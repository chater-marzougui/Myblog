package global

var userID = 1
var authenticated = false

func GetUserID() int {
	return userID
}

func ModifyUser(newUser int) {
	userID = newUser
}

func IsAuthenticated() bool {
	return authenticated
}

func SetAuthenticated() {
	authenticated = true
}
