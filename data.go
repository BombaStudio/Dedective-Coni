package main

var key string

type Room struct {
	RoomID   string        `json:"roomID"`
	Scenario string        `json:"scenario"`
	Sus      bool          `json:"sus"`
	Messages []MessageData `json:"messages"`
}

// MessageData represents a message in a room.
type MessageData struct {
	ClientID string `json:"clientID"`
	Message  string `json:"message"`
}

var roomsTemplate = []Room{
	{
		RoomID:   "1234",
		Scenario: "you killed my wife",
		Messages: []MessageData{
			{
				ClientID: "31",
				Message:  "I'm sorry, I didn't mean to kill your wife",
			},
		},
	},
}

type Status struct {
	StatusEvent string `json:"statusEvent"`
	Command     string `json:"command"`
	Runned      bool   `json:"runned"`
}

var rooms = []Room{
	{
		RoomID:   "test",
		Scenario: "you killed my wife",
		Messages: []MessageData{
			{
				ClientID: "31",
				Message:  "I'm sorry, I didn't mean to kill your wife",
			},
		},
	},
}

func addNewRoom(
	roomID string,
	scenario string,
	sus bool,
	messages []MessageData,
) Room {
	room := Room{
		RoomID:   roomID,
		Scenario: scenario,
		Sus:      sus,
		Messages: messages,
	}
	rooms = append(rooms, room)
	return room
}

func getChat(roomID string) []MessageData {
	var data []MessageData
	for _, room := range rooms {
		if room.RoomID == roomID {
			data = room.Messages
			break
		}
	}
	return data
}
func addChat(roomID string, message MessageData) Status {
	var data Status = Status{
		StatusEvent: "addChat",
		Command:     message.Message,
		Runned:      false,
	}
	var i = 0
	for _, room := range rooms {
		if room.RoomID == roomID {
			rooms[i].Messages = append(rooms[i].Messages, message)
			data.Runned = true
			break
		}
		i++
	}
	return data
}
