package main

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	fdao "bonex-middleware/dao/firebase"
	mdao "bonex-middleware/dao/mysql"
	rdao "bonex-middleware/dao/redis"
	bdao "bonex-middleware/dao/stellar"
	"bonex-middleware/log"
	"bonex-middleware/services/api"
	"bonex-middleware/services/faucet"
	"bonex-middleware/services/messaging"
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

	bDao := bdao.NewStellar(cfg)

	d := dao.New(redisDao, dbDao, fbDao, bDao)

	faucetModule := faucet.New(d, cfg)
	err = faucetModule.PromtKey()
	if err != nil {
		log.Fatalf("Cannot promt necessary keys to run faucet: %s", err.Error())
	}

	apiModule := api.New(d, cfg, faucetModule)
	msgModule := messaging.New(d, cfg)

	runModules(apiModule, faucetModule, msgModule)

	log.Infof("Exiting")
	os.Exit(0)
}
