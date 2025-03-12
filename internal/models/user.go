package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Signed   bool   `json:"signed"` // Signed будет равным true, если пользователь подписан на рассылку; иначе false.
}
