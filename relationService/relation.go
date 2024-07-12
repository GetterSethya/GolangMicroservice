package main

type Relation struct {
	Id         string `json:"id"`
	TargetId   string `json:"targetId"`
	FollowerId string `json:"followerId"`

	CreatedAt int64       `json:"createdAt"`
	UpdatedAt int64       `json:"updatedAt"`
	DeletedAt interface{} `json:"deletedAt"`
}
