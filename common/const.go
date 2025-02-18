package common

const (
	CurrentUser = "current_user"
)

type TokenPayload struct {
	UId   int    `json:"user_id"`
	URole string `json:"role"`
}

func (t TokenPayload) UserId() int {
	return t.UId
}

func (t TokenPayload) Role() string {
	return t.URole
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
