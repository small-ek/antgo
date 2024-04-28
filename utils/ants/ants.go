package ants

import (
	"github.com/panjf2000/ants/v2"
)

// NewPool accepts the tasks and process them concurrently,
// it limits the total of goroutines to a given number by recycling goroutines.
var NewPool *ants.Pool

func InitPool(count int) {
	var err error
	if count > 0 {
		NewPool, err = ants.NewPool(count, ants.WithPreAlloc(true))
		if err != nil {
			panic(err)
		}
	}

	return
}
