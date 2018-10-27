package main

import (
    "bonex-middleware/config"
    dao "bonex-middleware/dao/mysql"
    "bonex-middleware/log"
    "bonex-middleware/services/api"
    "flag"
    "os"
)

func main() {
    configFile := flag.String("conf", "./config.json", "Pathname to config file")

    flag.Parse()

    cfg, err := config.NewFromFile(configFile)
    if err != nil {
        log.Errorf("Cannot decode config: %s", err.Error())
        return
    }

    // Setup log
    log.Init(cfg.LogLevel)

    // Init dao
    d, err := dao.New(cfg, log.GetInstance())
    if err != nil {
        log.Fatalf("Cannot init DAO: %s", err.Error())
    }

    apiModule := api.New(d, cfg)

    runModules(apiModule)

    log.Infof("Exiting")
    os.Exit(0)
}
