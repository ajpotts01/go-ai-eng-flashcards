package models

// Message represents a single entry in a conversation, with a role and content.
// This is a core data structure used by the service and handler.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
