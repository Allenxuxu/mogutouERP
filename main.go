package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Allenxuxu/mogutouERP/models"
	"github.com/Allenxuxu/mogutouERP/pkg/token"
	"github.com/gin-gonic/gin"
	config "github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/pkg/browser"
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

	gin.DisableConsoleColor()
	// gin.SetMode(gin.ReleaseMode)
	r := initRouter()
	go func() {
		time.Sleep(time.Second)
		browser.OpenURL("http://127.0.0.1:1988/ui")
		fmt.Println("Open: http://127.0.0.1:1988/ui")
	}()

	r.Run(":1988") // listen and serve on 0.0.0.0:8080
}
