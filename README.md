# 🐦‍⬛ CACHE

## Description

This package provides a service cache functionality based on
generic keys and values. The package also provides
functionality for automatic removal of expired records.
Additionally, the user can set the cache size.

## Usage Instructions

For example, to create a new cache that can store string
values and integer keys, add the following code:

```go
newCache, err := cache.New[int, string](&cache.Config{
	TTL:  150 * time.Millisecond,
	Size: 5,
})
```

The `TTL` parameter defines the lifetime of a record,
starting from the moment it is added to the cache. Once it
is exhausted, the record will be removed from the cache.
However, the user can also allow unlimited lifetime for
records in the cache. In this case, the `TTL` parameter
should either not be specified or be equal to `0`.

The `Size` parameter defines the cache size. When trying to
add a new record to a full cache, the oldest record is
removed from the cache and replaced with the added record.
This approach allows controlling the cache size. However,
the user can also allow an unlimited cache size. In this
case, the `Size` parameter should be equal to `0`.

This minimalist design allows the user to store arbitrary
keys and data in the cache. To write an element to the
cache, add the following code:

```go
newCache.Set(1, "value")
```

And to get the value by key, execute:

```go
v, k := newCache.Get(1)
```

Note that the first returned parameter is the value, and the
second is a flag indicating the existence of the element in
the cache.
