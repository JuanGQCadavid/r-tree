package mathstuff

import (
	"math"

	"github.com/JuanGQCadavid/r-tree/app/core/domain"
)

// ChatGPT
// haversine calculates the great-circle distance between two points
// on the Earth's surface given their latitude and longitude.
// The result is returned in kilometers.
func Haversine(point1, point2 *domain.LatLon) float64 {
	const earthRadius = 6371.0 // Earth radius in kilometers

	// Convert degrees to radians
	toRadians := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	lat1 := toRadians(point1.Lat)
	lon1 := toRadians(point1.Lon)
	lat2 := toRadians(point2.Lat)
	lon2 := toRadians(point2.Lon)

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	return earthRadius * c
}

func CalculateArea(limitA, limitB *domain.LatLon) float64 {
	return math.Abs(limitA.Lat-limitB.Lat) * math.Abs(limitA.Lon-limitB.Lon)
}

// New cords calculte the coordinates that will encompas the new point.
// Returns the A,B new limits, with the delta of the new area added and the total area
func NewCoords(limitA, limitB, newPoint *domain.LatLon) (*domain.LatLon, *domain.LatLon, float64, float64) {

	minLat := math.Min(math.Min(limitA.Lat, limitB.Lat), newPoint.Lat)
	minLon := math.Min(math.Min(limitA.Lon, limitB.Lon), newPoint.Lon)

	maxLat := math.Max(math.Max(limitA.Lat, limitB.Lat), newPoint.Lat)
	maxLon := math.Max(math.Max(limitA.Lon, limitB.Lon), newPoint.Lon)

	newLimitA := &domain.LatLon{
		Lat: minLat,
		Lon: minLon,
	}

	newLimitB := &domain.LatLon{
		Lat: maxLat,
		Lon: maxLon,
	}

	newArea := CalculateArea(newLimitA, newLimitB)
	oldArea := CalculateArea(limitA, limitB)

	return newLimitA, newLimitB, newArea - oldArea, newArea
}

// New cords calculte the coordinates that will encompas the new point.
// Returns the A,B new limits, with the delta of the new area added and the total area
func NewCoordsV2[T any](limitA, limitB, newPoint *domain.Location[T]) (*domain.LatLon, *domain.LatLon, float64, float64) {
	lmtA, lmtB := MinSquare[T](limitA, limitB)
	newLimitA, newLimitB := MinSquare[T](&domain.Location[T]{
		LimitA: lmtA,
		LimitB: lmtB,
	}, newPoint)

	newArea := CalculateArea(newLimitA, newLimitB)
	oldArea := CalculateAreaV2[T](limitA, limitB)

	return newLimitA, newLimitB, newArea - oldArea, newArea
}

func MinSquare[T any](limitA, limitB *domain.Location[T]) (*domain.LatLon, *domain.LatLon) {
	minALat := math.Min(limitA.LimitA.Lat, limitA.LimitB.Lat)
	minBLat := math.Min(limitB.LimitA.Lat, limitB.LimitB.Lat)
	minABLat := math.Min(minALat, minBLat)

	minALon := math.Min(limitA.LimitA.Lon, limitA.LimitB.Lon)
	minBLon := math.Min(limitB.LimitA.Lon, limitB.LimitB.Lon)
	minABLon := math.Min(minALon, minBLon)

	maxALon := math.Max(limitA.LimitA.Lon, limitA.LimitB.Lon)
	maxBLon := math.Max(limitB.LimitA.Lon, limitB.LimitB.Lon)
	maxABLon := math.Max(maxALon, maxBLon)

	maxALat := math.Max(limitA.LimitA.Lat, limitA.LimitB.Lat)
	maxBLat := math.Max(limitB.LimitA.Lat, limitB.LimitB.Lat)
	maxABLat := math.Max(maxALat, maxBLat)

	return &domain.LatLon{
			Lat: minABLat,
			Lon: minABLon,
		}, &domain.LatLon{
			Lat: maxABLon,
			Lon: maxABLat,
		}
}

func CalculateAreaV2[T any](limitA, limitB *domain.Location[T]) float64 {
	lmtA, lmtB := MinSquare[T](limitA, limitB)
	return math.Abs(lmtA.Lat-lmtB.Lat) * math.Abs(lmtA.Lon-lmtB.Lon)
}

func IsPointCovered[T any](loc *domain.Location[T], point *domain.LatLon) bool {
	if loc.LimitA.Lat <= point.Lat && point.Lat <= loc.LimitB.Lat {
		if loc.LimitA.Lon <= point.Lon && point.Lon <= loc.LimitB.Lon {
			return true
		}
	}

	return false
}

func GetMRB[T any](node []*domain.Location[T]) (*domain.LatLon, *domain.LatLon) {
	if len(node) == 0 {
		return nil, nil
	}

	var (
		l_1_coords = node[0]
	)

	for i := 1; i < len(node); i++ {
		a, b := MinSquare[T](node[i], l_1_coords)
		l_1_coords = &domain.Location[T]{
			LimitA: a,
			LimitB: b,
		}
	}

	return l_1_coords.LimitA, l_1_coords.LimitB
}
