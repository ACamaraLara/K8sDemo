package dataTypes

type ActionType string

const (
	Signup ActionType = "SignUp"
	Login  ActionType = "Login"
)

type User struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Action    ActionType `json:"action"`
}
