package kafka

type UserRegisteredEvent struct {
	UserID string `json:"user_id"`
	Login  string `json:"login"`
	Time   string `json:"time"`
}
