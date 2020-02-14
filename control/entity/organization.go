package entity

import (
	"context"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"

	"github.com/beneath-core/internal/hub"
)

// Organization represents the entity that manages billing on behalf of its users
type Organization struct {
	OrganizationID uuid.UUID `sql:",pk,type:uuid,default:uuid_generate_v4()"`
	Name           string    `sql:",unique,notnull",validate:"required,gte=1,lte=40"`
	Personal       bool      `sql:",notnull"`
	CreatedOn      time.Time `sql:",default:now()"`
	UpdatedOn      time.Time `sql:",default:now()"`
	DeletedOn      time.Time
	Projects       []*Project
	Services       []*Service
	Users          []*User
}

var (
	// used for validation
	organizationNameRegex *regexp.Regexp
)

func init() {
	// configure validation
	organizationNameRegex = regexp.MustCompile("^[_a-z][_a-z0-9]*$")
	GetValidator().RegisterStructValidation(organizationValidation, Organization{})
}

// custom stream validation
func organizationValidation(sl validator.StructLevel) {
	s := sl.Current().Interface().(Organization)

	if !organizationNameRegex.MatchString(s.Name) {
		sl.ReportError(s.Name, "Name", "", "alphanumericorunderscore", "")
	}
}

// FindOrganization finds a organization by ID
func FindOrganization(ctx context.Context, organizationID uuid.UUID) *Organization {
	organization := &Organization{
		OrganizationID: organizationID,
	}
	err := hub.DB.ModelContext(ctx, organization).WherePK().Column("organization.*", "Projects", "Services", "Users").Select()
	if !AssertFoundOne(err) {
		return nil
	}
	return organization
}

// FindOrganizationByName finds a organization by name
func FindOrganizationByName(ctx context.Context, name string) *Organization {
	organization := &Organization{}
	err := hub.DB.ModelContext(ctx, organization).
		Where("lower(name) = lower(?)", name).
		Column("organization.*", "Projects", "Services", "Users").
		Select()
	if !AssertFoundOne(err) {
		return nil
	}
	return organization
}

// FindAllOrganizations returns all active organizations
func FindAllOrganizations(ctx context.Context) []*Organization {
	var organizations []*Organization
	err := hub.DB.ModelContext(ctx, &organizations).
		Column("organization.*").
		Select()
	if err != nil {
		panic(err)
	}
	return organizations
}

// UpdatePersonalStatus updates an organization's "personal" status
// this is called when a user changes plans to/from an Enterprise plan
func (o *Organization) UpdatePersonalStatus(ctx context.Context, personal bool) error {
	org := &Organization{
		OrganizationID: o.OrganizationID,
		Personal:       personal,
		UpdatedOn:      time.Now(),
	}

	_, err := hub.DB.ModelContext(ctx, org).
		Column("personal", "updated_on").
		WherePK().
		Update()
	if err != nil {
		return err
	}

	// TODO: if upgrading to a non-personal plan, trigger an email to the user:
	// -- if prevBillingPlan.Personal && !newBillingPlan.Personal, need to prompt the user to
	// -- 1) change the name of their organization, so its the enterprise's name, not the original user's name
	// -- 2) add teammates (which will call "InviteUser"; then the user will have to accept the invite, which will call "JoinOrganization")

	return nil
}

// InviteUser gives a user permission to join the organization
// the user must then "JoinOrganization" to officially change their membership
func (o *Organization) InviteUser(ctx context.Context, userID uuid.UUID, view bool, admin bool) error {
	// clear cache
	getUserOrganizationPermissionsCache().Clear(ctx, userID, o.OrganizationID)

	perms := &PermissionsUsersOrganizations{
		UserID:         userID,
		OrganizationID: o.OrganizationID,
		View:           view,
		Admin:          admin,
	}

	_, err := hub.DB.ModelContext(ctx, perms).OnConflict("(user_id, organization_id) DO UPDATE").Insert()
	if err != nil {
		return err
	}

	// done
	return nil
}

// RemoveUser removes a member from the organization
func (o *Organization) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	// TODO: only if not last user (there's a check in resolver, but it should be part of db tx)

	// clear cache
	getUserOrganizationPermissionsCache().Clear(ctx, userID, o.OrganizationID)

	// delete
	return hub.DB.WithContext(ctx).Delete(&PermissionsUsersOrganizations{
		UserID:         userID,
		OrganizationID: o.OrganizationID,
	})
}

// ChangeName changes an organization's name
func (o *Organization) ChangeName(ctx context.Context, name string) (*Organization, error) {
	o.Name = name
	o.UpdatedOn = time.Now()

	_, err := hub.DB.ModelContext(ctx, o).
		Column("name").
		WherePK().
		Update()
	if err != nil {
		return nil, err
	}

	return o, nil
}

// ChangeUserPermissions changes a user's permissions within the organization
// Q: this won't fuck up the object "p" that I have, right?
// A: No, it just invalidates it from the cache so other processes won't be served it
// A2: But, we should actually move it to after the DB update for consistency
// Q: should this really be on the Organization struct (vs the Permissions struct)? it leads to a double "FindPermissions" call, as the resolver calls it right before calling this function
func (o *Organization) ChangeUserPermissions(ctx context.Context, userID uuid.UUID, view *bool, admin *bool) (*PermissionsUsersOrganizations, error) {
	permissions := FindPermissionsUsersOrganizations(ctx, userID, o.OrganizationID)

	getUserOrganizationPermissionsCache().Clear(ctx, permissions.UserID, permissions.OrganizationID)

	if view != nil {
		permissions.View = *view
	}
	if admin != nil {
		permissions.Admin = *admin
	}

	_, err := hub.DB.ModelContext(ctx, permissions).
		Column("view", "admin").
		WherePK().
		Update()

	return permissions, err
}

// ChangeUserQuotas allows an admin to configure a member's quotas
// Q: do I have to invalidate the getSecretCache?
// Q: should this really be on the Organization struct (vs the User struct)? it leads to a double "FindUser" call, as the resolver calls it right before calling this function
func (o *Organization) ChangeUserQuotas(ctx context.Context, userID uuid.UUID, readQuota *int, writeQuota *int) (*User, error) {
	user := FindUser(ctx, userID)

	if readQuota != nil {
		user.ReadQuota = int64(*readQuota)
	}
	if writeQuota != nil {
		user.WriteQuota = int64(*writeQuota)
	}

	user.UpdatedOn = time.Now()

	_, err := hub.DB.ModelContext(ctx, user).
		Column("read_quota", "write_quota", "updated_on").
		WherePK().
		Update()
	if err != nil {
		return nil, err
	}

	return user, err
}

// Delete deletes the organization
// this will fail if there are any users, services, or projects that are still tied to the organization
func (o *Organization) Delete(ctx context.Context) error {
	err := hub.DB.WithContext(ctx).Delete(o)
	if err != nil {
		return err
	}

	return nil
}

// FindOrganizationPermissions retrieves all users' permissions for a given organization
func FindOrganizationPermissions(ctx context.Context, organizationID uuid.UUID) []*PermissionsUsersOrganizations {
	var permissions []*PermissionsUsersOrganizations
	err := hub.DB.ModelContext(ctx, &permissions).Where("permissions_users_organizations.organization_id = ?", organizationID).Column("permissions_users_organizations.*", "User", "Organization").Select()
	if err != nil {
		panic(err)
	}

	return permissions
}