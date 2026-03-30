package sqlite

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"game-catalog-backend/internal/game"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string, seedDemoData bool) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}

	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	if seedDemoData {
		if err := store.seed(); err != nil {
			db.Close()
			return nil, err
		}
	}

	return store, nil
}

func (s *Store) initSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			genre TEXT NOT NULL,
			platform TEXT NOT NULL,
			release_year INTEGER NOT NULL,
			rating INTEGER NOT NULL,
			status TEXT NOT NULL,
			image_path TEXT NOT NULL DEFAULT ''
		)
	`)
	if err != nil {
		return err
	}

	return s.ensureImagePathColumn()
}

func (s *Store) ensureImagePathColumn() error {
	rows, err := s.db.Query(`PRAGMA table_info(games)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var hasImagePath bool
	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal sql.NullString
			primaryKey int
		)

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &primaryKey); err != nil {
			return err
		}

		if name == "image_path" {
			hasImagePath = true
			break
		}
	}

	if hasImagePath {
		return nil
	}

	_, err = s.db.Exec(`ALTER TABLE games ADD COLUMN image_path TEXT NOT NULL DEFAULT ''`)
	return err
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
	result, err := s.db.Exec(`
		INSERT INTO games (title, genre, platform, release_year, rating, status, image_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, entity.Title, entity.Genre, entity.Platform, entity.ReleaseYear, entity.Rating, entity.Status, entity.ImagePath)
	if err != nil {
		return game.Game{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return game.Game{}, err
	}

	return s.GetByID(id)
}

func (s *Store) List(filters game.Filters) ([]game.Game, error) {
	rows, err := s.db.Query(`
		SELECT id, title, genre, platform, release_year, rating, status, image_path
		FROM games
		WHERE (? = '' OR status = ?)
		  AND (? = '' OR genre = ?)
		ORDER BY id DESC
	`, filters.Status, filters.Status, filters.Genre, filters.Genre)
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
		SELECT id, title, genre, platform, release_year, rating, status, image_path
		FROM games
		WHERE id = ?
	`, id).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Genre,
		&entity.Platform,
		&entity.ReleaseYear,
		&entity.Rating,
		&entity.Status,
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
	result, err := s.db.Exec(`
		UPDATE games
		SET title = ?, genre = ?, platform = ?, release_year = ?, rating = ?, status = ?, image_path = ?
		WHERE id = ?
	`, entity.Title, entity.Genre, entity.Platform, entity.ReleaseYear, entity.Rating, entity.Status, entity.ImagePath, id)
	if err != nil {
		return game.Game{}, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return game.Game{}, err
	}

	if affected == 0 {
		return game.Game{}, game.ErrNotFound
	}

	return s.GetByID(id)
}

func (s *Store) Delete(id int64) error {
	result, err := s.db.Exec(`DELETE FROM games WHERE id = ?`, id)
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
