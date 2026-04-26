package postgres

import (
	"database/sql"
	"errors"

	"game-catalog-backend/internal/game"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New(databaseURL string, seedDemoData bool) (*Store, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	store := &Store{db: db}

	if err := store.runMigrations(); err != nil {
		_ = db.Close()
		return nil, err
	}

	if seedDemoData {
		if err := store.seed(); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	return store, nil
}

func (s *Store) seed() error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM games`).Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	for _, seededGame := range game.SeedGames {
		if _, err := s.Create(seededGame); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Create(entity game.Game) (game.Game, error) {
	var created game.Game
	err := s.db.QueryRow(`
		INSERT INTO games (title, genre, platform, release_year, rating, status, description, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, genre, platform, release_year, rating, status, description, image_path
	`, entity.Title, entity.Genre, entity.Platform, entity.ReleaseYear, entity.Rating, entity.Status, entity.Description, entity.ImagePath).Scan(
		&created.ID,
		&created.Title,
		&created.Genre,
		&created.Platform,
		&created.ReleaseYear,
		&created.Rating,
		&created.Status,
		&created.Description,
		&created.ImagePath,
	)
	if err != nil {
		return game.Game{}, err
	}

	return created, nil
}

func (s *Store) List(filters game.Filters) ([]game.Game, error) {
	rows, err := s.db.Query(`
		SELECT id, title, genre, platform, release_year, rating, status, description, image_path
		FROM games
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR genre = $2)
		ORDER BY id DESC
	`, filters.Status, filters.Genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []game.Game
	for rows.Next() {
		var entity game.Game
		if err := rows.Scan(
			&entity.ID,
			&entity.Title,
			&entity.Genre,
			&entity.Platform,
			&entity.ReleaseYear,
			&entity.Rating,
			&entity.Status,
			&entity.Description,
			&entity.ImagePath,
		); err != nil {
			return nil, err
		}
		games = append(games, entity)
	}

	return games, rows.Err()
}

func (s *Store) GetByID(id int64) (game.Game, error) {
	var entity game.Game
	err := s.db.QueryRow(`
		SELECT id, title, genre, platform, release_year, rating, status, description, image_path
		FROM games
		WHERE id = $1
	`, id).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Genre,
		&entity.Platform,
		&entity.ReleaseYear,
		&entity.Rating,
		&entity.Status,
		&entity.Description,
		&entity.ImagePath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return game.Game{}, game.ErrNotFound
		}
		return game.Game{}, err
	}

	return entity, nil
}

func (s *Store) Update(id int64, entity game.Game) (game.Game, error) {
	var updated game.Game
	err := s.db.QueryRow(`
		UPDATE games
		SET title = $1, genre = $2, platform = $3, release_year = $4, rating = $5, status = $6, description = $7, image_path = $8
		WHERE id = $9
		RETURNING id, title, genre, platform, release_year, rating, status, description, image_path
	`, entity.Title, entity.Genre, entity.Platform, entity.ReleaseYear, entity.Rating, entity.Status, entity.Description, entity.ImagePath, id).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Genre,
		&updated.Platform,
		&updated.ReleaseYear,
		&updated.Rating,
		&updated.Status,
		&updated.Description,
		&updated.ImagePath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return game.Game{}, game.ErrNotFound
		}
		return game.Game{}, err
	}

	return updated, nil
}

func (s *Store) Delete(id int64) error {
	result, err := s.db.Exec(`DELETE FROM games WHERE id = $1`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return game.ErrNotFound
	}

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
