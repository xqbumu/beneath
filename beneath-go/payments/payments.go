package payments

import (
	"github.com/beneath-core/beneath-go/payments/driver"
	"github.com/beneath-core/beneath-go/payments/driver/anarchism"
	"github.com/beneath-core/beneath-go/payments/driver/stripecard"
	"github.com/beneath-core/beneath-go/payments/driver/stripewire"
)

// InitDrivers initializes all of the payments drivers
func InitDrivers(drivers []string) map[string]driver.PaymentsDriver {
	payments := make(map[string]driver.PaymentsDriver)
	for _, driver := range drivers {
		switch driver {
		case "stripecard":
			sc := stripecard.New()
			payments[driver] = &sc
		case "stripewire":
			sw := stripewire.New()
			payments[driver] = &sw
		case "anarchism":
			a := anarchism.New()
			payments[driver] = &a
		default:
			panic("unrecognized payments driver")
		}
	}
	return payments
}