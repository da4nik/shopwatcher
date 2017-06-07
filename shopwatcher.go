package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/da4nik/shopwatcher/integrations"
	_ "github.com/da4nik/shopwatcher/integrations/telegram"
	_ "github.com/joho/godotenv/autoload"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/shopwatcher/core"
	"github.com/da4nik/shopwatcher/db"
)

var (
	// Version - release version
	Version string

	// BuildTime - build time :)
	BuildTime string
)

func main() {
	if Version != "" || BuildTime != "" {
		fmt.Printf("Shopwatcher build version %s, build time %s\n\n", Version, BuildTime)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.DebugLevel)

	defer db.Close()

	integrations.Start()
	defer integrations.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	core.StartScheduler(ctx)
	defer cancel()

	// prod := types.Product{
	// 	URL:  "https://www.wildberries.ru/catalog/3199841/detail.aspx",
	// 	Name: "Some name",
	// }
	//
	// if err := db.SaveProduct(prod); err != nil {
	// 	logrus.Debugf("Err: %s\n", err.Error())
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
