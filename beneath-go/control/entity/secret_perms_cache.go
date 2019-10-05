package entity

import (
	"bytes"
	"context"
	"encoding/gob"
	"reflect"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-redis/cache/v7"
	uuid "github.com/satori/go.uuid"

	"github.com/beneath-core/beneath-go/db"
)

// PermissionsCache encapsulates the secret's permissions and how to access the permissions in short-term and long-term memory
type PermissionsCache struct {
	codec     *cache.Codec
	prototype reflect.Type
	query     string
}

var (
	userOrganizationPermissions *PermissionsCache
	userProjectPermissions      *PermissionsCache
	serviceStreamPermissions    *PermissionsCache
)

var permsCacheConfig = struct {
	cacheTime    time.Duration
	cacheLRUTime time.Duration
	cacheLRUSize int
	redisKeyFn   func(ownerID uuid.UUID, resourceID uuid.UUID) string
}{
	cacheTime:    time.Hour,
	cacheLRUTime: 1 * time.Minute,
	cacheLRUSize: 20000,
	redisKeyFn: func(ownerID uuid.UUID, resourceID uuid.UUID) string {
		return string(append(ownerID.Bytes(), resourceID.Bytes()...))
	},
}

func getUserOrganizationPermissionsCache() *PermissionsCache {
	if userOrganizationPermissions == nil {
		userOrganizationPermissions = NewPermissionsCache(OrganizationPermissions{}, `
			select p.view, p.admin
			from permissions_users_organizations p
			where p.user_id = ? and p.organization_id = ?
		`)
	}

	return userOrganizationPermissions
}
func getUserProjectPermissionsCache() *PermissionsCache {
	if userProjectPermissions == nil {
		userProjectPermissions = NewPermissionsCache(ProjectPermissions{}, `
			select p.view, p.create, p.admin
			from permissions_users_projects p
			where p.user_id = ? and p.project_id = ?
		`)
	}

	return userProjectPermissions
}
func getServiceStreamPermissionsCache() *PermissionsCache {
	if serviceStreamPermissions == nil {
		serviceStreamPermissions = NewPermissionsCache(StreamPermissions{}, `
			select p.read, p.write
			from permissions_services_streams p
			where p.service_id = ? and p.stream_id = ?
		`)
	}

	return serviceStreamPermissions
}

// CachedUserOrganizationPermissions returns organization permissions for a given owner-resource combo
func CachedUserOrganizationPermissions(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) OrganizationPermissions {
	return getUserOrganizationPermissionsCache().Get(ctx, userID, organizationID).(OrganizationPermissions)
}

// CachedUserProjectPermissions returns project permissions for a given owner-resource combo
func CachedUserProjectPermissions(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) ProjectPermissions {
	return getUserProjectPermissionsCache().Get(ctx, userID, projectID).(ProjectPermissions)
}

// CachedServiceStreamPermissions returns stream permissions for a given owner-resource combo
func CachedServiceStreamPermissions(ctx context.Context, serviceID uuid.UUID, streamID uuid.UUID) StreamPermissions {
	return getServiceStreamPermissionsCache().Get(ctx, serviceID, streamID).(StreamPermissions)
}

// NewPermissionsCache initializes a PermissionCache object for a given prototype (organization/project/stream)
func NewPermissionsCache(prototype interface{}, query string) *PermissionsCache {
	pm := &PermissionsCache{}
	pm.prototype = reflect.TypeOf(prototype)
	pm.query = query
	pm.codec = &cache.Codec{
		Redis:     db.Redis,
		Marshal:   pm.marshal,
		Unmarshal: pm.unmarshal,
	}
	pm.codec.UseLocalCache(permsCacheConfig.cacheLRUSize, permsCacheConfig.cacheLRUTime)
	return pm
}

// Get fetches permissions by applying the cached query to the given parameters
func (c *PermissionsCache) Get(ctx context.Context, ownerID uuid.UUID, resourceID uuid.UUID) interface{} {
	res := reflect.New(c.prototype)
	err := c.codec.Once(&cache.Item{
		Key:        permsCacheConfig.redisKeyFn(ownerID, resourceID),
		Object:     res.Interface(),
		Expiration: permsCacheConfig.cacheTime,
		Func:       c.getterFunc(ctx, ownerID, resourceID),
	})
	if err != nil {
		panic(err)
	}
	return res.Elem().Interface()
}

func (c PermissionsCache) marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c PermissionsCache) unmarshal(b []byte, v interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func (c PermissionsCache) getterFunc(ctx context.Context, ownerID uuid.UUID, resourceID uuid.UUID) func() (interface{}, error) {
	return func() (interface{}, error) {
		res := reflect.New(c.prototype)
		_, err := db.DB.QueryContext(ctx, res.Interface(), c.query, ownerID, resourceID)
		if err != nil && err != pg.ErrNoRows {
			panic(err)
		}
		return res.Elem().Interface(), nil
	}
}
