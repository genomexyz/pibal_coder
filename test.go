package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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
			if login_state {
				return "{\"status\": \"login\", \"username\": \"\", \"display_name\": \"\"}"
			} else {
				return "{\"status\": \"not_login\", \"username\": \"\", \"display_name\": \"\"}"
			}
		} else if s == "akses_login" {
			if login_state {
				return "{\"status\": \"login\", \"username\": \"\", \"display_name\": \"\"}"
			} else {
				return "{\"status\": \"not_login\", \"username\": \"\", \"display_name\": \"\"}"
			}
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
					var sending_data bson.M
					sending_data["status"] = "login_success"
					sending_data["username"] = username_login
					sending_data["display_name"] = display_name_login
					json_send, err := json.Marshal(sending_data)
					if err != nil {
						log.Fatal(err)
					}
					return json_send
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
	//w.OpenDevTools()

	// Blocking pattern
	a.Wait()
}
