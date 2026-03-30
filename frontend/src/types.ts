export type GameStatus = "planned" | "playing" | "completed";

export interface Game {
  id: number;
  title: string;
  genre: string;
  platform: string;
  releaseYear: number;
  rating: number;
  status: GameStatus;
  description: string;
  imagePath: string;
}

export interface GameFormData {
  title: string;
  genre: string;
  platform: string;
  releaseYear: string;
  rating: string;
  status: GameStatus;
  description: string;
  imagePath: string;
}
