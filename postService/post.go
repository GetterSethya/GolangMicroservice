package main

type Post struct {
	Id           string `json:"id"`
	Image        string `json:"image"`
	Body         string `json:"body"`
	IdUser       string `json:"idUser"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	TotalLikes   int64  `json:"totalLikes"`
	TotalReplies int64  `json:"totalReplies"`

	CreatedAt int64       `json:"createdAt"`
	UpdatedAt int64       `json:"updatedAt"`
	DeletedAt interface{} `json:"deletedAt"`
}
