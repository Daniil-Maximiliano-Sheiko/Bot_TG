package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
	Ok     bool          `json:"ok"`
	Result []UpdateStuct `json:"result"`
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

type UpdateStuct struct {
	Id                int     `json:"update_id"`
	Message           Message `json:"message"`
	EditedMessage     Message `json:"edited_message"`
	ChannelPost       Message `json:"channel_post"`
	EditedChannelPost Message `json:"edited_channel_post"`
}

const apiUrl = "https://api.telegram.org/" + "bot5504725655:AAHXPRXyT51v9bCRrrvAAdQRVZrBlNu5O2Y"

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	var R MainStruct

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

var appeal = "БОООт"

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

	// if len(v.Result) > 0 {
	// 	ev := v.Result[len(v.Result)-1]
	// 	// for _, ev := range v.Result {
	// 	txt := ev.Message.Text
	// 	if txt == "Привет" {
	// 		txtmsg := SendMessage{
	// 			ChId: ev.Message.Chat.Id,
	// 			Text: "И тебя туда же",
	// 			ProtectContent: true,
	// 		}

	// 		bytemsg, _ := json.Marshal(txtmsg)
	// 		_, err = http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return lastId
	// 		} else {
	// 			return ev.Id + 1
	// 		}
	// 	}
	// }
	if len(v.Result) > 0 {
		ev := v.Result[len(v.Result)-1]
		txt := ev.Message.Text
		if txt == "Ответ" {
			txtmsg := SendMessage{
				ChId:             ev.Message.Chat.Id,
				Text:             "Как ответить?",
				ReplyToMessageId: ev.Message.Id,
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
	}
	if len(v.Result) > 0 {
		ev := v.Result[len(v.Result)-1]
		txt := ev.Message.Text
		if txt == "Как дела?" {
			txtmsg := SendMessage{
				ChId:             ev.Message.Chat.Id,
				Text:             "Нормально",
				ReplyToMessageId: ev.Message.Id,
				ProtectContent: true,
			}

			bytemsg, _ := json.Marshal(txtmsg)
			_, err = http.Post(apiUrl+"/sendMessage", "application/json", bytes.NewReader(bytemsg))
			if err != nil {
				fmt.Println(err)
				return lastId
			} else {
				return ev.Id + 1
			}
			if strings.Split(txt, ", ")[0] == appeal {

				switch strings.Split(strings.Split(txt, ", ")[1], ": ")[0] {
				case "Привет":
					{
						return Privet(lastId, ev)
					}
				case "сгенерируй число":
					{
						return RandGen(lastId, ev, txt)
					}
				case "измени обращение на":
					{
						if strings.Contains(txt, ": ") {
							return ChangeName(lastId, ev, txt)
						} else {
							fmt.Println("error")
						}
					}
				}
	
			}
		}
		return lastId
	}

func Privet(lastId int, ev UpdateStruct) int {
		txtmsg := SendMessage{
			ChId: ev.Message.Chat.Id,
			Text: "Привет",
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
