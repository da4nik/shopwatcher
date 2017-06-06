package core

import (
	"context"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/shopwatcher/db"
	"github.com/da4nik/shopwatcher/parsers"
	"github.com/da4nik/shopwatcher/types"
)

// StartScheduler starts scheduler
func StartScheduler(ctx context.Context) {
	go scheduler(ctx)
}

func scheduler(ctx context.Context) {
	log := logger()
	log.Infoln("Scheduler started.")

	var wg sync.WaitGroup
	working := true
	for working {
		// Parsing all products concurrently
		products, _ := db.AllProducts()
		log.Infof("Processing %d product(s)\n", len(products))
		for _, product := range products {
			wg.Add(1)
			go parseProduct(product, &wg)
		}
		wg.Wait()

		log.Infoln("Processing done, sleeping for 1 minute.")
		select {
		case <-ctx.Done():
			working = false
		case <-time.After(time.Minute * 1):
		}
	}
	log.Infoln("Scheduler finished.")
}

// parseProduct parse new version of product
func parseProduct(product types.Product, wg *sync.WaitGroup) {
	defer wg.Done()
	log := logger().WithField("function", "parseProduct")

	log.Debugf("Parsing product: %s (%s)", product.Name, product.URL)

	newProduct, err := parsers.Parse(product.URL)
	if err != nil {
		return
	}

	if product.Equal(newProduct) {
		return
	}

	db.SaveProduct(newProduct)
	// TODO: Emit product change event
}

func logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"module": "core",
	})
}
