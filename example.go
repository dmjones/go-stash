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
