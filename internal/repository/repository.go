package repository

// Repository provides access to the message store.
type Repository interface {
	CreateMessage(title, body string) (string, error)
	ListMessages() (string, error)
	CreateReply(messageID, reply string) error
	ListReplies(messageID string) (string, error)
	Ping() (string, error)
}
