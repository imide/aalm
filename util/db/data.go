package db

import "time"

type Player struct {
	DiscordID         string    `bson:"discordId"`
	TeamPlaying       string    `bson:"teamPlaying"`
	Stars             float32   `bson:"stars"`
	Position          string    `bson:"position"`
	SeasonsPlayed     int       `bson:"seasonsPlayed"`
	Contracted        bool      `bson:"contracted"`
	SeasonsContracted int       `bson:"seasonsContracted"`
	IsSuspended       bool      `bson:"isSuspended"`
	SuspensionExpires time.Time `bson:"suspensionExpires"`
}

type Team struct {
	Name           string   `bson:"name"`
	Logo           string   `bson:"logo"`
	TeamOwner      string   `bson:"teamOwner"`
	Coach          []string `bson:"coach"`
	Players        []string `bson:"players"`
	PlayerMax      int      `bson:"playerMax"`
	RoleID         string   `bson:"roleId"`
	Wins           uint     `bson:"wins"`
	Losses         uint     `bson:"losses"`
	StarsRecruited float32  `bson:"starsRecruited"`
	MaxStars       float32  `bson:"maxStars"`
}
