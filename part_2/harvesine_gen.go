package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"performance_aware/pkg/harvesine"
	"strconv"
)

func main() {
    var discard bool
    flag.BoolVar(&discard, "discard", false, "discard generated output")
    flag.Parse()
    if len(flag.Args()) == 2 {
	seed, err := strconv.ParseInt(flag.Args()[0], 10, 64)
	count, err := strconv.ParseInt(flag.Args()[1], 10, 64)
	if err != nil {
	    fmt.Println("Usage: [random seed] [number of coordinate pairs to generate]")
	}
	pairs, distances := CreateHarvesine(seed, count)

	if !discard {
	    f, err := os.Create(path.Join("./", "output.json"))
	    if err != nil {
		log.Fatal(err)
	    }
	    defer f.Close()

	    writer := bufio.NewWriter(f)
	    defer writer.Flush()
	    enc := json.NewEncoder(writer)
	    enc.SetIndent("", " ")
	    if err := enc.Encode(Output{pairs}); err != nil {
		log.Fatal(err)
	    }
	}

	avg := Average(distances)
	fmt.Printf("average: %f\n", avg)
    } else {
	fmt.Println("Usage: [random seed] [number of coordinate pairs to generate]")
    }

}

func CreateHarvesine(seed int64, nrOfPairs int64) ([]Pair, []float64){
    rand := rand.New(rand.NewSource(seed))
    pairs := make([]Pair, nrOfPairs)
    clusters := rand.Intn((100 + int(seed)) % 1024)
    var distances []float64
    pos := 0
    for i := 0; i < clusters; i++ {
	xMin := rand.Float64()*360 - 180
	xMax := rand.Float64()*360 - 180
	if xMin > xMax {
	    xMax, xMin = xMin, xMax
	}
	yMin := rand.Float64()*180 - 90 
	yMax := rand.Float64()*180 - 90 
	if yMin > yMax {
	    yMax, yMin = yMin, yMax
	}
	items := (len(pairs) - pos) / (clusters - i)
	for j := 0; j < items; j++ {
	    posI := pos + j	
	    pairs[posI].X0 = rand.Float64()*xMax - xMin
	    pairs[posI].Y0 = rand.Float64()*yMax - yMin
	    pairs[posI].X1 = rand.Float64()*xMax - xMin
	    pairs[posI].Y1 = rand.Float64()*yMax - yMin
	    item := pairs[posI]
	    harvesine := harvesine.ReferenceHarvesine(item.X0, item.Y0, item.X1, item.Y1, 6372.8) 
	    distances = append(distances, harvesine)
	}
	pos += items
    }
    return pairs, distances
}

func Average(xx []float64) float64 {
	var avg float64
	for _, x := range xx {
		avg += x
	}
	return avg / float64(len(xx))
}

type Pair struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

type Output struct {
	Pairs []Pair `json:"pairs"`
}
