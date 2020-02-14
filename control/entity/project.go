package entity

import (
	"context"
	"regexp"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/cache/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/go-playground/validator.v9"

	"github.com/beneath-core/internal/hub"
)

// Project represents a Beneath project
type Project struct {
	ProjectID      uuid.UUID `sql:",pk,type:uuid,default:uuid_generate_v4()"`
	Name           string    `sql:",notnull",validate:"required,gte=3,lte=16"`
	DisplayName    string    `validate:"omitempty,lte=40"`
	Site           string    `validate:"omitempty,url,lte=255"`
	Description    string    `validate:"omitempty,lte=255"`
	PhotoURL       string    `validate:"omitempty,url,lte=255"`
	Public         bool      `sql:",notnull,default:true"`
	Locked         bool      `sql:",notnull,default:false"`
	OrganizationID uuid.UUID `sql:",on_delete:restrict,notnull,type:uuid"`
	Organization   *Organization
	CreatedOn      time.Time `sql:",default:now()"`
	UpdatedOn      time.Time `sql:",default:now()"`
	DeletedOn      time.Time
	Streams        []*Stream
	Models         []*Model
	Users          []*User `pg:"many2many:permissions_users_projects,fk:project_id,joinFK:user_id"`
}

var (
	// regex used in validation
	projectNameRegex *regexp.Regexp

	// redis cache for project data
	projectCache *cache.Codec
)

func init() {
	projectNameRegex = regexp.MustCompile("^[_a-z][_a-z0-9]*$")
	GetValidator().RegisterStructValidation(validateProject, Project{})
}

// custom project validation
func validateProject(sl validator.StructLevel) {
	p := sl.Current().Interface().(Project)

	if !projectNameRegex.MatchString(p.Name) {
		sl.ReportError(p.Name, "Name", "", "alphanumericorunderscore", "")
	}
}

// FindProject finds a project by ID
func FindProject(ctx context.Context, projectID uuid.UUID) *Project {
	project := &Project{
		ProjectID: projectID,
	}
	err := hub.DB.ModelContext(ctx, project).WherePK().Column("project.*", "Streams", "Users").Select()
	if !AssertFoundOne(err) {
		return nil
	}
	return project
}

// FindProjects returns a sample of projects
func FindProjects(ctx context.Context) []*Project {
	var projects []*Project
	err := hub.DB.ModelContext(ctx, &projects).Where("project.public = true").Limit(200).Order("name").Select()
	if err != nil {
		panic(err)
	}
	return projects
}

// FindProjectByName finds a project by name
func FindProjectByName(ctx context.Context, name string) *Project {
	project := &Project{}
	err := hub.DB.ModelContext(ctx, project).
		Where("lower(project.name) = lower(?)", name).
		Column("project.*", "Streams", "Users").
		Select()
	if !AssertFoundOne(err) {
		return nil
	}
	return project
}

// GetProjectID implements engine/driver.Project
func (p *Project) GetProjectID() uuid.UUID {
	return p.ProjectID
}

// GetProjectName implements engine/driver.Project
func (p *Project) GetProjectName() string {
	return p.Name
}

// GetPublic implements engine/driver.Project
func (p *Project) GetPublic() bool {
	return p.Public
}

// CreateWithUser creates a project and makes user a member
func (p *Project) CreateWithUser(ctx context.Context, userID uuid.UUID, view bool, create bool, admin bool) error {
	// validate
	err := GetValidator().Struct(p)
	if err != nil {
		return err
	}

	// create project and PermissionsUsersProjects in one transaction
	return hub.DB.WithContext(ctx).RunInTransaction(func(tx *pg.Tx) error {
		// insert project
		_, err := tx.Model(p).Insert()
		if err != nil {
			return err
		}

		// connect project to userID
		err = tx.Insert(&PermissionsUsersProjects{
			UserID:    userID,
			ProjectID: p.ProjectID,
			View:      view,
			Create:    create,
			Admin:     admin,
		})
		if err != nil {
			return err
		}

		err = hub.Engine.RegisterProject(ctx, p)
		if err != nil {
			return err
		}

		return nil
	})
}

// AddUser makes user a member of project
func (p *Project) AddUser(ctx context.Context, userID uuid.UUID, view bool, create bool, admin bool) error {
	return hub.DB.WithContext(ctx).Insert(&PermissionsUsersProjects{
		UserID:    userID,
		ProjectID: p.ProjectID,
		View:      view,
		Create:    create,
		Admin:     admin,
	})
}

// RemoveUser removes a member from the project
func (p *Project) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	// clear cache
	getUserProjectPermissionsCache().Clear(ctx, userID, p.ProjectID)

	// TODO only if not last user (there's a check in resolver, but it should be part of db tx)
	return hub.DB.WithContext(ctx).Delete(&PermissionsUsersProjects{
		UserID:    userID,
		ProjectID: p.ProjectID,
	})
}

// UpdateDetails updates projects user-facing details
func (p *Project) UpdateDetails(ctx context.Context, displayName *string, public *bool, site *string, description *string, photoURL *string) error {
	// set fields
	if displayName != nil {
		p.DisplayName = *displayName
	}
	if public != nil {
		p.Public = *public
	}
	if site != nil {
		p.Site = *site
	}
	if description != nil {
		p.Description = *description
	}
	if photoURL != nil {
		p.PhotoURL = *photoURL
	}

	// validate
	err := GetValidator().Struct(p)
	if err != nil {
		return err
	}

	// update in tx with call to bigquery
	return hub.DB.WithContext(ctx).RunInTransaction(func(tx *pg.Tx) error {
		p.UpdatedOn = time.Now()
		_, err = hub.DB.WithContext(ctx).Model(p).
			Column("display_name", "public", "site", "description", "photo_url", "updated_on").
			WherePK().
			Update()
		if err != nil {
			return err
		}

		// update in warehouse
		err = hub.Engine.RegisterProject(ctx, p)
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateOrganization changes a project's organization
func (p *Project) UpdateOrganization(ctx context.Context, organizationID uuid.UUID) (*Project, error) {
	p.OrganizationID = organizationID
	p.UpdatedOn = time.Now()

	_, err := hub.DB.ModelContext(ctx, p).Column("organization_id", "updated_on").WherePK().Update()

	// re-find project so we can return the name of the new organization
	err = hub.DB.ModelContext(ctx, p).
		WherePK().
		Column("project.*", "Organization").
		Select()
	if !AssertFoundOne(err) {
		return nil, err
	}

	return p, err
}

// SetLock sets a project's "locked" status
func (p *Project) SetLock(ctx context.Context, isLocked bool) error {
	p.Locked = isLocked
	p.UpdatedOn = time.Now()

	_, err := hub.DB.ModelContext(ctx, p).
		Column("locked", "updated_on").
		WherePK().
		Update()
	if err != nil {
		return err
	}

	return nil
}

// Delete safely deletes the project (fails if the project still has content)
func (p *Project) Delete(ctx context.Context) error {
	return hub.DB.WithContext(ctx).RunInTransaction(func(tx *pg.Tx) error {
		err := hub.DB.WithContext(ctx).Delete(p)
		if err != nil {
			return err
		}

		err = hub.Engine.RemoveProject(ctx, p)
		if err != nil {
			return err
		}

		return nil
	})
}

func getProjectCache() *cache.Codec {
	if projectCache == nil {
		projectCache = &cache.Codec{
			Redis:     hub.Redis,
			Marshal:   msgpack.Marshal,
			Unmarshal: msgpack.Unmarshal,
		}
	}
	return projectCache
}