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

package stash

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
	"time"
)

type struct1 struct {
	Foo string
	Bar bool
	Baz []byte
}

type struct2 struct {
	Foo string
	S1  struct1
}

func makeTempFilename() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%stest-%d", os.TempDir(), rand.Int())
}

func TestEmptyFileGetsCreated(t *testing.T) {
	filename := makeTempFilename()
	defer os.Remove(filename)

	_, err := NewStash(filename, true)
	assert.Nil(t, err)

	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestFlush(t *testing.T) {
	filename := makeTempFilename()
	defer os.Remove(filename)

	s, err := NewStash(filename, false)
	assert.Nil(t, err)

	s.Save("foo", "bar")

	_, err = os.Stat(filename)
	assert.True(t, os.IsNotExist(err))

	s.Flush()
	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestEmptyFileWriteThenRead(t *testing.T) {
	filename := makeTempFilename()
	defer os.Remove(filename)

	jd, err := NewStash(filename, true)
	require.Nil(t, err)

	const key1 = "KEY1"
	const key2 = "KEY2"

	s1 := struct1{
		Bar: true,
		Baz: []byte("testing123"),
		Foo: "Hello, World!",
	}

	s2 := "Hello"

	jd.Save(key1, s1)
	jd.Save(key2, s2)

	jd2, err := NewStash(filename, true)
	require.Nil(t, err)

	var s1x struct1
	err = jd2.Read(key1, &s1x)
	require.Nil(t, err)

	var s2x string
	err = jd2.Read(key2, &s2x)
	require.Nil(t, err)

	assert.Equal(t, s1, s1x)
	assert.Equal(t, s2, s2x)
}
