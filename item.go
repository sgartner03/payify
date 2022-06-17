package main

type Item struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Amount   float32 `json:"amount"`
	Unit     string  `json:"unit"`
	Username string  `json:"username"`
}

func New(id int, name string, price float32, amount float32, unit string, username string) Item {
	return Item{id, name, price, amount, unit, username}
}

func (item Item) TotalPrice() float32 {
	return item.Price * item.Amount
}
