package main

import (
    "bonex-middleware/config"
    mdao "bonex-middleware/dao/mysql"
    rdao "bonex-middleware/dao/redis"
    "bonex-middleware/dao"
    "bonex-middleware/log"
    "bonex-middleware/services/api"
    "flag"
    "os"
    "bonex-middleware/services/faucet"
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

    d := dao.New(redisDao, dbDao)

    s := faucet.NewStellar(cfg)
    faucetModule := faucet.New(d, cfg, s)

    apiModule := api.New(d, cfg, faucetModule)



    runModules(apiModule, faucetModule)

    log.Infof("Exiting")
    os.Exit(0)
}
