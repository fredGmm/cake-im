package module


type Message struct {
	//Username string `json:"username"`
	Message  string `json:"message"`
	ToUserId string `json:"to_user_id"`
	UserId   string `json:"user_id"`
	GroupId  string `json:"group_id"`
}
