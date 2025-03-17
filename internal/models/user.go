package models

type User struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Subscribed bool   `json:"subscribed"` // Subscribed будет равным true, если пользователь подписан на рассылку; иначе false.
}
