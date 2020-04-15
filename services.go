package main

import (
	"context"
	"fmt"
	"github.com/gocarina/gocsv"
	"googlemaps.github.io/maps"
	"log"
	"os"
)

var GoogleMapsApiKey string = os.Getenv("GOOGLE_MAPS_API_KEY")

func ReadCepFromCsv(filename string) ([]*Address, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []*Address{}, err
	}
	defer f.Close()
	addresses := []*Address{}
	if err := gocsv.UnmarshalFile(f, &addresses); err != nil {
		return []*Address{}, err
	}
	return addresses, nil
}

func ConvertCsvWithCepToLatitudeLongitude(addresses []*Address) ([]*Address, error) {
	c, err := maps.NewClient(maps.WithAPIKey(GoogleMapsApiKey))
	if err != nil {
		log.Printf("fatal error: %s", err)
		return nil, err
	}
	for _, address := range addresses {
		geocodeRequest := &maps.GeocodingRequest{
			Address: address.CEP,
			Region: "BR",
		}
		geocodeResult, err := c.Geocode(context.Background(), geocodeRequest)
		if err != nil {
			log.Printf("fatal error: %s", err)
			return nil, err
		}
		location := geocodeResult[0].Geometry.Location
		latitude := location.Lat
		longitude := location.Lng
		address.Latitude = fmt.Sprintf("%f", latitude)
		address.Longitude = fmt.Sprintf("%f", longitude)
	}
	return addresses, nil
}

func WriteCsv(addresses []*Address) error {
	clientsFile, err := os.OpenFile("output.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	_, err = gocsv.MarshalString(&addresses)
	if err != nil {
		log.Printf("fatal error: %s", err)
		return err
	}
	defer clientsFile.Close()
	err = gocsv.MarshalFile(&addresses, clientsFile)
	if err != nil {
		log.Printf("fatal error: %s", err)
		return err
	}
	return nil
}
