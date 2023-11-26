package brick

import (
	"encoding/json"

	"github.com/simon-siggaard/lego/pkg/brick/cache"
)

// UnmarshalCachedOr unmarshals the cached value of key or calls orFunc and caches and returns the result.
func UnmarshalCachedOr[T any](
	rdb *cache.RedisClient,
	key string,
	orFunc func() (T, error),
) (T, error) {
	var t T

	cached, err := rdb.Get(key)
	if err != nil {
		sets, err := orFunc()
		if err != nil {
			return t, err
		}

		bs, err := json.Marshal(sets)
		if err != nil {
			return t, err
		}

		err = rdb.Set(key, bs)
		if err != nil {
			return t, err
		}

		return sets, nil
	}

	err = json.Unmarshal(cached, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}
