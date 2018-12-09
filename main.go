package main

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	fdao "bonex-middleware/dao/firebase"
	mdao "bonex-middleware/dao/mysql"
	rdao "bonex-middleware/dao/redis"
	"bonex-middleware/log"
	"bonex-middleware/services/api"
	"bonex-middleware/services/faucet"
	"flag"
	"os"
)

func main() {
	configFile := flag.String("conf", "./config.json", "Pathname to config file")
	//var DisableFaucet bool
	//flag.BoolVar(&DisableFaucet, "disable-faucet", false, "Publish withdraws")

	flag.Parse()

	cfg, err := config.NewFromFile(configFile)
	if err != nil {
		log.Errorf("Cannot decode config: %s", err.Error())
		return
	}

	// Setup log
	log.Init(cfg.LogLevel)

	// Init dao
	dbDao, err := mdao.NewMysql(cfg, log.GetInstance())
	if err != nil {
		log.Fatalf("Cannot init DB DAO: %s", err.Error())
	}

	redisDao, err := rdao.NewRedis(cfg)
	if err != nil {
		log.Fatalf("Cannot init Redis DAO: %s", err.Error())
	}

	fbDao, err := fdao.NewFirebase(cfg)
	if err != nil {
		log.Fatalf("Cannot init Firebase DAO: %s", err.Error())
	}

	d := dao.New(redisDao, dbDao, fbDao)

	s := faucet.NewStellar(cfg)
	faucetModule := faucet.New(d, cfg, s)
	err = faucetModule.PromtKey()
	if err != nil {
		log.Fatalf("Cannot promt necessary keys to run faucet: %s", err.Error())
	}

	apiModule := api.New(d, cfg, faucetModule)

	runModules(apiModule, faucetModule)

	log.Infof("Exiting")
	os.Exit(0)
}
