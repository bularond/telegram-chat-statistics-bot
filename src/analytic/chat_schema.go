package analytic

import (
	"encoding/json"
	"fmt"
)

type Chat struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Id       int       `json:"Id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Id               int    `json:"Id"`
	Type             string `json:"type"`
	Date             string `json:"date"`
	DateUnixtime     string `json:"date_unixtime"`
	From             string `json:"from"`
	FromId           string `json:"from_id"`
	Edited           string `json:"edited"`
	EditedUnixtime   string `json:"edited_unixtime"`
	ForwardedFrom    string `json:"forwarded_from"`
	ReplyToMessageId int    `json:"reply_to_message_id"`

	Photo  string `json:"photo"`
	Width  int    `json:"width"`
	Height int    `json:"height"`

	File            string `json:"file"`
	Thumbnail       string `json:"thumbnail"`
	MediaType       string `json:"media_type"`
	Performer       string `json:"performer"`
	Title           string `json:"title"`
	MimeType        string `json:"mime_type"`
	DurationSeconds int    `json:"duration_seconds"`

	LocationInformation       Location `json:"location_information"`
	LiveLocationPeriodSeconds int      `json:"live_location_period_seconds"`

	Actor         string `json:"actor"`
	ActorId       string `json:"actor_id"`
	Action        string `json:"action"`
	DiscardReason string `json:"discard_reason"`
	Emoticon      string `json:"emoticon"`

	StickerEmoji string `json:"sticker_emoji"`

	RawText json.RawMessage `json:"text"`
	Text    string
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ParseJson(rawData []byte) (chat *Chat, err error) {
	err = json.Unmarshal(rawData, &chat)
	if err != nil {
		return
	}
	for i := range chat.Messages {
		err = chat.Messages[i].setText()
		if err != nil {
			return
		}
		chat.Messages[i].RawText = nil
	}
	return
}

func (m *Message) setText() (err error) {
	if m.RawText[0] == '"' {
		err = json.Unmarshal(m.RawText, &m.Text)
		return
	} else if m.RawText[0] != '[' {
		err = fmt.Errorf("unknown type of message text: %v", m)
		return
	}

	var textParts []json.RawMessage
	err = json.Unmarshal(m.RawText, &textParts)
	if err != nil {
		return
	}

	for _, partText := range textParts {
		if partText[0] == '"' {
			tmpText := ""
			err = json.Unmarshal(partText, &tmpText)
			if err != nil {
				return
			}
			m.Text += tmpText
		} else if partText[0] == '{' {
			var tmpText struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}
			err = json.Unmarshal(partText, &tmpText)
			if err != nil {
				return
			}
			m.Text += tmpText.Text
		} else {
			err = fmt.Errorf("unknown type of message text: %v", m)
			return
		}
	}
	return
}
