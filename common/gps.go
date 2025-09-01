package common

import (
    "math/rand"
)

type Location struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

// RandomLocation returns a random GPS location in Mumbai
func RandomLocation() Location {
    // Mumbai approx bounds
    minLat, maxLat := 18.89, 19.30
    minLng, maxLng := 72.77, 72.99
    lat := minLat + rand.Float64()*(maxLat-minLat)
    lng := minLng + rand.Float64()*(maxLng-minLng)
    return Location{Lat: lat, Lng: lng}
}
