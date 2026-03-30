import type { Game } from "../types";
import { resolveAssetUrl } from "../api";

interface GameDetailsProps {
  game: Game | null;
  relatedGames: Game[];
  loading: boolean;
  onOpen: (id: number) => void;
  onClose: () => void;
}

function GameDetails({ game, relatedGames, loading, onOpen, onClose }: GameDetailsProps) {
  const heroBackground = game?.imagePath ? resolveAssetUrl(game.imagePath) : "";

  return (
    <section className="showcase-shell" data-testid="game-details-panel">
      <article
        className={`showcase-main${heroBackground ? " showcase-has-image" : ""}`}
        style={heroBackground ? { backgroundImage: `linear-gradient(90deg, rgba(7, 12, 14, 0.82), rgba(7, 12, 14, 0.28)), url(${heroBackground})` } : undefined}
      >
        <div className="showcase-nav">
          <div className="genre-pill active">Genre</div>
          <div className="genre-pill">Action</div>
          <div className="genre-pill">Adventure</div>
          <div className="genre-pill">Strategy</div>
        </div>

        <div className="showcase-content">
          {loading ? (
            <p className="panel-text">Загрузка карточки...</p>
          ) : !game ? (
            <div>
              <p className="section-kicker">Featured</p>
              <h2>Выберите игру</h2>
              <p className="panel-text">Откройте запись из коллекции, чтобы показать её как главный экран.</p>
            </div>
          ) : (
            <>
              <p className="price-tag">Score {game.rating}/10</p>
              <h2>{game.title}</h2>
              <div className="showcase-line" />
              <div className="showcase-meta">
                <span>{game.platform}</span>
                <span>{game.status}</span>
                <span>{game.releaseYear}</span>
              </div>
              <p className="showcase-description">
                Жанр: {game.genre}. Это главная карточка каталога с крупной обложкой, статусом и быстрой
                навигацией для демонстрации на защите.
              </p>
              <div className="showcase-actions">
                <button type="button" className="primary-cta" onClick={onClose}>
                  Скрыть фокус
                </button>
                {game.imagePath && (
                  <a className="ghost-button inline-link" href={resolveAssetUrl(game.imagePath)} target="_blank" rel="noreferrer">
                    Открыть постер
                  </a>
                )}
              </div>
            </>
          )}
        </div>
      </article>

      <aside className="showcase-side panel">
        <div className="panel-header">
          <div>
            <p className="section-kicker">Scenes</p>
            <h2>Быстрый выбор</h2>
          </div>
        </div>

        <div className="scene-list">
          {relatedGames.length === 0 ? (
            <p className="panel-text">Добавьте ещё несколько игр, чтобы заполнить сцену справа.</p>
          ) : (
            relatedGames.map((item) => (
              <button key={item.id} type="button" className="scene-card" onClick={() => onOpen(item.id)}>
                {item.imagePath ? (
                  <img src={resolveAssetUrl(item.imagePath)} alt={item.title} />
                ) : (
                  <div className="scene-fallback">{item.title}</div>
                )}
                <span>{item.title}</span>
              </button>
            ))
          )}
        </div>
      </aside>
    </section>
  );
}

export default GameDetails;
