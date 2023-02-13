package environment

type ProductionEnvironment struct {
	router IRouter
}

func (p ProductionEnvironment) GetRouter() IRouter {
	//TODO implement me
	panic("implement me")
}

func NewProductionEnvironment() IEnvironment {
	return &ProductionEnvironment{}
}
