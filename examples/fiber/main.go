package main

import (
	"github.com/small-ek/antgo/examples/fiber/boot"
	_ "github.com/small-ek/antgo/frame/serve/fiber"
)

func main() {
	boot.Serve()
	//result := model.Admin{}
	//ant.Db("mysql2").Table("admin").Find(&result)
	//ant.Log().Info("result", zap.String("12", conv.String(result)))

	//eng := ant.New(config).Serve(app)

	//alog.Info("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}
	//defer eng.Close()
}
