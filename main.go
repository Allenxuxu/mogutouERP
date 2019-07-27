package main

import (
	"flag"
	"log"

	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/pkg/token"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	config "github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

func main() {
	path := flag.String("c", "/etc/conf", "配置文件夹路径")
	flag.Parse()

	token.InitConfig(*path+"/jwt.json", "jwt-key")
	//read config
	fileSource := file.NewSource(
		file.WithPath(*path + "/conf.json"),
	)
	conf := config.NewConfig()
	err := conf.Load(fileSource)
	if err != nil {
		log.Fatal(err)
	}

	var info models.DBInfo
	err = conf.Get("mysql").Scan(&info)
	if err != nil {
		log.Fatal(err)
	}
	models.Init(&info)

	var listenAddr string
	err = conf.Get("listen").Scan(&listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	gin.DisableConsoleColor()
	// gin.SetMode(gin.ReleaseMode)
	r := initRouter()
	r.Run(listenAddr) // listen and serve on 0.0.0.0:8080
}
