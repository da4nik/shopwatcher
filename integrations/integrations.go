package integrations

import (
	"context"

	"github.com/da4nik/shopwatcher/db"
	"github.com/da4nik/shopwatcher/integrations/telegram"
	"github.com/da4nik/shopwatcher/types"
)

var (
	cancelFuncs []context.CancelFunc
	outChans    []chan types.Product
	inChan      = make(chan types.Product, 10)
	done        = make(chan bool, 1)
)

// Start starts all integrations
func Start() {
	ctx, telegramCancel := context.WithCancel(context.Background())
	telegramChan := telegram.Start(ctx, inChan)

	cancelFuncs = append(cancelFuncs, telegramCancel)
	outChans = append(outChans, telegramChan)

	go listen()
}

// Stop stops all integrations
func Stop() {
	for _, cancel := range cancelFuncs {
		cancel()
	}
	done <- true
}

// Notify notifies about product changes
func Notify(product types.Product) {
	for _, out := range outChans {
		out <- product
	}
}

func listen() {
	for {
		select {
		case product := <-inChan:
			// TODO: merge users on save
			db.SaveProduct(product)
		case <-done:
			return
		}
	}
}
