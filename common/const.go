package common

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
