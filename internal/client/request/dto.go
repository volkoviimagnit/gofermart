package request

type ItemDTO struct {
	Match      string  `json:"match"`
	Reward     float64 `json:"reward"`
	RewardType string  `json:"reward_type"`
}

type OrderItemDTO struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type OrderDTO struct {
	Number      int            `json:"order,string"`
	Goods       []OrderItemDTO `json:"goods"`
	WaitAccrual float64        `json:"-"`
}
