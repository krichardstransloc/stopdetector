package main

import (
	"time"
)

func main() {
	beyu := &Position{
		Latitude:  35.996656,
		Longitude: -78.903862,
	}

	toro := &Position{
		Latitude:  35.997048,
		Longitude: -78.903628,
	}

	parlour := &Position{
		Latitude:  35.996644,
		Longitude: -78.902220,
	}

	pompieri := &Position{
		Latitude:  35.995972,
		Longitude: -78.9003702,
	}

	// poll ops to update stop cache. for now, use fake data
	s1 := &Stop{
		StopID:    1,
		PatternID: 2,
		Position:  beyu,
		Sequence:  1,
	}
	s2 := &Stop{
		StopID:    2,
		PatternID: 2,
		Position:  toro,
		Sequence:  2,
	}
	s3 := &Stop{
		StopID:    3,
		PatternID: 2,
		Position:  parlour,
		Sequence:  3,
	}
	stopCache := make(StopCache)
	stopCache[2] = StopSlice{s1, s2, s3}

	nowSeconds := int(time.Now().Unix())

	vehicleStop := &VehicleStop{
		Stop:      s2,
		Timestamp: nowSeconds - 200,
	}

	// pull from rmq to update vehicle. for now, use fake data
	vehicle := &Vehicle{
		VehicleID: 101,
		PatternID: 2,
		LastStop:  vehicleStop,
		Updates:   NewUpdateQueue(10),
	}

	vehicleCache := make(VehicleCache)
	vehicleCache[vehicle.VehicleID] = vehicle

	p1 := &PositionUpdate{
		PatternID: 2,
		Position:  beyu,
		Timestamp: nowSeconds - 100,
	}
	p2 := &PositionUpdate{
		PatternID: 2,
		Position:  toro,
		Timestamp: nowSeconds,
	}
	p3 := &PositionUpdate{
		PatternID: 2,
		Position:  pompieri,
		Timestamp: nowSeconds + 20,
	}

	vehicle.AddUpdate(p1)
	vehicle.AddUpdate(p2)
	vehicle.AddUpdate(p3)

	Process(vehicleCache, stopCache)
}
