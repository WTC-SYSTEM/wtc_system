package user

type User struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUserHashedDTO) NewUser() *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: string(u.Password),
	}
}

func (u *CreateUserDTO) Hashed(hashedPassword []byte) *CreateUserHashedDTO {
	return &CreateUserHashedDTO{
		Username: u.Username,
		Email:    u.Email,
		Password: hashedPassword,
	}
}

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserHashedDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type GetUserByEmailAndPasswordDTO struct {
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
	RegSuccess     string = "REG_SUCCESS"
	GetUserSuccess string = "GET_USER_SUCCESS"
)
