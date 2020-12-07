package main

import "math"

const (
	meanRadius = 6371000.7900
	piOver180  = math.Pi / 180
)

func radians(degrees float64) float64 { return degrees * piOver180 }

// Position is a Lat/Lng pair
type Position struct {
	Latitude  float64
	Longitude float64
}

// Distance returns the distance (haversine) between two positions
func (me *Position) Distance(p *Position) float64 {
	const diameter = 2 * meanRadius
	lat1 := radians(me.Latitude)
	lat2 := radians(p.Latitude)
	latH := math.Sin((lat1 - lat2) / 2)
	latH *= latH
	lonH := math.Sin(radians(me.Longitude-p.Longitude) / 2)
	lonH *= lonH
	tmp := latH + math.Cos(lat1)*math.Cos(lat2)*lonH
	return diameter * math.Asin(math.Sqrt(tmp))
}
