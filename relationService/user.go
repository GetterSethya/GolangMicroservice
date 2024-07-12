package main

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`

	CreatedAt int64       `json:"createdAt"`
	UpdatedAt int64       `json:"updatedAt"`
	DeletedAt interface{} `json:"deletedAt"`
}
