import type { GameFormData, GameStatus } from "./types";

export const STATUS_OPTIONS: Array<{ value: GameStatus; label: GameStatus }> = [
  { value: "planned", label: "planned" },
  { value: "playing", label: "playing" },
  { value: "completed", label: "completed" },
];

export const DEFAULT_FORM_STATE: GameFormData = {
  title: "",
  genre: "",
  platform: "",
  releaseYear: "",
  rating: "",
  status: "planned",
  description: "",
  imagePath: "",
};
