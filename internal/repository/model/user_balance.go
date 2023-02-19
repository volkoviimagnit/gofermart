package model

type UserBalance struct {
	UserID    string
	Current   float64
	Withdrawn float64
}

func NewUserBalance(userID string, current float64, withdrawn float64) *UserBalance {
	return &UserBalance{UserID: userID, Current: current, Withdrawn: withdrawn}
}

func (ub *UserBalance) GetCurrent() float64 {
	return ub.Current
}

func (ub *UserBalance) GetWithdrawn() float64 {
	return ub.Withdrawn
}

func (ub *UserBalance) SetCurrent(current float64) {
	ub.Current = current
}

func (ub *UserBalance) SetWithdrawn(withdrawn float64) {
	ub.Withdrawn = withdrawn
}
