package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
	"github.com/JuanGQCadavid/r-tree/app/core/rtree"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
	tree   = rtree.NewRTree[string](3)
)

type HttpResponse struct {
	Response []*domain.Location[string] `json:"results,omitempty"`
}

const (
	fileName = "tartu_stops.csv"
)

func init() {
	populateRTree(tree, fileName)

	router.GET("/", GetPoints)          // OK
	router.GET("/fake/", GetPointsFake) // OK
	router.Run(":8001")
}

func populateRTree(tree *rtree.RTree[string], fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatalln("Error reading records ", err.Error())
	}

	records = records[1:] // Deleting header

	for _, eachrecord := range records {
		if len(eachrecord) < 4 {
			continue
		}
		var (
			id  = eachrecord[0]
			lat = eachrecord[2]
			lon = eachrecord[3]
		)
		fmt.Println(id, lat, lon)

		latF, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			log.Println("There is an error while parsing the Lat -> ", lat, " err: ", err.Error())
		}

		lngF, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			log.Println("There is an error while parsing the Lat -> ", lat, " err: ", err.Error())
		}

		tree.InsertLocation(
			&domain.LatLon{
				Lat: latF,
				Lon: lngF,
			},
			id,
			tree.Root,
		)
	}

}

func GetPointsFake(context *gin.Context) {
	results := []*domain.Location[string]{
		&domain.Location[string]{
			Value: "122420",
			LimitA: &domain.LatLon{
				Lat: 58.35881504,
				Lon: 26.67327622,
			},
			LimitB: &domain.LatLon{
				Lat: 58.35881504,
				Lon: 26.67327622,
			},
		},
		&domain.Location[string]{
			Value: "161498",
			LimitA: &domain.LatLon{
				Lat: 58.34987601,
				Lon: 26.76803853,
			},
			LimitB: &domain.LatLon{
				Lat: 58.34987601,
				Lon: 26.76803853,
			},
		},
		&domain.Location[string]{
			Value: "125530",
			LimitA: &domain.LatLon{
				Lat: 58.35889909,
				Lon: 26.67292438,
			},
			LimitB: &domain.LatLon{
				Lat: 58.35889909,
				Lon: 26.67292438,
			},
		},
	}

	context.JSON(http.StatusOK, &HttpResponse{
		Response: results,
	})

}

func GetPoints(context *gin.Context) {
	var (
		lat = context.DefaultQuery("lat", "0")
		lng = context.DefaultQuery("lng", "0")
	)

	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Println("There is an error while parsing the Lat -> ", lat, " err: ", err.Error())
	}

	lngF, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		log.Println("There is an error while parsing the Lat -> ", lat, " err: ", err.Error())
	}

	results := tree.Search(&domain.LatLon{
		Lat: latF,
		Lon: lngF,
	}, tree.Root)

	// tree.Search(&domain.LatLon{
	// 	Lat: latF,
	// 	Lon: lngF,
	// }, tree.Root)

	context.JSON(http.StatusOK, &HttpResponse{
		Response: results,
	})

}

func main() {

}
