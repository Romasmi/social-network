package main

import (
	"fmt"
	"github.com/Romasmi/social-network/internal/app"
)

func main() {
	appInstance, err := app.CreateApp()
	if err != nil {
		fmt.Printf("error while app init: %v", err)
	}
	defer appInstance.OnStop()
	appInstance.Run()

	// TODO Add gracefull shutdown
}
