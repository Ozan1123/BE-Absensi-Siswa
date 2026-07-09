package utils

import (
	"math"
)

// SCHOOL_LATITUDE and SCHOOL_LONGITUDE are the coordinates of the school
const (
	SCHOOL_LATITUDE       = -6.467184782305172
	SCHOOL_LONGITUDE      = 106.8646612685341
	ALLOWED_RADIUS_METERS = 100.0
)

// GetDistanceMeters calculates the distance in meters between two GPS coordinates using Haversine formula
func GetDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return R * c
}

// IsInsideSchool checks if the student coordinates are within the allowed radius
func IsInsideSchool(lat, lon float64) bool {
	distance := GetDistanceMeters(lat, lon, SCHOOL_LATITUDE, SCHOOL_LONGITUDE)
	return distance <= ALLOWED_RADIUS_METERS
}
