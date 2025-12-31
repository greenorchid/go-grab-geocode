package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// test JSON parsing
func TestLocationJSONParsing(t *testing.T) {
	data := `[
		{"display_name":"Enfield, IE","lat":"53.4161821","lon":"-6.8341687"}
	]`

	var results []Location
	err := json.Unmarshal([]byte(data), &results)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	loc := results[0]
	if loc.DisplayName != "Enfield, IE" {
		t.Errorf("unexpected display name: %s", loc.DisplayName)
	}
	if loc.Lat != "53.4161821" || loc.Lon != "-6.8341687" {
		t.Errorf("unexpected coordinates: %s, %s", loc.Lat, loc.Lon)
	}
}

// helper function to build the Nominatim URL
func buildNominatimURL(query string, limit int, country string) string {
	endpoint := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", fmt.Sprint(limit))
	if country != "" {
		params.Set("countrycodes", country)
	}
	return endpoint + "?" + params.Encode()
}

// test URL building
func TestBuildNominatimURL(t *testing.T) {
	urlStr := buildNominatimURL("Enfield", 5, "ie")

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	q := u.Query()
	if q.Get("q") != "Enfield" {
		t.Errorf("expected q=Enfield, got %s", q.Get("q"))
	}
	if q.Get("format") != "json" {
		t.Errorf("expected format=json, got %s", q.Get("format"))
	}
	if q.Get("limit") != "5" {
		t.Errorf("expected limit=5, got %s", q.Get("limit"))
	}
	if q.Get("countrycodes") != "ie" {
		t.Errorf("expected countrycodes=ie, got %s", q.Get("countrycodes"))
	}
}

// test concise output logic
func TestConciseOutput(t *testing.T) {
	results := []Location{
		{DisplayName: "Enfield, IE", Lat: "53.4161821", Lon: "-6.8341687"},
		{DisplayName: "Enfield, GB", Lat: "51.6520851", Lon: "-0.0810175"},
	}

	// simulate concise flag
	concise := true
	var output string
	if concise {
		output = fmt.Sprintf("%s %s", results[0].Lat, results[0].Lon)
	}

	expected := "53.4161821 -6.8341687"
	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}
}

// test HTTP request/response using a mock server
func TestHTTPCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "q=CEnfield") {
			t.Errorf("expected query to contain q=CEnfield")
		}
		w.Write([]byte(`[{"display_name":"Enfield, IE","lat":"53.4161821","lon":"-6.8341687"}]`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "?q=CEnfield&format=json&limit=1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var results []Location
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}

	if results[0].DisplayName != "Enfield, IE" {
		t.Errorf("unexpected display name: %s", results[0].DisplayName)
	}
	if results[0].Lat != "53.4161821" || results[0].Lon != "-6.8341687" {
		t.Errorf("unexpected coordinates: %s, %s", results[0].Lat, results[0].Lon)
	}
}
