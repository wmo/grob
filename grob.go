package main

import (
	"fmt"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	debug := false // switch to true for a bit of extra information

	// grok the command line args
	filename, dist, lat, lon, pattern := grok_args(os.Args)

	if debug {
		fmt.Printf("# d,lat,lon,'pattern': %.3f, %f,%f,'%s'\n", dist, lat, lon, pattern)
	}

	// start reading the file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				handleNode(*v, dist, lat, lon, pattern)
				//handleNode(*v, lat, lon, dist)
				nc++
			case *osmpbf.Way:
				// Process Way v.
				wc++
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}
	if debug {
		fmt.Printf("# Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)
	}
}

func grok_args(args []string) (filename string, dist float64, lat float64, lon float64, pattern string) {

	var err error

	// default values
	filename = ""
	lat = 0.0
	lon = 0.0
	dist = -1.0 // negative distance means: no within-distance checking
	pattern = ""

	// show help if not enough args
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, `Usage: 

%s osm-file [-d max-dist lat lon] [pattern]

    eg. %s england-latest.osm mexican
    eg. %s central-america-latest.osm -d 10 12.1166 -68.9333 > willemstad10k.csv

The unit for maximum distance is km. 
`, os.Args[0], os.Args[0], os.Args[0])
		os.Exit(1)
	}

	// first arg is filename
	filename = os.Args[1]
	if !fileExists(filename) {
		fmt.Fprintf(os.Stderr, "File does not exist: %s\n", os.Args[1])
		os.Exit(1)
	}

	// distance,lat,lon args
	pattern_idx := 2
	if args[2] == "-d" {
		dist, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Illegal value for max-distance: %s\n", os.Args[3])
			os.Exit(1)
		}

		lat, err = strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Illegal value for latitude: %s\n", os.Args[4])
			os.Exit(1)
		}

		lon, err = strconv.ParseFloat(args[5], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Illegal value for longitude: %s\n", os.Args[5])
			os.Exit(1)
		}
		pattern_idx = 6
	}

	// pattern arg
	if len(args) >= (pattern_idx + 1) {
		pattern = strings.ToLower(args[pattern_idx]) // SEARCH IN LOWER CASE!
	}
	return
}

/* approximately calculate the distance between 2 points
   from: http://www.movable-type.co.uk/scripts/latlong.html
   note: φ=lat λ=lon  in RADIANS!
   var x = (λ2-λ1) * Math.cos((φ1+φ2)/2);
   var y = (φ2-φ1);
   var d = Math.sqrt(x*x + y*y) * R;
*/
func rough_distance(lat1, lon1, lat2, lon2 float64) float64 {

	// convert to radians
	lat1 = lat1 * math.Pi / 180.0
	lon1 = lon1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0
	lon2 = lon2 * math.Pi / 180.0

	r := 6371.0 // km
	x := (lon2 - lon1) * math.Cos((lat1+lat2)/2)
	y := (lat2 - lat1)
	d := math.Sqrt(x*x+y*y) * r
	return d
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func handleNode(nd osmpbf.Node, dist float64, lat float64, lon float64, contains_pattern string) {
	estim_distance := 0.0
	bingo_distance := true
	if dist >= 0.0 { // do we need to examine the distance?
		estim_distance = rough_distance(lat, lon, nd.Lat, nd.Lon)
		bingo_distance = (estim_distance < dist)
	}

	bingo_pattern := true
	if len(contains_pattern) > 0 { // do we need to examine the pattern?
		bingo_pattern = false
		for k, v := range nd.Tags {
			bingo_pattern = strings.Contains(strings.ToLower(k), contains_pattern) ||
				strings.Contains(strings.ToLower(v), contains_pattern)
			if bingo_pattern {
				break // out of the loop
			}
		}
	}

	if bingo_distance && bingo_pattern {
		// turn the Tags map into a k:v string
		tgs := ""
		for k, v := range nd.Tags {
			tgs = tgs + " " + k + ":" + v
		}

		if dist >= 0.0 {
			fmt.Printf("%f, %f, %s #,%.2f\n", nd.Lat, nd.Lon, tgs, estim_distance)
		} else {
			fmt.Printf("%f, %f, %s\n", nd.Lat, nd.Lon, tgs)
		}
	}
}
