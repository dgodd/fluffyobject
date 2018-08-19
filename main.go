package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/r3labs/sse"
)

const connStr = "postgres://localhost/fluffyobject?sslmode=disable"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func main() {
	server := &sse.Server{
		BufferSize: sse.DefaultBufferSize,
		AutoStream: false,
		AutoReplay: false,
		Streams:    make(map[string]*sse.Stream),
	}
	server.CreateStream("messages")

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", server.HTTPHandler)
	mux.HandleFunc("/api/senddata", func(w http.ResponseWriter, r *http.Request) {
		db.QueryRow("NOTIFY data_changed")
	})
	mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))

	go func() {
		notificationChan := make(chan *pq.Notification)
		l, err := pq.NewListenerConn(connStr, notificationChan)
		if err != nil {
			panic(err)
		}
		defer l.Close()
		if ok, err := l.Listen("data_changed"); !ok || err != nil {
			panic(err)
		}

		for {
			select {
			case <-notificationChan:
				users, err := AllUsers()
				if err != nil {
					panic(err)
				}
				if err := Publish(server, "users", users); err != nil {
					panic(err)
				}
				objects, err := AllObjects()
				if err != nil {
					panic(err)
				}
				if err := Publish(server, "objects", objects); err != nil {
					panic(err)
				}
				fmt.Printf("Users: %#v\nObjects: %#v\n", users, objects)
			}
		}
	}()

	http.ListenAndServe(":8080", mux)
}

func Publish(server *sse.Server, event string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	server.Publish("messages", &sse.Event{
		Event: []byte(event),
		Data:  b,
	})
	return nil
}

type User struct {
	ID    string
	Email string
}

func AllUsers() ([]User, error) {
	rows, err := db.Query("select id, email from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, err
}

type Object struct {
	ID    string
	Name  string
	Image string
}

func AllObjects() ([]Object, error) {
	rows, err := db.Query("select id, name, image from objects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var objects []Object
	for rows.Next() {
		object := Object{}
		if err := rows.Scan(&object.ID, &object.Name, &object.Image); err != nil {
			log.Fatal(err)
		}
		objects = append(objects, object)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return objects, err
}
