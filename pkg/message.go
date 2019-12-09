package pkg

type Message struct {
	UserId      string   `json:"user_id"`
	Name        string   `json:"name"`
	Avatar      string    `json:"avatar"`
	Text        string   `json:"text"`
	GroupId     string   `json:"group_id"`
	ToUserId    string   `json:"to_user_id"`
	Type        string   `json:"type"` // BroadcastRadio  , single
	OnlineCount int      `json:"online_count"`
	Images      []string `json:"images"`
	Time        string   `json:"time"`
}
