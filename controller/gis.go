package controller

import (
	"log"
	"net/http"
	"strconv"
    "github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"go.mongodb.org/mongo-driver/bson"
    "github.com/gocroot/config"
)


func NearbyRoadHandler(respw http.ResponseWriter, req *http.Request) {
    latitude, errLat := strconv.ParseFloat(req.URL.Query().Get("latitude"), 64)
    longitude, errLng := strconv.ParseFloat(req.URL.Query().Get("longitude"), 64)

    if errLat != nil || errLng != nil || latitude == 0 || longitude == 0 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
            "error": "Invalid or missing latitude and longitude",
        })
        return
    }

    filter := bson.M{
        "geometry": bson.M{
            "$near": bson.M{
                "$geometry": bson.M{
                    "type":        "Point",
                    "coordinates": []float64{longitude, latitude},
                },
                "$maxDistance": 1000, // Jarak maksimum dalam meter
            },
        },
    }

    // Menambahkan parameter limit (misalnya, 10)
    limit := 10

    roads, err := atdb.GetOneManyDocs[bson.M](config.Mongoconn, "jalan", filter, limit)
    if err != nil {
        log.Println("Error fetching roads from MongoDB:", err)
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{
            "error": "Failed to fetch roads from MongoDB",
        })
        return
    }

    if len(roads) == 0 {
        at.WriteJSON(respw, http.StatusNotFound, map[string]string{
            "error": "No roads found near the specified location",
        })
        return
    }

    at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
        "status": "Success",
        "roads":  roads,
    })
}
