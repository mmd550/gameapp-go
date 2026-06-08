package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type GetProfileRequest struct {
	UserId uint
}

type GetProfileResponse struct {
	Name string `json:"name"`
}

type RegisterRequest struct {
	Name        string `json:"name"         validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,persianphone"`
	Password    string `json:"password"     validate:"required,min=8"`
}

type RegisteredUser struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Id          uint   `json:"id"`
}
type RegisterResponse struct {
	User RegisteredUser `json:"user"`
}
