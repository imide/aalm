package bloxlink

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handleBloxlinkError(err error) (string, error) {
	log.Println("Bloxlink connection error: ", err)
	return "Bloxlink Connection Failed", err
}

func GetRobloxId(userId string) (string, error) {
	url := fmt.Sprintf("%s/%s/discord-to-roblox/%s",
		"https://api.blox.link/v4/public/guilds",
		os.Getenv("GUILD_ID"), userId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return handleBloxlinkError(err)
	}

	req.Header.Add("Authorization", os.Getenv("BLOXLINK_TOKEN"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handleBloxlinkError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleBloxlinkError(err)
	}

	return string(body), nil
}

//todo: whatever else is needed cause i dont remember lol
