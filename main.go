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
	mux.HandleFunc("/api/object_users", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var data ObjectUser
		json.NewDecoder(r.Body).Decode(&data)
		fmt.Printf("ObjectUsers: %d => %d\n", data.ObjectID, data.UserID)
		if _, err := db.Exec("INSERT INTO object_users(object_id,user_id) VALUES($1,$2) ON CONFLICT(object_id,user_id) DO UPDATE SET updated_at=now();", data.ObjectID, data.UserID); err != nil {
			panic(err)
		}
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

	fmt.Println("Run server on port 8080")
	http.ListenAndServe(":8080", logHandler(mux))
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Before: %s", r.URL)
		h.ServeHTTP(w, r) // call original
		log.Printf("After: %s", r.URL)
	})
}

type ObjectUser struct {
	ObjectID int `json:"object_id"`
	UserID   int `json:"user_id"`
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
	ID    int
	Email string
}

func AllUsers() (map[int]User, error) {
	rows, err := db.Query("select id, email from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make(map[int]User)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			log.Fatal(err)
		}
		users[user.ID] = user
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, err
}

type Object struct {
	ID    int
	Name  string
	Image string
}

func AllObjects() (map[int]Object, error) {
	rows, err := db.Query("select id, name, image from objects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	objects := make(map[int]Object)
	for rows.Next() {
		object := Object{}
		if err := rows.Scan(&object.ID, &object.Name, &object.Image); err != nil {
			log.Fatal(err)
		}
		objects[object.ID] = object
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return objects, err
}
