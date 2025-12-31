package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Location struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

func main() {
	limit := flag.Int("limit", 5, "number of results") // is a sensible limits
	country := flag.String("country", "", "country code (e.g. ie, gb, us)")
	concise := flag.Bool("concise", false, "print only first result coordinates")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("usage: geofinder [options] <location>")
		os.Exit(1)
	}

	query := strings.Join(flag.Args(), " ")
	// simplest api to consume, no creds ;-)
	endpoint := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", fmt.Sprint(*limit))

	if *country != "" {
		params.Set("countrycodes", *country)
	}

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}

	req.Header.Set("User-Agent", "geofinder/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http error:", err)
		return
	}
	defer resp.Body.Close()

	var results []Location
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		fmt.Println("decode error:", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("no results found")
		return
	}

	// warn, for first result only
	if *concise {
		fmt.Printf("%s %s\n", results[0].Lat, results[0].Lon)
		return
	}

	for i, loc := range results {
		fmt.Printf("%d. %s\n", i+1, loc.DisplayName)
		fmt.Printf("   Lat: %s, Lon: %s\n\n", loc.Lat, loc.Lon)
	}
}
