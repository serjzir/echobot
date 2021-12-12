package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// точка входа
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	botToken := os.Getenv("API-TOKEN")
	// https://api.telegram.org/bot<token>/METHOD_NAME
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offsset := 0
	for ;; {
		updates, err := getUpdates(botUrl, offsset)
		if err != nil {
			log.Printf("something wet wrong: %v", err)
		}
		for _, update := range updates {
			err = respond(botUrl, update)
			offsset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

// get update
func getUpdates(botUrl string, offsset int) ([]Update, error){
	res, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offsset))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// request for update
func respond(botUrl string, update Update) (error){
	var botMsg BotMsg
	botMsg.ChatId = update.Message.Chat.ChatId
	botMsg.Text = update.Message.Text
	buf, err := json.Marshal(botMsg)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}