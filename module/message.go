package module


type Message struct {
	UserId   string `json:"user_id"`
	//Username string `json:"username"`
	Message  string `json:"message"`
	GroupId  string `json:"group_id"`
	ToUserId string `json:"to_user_id"`
	Type string `json:"type"` // BroadcastRadio  , single
}
