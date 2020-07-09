package user

type SignUpRequest struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

