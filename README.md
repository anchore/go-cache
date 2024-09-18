# go-cache

A basic caching library.

This library provides the interfaces used for persistent caching by the Anchore Go CLI tools, including
Syft and Grype.

Each application should provide means for configuring a `cache.Manager` and
individual caches should be obtained from the `cache.Manager`.

## Usage

To cache specific Go data types, the easiest method is using a `cache.Resolver` with the with automatic type version
based on a `cache.HashType` call, which automatically creates a unique version qualifier based on the _structure_ of the provided type.
This way, if the structure changes in any way it will end up with a new version key which will invalidate older cache entries
and result in populating the cache based on this new key.
The `cache.Resolver` will store items using the `json` package to serialize/deserialize values, so to save space
it is encouraged  to use `omitempty`. For example:

```go
type myCacheItem struct {
	Name string `json:"name",omitempty`
}
```

If it is common that checking for an item will result in errors, and you do not want to re-run the resolve function
when errors are encountered, instead of using `GetResolver`, you can use `GetResolverCachingErrors`, which is useful
for things such as resolving artifacts over a network, where a number of them will not be resolved, and you do not want
to continue to have the expense of running the network resolution. This should be used when it is acceptable a network
outage and cached errors is an acceptable risk.

An example can be seen in the [Syft golang cataloger](https://github.com/anchore/syft/blob/main/syft/pkg/cataloger/golang/licenses.go) fetching remote licenses.

## Example

```golang
package appcache

import "github.com/anchore/go-cache"

// default to a bypassed cache, so unit tests are easy to deal with
var manager = cache.NewBypassed()

// set the cache after any necessary configuration
func InitCache(dir string, ttl time.Duration) {
	manager = cache.NewFromDir(globalLogger, dir, ttl)
}

// utility function with an app-scoped global cache
func NewResolver[T any](name string) cache.Resolver[T] {
	return cache.NewResolver[T](manager.GetCache(name, cache.HashType[T]()))
}
```

```golang
package app

import "appcache"

type myDataType struct {
	Value string `json:"value,omitempty"`
}

var cacheResolver = appcache.NewResolver[myDataType]("my-top-level-cache-name")

func functionWhichCaches(someParam string) myDataType {
	return cacheResolver.Resolve(someParam, func() myDataType {
		// do some work and return 
		return myDataType{ ... }
    })
}
```
