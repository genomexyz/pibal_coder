package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/bson"
)

type userpass_check struct {
	username string
	password string
}

func main() {
	//state
	login_state := false
	username_login := ""
	display_name_login := ""

	//open db
	db, err := sql.Open("sqlite3", "./pibal.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println(db)

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:            "pibal_coder V0",
		BaseDirectoryPath:  "example",
		AppIconDefaultPath: "bg.png",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(600),
		Width:  astikit.IntPtr(600),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	//fmt.Println(w_mini)
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		//for feature akses admin
		if s == "akses_admin" {
			if login_state && username_login == "admin" {
				sending_data := bson.M{}
				sending_data["status"] = "login"
				sending_data["username"] = username_login
				sending_data["display_name"] = display_name_login

				// Perform a SELECT statement
				rows, err := db.Query("SELECT username, password, display_name FROM userpass")
				if err != nil {
					fmt.Println(err)
					sending_data["status_data"] = "err"
					json_send, err := json.Marshal(sending_data)
					if err != nil {
						log.Fatal(err)
					}
					json_send_str := string(json_send)
					return json_send_str
				}
				defer rows.Close()

				// Iterate over the rows returned by the SELECT statement
				cnt_row := 0
				for rows.Next() {
					sending_data["data-"+strconv.Itoa(cnt_row)] = bson.M{}
					var single_username string
					var single_password string
					var single_display_name string

					err = rows.Scan(&single_username, &single_password, &single_display_name)
					if err != nil {
						fmt.Println("got error in row: ", err)
						continue
					}
					sending_data["data-"+strconv.Itoa(cnt_row)].(bson.M)["username"] = single_username
					sending_data["data-"+strconv.Itoa(cnt_row)].(bson.M)["password"] = single_password
					sending_data["data-"+strconv.Itoa(cnt_row)].(bson.M)["display_name"] = single_display_name
					cnt_row++
				}

				json_send, err := json.Marshal(sending_data)
				if err != nil {
					log.Fatal(err)
				}
				json_send_str := string(json_send)
				return json_send_str
			} else {
				return "{\"status\": \"not_login\", \"username\": \"\", \"display_name\": \"\"}"
			}
		} else if s == "akses_login" {
			if login_state {
				return "{\"status\": \"login\", \"username\": \"\", \"display_name\": \"\"}"
			} else {
				return "{\"status\": \"not_login\", \"username\": \"\", \"display_name\": \"\"}"
			}
		} else if s == "open_input_pibal" {
			sending_data := bson.M{}
			sending_data["status"] = "not_login"
			if login_state {
				sending_data["permission"] = "granted"
				sending_data["status"] = "login"
			} else {
				sending_data["permission"] = "not granted"
			}
			json_send, err := json.Marshal(sending_data)
			if err != nil {
				log.Fatal(err)
			}
			return string(json_send)
		} else {
			var data bson.M
			json.Unmarshal([]byte(s), &data)
			fmt.Println("USERPASS YANG DIDAPATKAN", data)

			if data["action"] == "login" {
				user := data["username"]
				password := data["password"]
				rows, err := db.Query("select username, display_name from userpass where username = $1 and password = $2", user, password)
				if err != nil {
					fmt.Println("error select table:", err)
				}
				defer rows.Close()

				for rows.Next() {

					err = rows.Scan(&username_login, &display_name_login)

					if err != nil {
						log.Fatal(err)
					}
					break
				}
				if username_login != "" {
					login_state = true
					sending_data := bson.M{}
					sending_data["status"] = "login_success"
					sending_data["username"] = username_login
					sending_data["display_name"] = display_name_login
					json_send, err := json.Marshal(sending_data)
					if err != nil {
						log.Fatal(err)
					}
					return string(json_send)
				} else {
					sending_data := bson.M{}
					sending_data["status"] = "login_failed"
					sending_data["username"] = username_login
					sending_data["display_name"] = display_name_login
					json_send, err := json.Marshal(sending_data)
					if err != nil {
						log.Fatal(err)
					}
					return string(json_send)
				}

			}
		}

		// Process message
		fmt.Println("dapat pesan", s)
		return s
	})

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// Open dev tools
	fmt.Println("TRYING OPEN DEV TOOLS")
	// Open the developer tools window
	if err := w.OpenDevTools(); err != nil {
		log.Fatal(err)
	}

	// Blocking pattern
	a.Wait()
}
