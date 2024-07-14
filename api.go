package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getRooms(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedRoomID := vars["roomID"]

	for _, room := range rooms {
		if room.RoomID == requestedRoomID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(room)
			return
		}
	}

	// If room not found, return an error
	http.NotFound(w, r)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedRoomID := vars["roomID"]

	for _, room := range rooms {
		if room.RoomID == requestedRoomID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(room.Messages)
			return
		}
	}

	// If room not found, return an error
	http.NotFound(w, r)
}

func addRoom(w http.ResponseWriter, r *http.Request) {
	var sus bool = RandBool()
	var scenario string = generateText("", key)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addNewRoom(
		uuid.New().String(),
		scenario,
		sus,
		[]MessageData{},
	))
}

func api() {
	keyy := flag.String("key", "", "")
	flag.Parse()
	key = string(*keyy)
	router := mux.NewRouter()

	// Define your routes

	router.PathPrefix("/src").Handler(http.StripPrefix("/src", http.FileServer(http.Dir("./templates/src"))))

	router.PathPrefix("/src").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/src")
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "./templates")
		t, err := template.ParseGlob("./templates/index.gohtml")
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		} else {
			t.Execute(w, rooms)
		}

	})
	router.HandleFunc("/room/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedRoomID := vars["roomID"]
		t, err := template.ParseGlob("./templates/index.gohtml")
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		} else {
			t.Execute(w, requestedRoomID)
		}
	})
	router.HandleFunc("/add_room", addRoom).Methods("GET")
	router.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	})
	router.HandleFunc("/rooms/{roomID}", getRooms).Methods("GET")
	router.HandleFunc("/rooms/{roomID}/messages", getMessages).Methods("GET")
	router.HandleFunc("/rooms/{roomID}/lastMessages", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedRoomID := vars["roomID"]

		for _, room := range rooms {
			if room.RoomID == requestedRoomID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(room.Messages[len(room.Messages)-1])
				return
			}
		}

		// If room not found, return an error
		http.NotFound(w, r)
	}).Methods("GET")
	router.HandleFunc("/rooms/{roomID}/add_message/{clientID}/{message}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var roomID = vars["roomID"]
		var clientID = vars["clientID"]
		var message = vars["message"]
		addChat(roomID, MessageData{
			ClientID: clientID,
			Message:  message,
		})
		add_message_with_ai(roomID, key)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getChat(roomID))
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
