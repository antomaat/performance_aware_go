package main

import (
	"fmt"
	"performance_aware/pkg/harvesine"
	"performance_aware/pkg/json"
)

func main() {
    fmt.Printf("hello world\n")

    r := harvesine.ReferenceHarvesine(5, 2, 4, 8, 6372.8)
    fmt.Printf("harvesine: %f\n", r)

    json.CreateHarversineJson()
}

