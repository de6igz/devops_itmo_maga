package game

type Game struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Platform    string `json:"platform"`
	ReleaseYear int    `json:"releaseYear"`
	Rating      int    `json:"rating"`
	Status      string `json:"status"`
	ImagePath   string `json:"imagePath"`
}

type Filters struct {
	Status string
	Genre  string
}

var SeedGames = []Game{
	{
		Title:       "The Witcher 3",
		Genre:       "RPG",
		Platform:    "PC",
		ReleaseYear: 2015,
		Rating:      10,
		Status:      "completed",
		ImagePath:   "",
	},
	{
		Title:       "Hades",
		Genre:       "Roguelike",
		Platform:    "Switch",
		ReleaseYear: 2020,
		Rating:      9,
		Status:      "playing",
		ImagePath:   "",
	},
	{
		Title:       "Alan Wake 2",
		Genre:       "Action",
		Platform:    "PS5",
		ReleaseYear: 2023,
		Rating:      8,
		Status:      "planned",
		ImagePath:   "",
	},
}
