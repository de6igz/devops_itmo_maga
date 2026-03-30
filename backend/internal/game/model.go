package game

type Game struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Platform    string `json:"platform"`
	ReleaseYear int    `json:"releaseYear"`
	Rating      int    `json:"rating"`
	Status      string `json:"status"`
	Description string `json:"description"`
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
		Description: "Глубокая сюжетная RPG в мрачном фэнтези-мире с сильными побочными квестами, исследованием открытого мира и запоминающимися персонажами.",
		ImagePath:   "",
	},
	{
		Title:       "Hades",
		Genre:       "Roguelike",
		Platform:    "Switch",
		ReleaseYear: 2020,
		Rating:      9,
		Status:      "playing",
		Description: "Динамичный рогалик про побег из подземного царства, где быстрые бои, стильный арт и постоянное развитие героя работают вместе.",
		ImagePath:   "",
	},
	{
		Title:       "Alan Wake 2",
		Genre:       "Action",
		Platform:    "PS5",
		ReleaseYear: 2023,
		Rating:      8,
		Status:      "planned",
		Description: "Психологический survival horror с густой атмосферой, кинематографичной подачей и мрачным расследованием в духе триллера.",
		ImagePath:   "",
	},
}
