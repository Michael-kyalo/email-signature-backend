package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Hash passwords
}

type Click struct {
	ID        string `json:"id"`
	LinkID    string `json:"link_id"`
	Timestamp string `json:"timestamp"`
}
