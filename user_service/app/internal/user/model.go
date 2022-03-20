package user

type User struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUserDTO) NewUser() *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetUserByEmailAndPasswordDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	ID       string `json:"ID" validate:"required"`
	Username string `json:"username,omitempty" `
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

const (
	RegSuccess string = "REG_SUCCESS"
)
