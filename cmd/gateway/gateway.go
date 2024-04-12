package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mzhn-sochi/gateway/internal/app"
)

func main() {

	a := app.InitApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		if err := a.Run(); err != nil {
			log.Printf("app crashed: %s\n", err.Error())
		}
	}()

	sig := <-sigChan

	log.Printf("caught sig: %+v. graceful shutdown", sig)
	a.Shutdown()
}
