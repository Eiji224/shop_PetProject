package dto

type UserType string

const (
	TypeCustomer      UserType = "customer"
	TypeSeller        UserType = "seller"
	TypeAdministrator UserType = "administrator"
)

func (t UserType) String() string {
	return string(t)
}

func (t UserType) IsValid() bool {
	switch t {
	case TypeCustomer, TypeSeller, TypeAdministrator:
		return true
	default:
		return false
	}
}
