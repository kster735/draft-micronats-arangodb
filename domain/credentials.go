package domain

type Credentials struct {
	UID      string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
