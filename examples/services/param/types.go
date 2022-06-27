package param

type Company struct {
	ID *int `json:"app_id" header:"X-Companyid"`
}

type Item struct {
	Func func(string) error `json:"-"`
	ID   int                `json:"id" param:"id"`
}

type paramService struct {
	Company Company
	Item    Item
}
