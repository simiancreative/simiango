package crud

type Product struct {
	ID int
	ProductProperties
}
type ProductProperties struct {
	Name  string
	Depth int
}
type Products []Product
