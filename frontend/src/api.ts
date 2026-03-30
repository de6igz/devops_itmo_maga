import type { Game, GameFormData, GameStatus } from "./types";

const API_BASE = import.meta.env.VITE_API_URL || "";

export interface GameFilters {
  status: "" | GameStatus;
  genre: string;
}

type GamePayload = Omit<Game, "id">;

async function apiRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({ error: "Ошибка запроса." }));
    throw new Error(errorBody.error || "Ошибка запроса.");
  }

  if (response.status === 204) {
    return null as T;
  }

  return response.json() as Promise<T>;
}

export function formDataToPayload(formData: GameFormData): GamePayload {
  return {
    title: formData.title,
    genre: formData.genre,
    platform: formData.platform,
    releaseYear: Number(formData.releaseYear),
    rating: Number(formData.rating),
    status: formData.status,
    description: formData.description,
    imagePath: formData.imagePath,
  };
}

export function resolveAssetUrl(imagePath: string): string {
  if (!imagePath) {
    return "";
  }

  return `${API_BASE}${imagePath}`;
}

export function getGames(filters: GameFilters): Promise<Game[]> {
  const query = new URLSearchParams();

  if (filters.status) {
    query.set("status", filters.status);
  }

  if (filters.genre) {
    query.set("genre", filters.genre);
  }

  const queryString = query.toString();
  const path = queryString ? `/api/games?${queryString}` : "/api/games";

  return apiRequest<Game[]>(path);
}

export function getGame(id: number): Promise<Game> {
  return apiRequest<Game>(`/api/games/${id}`);
}

export function createGame(payload: GamePayload): Promise<Game> {
  return apiRequest<Game>("/api/games", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateGame(id: number, payload: GamePayload): Promise<Game> {
  return apiRequest<Game>(`/api/games/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export function deleteGame(id: number): Promise<null> {
  return apiRequest<null>(`/api/games/${id}`, {
    method: "DELETE",
  });
}

export async function uploadGameImage(file: File): Promise<string> {
  const body = new FormData();
  body.append("image", file);

  const response = await fetch(`${API_BASE}/api/uploads/image`, {
    method: "POST",
    body,
  });

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({ error: "Ошибка загрузки изображения." }));
    throw new Error(errorBody.error || "Ошибка загрузки изображения.");
  }

  const payload = (await response.json()) as { path: string };
  return payload.path;
}
