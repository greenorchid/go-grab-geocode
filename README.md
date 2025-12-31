# go-grab-geocode

A small Go command-line tool that looks up latitude and longitude for a given place name.

It uses OpenStreetMap Nominatim and returns multiple matches for ambiguous locations.

## Usage

```bash
./go-grab-geocode "Enfield"
./go-grab-geocode --country=ie "Enfield" 
./go-grab-geocode --limit=1 "Enfield Ireland"
./go-grab-geocode --concise "Enfield Ireland"