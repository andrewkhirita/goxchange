package types

type Currency struct {
	Code  string
	Name  string
	Rates map[string]float64
}
