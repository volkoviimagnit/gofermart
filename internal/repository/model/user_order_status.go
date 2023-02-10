package model

type UserOrderStatus string

const UserOrderStatusNew UserOrderStatus = "NEW"
const UserOrderStatusProcessing UserOrderStatus = "PROCESSING"
const UserOrderStatusInvalid UserOrderStatus = "INVALID"
const UserOrderStatusProcessed UserOrderStatus = "PROCESSED"

func (e UserOrderStatus) IsSuitable() bool {
	return e == UserOrderStatusProcessed
}

func (e UserOrderStatus) String() string {
	switch e {
	case UserOrderStatusNew:
		return "NEW"
	case UserOrderStatusProcessing:
		return "PROCESSING"
	case UserOrderStatusInvalid:
		return "INVALID"
	case UserOrderStatusProcessed:
		return "PROCESSED"
	default:
		return ""
	}
}
