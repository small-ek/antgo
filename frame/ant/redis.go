package ant

import (
	"github.com/small-ek/antgo/db/aredis"
	"github.com/small-ek/antgo/os/config"
)

// initRedis
func initRedis() {
	address := config.GetString("redis.address")
	db := config.GetInt("redis.db")

	if address != "" && db >= 0 {
		password := config.GetString("redis.password")
		aredis.New(address, password, db)
	}
}

// Redis
func Redis() *aredis.ClientRedis {
	return aredis.Client
}
