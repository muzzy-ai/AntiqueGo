package models

const (
	TaxPercent=10
)

func getTaxPercent() float64 {
	return float64(TaxPercent)/100.0
}

func GetTaxAmount(price float64) float64{
	return price * getTaxPercent()
}
