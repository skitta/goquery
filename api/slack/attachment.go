// This example will post a detailed message attachment as though it were sent
// from a task management service:
// {
//    "attachments":[
//       {
//          "fallback":"fallback",
//          "pretext":"text",
//          "color":"#D00000",
//          "fields":[
//             {
//                "title":"Notes",
//                "value":"This is much easier than I thought it would be.",
//                "short":false
//             }
//          ]
//       }
//    ]
// }
// More detail see https://api.slack.com/docs/message-attachments

package slack

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Message serve a slack attachment message
type Message struct {
	Attachments []*Attachment `json:"attachments"`
}

// Attachment serve a slack attachment json
type Attachment struct {
	Fallback string  `json:"fallback"`
	Text     string  `json:"text"`
	Pretext  string  `json:"pretext"`
	Color    string  `json:"color"`
	Fields   []Field `json:"fields"`
}

// Field serve a field in attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

var initAttachment *Attachment

func init() {
	if initAttachment == nil {
		a := Attachment{}
		a.Color = "#4169E1" // RoyalBlue
		a.Fallback = "An Error Occured!"
		initAttachment = &a
	}
}

// NewMessage return a Message type which point to a init attachment
func NewMessage(pretext, text string) Message {
	initAttachment.Text = text
	initAttachment.Pretext = pretext
	initAttachment.Fallback = pretext
	return Message{[]*Attachment{initAttachment}}
}

// AddColor change attachment's display color
func (m *Message) AddColor(color string) {
	initAttachment.Color = color
}

// AddField modify a field and give it to an attachment
func (m *Message) AddField(title, value string, short bool) {
	var f = Field{title, value, short}
	initAttachment.Fields = []Field{f}
}

// Post a Message to slack webhook
func Post(url string, m Message) error {
	body, _ := json.Marshal(m)
	payload := strings.NewReader(string(body))
	_, err := http.Post(url, "application/json", payload)
	if err != nil {
		return err
	}
	return nil
}
