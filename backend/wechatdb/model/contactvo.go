package model

type ContactVo struct {
	UserName      string   `json:"userName"`
	Alias         string   `json:"alias"`
	Remark        string   `json:"remark"`
	NickName      string   `json:"nickName"`
	BigHeadUrl    string   `json:"big_head_url"`
	SmallHeadUrl  string   `json:"small_head_url"`
	Labels        []string `json:"labels"`
	Phone         []string `json:"phones"`
	Gender        string   `json:"gender"`
	Signature     string   `json:"signature"`
	Country       string   `json:"country"`
	Province      string   `json:"province"`
	City          string   `json:"city"`
	ChatroomCount int      `json:"chatroom_count"`
}
