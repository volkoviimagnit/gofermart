package request

type UserBalanceWithdrawDTO struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

func NewUserBalanceWithdrawDTO() *UserBalanceWithdrawDTO {
	return &UserBalanceWithdrawDTO{}
}

func (dto *UserBalanceWithdrawDTO) GetOrderNumber() string {
	return dto.OrderNumber
}

func (dto *UserBalanceWithdrawDTO) GetSum() float64 {
	return dto.Sum
}
