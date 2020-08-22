package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &webhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	fmt.Println(body.Message.Text)

	// If the text contains marco, call the `startGame` function, which
	// is defined below
	if err := startGame(body.Message.Chat.ID); err != nil {
		fmt.Println("error in starting a game:", err)
		return
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

// The below code deals with the process of sending a response message
// to the user

// Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type SendMessageReqBody struct {
	ChatID      int64                `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

// startGame takes a chatID and sends "polo" to them
func startGame(chatID int64) error {

	keys := [3]string{"Name", "Language", "API"}
	values := [3]string{"John", "Python", "pyTelegramBotAPI"}
	markup := InlineKeyboardMarkup{}
	for i := 0; i < 3; i++ {
		markup.InlineKeyboard[i][0] = InlineKeyboardButton{Text: keys[i]}
		markup.InlineKeyboard[i][1] = InlineKeyboardButton{Text: values[i]}
	}

	// Create the request body struct
	reqBody := &SendMessageReqBody{
		ChatID:      chatID,
		Text:        "Polo!!",
		ReplyMarkup: markup,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot"+os.Getenv("BOT_TOKEN")+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

// main funtion starts our server on a port
func main() {
	http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(Handler))
}
