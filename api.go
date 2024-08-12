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
	var scenario string = generateText("VAR1, VAR2 and VAR3 is variable that you decided and create a criminal scenario that VAR1 is the name of suspect, VAR2 is the A brief history of the crime and incident of which the suspect is accused, VAR3 is evidence to incriminate the suspect. And write them like:\nName: VAR1\nEvent: VAR2\nEvidence: VAR3", "AIzaSyCTB075sQmU6Yh76nZZgYO_sV-bMzKANfg")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addNewRoom(
		uuid.New().String(),
		scenario,
		sus,
		[]MessageData{},
	))
}
func removeRoom(slice []Room, roomIDToRemove string) {
	for i, room := range rooms {
		if room.RoomID == roomIDToRemove {
			rooms = append(rooms[:i], rooms[i+1:]...)
			break
		}
	}
}
func deleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedRoomID := vars["roomID"]

	removeRoom(rooms, requestedRoomID)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(w, "{'roomID':'"+requestedRoomID+"','message':'room removed'}")
}
func api() {
	//keyy := flag.String("key", "", "")
	flag.Parse()
	//key = string(*keyy)
	//key := "AIzaSyCTB075sQmU6Yh76nZZgYO_sV-bMzKANfg"
	router := mux.NewRouter()

	// Define your routes

	router.PathPrefix("/src").Handler(http.StripPrefix("/src", http.FileServer(http.Dir("./templates/src"))))

	router.PathPrefix("/src").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/src")
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "./templates")
		t, err := template.ParseGlob("./templates/main.html")
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
		t, err := template.ParseGlob("./templates/index.html")
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		} else {
			for _, room := range rooms {
				if room.RoomID == requestedRoomID {
					t.Execute(w, requestedRoomID)
					return
				}
			}
			http.NotFound(w, r)
		}
	})
	router.HandleFunc("/add_room", addRoom).Methods("GET")
	router.HandleFunc("/delete_room/{roomID}", deleteRoom).Methods("GET")
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
		add_message_with_ai(roomID, "AIzaSyCTB075sQmU6Yh76nZZgYO_sV-bMzKANfg")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getChat(roomID))
	})
	log.Fatal(http.ListenAndServe(":8080", router))
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
