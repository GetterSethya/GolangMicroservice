package main

type RegisterUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refreshToken"`
}
