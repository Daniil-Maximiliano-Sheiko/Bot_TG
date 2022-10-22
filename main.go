package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	go UpdateLoop()
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	http.ListenAndServe("localhost:8080", router)
	//...
}

type MainStruct struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	Id       int    `ison:"id"`
	Is_bot   bool   `json"id_bot"`
	Name     string `json:"first_name"`
	Username string `json:"username"`
	Join     bool   `json:"can_join_groups"`
	Read     bool   `json:"can_read_all_group_messages"`
	Support  bool   `json:supports_inline_queries"`
}

type Updateresponse struct {
	Ok     bool           `json:"ok"`
	Result []UpdateStruct `json:"result"`
}

type User struct {
	Id       int    `json:"id"`
	Is_bot   bool   `json:"is_bot"`
	Username string `json:"username"`
	IsPrem   bool   `json:"is_prem"`
}

type Chat struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type Message struct {
	Id   int    `json:"message_id"`
	User User   `json:"from"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type SendMessage struct {
	ChId             int    `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageId int    `json:"reply_to_message_id"`
	ProtectContent   bool   `json:"protect_content"`
}

type UpdateStruct struct {
	Id                int     `json:"update_id"`
	Message           Message `json:"message"`
	EditedMessage     Message `json:"edited_message"`
	ChannelPost       Message `json:"channel_post"`
	EditedChannelPost Message `json:"edited_channel_post"`
}

const apiUrl = "https://api.telegram.org/" + "bot5504725655:AAHXPRXyT51v9bCRrrvAAdQRVZrBlNu5O2Y"

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var R MainStruct

	Ping()

	resp, err := http.Get(apiUrl + "/getMe")

	if err != nil {
		fmt.Println(err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	err = json.Unmarshal(respBody, &R) // заполнили перемнную р
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	respReady, err := json.Marshal(R.Result)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(respReady))

	println("Данные прочитаны")

	w.Write([]byte("Вывод успешно произведён!"))
}

var appeal = ("Мой бот")

func UpdateLoop() {
	lastId := 0
	for {
		lastId = Update(lastId)
		time.Sleep(1 * time.Second)
	}
}

func Update(lastId int) int {
	raw, err := http.Get(apiUrl + "/getUpdates?offset=" + strconv.Itoa(lastId))
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(raw.Body)

	var v Updateresponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		panic(err)
	}
	if len(v.Result) > 0 {
		ev := v.Result[len(v.Result)-1]
		txt := strings.ToLower(ev.Message.Text)
		if txt == "/privet" {
			txtmsg := SendMessage{
				ChId: ev.Message.Chat.Id,
				Text: "Хай и "  + strconv.Itoa(ev.Message.Chat.Id),
			}

			bytemsg, _ := json.Marshal(txtmsg)
			_, err = http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
			if err != nil {
				fmt.Println(err)
				return lastId
			} else {
				return ev.Id + 1
			}

		}

		if strings.Split(txt, ", ")[0] == appeal {
			switch strings.Split(strings.Split(txt, ", ")[1], ": ")[0] {
			case "Привет":
				{
					return Anek(lastId, ev)
				}
			case "Ответ":
				{
					return Otvet(lastId, ev)
				}
			}
		}
	}
	return lastId
}

func Anek(lastId int, ev UpdateStruct) int {
	txtmsg := SendMessage{
		ChId:           ev.Message.Chat.Id,
		Text:           "Привет, вот твой ID: " + strconv.Itoa(ev.Message.Chat.Id),
		ProtectContent: true,
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	} else {
		return ev.Id + 1
	}
}

func Otvet(lastId int, ev UpdateStruct) int {
	txtmsg := SendMessage{
		ChId:             ev.Message.Chat.Id,
		Text:             "Как ответить?",
		ReplyToMessageId: ev.Message.Id,
		ProtectContent:   true,
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
		return lastId
	} else {
		return ev.Id + 1
	}
}

// func ChangeName(lastId int, ev UpdateStruct, txt string) int {
// 	newap := strings.Split(txt, "измени обращение на: ")
// 	appeal = newap[1]
// 	fmt.Println(appeal)
// 	txtmsg := SendMessage{
// 		ChId:           ev.Message.Chat.Id,
// 		Text:           "Обращение изменено на: " + appeal,
// 		ProtectContent: true,
// 	}

// 	bytemsg, _ := json.Marshal(txtmsg)
// 	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))

// 	if err != nil {
// 		fmt.Println(err)
// 		return lastId
// 	} else {
// 		return ev.Id + 1
// 	}
// }

func Ping() {
	txtmsg := SendMessage{
		ChId: 520669485,
		Text: "Страница посещена",
	}

	bytemsg, _ := json.Marshal(txtmsg)
	_, err := http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	if err != nil {
		fmt.Println(err)
	}
}
