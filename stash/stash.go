// Copyright 2017 Duncan Jones

// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following
// conditions:

// The above copyright notice and this permission notice shall be included in all copies
// or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package stash provides a basic in-memory data store backed by a file on disk.
package stash

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"sync"
)

const (
	version1 = 1
)

// UnknownVersionError indicates an unsupported version number tag was found in the data
type UnknownVersionError struct {
	badVersion int
}

func (e UnknownVersionError) Error() string {
	return fmt.Sprintf("unsupported version number %d", e.badVersion)
}

// NoSuchKeyError indicates that a key does not exist in the database
type NoSuchKeyError struct {
	s string
}

func (e NoSuchKeyError) Error() string {
	return fmt.Sprintf("no such key: %s", e.s)
}

// Stash is a simple in-memory data store, backed by a file on disk. Create a Stash by calling
// the NewStash factory method. It is safe for multiple goroutines to call a Stash's methods
// concurrently.
type Stash struct {
	mutex     *sync.Mutex // protects access to the file
	file      string
	version   int
	autoFlush bool
	data      interface{}
}

// container is used when writing to disk, to store the data format version
// alongside the marshalled data.
type container struct {
	Version int
	Data    json.RawMessage
}

// v1Data is the version 1 data format - a simple map of strings to marshalled JSON data.
type v1Data map[string]json.RawMessage

// Save associates the value with the key in the data store, overwriting
// any previous value. If auto-flush is enabled, each call to Save will
// be persisted to disk immediately. Otherwise, Flush must be called.
//
// Values are stored using JSON marshalling, which means unexported fields
// will not be saved. See the documentation for the json package for more
// information.
func (s *Stash) Save(key string, value interface{}) error {
	switch s.version {
	case version1:
		marshalledData, err := json.Marshal(value)
		if err != nil {
			return errors.Wrap(err, "error marshalling value")
		}
		s.mutex.Lock()
		s.data.(v1Data)[key] = marshalledData
		s.mutex.Unlock()

		if s.autoFlush {
			return s.Flush()
		} else {
			return nil
		}
	default:
		return UnknownVersionError{s.version}
	}
}

// Read will store the value associated with the key into the
// variable pointed to by ptr.
//
//   var foo MyStruct
//   err = jd2.Read("myKey", &foo)
//   if err != nil {
//     ...
//   }
func (s *Stash) Read(key string, ptr interface{}) error {
	switch s.version {
	case version1:
		data := s.data.(v1Data)
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if item, ok := data[key]; ok {
			return json.Unmarshal(item, ptr)
		} else {
			return NoSuchKeyError{""}
		}

	default:
		return UnknownVersionError{s.version}
	}
}

// Flush writes the content of the in-memory database to disk. There
// is no need to call Flush if auto-flushing is enabled.
func (s *Stash) Flush() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	jsonData, err := json.Marshal(s.data)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal data")
	}

	container := container{Version: s.version, Data: jsonData}
	jsonFileData, err := json.Marshal(container)

	err = ioutil.WriteFile(s.file, jsonFileData, 0600)
	return errors.WithMessage(err, fmt.Sprintf("failed to write database to '%s'", s.file))
}

// readFromDisk reads the contents of jd.file into memory. This function will
// return an error if the file is not a Stash file.
func (s *Stash) readFromDisk() error {
	data, err := ioutil.ReadFile(s.file)
	if err != nil {
		return err
	}

	var container container
	err = json.Unmarshal(data, &container)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal outer data structure")
	}

	s.version = container.Version

	switch s.version {
	case version1:
		v1data := v1Data{}
		err = json.Unmarshal(container.Data, &v1data)
		if err != nil {
			return errors.Wrap(err, "failed to unwrap v1 data")
		}
		s.data = v1data
		return nil
	default:
		return UnknownVersionError{s.version}
	}
}

// NewStash constructs a new Stash, backed by the specified file on disk. If autoFlush is
// enabled, every call to Save will be automatically followed by a call to Flush, which writes
// the data store to disk.
//
// If filename points at an existing file, it is assumed to be a Stash file and is
// read into memory. If the file does not yet exist and autoFlush is enabled, an empty
// data store will be written to disk.
func NewStash(filename string, autoFlush bool) (*Stash, error) {
	result := Stash{file: filename, mutex: &sync.Mutex{}, autoFlush: autoFlush}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// new database
		result.version = version1
		result.data = v1Data(make(map[string]json.RawMessage))
		if autoFlush {
			return &result, result.Flush()
		} else {
			return &result, nil
		}
	} else {
		// existing database
		return &result, result.readFromDisk()
	}
}
