package response

type UserBalanceDTO struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewUserBalanceDTO(current float64, withdrawn float64) *UserBalanceDTO {
	return &UserBalanceDTO{Current: current, Withdrawn: withdrawn}
}
