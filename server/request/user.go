package request

// UserRequest ...
type UserRequest struct {
	Name     string `json:"name" `
	Email    string `json:"Email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest ...
type UserUpdateRequest struct {
	Name  string `json:"name" `
	Email string `json:"email" validate:"required"`
}

// UserUploadImageRequest ...
type UserUploadImageRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// UserLoginRequest ....
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}
