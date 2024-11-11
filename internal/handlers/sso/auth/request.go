package authhandler

type UpdateUserInfoReq struct {
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	TgLink  string `json:"tg_link,omitempty"`
	IsAdmin bool   `json:"is_admin,omitempty"`
}
