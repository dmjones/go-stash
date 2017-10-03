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

package main

import (
	"fmt"
	"github.com/dmjones500/go-stash/stash"
	"math/rand"
	"os"
	"time"
)

type Transaction struct {
	Timestamp int64
	Amount    int
}

type Account struct {
	Balance      int
	Name         string
	Transactions []Transaction
}

func main() {
	transactions := []Transaction{
		Transaction{time.Now().Unix(), 999},
		Transaction{time.Now().Unix(), 1295},
	}

	account := Account{145187, "checking", transactions}

	filename := makeTempFilename()
	defer os.Remove(filename)

	stash, err := stash.NewStash(filename, true)
	if err != nil {
		panic(err)
	}

	err = stash.Save("accountData", account)
	if err != nil {
		panic(err)
	}

	var anotherAccount Account
	err = stash.Read("accountData", &anotherAccount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", anotherAccount)
}

func makeTempFilename() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%stest-%d", os.TempDir(), rand.Int())
}
