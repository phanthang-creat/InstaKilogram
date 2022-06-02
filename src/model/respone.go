package model

type LoginRespone struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Message      string `json:"message"`
	User         struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Age      string `json:"age"`
		Avatar   string `json:"avatar"`
		Name     string `json:"name"`
	} `json:"user"`
}
