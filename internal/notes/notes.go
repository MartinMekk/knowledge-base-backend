package notes

import "time"

type Note struct {
	ID      string    `json:"id"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}

type Tag struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
