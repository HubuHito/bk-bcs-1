package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-webconsole/console/web"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-webconsole/handler"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-webconsole/route"

	"github.com/Tencent/bk-bcs/bcs-common/common/ssl"
	yaml "github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	etcd "github.com/asim/go-micro/plugins/registry/etcd/v4"
	mhttp "github.com/asim/go-micro/plugins/server/http/v4"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/urfave/cli/v2"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	service = "bcs-webconsole"
	version = "latest"
)

func main() {
	// Create service
	// srv := micro.NewService(
	// 	micro.Name(service),
	// 	micro.Version(version),
	// )

	var configPath string

	// new yaml encoder
	enc := yaml.NewEncoder()
	// new config
	conf, _ := config.NewConfig(
		config.WithReader(
			json.NewReader( // json reader for internal config merge
				reader.WithEncoder(enc),
			),
		),
	)

	confFlags := micro.Flags(
		&cli.StringFlag{
			Name:        "bcs-conf",
			Usage:       "bcs-config path",
			Required:    true,
			Destination: &configPath,
		},
	)

	confAction := micro.Action(func(c *cli.Context) error {
		logger.Info("load conf from ", configPath)
		if err := conf.Load(file.NewSource(file.WithPath(configPath))); err != nil {
			return err
		}
		return nil
	})

	// etcd 服务注册
	endpoints := conf.Get("endpoints", "ca").String("127.0.0.1:2379")
	etcdRegistry := etcd.NewRegistry(registry.Addrs(strings.Split(endpoints, ",")...))

	ca := conf.Get("etcd", "ca").String("")
	cert := conf.Get("cert", "ca").String("")
	key := conf.Get("key", "ca").String("")
	if ca != "" && cert != "" {
		tlsConfig, err := ssl.ClientTslConfVerity(ca, cert, key, "")
		if err != nil {
			logger.Fatal(err)
		}
		etcdRegistry.Init(registry.TLSConfig(tlsConfig))
	}

	etcdR := micro.Registry(etcdRegistry)

	srv := micro.NewService(micro.Server(mhttp.NewServer()), confFlags)

	opts := []micro.Option{
		micro.Name(service),
		micro.Version(version),
		confAction,
		etcdR,
	}
	srv.Init(opts...)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// 注册模板和静态资源
	router.SetHTMLTemplate(web.WebTemplate())
	router.StaticFS("/web_console/web/static", http.FS(web.WebStatic()))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", conf.Get("redis", "host").String("127.0.0.1"), conf.Get("redis", "port").Int(6379)),
		Password: "",
		DB:       conf.Get("redis", "db").Int(0),
	})

	handlerOpts := &route.Options{
		Client: srv.Client(), Config: conf, Router: router,
		RedisClient: redisClient,
	}

	if err := handler.Register(handlerOpts); err != nil {
		logger.Fatal(err)
	}
	if err := micro.RegisterHandler(srv.Server(), router); err != nil {
		logger.Fatal(err)
	}
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	// Register handler
	// pb.RegisterBcsWebconsoleHandler(srv.Server(), new(handler.BcsWebconsole))

	// // Run service
	// if err := srv.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
