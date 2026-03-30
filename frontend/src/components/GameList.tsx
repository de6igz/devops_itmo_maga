import type { Game } from "../types";
import { resolveAssetUrl } from "../api";

interface GameListProps {
  games: Game[];
  loading: boolean;
  onOpen: (id: number) => void;
  onEdit: (game: Game) => void;
  onDelete: (id: number) => void;
}

function GameList({ games, loading, onOpen, onEdit, onDelete }: GameListProps) {
  return (
    <section className="panel library-panel">
      <div className="panel-header">
        <div>
          <p className="section-kicker">Library</p>
          <h2>Коллекция</h2>
          <p className="panel-text">Каталог для просмотра, редактирования и удаления записей.</p>
        </div>
      </div>

      {loading ? (
        <p className="panel-text">Загрузка...</p>
      ) : games.length === 0 ? (
        <p className="panel-text">Игры не найдены. Добавьте первую запись через форму.</p>
      ) : (
        <div className="game-list">
          {games.map((game) => (
            <article className="game-card" key={game.id} onClick={() => onOpen(game.id)}>
              <div className="card-cover">
                {game.imagePath ? (
                  <img src={resolveAssetUrl(game.imagePath)} alt={game.title} />
                ) : (
                  <div className="cover-fallback">{game.title.slice(0, 2).toUpperCase()}</div>
                )}
              </div>

              <div className="card-copy">
                <div>
                  <h3>{game.title}</h3>
                  <p>{game.genre}</p>
                </div>

                <div className="game-meta">
                  <span>{game.platform}</span>
                  <span>{game.releaseYear}</span>
                  <span>Оценка: {game.rating}</span>
                  <span className={`status-badge status-${game.status}`}>{game.status}</span>
                </div>
              </div>

              <div className="card-actions">
                <button type="button" className="ghost-button" onClick={() => onOpen(game.id)}>
                  Открыть
                </button>
                <button type="button" className="ghost-button" onClick={(event) => {
                  event.stopPropagation();
                  onEdit(game);
                }}>
                  Редактировать
                </button>
                <button type="button" className="danger-button" onClick={(event) => {
                  event.stopPropagation();
                  onDelete(game.id);
                }}>
                  Удалить
                </button>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}

export default GameList;
