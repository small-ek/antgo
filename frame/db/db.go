package db

import (
	"github.com/small-ek/antgo/os/config"
	"log"
)

type Db struct {
	Name     string
	Type     string
	hostname string
	port     string
	username string
	password string
	database string
	params   string
	log      bool
}

func InitDb() {
	cfg := config.Decode()
	connections := cfg.Get("connections").Maps()
	log.Println(connections)
}
