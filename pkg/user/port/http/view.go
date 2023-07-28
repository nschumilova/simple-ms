package http

type UserBrief struct {
	UserName  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type UserFull struct {
	Id        uint   `json:"id"`
	UserName  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type ErrorDescription struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
