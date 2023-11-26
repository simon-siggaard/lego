package user

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

func (s CachedStore) Summary(username string) (brick.User, error) {
	userSummaryURL := brick.Domain + "/api/user/by-username/" + username

	user, err := brick.UnmarshalCachedOr(
		s.redisCache,
		userSummaryURL,
		func() (brick.User, error) {
			return s.store.Summary(username)
		},
	)
	if err != nil {
		return brick.User{}, err
	}

	return user, nil
}

func (s CachedStore) Details(id string) (brick.User, error) {
	userDetailsURL := brick.Domain + "/api/user/by-id/" + id

	user, err := brick.UnmarshalCachedOr(
		s.redisCache,
		userDetailsURL,
		func() (brick.User, error) {
			return s.store.Details(id)
		},
	)
	if err != nil {
		return brick.User{}, err
	}

	return user, nil
}

func (s CachedStore) All() ([]brick.User, error) {
	summaryURL := brick.Domain + "/api/users"
	detailsURL := brick.Domain + "/api/user/by-id"

	users, err := brick.UnmarshalCachedOr(
		s.redisCache,
		summaryURL,
		func() ([]brick.User, error) {
			return s.store.All()
		},
	)
	if err != nil {
		return nil, err
	}

	for n, user := range users {
		user := user // no longer necessary in Go 1.22
		url := detailsURL + "/" + user.ID

		user, err = brick.UnmarshalCachedOr(
			s.redisCache,
			url,
			func() (brick.User, error) {
				return s.store.Details(user.ID)
			},
		)
		if err != nil {
			return nil, err
		}

		users[n] = user
	}

	return users, nil
}
