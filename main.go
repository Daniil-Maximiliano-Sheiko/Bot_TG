package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

func main() {
	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"sqlite3_mod_regexp",
			},
		})

	go UpdateLoop()
	router := mux.NewRouter()
	router.HandleFunc("/api", IndexHandler)
	router.HandleFunc("/botname", NameHandler)
	router.HandleFunc("/eventId", EvIdHandler)
	router.HandleFunc("/lastId", LastIdHandler)
	router.HandleFunc("/login", LogHandler)
	router.HandleFunc("/register", RegHandler)
	router.PathPrefix("/").Handler(http.FileServer((http.Dir("./static/"))))
	http.ListenAndServe("localhost:8000", router)
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

func NameHandler(w http.ResponseWriter, _ *http.Request) {
    db, err := sql.Open("sqlite3", "file:database.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    var gotname string
    var resp sql.NullString // для результата
    err = db.QueryRow("SELECT name FROM bot_status").Scan(&resp)
    if err != nil {
        fmt.Println(err)
    }
    if resp.Valid { // если результат валид
        gotname = resp.String // берём оттуда обычный string
    }
    w.Write([]byte(gotname))
}

func EvIdHandler(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "db.sql")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    var goteventid string
    var resp sql.NullString // для результата
    err = db.QueryRow("SELECT id FROM bot_status").Scan(&resp)
    if err != nil {
        fmt.Println(err)
    }
    if resp.Valid { // если результат валид
        goteventid = resp.String // берём оттуда обычный string
    }
    w.Write([]byte(goteventid))
}

func LastIdHandler(w http.ResponseWriter, _ *http.Request) {
	db, err := sql.Open("sqlite3", "db.sql")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    var gotlastid string
    var resp sql.NullString // для результата
    err = db.QueryRow("SELECT lastid FROM bot_status").Scan(&resp)
    if err != nil {
        fmt.Println(err)
    }
    if resp.Valid { // если результат валид
        gotlastid = resp.String // берём оттуда обычный string
    }
    w.Write([]byte(gotlastid))
}

func AuthCheck(w http.ResponseWriter, _ *http.Request) {}

func LogHandler(w http.ResponseWriter, _ *http.Request) {}

func RegHandler(w http.ResponseWriter, _ *http.Request) {}

var appeal = ("GO")

func UpdateLoop() {
    db, err := sql.Open("sqlite3", "tgbot.db")
    if err != nil {
        panic(err)
    }
    defer db.Close() //закрывает коннект при закрытии программы
    lastId := -1
    for {
        newId := Update(lastId)
        if lastId != newId {
            lastId = newId
            db.Exec(`UPDATE `) // новый lastid в таблицу bot_status
        }
        time.Sleep(50 * time.Millisecond)
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
			case "привет":
				{
					return Anek(lastId, ev)
				}
			case "ответ":
				{
					return Otvet(lastId, ev)
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

func ChangeName(lastId int, ev UpdateStruct, txt string) int {
	newap := strings.Split(txt, "измени обращение на: ")
	appeal = newap[1]
	fmt.Println(appeal)
	txtmsg := SendMessage{
		ChId:           ev.Message.Chat.Id,
		Text:           "Обращение изменено на: " + appeal,
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
