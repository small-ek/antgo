package main

import (
	"github.com/small-ek/antgo/examples/fiber/boot"
	"github.com/small-ek/antgo/examples/gin/model"
	"github.com/small-ek/antgo/frame/ant"
	_ "github.com/small-ek/antgo/frame/serve/fiber"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
)

func main() {
	boot.Serve()
	result := model.Admin{}
	ant.Db("mysql2").Table("s_admin").Find(&result)
	ant.Log().Info("result", zap.String("12", conv.String(result)))

	//eng := ant.New(config).Serve(app)

	//alog.Info("main", zap.Any("result", result))
	//tt := Test{Name: "22121"}
	//for i := 0; i < 10; i++ {
	//	ant.Log().Info("222121212=============================" + conv.String(i))
	//}
	//defer eng.Close()
}
