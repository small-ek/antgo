package ants

import (
	"github.com/panjf2000/ants/v2"
	"github.com/small-ek/antgo/os/config"
)

// NewPool
var NewPool *ants.Pool

func InitPool(count ...int) {
	var err error
	var maxPoolCount = config.GetInt("system.max_pool_count")

	if len(count) > 0 {
		maxPoolCount = count[0]
	}
	NewPool, err = ants.NewPool(maxPoolCount, ants.WithPreAlloc(true))

	if err != nil {
		panic(err)
	}
	return
}
