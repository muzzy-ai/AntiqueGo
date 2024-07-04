package models


type Model struct {
	Model interface{}
}

func RegisterModel()[]Model{
	return []Model{
		{Model:User{}},
		{Model:Address{}},
        {Model:Product{}},
        {Model:Category{}},
        {Model:Section{}},
        {Model:ProductImage{}},
		{Model:Order{}},
        {Model:OrderItem{}},
        {Model:Shipment{}},
        {Model:Payment{}},
		{Model:OrderCustomer{}},
		{Model:Cart{}},
        {Model:CartItem{}},
	}
}