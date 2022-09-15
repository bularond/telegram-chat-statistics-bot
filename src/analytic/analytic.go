package analytic

import (
	"strconv"
	"strings"
	"time"
)

type Person struct {
	Name string
	Id   string

	MessageCount int
	WordsCount   int
	CharCount    int
}

type ChatStatistics struct {
	Chat *Chat

	MessageCount int
	WordsCount   int
	CharCount    int

	Persons map[string]*Person

	WorsdMap map[string]int
	DateMap  map[time.Time]int
}

func newChatStatistics(chat *Chat) ChatStatistics {
	return ChatStatistics{
		Chat: chat,

		MessageCount: 0,
		WordsCount:   0,
		CharCount:    0,

		Persons:  make(map[string]*Person),
		WorsdMap: make(map[string]int),
		DateMap:  make(map[time.Time]int),
	}
}

func newPerson(name string, id string) *Person {
	return &Person{
		Name: name,
		Id:   id,

		MessageCount: 0,
		WordsCount:   0,
		CharCount:    0,
	}
}

func GetChatStatistics(chat *Chat) (*ChatStatistics, error) {
	stats := newChatStatistics(chat)

	for _, message := range chat.Messages {
		err := stats.handleMessage(&message)
		if err != nil {
			return nil, err
		}
	}

	if _, exists := stats.Persons[""]; exists {
		delete(stats.Persons, "")
	}

	return &stats, nil
}

func (cs *ChatStatistics) handleMessage(message *Message) (err error) {
	person := cs.getPerson(message.From, message.FromId)
	cs.MessageCount += 1
	person.MessageCount += 1

	text := strings.ToLower(message.Text)
	cs.CharCount += len(text)
	person.CharCount += len(text)

	words := strings.Fields(text)
	cs.WordsCount += len(words)
	person.WordsCount += len(words)

	for _, word := range words {
		cs.WorsdMap[word] += 1
	}

	messageDate, err := message.getDayOfMessage()
	if err != nil {
		return
	}
	cs.DateMap[messageDate] += 1

	return nil
}

func (m *Message) getDayOfMessage() (time.Time, error) {
	timestamp, err := strconv.ParseInt(m.DateUnixtime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	messageTime := time.Unix(timestamp, 0)
	timestamp -= int64(messageTime.Hour()*60*60 + messageTime.Minute()*60 + messageTime.Second())

	return time.Unix(timestamp, 0), nil
}

func (cs *ChatStatistics) getPerson(name string, id string) *Person {
	if person, exists := cs.Persons[id]; exists {
		return person
	} else {
		person = newPerson(name, id)
		cs.Persons[id] = person
		return person
	}
}
