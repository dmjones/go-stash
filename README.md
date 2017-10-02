# go-stash

Go-stash provides an in-memory data store, backed by a file on disk. It's designed for simple object storage when you don't need the overhead of a proper database.


Stash instances are safe for concurrent use by multiple goroutines. Data is serialized to disk using `json.Marshal()`, please see [the JSON docs](https://golang.org/pkg/encoding/json/#Marshal) to understand how values are encoded to JSON.

## Installation

```bash
$ go get -u github.com/dmjones500/go-stash/stash
```

## Documentation

https://godoc.org/github.com/dmjones500/go-stash/stash

## Quick Start

Create a new Stash, which auto-flushes to disk:

```Go
stash, err := stash.NewStash(filename, true)
if err != nil {
    // handle error
}
```

Save a complex structure:

```Go
err = stash.Save("accountData", account)
```

Retrieve it to another variable:

```Go
var anotherAccount Account
err = stash.Read("accountData", &anotherAccount)
```
	
If you don't enable auto-flush, then call `Flush()` when you want to persist all data to disk:

```Go
err = stash.Flush()
```

## License

Go-stash is licensed under the [MIT License](https://opensource.org/licenses/MIT).
