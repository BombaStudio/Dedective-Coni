package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func generateText(prmpt string, key string) string {

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		fmt.Println(1)
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prmpt))
	if err != nil {
		fmt.Println(2)
		log.Fatal(err)
	}
	fmt.Println(3)
	fmt.Println(printResponse(resp))
	fmt.Println("Totel token count is " + fmt.Sprint(resp.UsageMetadata.TotalTokenCount))
	fmt.Println("max token count is " + fmt.Sprint(resp.UsageMetadata.PromptTokenCount))
	return printResponse(resp)
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var output string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			//fmt.Println(cand.Content)
			for _, part := range cand.Content.Parts {
				//fmt.Println(part)
				output += fmt.Sprint(part)
			}
		}
	}
	return output

}

func add_message_with_ai(roomID string, key string) {

	var i = 0
	for _, room := range rooms {
		if room.RoomID == roomID {
			msg, err := json.Marshal(rooms[i])
			if err != nil {
				break
			}
			fmt.Println(string(msg))
			rooms[i].Messages = append(rooms[i].Messages, MessageData{
				ClientID: "AI",
				Message:  generateText("You are a suspect in an interrogation room and you are accused of 'scenario', which I will mention at the end of the text. And our conversations during this query are recorded in the JSON parameter in the 'messages' variable. Among the JSON parameters in the 'messages' variable, if 'you' is written in the 'clientID' parameter, you have said what is written in the 'message' parameter before, but if 'me' is written in the 'clientID' parameter, I have said what is written in the 'message' parameter before. It will base you on the 'clientID' and 'message' parameters from the JSON parameters in the 'messages' variable, and you will look at the 'sus' boolean value and answer the question as a guilty person if it is true, or as an innocent if it is false, if I asked you a question. But if I don't ask anything, you will try to convince me that you are innocent."+string(msg), key),
			})
			break
		}
		i++
	}

}
