package util

const (
	// USD stands for US Dollar
	USD = "USD"
	// EUR stands for Euro
	EUR = "EUR"
	// CAD stands for Canadian Dollar
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
