package set

import (
	"github.com/simon-siggaard/lego/pkg/brick"
	"github.com/simon-siggaard/lego/pkg/brick/cache"
)

// CachedStore is a cached implementation of Store.
type CachedStore struct {
	redisCache *cache.RedisClient
	store      Store
}

// NewCachedStore returns a new cached LEGO set store.
func NewCachedStore(store Store) CachedStore {
	return CachedStore{
		redisCache: cache.NewRedisClient(),
		store:      store,
	}
}

// Summaries returns LEGO set summaries.
func (s CachedStore) Summaries() ([]brick.Set, error) {
	summaryURL := brick.Domain + "/api/sets"

	sets, err := brick.UnmarshalCachedOr(
		s.redisCache,
		summaryURL,
		func() ([]brick.Set, error) {
			return s.store.Summaries()
		},
	)
	if err != nil {
		return nil, err
	}

	return sets, nil
}

// Details returns the details of a LEGO set.
func (s CachedStore) Details(id string) (brick.Set, error) {
	detailsURL := brick.Domain + "/api/set/by-id/" + id

	set, err := brick.UnmarshalCachedOr(
		s.redisCache,
		detailsURL,
		func() (brick.Set, error) {
			return s.store.Details(id)
		},
	)
	if err != nil {
		return brick.Set{}, err
	}

	return set, nil
}
