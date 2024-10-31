package dataTypes

type ActionType string

const (
	Signup ActionType = "SignUp"
	Login  ActionType = "Login"
)

type User struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}
