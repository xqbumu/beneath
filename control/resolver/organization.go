package resolver

import (
	"context"

	"github.com/beneath-core/control/entity"
	"github.com/beneath-core/control/gql"
	"github.com/beneath-core/internal/middleware"
	uuid "github.com/satori/go.uuid"
	"github.com/vektah/gqlparser/gqlerror"
)

// Organization returns the gql.OrganizationResolver
func (r *Resolver) Organization() gql.OrganizationResolver {
	return &organizationResolver{r}
}

type organizationResolver struct{ *Resolver }

func (r *organizationResolver) OrganizationID(ctx context.Context, obj *entity.Organization) (string, error) {
	return obj.OrganizationID.String(), nil
}

func (r *queryResolver) OrganizationByName(ctx context.Context, name string) (*entity.Organization, error) {
	organization := entity.FindOrganizationByName(ctx, name)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", name)
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organization.OrganizationID)
	if !perms.View {
		return nil, gqlerror.Errorf("You are not allowed to view organization %s", name)
	}

	return organization, nil
}

func (r *mutationResolver) InviteUserToOrganization(ctx context.Context, username string, organizationID uuid.UUID, view bool, admin bool) (*entity.User, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	if organization.Personal == true {
		return nil, gqlerror.Errorf("Upgrade to an Enterprise plan to add users to your organization")
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.Admin {
		return nil, gqlerror.Errorf("Not allowed to perform admin functions on organization %s", organizationID.String())
	}

	user := entity.FindUserByUsername(ctx, username)
	if user == nil {
		return nil, gqlerror.Errorf("No user found with that username")
	}

	for _, u := range organization.Users {
		if u.UserID == user.UserID {
			return nil, gqlerror.Errorf("User is already a member of the organization")
		}
	}

	err := organization.InviteUser(ctx, user.UserID, view, admin)
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	// TODO: trigger an email to the invited user, so they can "accept invite to join organization" (which will switch the User.OrganizationID)

	return user, nil
}

func (r *mutationResolver) RemoveUserFromOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (bool, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return false, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.Admin {
		return false, gqlerror.Errorf("Not allowed to perform admin functions in organization %s", organizationID.String())
	}

	if len(organization.Users) < 2 {
		return false, gqlerror.Errorf("Can't remove last member of organization")
	}

	err := organization.RemoveUser(ctx, userID)
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	return true, nil
}

func (r *queryResolver) UsersOrganizationPermissions(ctx context.Context, organizationID uuid.UUID) ([]*entity.PermissionsUsersOrganizations, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.View {
		return nil, gqlerror.Errorf("You are not allowed to view organization %s", organizationID.String())
	}

	permissions := entity.FindOrganizationPermissions(ctx, organizationID)
	if permissions == nil {
		return nil, gqlerror.Errorf("Permissions not found for organization %s", organizationID.String())
	}

	return permissions, nil
}

func (r *mutationResolver) UpdateOrganizationName(ctx context.Context, organizationID uuid.UUID, name string) (*entity.Organization, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.Admin {
		return nil, gqlerror.Errorf("Not allowed to perform admin functions in organization %s", organizationID.String())
	}

	organization, err := organization.ChangeName(ctx, name)
	if err != nil {
		return nil, gqlerror.Errorf("Failed to update organization name")
	}

	return organization, nil
}

func (r *mutationResolver) UpdateUserOrganizationPermissions(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, view *bool, admin *bool) (*entity.PermissionsUsersOrganizations, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.Admin {
		return nil, gqlerror.Errorf("Not allowed to perform admin functions in organization %s", organizationID.String())
	}

	user := entity.FindUser(ctx, userID)
	if user == nil {
		return nil, gqlerror.Errorf("User %s not found", userID.String())
	}

	permissions := entity.FindPermissionsUsersOrganizations(ctx, userID, organizationID)
	if permissions == nil {
		return nil, gqlerror.Errorf("You must invite the user to the organization before editing its permissions. Permissions not found for organization %s and user %s", organizationID.String(), userID.String())
	}

	permissions, err := organization.ChangeUserPermissions(ctx, userID, view, admin)
	if err != nil {
		return nil, gqlerror.Errorf("Failed to update permissions")
	}

	return permissions, nil
}

func (r *mutationResolver) UpdateUserOrganizationQuotas(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, readQuota *int, writeQuota *int) (*entity.User, error) {
	organization := entity.FindOrganization(ctx, organizationID)
	if organization == nil {
		return nil, gqlerror.Errorf("Organization %s not found", organizationID.String())
	}

	secret := middleware.GetSecret(ctx)
	perms := secret.OrganizationPermissions(ctx, organizationID)
	if !perms.Admin {
		return nil, gqlerror.Errorf("Not allowed to perform admin functions in organization %s", organizationID.String())
	}

	user := entity.FindUser(ctx, userID)
	if user == nil {
		return nil, gqlerror.Errorf("User %s not found", userID.String())
	}

	user, err := organization.ChangeUserQuotas(ctx, userID, readQuota, writeQuota)
	if err != nil {
		return nil, gqlerror.Errorf("Failed to update the user's quotas")
	}

	return user, nil
}