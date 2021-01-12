package main

import (
	"flag"
	"github.com/w3liu/consensus/log"
	"go.uber.org/zap"
	"pbft_consenus/config"
	"pbft_consenus/state"
)

var (
	BUILD_TIME string
	GIT_HASH   string
	GO_VERSION string
)

var c = flag.String("c", "./config/config.toml", "配置文件路径，默认为./config/config.toml")

func main() {
	log.Info("init", zap.String("build", BUILD_TIME), zap.String("git", GIT_HASH), zap.String("go", GO_VERSION))
	flag.Parse()
	cfg, err := config.Init(*c)
	if err != nil {
		panic(err)
	}
	state.NewState(cfg).Start()
}
