package main

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/workflow/hotelreservation2"
)

func main() {
	ctx := context.Background()
	
	var geoDB backend.NoSQLDatabase
	geoService, _ := hotelreservation2.NewGeoServiceImpl(ctx, geoDB)

	var lat, lon float64
	geoService.Nearby(ctx, lat, lon)
}
