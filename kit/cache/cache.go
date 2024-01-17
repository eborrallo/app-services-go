package cache

// Cache defines the expected behaviour from a cache.
type Cache interface {
	// Get Fetch the value of a key.
	Get(Key) ([]byte, error)
	// Set Store a value for a key.
	Set(Key, Value)
}

type Key string
type Value interface{}

//go:generate mockery --case=snake --outpkg=cachemocks --output=cachemocks --name=Cache
