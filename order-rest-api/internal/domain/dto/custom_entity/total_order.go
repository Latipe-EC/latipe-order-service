package custom_entity

type TotalOrderInSystemInMonth struct {
	Month  int `json:"month"`
	Amount int `json:"amount"`
	Count  int `json:"count"`
}

type TotalOrderInSystemInHours struct {
	Hours  int `json:"hour"`
	Amount int `json:"amount"`
	Count  int `json:"count"`
}

type TotalOrderInSystemInDay struct {
	Day    int `json:"day"`
	Amount int `json:"amount"`
	Count  int `json:"count"`
}
