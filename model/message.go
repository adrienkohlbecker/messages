package model

import "time"

type Message struct {
	ID          string        `json:"id"`
	Sender      string        `json:"sender"`
	Content     string        `json:"content"`
	Timestamp   time.Time     `json:"timestamp"`
	Sent        bool          `json:"sent"`
	Attachments []*Attachment `json:"attachments"`
	Group       string        `json:"group"`
	Kind        string        `json:"kind"`
	RepliesTo   string        `json:"replies_to"`
}
