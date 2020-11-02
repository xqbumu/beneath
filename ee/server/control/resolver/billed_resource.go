package resolver

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.com/beneath-hq/beneath/control/entity"
	"gitlab.com/beneath-hq/beneath/control/gql"
	"gitlab.com/beneath-hq/beneath/services/middleware"
)

// BilledResource returns the gql.BilledResourceResolver
func (r *Resolver) BilledResource() gql.BilledResourceResolver {
	return &billedResourceResolver{r}
}

type billedResourceResolver struct{ *Resolver }

func (r *billedResourceResolver) EntityKind(ctx context.Context, obj *entity.BilledResource) (string, error) {
	return string(obj.EntityKind), nil
}

func (r *billedResourceResolver) Product(ctx context.Context, obj *entity.BilledResource) (string, error) {
	return string(obj.Product), nil
}

func (r *billedResourceResolver) Quantity(ctx context.Context, obj *entity.BilledResource) (float64, error) {
	return float64(obj.Quantity), nil
}

func (r *billedResourceResolver) Currency(ctx context.Context, obj *entity.BilledResource) (string, error) {
	return string(obj.Currency), nil
}

func (r *queryResolver) BilledResources(ctx context.Context, organizationID uuid.UUID, billingTime time.Time) ([]*entity.BilledResource, error) {
	secret := middleware.GetSecret(ctx)
	perms := r.Permissions.OrganizationPermissionsForSecret(ctx, secret, organizationID)
	if !perms.Admin {
		return nil, gqlerror.Errorf("Not allowed to perform admin functions on organization %s", organizationID.String())
	}

	billedResources := entity.FindBilledResources(ctx, organizationID, billingTime)
	if billedResources == nil {
		return nil, gqlerror.Errorf("Billed resources for organization %s and time %s not found", organizationID.String(), billingTime)
	}

	return billedResources, nil
}