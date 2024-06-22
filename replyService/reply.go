package main

type Reply struct {
	Id         string      `json:"id"`
	Body       string      `json:"body"`
	IdUser     string      `json:"idUser"`
	Username   string      `json:"username"`
	Name       string      `json:"name"`
	Profile    string      `json:"profile"`
	IdPost     string      `json:"idPost"`
	TotalChild int64       `json:"totalChild"`
	ParentId   interface{} `json:"parentId"`
	CreatedAt  int64       `json:"createdAt"`
	UpdatedAt  int64       `json:"updatedAt"`
	DeletedAt  interface{} `json:"deletedAt"`
}
