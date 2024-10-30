package main

import (
	"fmt"
	"golang/functions"
	"golang/services"
	"golang/structs"
	"time"
)

func main() {
	start := time.Now()
	game := functions.FetchFile[structs.GameProps]("test")

	/*res := */
	services.Calculate(&game, "4645")
	// functions.ToStringPretty(res)

	elapsed := time.Since(start)
	fmt.Printf("It took %s", elapsed)
}
