package main

import (
	"github.com/jgero/schlingel/api"
)

func main() {
	router := api.BuildRouter()
	router.Run(":8080")
}
