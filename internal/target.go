package internal

type EmailTarget struct {
	Email string `json:"email"`
}

func (e EmailTarget) Type() int {
	return TYPEEmail
}

func (e EmailTarget) Value() string {
	return e.Email
}

type PhoneTarget struct {
	Phone string `json:"phone"`
}

func (e PhoneTarget) Type() int {
	return TYPEPhone
}

func (e PhoneTarget) Value() string {
	return e.Phone
}

type IdTarget struct {
	Id string `json:"id"`
}

func (e IdTarget) Type() int {
	return TYPEId
}

func (e IdTarget) Value() string {
	return e.Id
}
