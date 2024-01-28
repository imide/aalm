package db

import "time"

type Player struct {
	ID          string           `bson:"_id"`
	TeamPlaying string           `bson:"team_id"`
	Stars       float32          `bson:"stars"`
	Contracted  bool             `bson:"contracted"`
	Suspension  *Suspension      `bson:"suspension"`
	Rings       []PlayerRingData `bson:"rings"`
}

type Team struct {
	ID            string       `bson:"_id"`
	Name          string       `bson:"name"`
	Logo          string       `bson:"logo"`
	EmojiID       string       `bson:"emojiId"`
	Owner         string       `bson:"teamOwner"`
	Coaches       []string     `bson:"coach"`
	Players       []PlayerInfo `bson:"players"`
	PlayerMax     int          `bson:"playerMax"`
	RoleID        string       `bson:"roleId"`
	Wins          uint         `bson:"wins"`
	Losses        uint         `bson:"losses"`
	MaxStars      float32      `bson:"maxStars"`
	DiscordInvite string       `bson:"discordInvite"`
}

// TODO: implement discord invites for teams in teams.go  and wherever else it's needed (i forgort lol)
type Ring struct {
	ID     string `bson:"_id"`
	RoleID string `bson:"roleId"`
	Name   string `bson:"name"`
	Desc   string `bson:"desc"`
}

type PlayerRingData struct {
	RingID      string `bson:"ring_id"`
	DateAwarded string `bson:"date_awarded"`
}

type Suspension struct {
	WhoSuspended string    `bson:"who_suspended"`
	Reason       string    `bson:"reason"`
	Started      time.Time `bson:"started"`
	Ends         int64     `bson:"ends"`
	Active       bool      `bson:"active"`
}

type PlayerInfo struct {
	ID    string  `bson:"player_id"`
	Stars float32 `bson:"stars"`
}
