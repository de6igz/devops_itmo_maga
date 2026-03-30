import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import App from "./App";
import type { Game } from "./types";

function createJsonResponse(data: unknown, status = 200) {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => data,
  } as Response;
}

function createFetchMock(initialGames: Game[]) {
  let games = [...initialGames];
  let nextId = games.length + 1;

  return vi.fn(async (input: RequestInfo | URL, init: RequestInit = {}) => {
    const url = new URL(typeof input === "string" ? input : input.toString(), "http://localhost");
    const method = init.method || "GET";
    const pathname = url.pathname;

    if (pathname === "/api/games" && method === "GET") {
      let result = [...games];

      const status = url.searchParams.get("status");
      const genre = url.searchParams.get("genre");

      if (status) {
        result = result.filter((game) => game.status === status);
      }

      if (genre) {
        result = result.filter((game) => game.genre.toLowerCase() === genre.toLowerCase());
      }

      return createJsonResponse(result);
    }

    if (pathname === "/api/games" && method === "POST") {
      const payload = JSON.parse(String(init.body)) as Omit<Game, "id">;
      const newGame: Game = { id: nextId++, ...payload };
      games = [newGame, ...games];
      return createJsonResponse(newGame, 201);
    }

    if (pathname === "/api/uploads/image" && method === "POST") {
      return createJsonResponse({ path: "/blob/uploaded-cover.png" }, 201);
    }

    const idMatch = pathname.match(/^\/api\/games\/(\d+)$/);

    if (!idMatch) {
      return createJsonResponse({ error: "Маршрут не найден." }, 404);
    }

    const id = Number(idMatch[1]);
    const game = games.find((item) => item.id === id);

    if (!game) {
      return createJsonResponse({ error: "Игра не найдена." }, 404);
    }

    if (method === "GET") {
      return createJsonResponse(game);
    }

    if (method === "PUT") {
      const payload = JSON.parse(String(init.body)) as Omit<Game, "id">;
      const updatedGame: Game = { id, ...payload };
      games = games.map((item) => (item.id === id ? updatedGame : item));
      return createJsonResponse(updatedGame);
    }

    if (method === "DELETE") {
      games = games.filter((item) => item.id !== id);
      return createJsonResponse(null, 204);
    }

    return createJsonResponse({ error: "Метод не поддерживается." }, 405);
  });
}

describe("App", () => {
  beforeEach(() => {
    vi.stubGlobal(
      "fetch",
      createFetchMock([
        {
          id: 1,
          title: "The Witcher 3",
          genre: "RPG",
          platform: "PC",
          releaseYear: 2015,
          rating: 10,
          status: "completed",
          imagePath: "/blob/witcher.png",
        },
        {
          id: 2,
          title: "Hades",
          genre: "Roguelike",
          platform: "Switch",
          releaseYear: 2020,
          rating: 9,
          status: "playing",
          imagePath: "/blob/hades.png",
        },
      ]),
    );
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it("renders the games list", async () => {
    render(<App />);

    expect(await screen.findAllByText("The Witcher 3")).toHaveLength(2);
    expect(screen.getAllByText("Hades").length).toBeGreaterThan(0);
  });

  it("renders the add game form", async () => {
    render(<App />);

    expect(await screen.findByRole("heading", { name: "Новая игра" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Например, The Last of Us Part II")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Добавить игру" })).toBeInTheDocument();
  });

  it("adds a new game", async () => {
    const user = userEvent.setup();
    render(<App />);

    expect(await screen.findAllByText("The Witcher 3")).toHaveLength(2);

    await user.type(screen.getByPlaceholderText("Например, The Last of Us Part II"), "Celeste");
    await user.type(screen.getByPlaceholderText("Action, RPG, Strategy"), "Platformer");
    await user.type(screen.getByPlaceholderText("PC, PS5, Xbox, Switch"), "Switch");
    await user.type(screen.getByLabelText("Год выпуска"), "2018");
    await user.type(screen.getByLabelText("Оценка"), "9");
    await user.selectOptions(screen.getByLabelText("Статус игры"), "completed");
    await user.click(screen.getByRole("button", { name: "Добавить игру" }));

    expect(await screen.findAllByText("Celeste")).toHaveLength(2);
  });

  it("deletes a game", async () => {
    const user = userEvent.setup();
    render(<App />);

    expect(await screen.findAllByText("The Witcher 3")).toHaveLength(2);

    await user.click(screen.getAllByRole("button", { name: "Удалить" })[0]);

    await waitFor(() => {
      expect(screen.queryAllByText("The Witcher 3")).toHaveLength(0);
    });
  });

  it("opens the game card for a user scenario", async () => {
    const user = userEvent.setup();
    render(<App />);

    expect(await screen.findAllByText("The Witcher 3")).toHaveLength(2);
    await user.click(screen.getAllByRole("button", { name: "Открыть" })[0]);

    const detailsPanel = screen.getByTestId("game-details-panel");

    expect(within(detailsPanel).getByRole("heading", { name: "The Witcher 3" })).toBeInTheDocument();
    expect(within(detailsPanel).getByText("completed")).toBeInTheDocument();
    expect(within(detailsPanel).getByText("PC")).toBeInTheDocument();
  });
});
