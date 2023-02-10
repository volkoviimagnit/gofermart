package response

type OrderDTO struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *float64    `json:"accrual"`
}

func (o *OrderDTO) GetNumber() string {
	return o.Order
}

func (o *OrderDTO) GetStatus() OrderStatus {
	return o.Status
}

func (o *OrderDTO) GetAccrual() *float64 {
	return o.Accrual
}

func (o *OrderDTO) IsTerminal() bool {
	return o.Status == AccrualStatusInvalid || o.Status == AccrualStatusProcessed
}
