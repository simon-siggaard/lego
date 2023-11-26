package set

import (
	"github.com/simon-siggaard/lego/pkg/brick"
	"github.com/simon-siggaard/lego/pkg/brick/cache"
)

type CachedStore struct {
	redisCache *cache.RedisClient
	store      Store
}

func NewCachedStore(store Store) CachedStore {
	return CachedStore{
		redisCache: cache.NewRedisClient(),
		store:      store,
	}
}

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
