package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	api "github.com/wayne011872/api-toolkit"
	"github.com/wayne011872/api-toolkit/mid"
	"github.com/wayne011872/goSterna/util"
	"github.com/wayne011872/log"
	"github.com/wayne011872/microservice"
	"github.com/wayne011872/microservice/di"
	myapi "github.com/wayne011872/myTodoBackend/api"
	"github.com/wayne011872/myTodoBackend/model"
	"github.com/wayne011872/pgx/conn"
	"github.com/wayne011872/wayneLib"
	authMid "github.com/wayne011872/wayneLib/auth/mid"
)

var (
	service = flag.String("service", "cli", "service (api, cli)")
	v       = flag.Bool("v", false, "version")

	Version   = "1.0.0"
	BuildTime = "2000-01-01T00:00:00+0800"
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println("Version: " + Version)
		fmt.Println("Build Time: " + BuildTime)
		return
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	envFile := path + "/.env"
	if util.FileExists(envFile) {
		godotenv.Load(envFile)
	}
	modelCfg, err := model.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}
	ms, err := microservice.New(modelCfg, &mydi{})
	if err != nil {
		panic(err)
	}

	serv := &myservice{ms}
	switch *service {
	case "cli":
		serv.runCli()
	default:
		microservice.RunService(
			serv.runApi,
		)
	}
}

func (serv *myservice) runApi(ctx context.Context) {
	cfg, err := api.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}
	if !cfg.IsMockAuth {
		authAddress := os.Getenv("GRPC_Auth_Address")
		if authAddress == "" {
			panic("missing auth address")
		}
		cfg.SetAuth(authMid.NewGinInterAuthMid(authAddress))
	}
	cfg.SetMiddles(
		mid.NewGinMiddle(di.GinMiddleHandler(serv.GetDI())),
		serv.GetModelCfgMgr())
	cfg.AddProms(conn.PgxOpsQueued)
	cfg.SetServerErrorHandler(wayneLib.ServerErrorHandler)
	cfg.SetAPIs(
		myapi.NewTodoItemAPI(),
	)
	log, err := serv.NewLog("api")
	if err != nil {
		panic(err)
	}
	cfg.Logger = log
	err = api.AutoGinApiRun(ctx, cfg)
	if err != nil {
		cfg.Logger.Fatalf("api run fail: %v", err)
	}
}

type myservice struct {
	microservice.MicroService[*model.Config, *mydi]
}

func (serv *myservice) runCli() {
	cfg, err := serv.NewCfg("cli")
	if err != nil {
		panic(err)
	}
	defer cfg.Close()

	if err != nil {
		panic(err)
	}
}

type mydi struct {
	di.CommonServiceDI
	*conn.PgxConf   `yaml:"postgres,omitempty"`
	*log.LoggerConf `yaml:"log,omitempty"`
}

func (d *mydi) IsConfEmpty() error {
	if d.PgxConf == nil {
		return errors.New("postgres config is empty")
	}

	if d.LoggerConf == nil {
		return errors.New("logger config is empty")
	}

	return nil
}
