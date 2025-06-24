package models

type UserForStatus struct {
	Id             int    `json:"id"`
	Login          string `json:"login"`
	Quantity       int    `json:"quantity"`
	Points         int    `json:"points"`
	ReferrerID     int    `json:"referrer_id"`
	CompletedTasks []int  `json:"completed_tasks"`
}

type UserForLeaderBoard struct {
	Login    string `json:"login"`
	Quantity int    `json:"quantity"`
	Points   int    `json:"points"`
}
