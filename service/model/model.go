package model

type ProductCreate struct {
	Name     string
	Price    float32
	Discount float32
	Store    string
}

type ProductUpdate struct {
	Id       int64
	Name     string
	Price    float32
	Discount float32
	Store    string
}
