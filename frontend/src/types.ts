export type GameStatus = "planned" | "playing" | "completed";

export interface Game {
  id: number;
  title: string;
  genre: string;
  platform: string;
  releaseYear: number;
  rating: number;
  status: GameStatus;
  imagePath: string;
}

export interface GameFormData {
  title: string;
  genre: string;
  platform: string;
  releaseYear: string;
  rating: string;
  status: GameStatus;
  imagePath: string;
}
