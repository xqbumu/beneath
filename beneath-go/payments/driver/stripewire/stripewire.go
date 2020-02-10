package stripewire

import (
	"net/http"

	"fmt"

	"github.com/beneath-core/beneath-go/control/entity"
	"github.com/beneath-core/beneath-go/core"
	"github.com/beneath-core/beneath-go/core/httputil"
	"github.com/beneath-core/beneath-go/core/log"
	"github.com/beneath-core/beneath-go/core/middleware"
	"github.com/beneath-core/beneath-go/core/timeutil"
	"github.com/beneath-core/beneath-go/payments/driver"
	"github.com/beneath-core/beneath-go/payments/driver/stripeutil"
	uuid "github.com/satori/go.uuid"
	stripe "github.com/stripe/stripe-go"
)

// StripeWire implements beneath.PaymentsDriver
type StripeWire struct{}

type configSpecification struct {
	StripeSecret        string `envconfig:"CONTROL_STRIPE_SECRET" required:"true"`
	PaymentsAdminSecret string `envconfig:"PAYMENTS_ADMIN_SECRET" required:"true"`
}

// New initializes a StripeWire object
func New() StripeWire {
	var config configSpecification
	core.LoadConfig("beneath", &config)
	stripeutil.InitStripe(config.StripeSecret)

	return StripeWire{}
}

// GetHTTPHandlers returns the necessary handlers to implement Stripe card payments
func (w *StripeWire) GetHTTPHandlers() map[string]httputil.AppHandler {
	return map[string]httputil.AppHandler{
		"initialize_customer": handleInitializeCustomer,
		// "webhook":               handleStripeWebhook,       // TODO: when a customer pays by wire, check to see if any important Stripe events are emitted via webhook
		"get_payment_details": handleGetPaymentDetails,
	}
}

// create/update a customer's billing info and Stripe registration
func handleInitializeCustomer(w http.ResponseWriter, req *http.Request) error {
	var config configSpecification
	core.LoadConfig("beneath", &config)

	organizationID, err := uuid.FromString(req.URL.Query().Get("organizationID"))
	if err != nil {
		return httputil.NewError(400, "couldn't get organizationID from the request")
	}

	organization := entity.FindOrganization(req.Context(), organizationID)
	if organization == nil {
		return httputil.NewError(400, "organization not found")
	}

	billingPlanID, err := uuid.FromString(req.URL.Query().Get("billingPlanID"))
	if err != nil {
		return httputil.NewError(400, "couldn't get billingPlanID from the request")
	}

	billingPlan := entity.FindBillingPlan(req.Context(), billingPlanID)
	if billingPlan == nil {
		return httputil.NewError(400, "billing plan not found")
	}

	emailAddress := req.URL.Query().Get("emailAddress")
	if emailAddress == "" {
		return httputil.NewError(400, "couldn't get emailAddress from the request")
	}

	// Beneath will call the function from an admin panel (after a customer discussion)
	secret := middleware.GetSecret(req.Context())
	if secret.GetSecretID().String() != config.PaymentsAdminSecret {
		return httputil.NewError(403, fmt.Sprintf("Enterprise plans require a Beneath Payments Admin to activate"))
	}

	// Our requests to Stripe differ whether or not the customer is already registered in Stripe
	var customer *stripe.Customer
	driverPayload := make(map[string]interface{})
	billingInfo := entity.FindBillingInfo(req.Context(), organization.OrganizationID)
	if billingInfo.DriverPayload["customer_id"] != nil {
		// customer is already registered with Stripe
		driverPayload["customer_id"] = billingInfo.DriverPayload["customer_id"]
		stripeutil.UpdateWireCustomer(driverPayload["customer_id"].(string), emailAddress)
	} else {
		// customer needs to be registered with stripe
		customer = stripeutil.CreateWireCustomer(organization.OrganizationID, organization.Name, emailAddress)
		driverPayload["customer_id"] = customer.ID
	}

	_, err = entity.UpdateBillingInfo(req.Context(), organization.OrganizationID, billingPlan.BillingPlanID, entity.StripeWireDriver, driverPayload)
	if err != nil {
		log.S.Errorf("Error updating billing info: %v\\n", err)
		return httputil.NewError(500, "error updating billing info: %v\\n", err)
	}

	return nil
}

func handleGetPaymentDetails(w http.ResponseWriter, req *http.Request) error {
	// TODO: is there anything we want to return to the front-end?
	// - bank account information where the wire should be sent
	// - state of recent payment (paid, X days remaining, Y days overdue)
	return nil
}

// IssueInvoiceForResources implements Payments interface
func (w *StripeWire) IssueInvoiceForResources(billingInfo driver.BillingInfo, billedResources []driver.BilledResource) error {
	if billingInfo.GetDriverPayload()["customer_id"] == nil {
		panic("stripe customer id is not set")
	}

	var seatCount int64
	var seatPrice int64
	var seatStartTime int64
	var seatEndTime int64

	for _, item := range billedResources {
		// only itemize the products that cost money (i.e. don't itemize the included Reads and Writes)
		// only itemize the products for this month's bill
		if (item.GetTotalPriceCents() != 0) && (item.GetBillingTime().UTC() == timeutil.BeginningOfThisPeriod(timeutil.PeriodMonth)) {
			// count seats, itemize everything else
			if item.GetProduct() == string(entity.SeatProduct) {
				seatCount++
				seatPrice = int64(item.GetTotalPriceCents())
				seatStartTime = item.GetStartTime().Unix()
				seatEndTime = item.GetEndTime().Unix()
			} else {
				stripeutil.NewInvoiceItemOther(billingInfo.GetDriverPayload()["customer_id"].(string), int64(item.GetTotalPriceCents()), string(billingInfo.GetBillingPlanCurrency()), item.GetStartTime().Unix(), item.GetEndTime().Unix(), stripeutil.PrettyDescription(item.GetProduct()))
			}
		}
	}

	// batch seats
	if seatCount > 0 {
		stripeutil.NewInvoiceItemSeats(billingInfo.GetDriverPayload()["customer_id"].(string), seatCount, seatPrice, string(billingInfo.GetBillingPlanCurrency()), seatStartTime, seatEndTime, stripeutil.PrettyDescription(string(entity.SeatProduct)))
	}

	inv := stripeutil.CreateInvoice(billingInfo.GetDriverPayload()["customer_id"].(string), billingInfo.GetPaymentsDriver())

	stripeutil.SendInvoice(inv.ID)

	return nil
}