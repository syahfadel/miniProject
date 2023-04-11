package entities

type GetFollowingResponse struct {
	FollowingID int    `json:"id"`
	UserName    string `json:"username"`
}

type GetFollowedResponse struct {
	FollowedID int    `json:"id"`
	UserName   string `json:"username"`
}
