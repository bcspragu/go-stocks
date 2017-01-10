package stock

// Retreiver can get the price of a stock
type Retreiver interface {
	Price(stock string) (float64, error)
}
