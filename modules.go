package main

import (
	"context"
	"bonex-middleware/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const gracefulStopTtl = 5

type (
	module interface {
		Run() error
		Title() string
		GracefulStop(context.Context) error
	}
)

func runModules(modules ...module) {
	// Run modules
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	if len(modules) > 0 {
		for _, m := range modules {
			go func(module module) {
				log.Infof("Starting module %s", module.Title())
				err := module.Run()
				if err != nil {
					log.Fatalf("Module fatal error: %s", err.Error())
				}
			}(m)
		}
	}

	// Wait for termination
	sig := <-gracefulStop
	log.Warnf("Caught sig: %+v", sig)

	stopModules(time.Second*gracefulStopTtl, modules...)
}

func stopModules(wait time.Duration, modules ...module) {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	wg := &sync.WaitGroup{}
	wg.Add(len(modules))
	for _, m := range modules {
		go func(module module) {
			stopModule(ctx, module)
			wg.Done()
		}(m)
	}

	doneChan := make(chan bool)
	go func() {
		wg.Wait()
		doneChan <- true
	}()

	select {
	case <-doneChan:
		log.Infof("all modules stopped gracefuly")
	case <-ctx.Done():
		log.Warnf("timeout exceeded: some modules might fail to stop: %s", ctx.Err().Error())
	}
}

func stopModule(ctx context.Context, module module) {
	log.Infof("Stoping module %s", module.Title())
	err := module.GracefulStop(ctx)
	if err != nil {
		log.Errorf("module %s returned an error %s", module.Title(), err.Error())
	}
}
