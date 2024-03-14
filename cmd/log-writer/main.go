package main

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/web"
)

func main() {
	router := web.SetupRouter()
	err := router.Run(":1233")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
