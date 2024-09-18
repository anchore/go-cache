package cache

import (
	"bytes"
	"encoding/json"

	"github.com/anchore/go-logger"
)

// Resolver interface provides a single Resolve method, which will return from cache
// or call the provided resolve function to get the value if not available in cache
type Resolver[T any] interface {
	// Resolve attempts to resolve the given key from cache and convert it to the type of the cache,
	// or calls the resolver function if unable to resolve a cached value
	Resolve(key string, resolver resolverFunc[T]) (T, error)
}

// NewResolver returns a cache resolver for persistent cached data across Syft runs, stored in a unique
// location based on the provided name and versioned by the type
func NewResolver[T any](cache Cache) Resolver[T] {
	return &cacheResolver[T]{
		cache: cache,
	}
}

const resolverKeySuffix = ".json"

type resolverFunc[T any] func() (T, error)

type cacheResolver[T any] struct {
	log   logger.Logger
	cache Cache
}

var _ interface {
	Resolver[int]
} = (*cacheResolver[int])(nil)

func (r *cacheResolver[T]) Resolve(key string, resolver resolverFunc[T]) (T, error) {
	key += resolverKeySuffix

	rdr, err := r.cache.Read(key)
	if rdr == nil || err != nil {
		return r.resolveAndCache(key, resolver)
	}
	defer func() {
		err = rdr.Close()
		if err != nil {
			r.log.Debug(err)
		}
	}()

	dec := json.NewDecoder(rdr)
	if dec == nil {
		r.log.Tracef("error getting cache json decoder for %v %v: %v", r.cache, key, err)
		return r.resolveAndCache(key, resolver)
	}
	var t T
	err = dec.Decode(&t)
	if err != nil {
		r.log.Tracef("error decoding cached entry for %v %v: %v", r.cache, key, err)
		return r.resolveAndCache(key, resolver)
	}
	// no error, able to resolve from cache
	return t, nil
}

func (r *cacheResolver[T]) resolveAndCache(key string, resolver func() (T, error)) (T, error) {
	t, err := resolver()
	if err != nil {
		return t, err
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	enc.SetEscapeHTML(false)
	err = enc.Encode(t)
	if err != nil {
		return t, err
	}
	err = r.cache.Write(key, &data)
	return t, err
}
