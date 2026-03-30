import { useEffect, useMemo, useState } from "react";
import type { ChangeEvent, FormEvent } from "react";

import {
  createGame,
  deleteGame,
  formDataToPayload,
  getGame,
  getGames,
  resolveAssetUrl,
  updateGame,
  uploadGameImage,
} from "./api";
import GameDetails from "./components/GameDetails";
import GameForm from "./components/GameForm";
import GameList from "./components/GameList";
import { DEFAULT_FORM_STATE, STATUS_OPTIONS } from "./constants";
import type { Game, GameFormData, GameStatus } from "./types";

function App() {
  const [games, setGames] = useState<Game[]>([]);
  const [selectedGame, setSelectedGame] = useState<Game | null>(null);
  const [formData, setFormData] = useState<GameFormData>(DEFAULT_FORM_STATE);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [loadingList, setLoadingList] = useState(true);
  const [loadingDetails, setLoadingDetails] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [selectedImageFile, setSelectedImageFile] = useState<File | null>(null);
  const [imagePreviewUrl, setImagePreviewUrl] = useState("");
  const [filters, setFilters] = useState<{ status: "" | GameStatus; genre: string }>({
    status: "",
    genre: "",
  });

  async function loadGames(nextFilters = filters) {
    setLoadingList(true);
    setErrorMessage("");

    try {
      const gamesList = await getGames(nextFilters);
      setGames(gamesList);

      if (gamesList.length === 0) {
        setSelectedGame(null);
        return;
      }

      setSelectedGame((current) => gamesList.find((item) => item.id === current?.id) || gamesList[0]);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : "Ошибка загрузки.");
    } finally {
      setLoadingList(false);
    }
  }

  useEffect(() => {
    void loadGames(filters);
  }, []);

  useEffect(() => {
    return () => {
      if (imagePreviewUrl.startsWith("blob:")) {
        URL.revokeObjectURL(imagePreviewUrl);
      }
    };
  }, [imagePreviewUrl]);

  const displayedGame = selectedGame ?? games[0] ?? null;
  const relatedGames = useMemo(
    () => games.filter((game) => game.id !== displayedGame?.id).slice(0, 3),
    [games, displayedGame?.id],
  );

  const featuredGenres = useMemo(() => {
    const uniqueGenres = Array.from(new Set(games.map((game) => game.genre))).slice(0, 4);
    return uniqueGenres.length > 0 ? uniqueGenres : ["Action", "Adventure", "Strategy"];
  }, [games]);

  async function handleOpenGame(id: number) {
    setLoadingDetails(true);
    setErrorMessage("");

    try {
      const game = await getGame(id);
      setSelectedGame(game);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : "Ошибка загрузки.");
    } finally {
      setLoadingDetails(false);
    }
  }

  function handleInputChange(event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) {
    const { name, value } = event.target;
    setFormData((currentState) => ({
      ...currentState,
      [name]: value,
    }));
  }

  function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0] ?? null;
    setSelectedImageFile(file);

    if (imagePreviewUrl.startsWith("blob:")) {
      URL.revokeObjectURL(imagePreviewUrl);
    }

    if (file) {
      setImagePreviewUrl(URL.createObjectURL(file));
      return;
    }

    setImagePreviewUrl(formData.imagePath ? resolveAssetUrl(formData.imagePath) : "");
  }

  function resetForm(nextState: GameFormData = DEFAULT_FORM_STATE) {
    setFormData(nextState);
    setEditingId(null);
    setSelectedImageFile(null);

    if (imagePreviewUrl.startsWith("blob:")) {
      URL.revokeObjectURL(imagePreviewUrl);
    }

    setImagePreviewUrl(nextState.imagePath ? resolveAssetUrl(nextState.imagePath) : "");
  }

  function handleEditGame(game: Game) {
    setEditingId(game.id);
    setSelectedImageFile(null);
    setFormData({
      title: game.title,
      genre: game.genre,
      platform: game.platform,
      releaseYear: String(game.releaseYear),
      rating: String(game.rating),
      status: game.status,
      imagePath: game.imagePath,
    });
    setImagePreviewUrl(game.imagePath ? resolveAssetUrl(game.imagePath) : "");
  }

  async function handleDeleteGame(id: number) {
    setErrorMessage("");

    try {
      await deleteGame(id);

      if (selectedGame?.id === id) {
        setSelectedGame(null);
      }

      if (editingId === id) {
        resetForm();
      }

      await loadGames(filters);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : "Ошибка удаления.");
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setErrorMessage("");

    try {
      let nextImagePath = formData.imagePath;

      if (selectedImageFile) {
        nextImagePath = await uploadGameImage(selectedImageFile);
      }

      const nextFormData = {
        ...formData,
        imagePath: nextImagePath,
      };

      if (editingId) {
        const updatedGame = await updateGame(editingId, formDataToPayload(nextFormData));
        setSelectedGame(updatedGame);
      } else {
        const createdGame = await createGame(formDataToPayload(nextFormData));
        setSelectedGame(createdGame);
      }

      resetForm();
      await loadGames(filters);
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : "Ошибка сохранения.");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleFilterChange(event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) {
    const { name, value } = event.target;
    const nextFilters = {
      ...filters,
      [name]: value,
    } as { status: "" | GameStatus; genre: string };

    setFilters(nextFilters);
    await loadGames(nextFilters);
  }

  async function handleGenreShortcut(genre: string) {
    const nextFilters = {
      ...filters,
      genre: filters.genre === genre ? "" : genre,
    };

    setFilters(nextFilters);
    await loadGames(nextFilters);
  }

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <p className="eyebrow">Game Spotlight</p>
          <h1>Каталог видеоигр</h1>
        </div>
        <div className="topbar-actions">
          <button type="button" className="ghost-button">
            Search
          </button>
          <button type="button" className="ghost-button">
            Wishlist
          </button>
          <button type="button" className="primary-cta">
            Signup / Login
          </button>
        </div>
      </header>

      {errorMessage && <div className="error-banner">{errorMessage}</div>}

      <section className="genre-strip panel">
        <div className="genre-shortcuts">
          {featuredGenres.map((genre) => (
            <button
              key={genre}
              type="button"
              className={`genre-chip${filters.genre === genre ? " active" : ""}`}
              onClick={() => void handleGenreShortcut(genre)}
            >
              {genre}
            </button>
          ))}
        </div>

        <div className="filter-grid compact">
          <label>
            Статус
            <select aria-label="Фильтр по статусу" name="status" value={filters.status} onChange={handleFilterChange}>
              <option value="">Все</option>
              {STATUS_OPTIONS.map((status) => (
                <option key={status.value} value={status.value}>
                  {status.label}
                </option>
              ))}
            </select>
          </label>

          <label>
            Жанр
            <input
              aria-label="Фильтр по жанру"
              name="genre"
              placeholder="Например, RPG"
              value={filters.genre}
              onChange={handleFilterChange}
            />
          </label>
        </div>
      </section>

      <GameDetails
        game={displayedGame}
        relatedGames={relatedGames}
        loading={loadingDetails}
        onOpen={handleOpenGame}
        onClose={() => setSelectedGame(null)}
      />

      <main className="dashboard-grid">
        <GameForm
          formData={formData}
          isEditing={editingId !== null}
          isSubmitting={submitting}
          imagePreviewUrl={imagePreviewUrl}
          onChange={handleInputChange}
          onFileChange={handleFileChange}
          onSubmit={handleSubmit}
          onCancel={resetForm}
        />

        <GameList
          games={games}
          loading={loadingList}
          onOpen={handleOpenGame}
          onEdit={handleEditGame}
          onDelete={handleDeleteGame}
        />
      </main>
    </div>
  );
}

export default App;
