package crud

type Product struct {
	ID int
	ProductProperties
}
type ProductProperties struct {
	Name string
}
type Products []Product
