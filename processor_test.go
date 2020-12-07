package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessor__Approaching_Stop(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	s1 := &Stop{
		StopID:    1,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.920733,
			Longitude: -78.902220,
		},
		Sequence: 1,
	}
	s2 := &Stop{
		StopID:    2,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.91937,
			Longitude: -78.9310093,
		},
		Sequence: 2,
	}

	stopCache := make(StopCache)
	stopCache[2] = StopSlice{s1, s2}

	vehicleStop := &VehicleStop{
		Stop:      s1,
		Timestamp: nowSeconds - 200,
	}

	vehicle := &Vehicle{
		VehicleID: 101,
		PatternID: 2,
		LastStop:  vehicleStop,
		Updates:   NewUpdateQueue(10),
	}

	vehicleCache := make(VehicleCache)
	vehicleCache[vehicle.VehicleID] = vehicle

	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.919948,
			Longitude: -78.930431,
		},
		Timestamp: nowSeconds - 20,
	})
	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.919003,
			Longitude: -78.930301,
		},
		Timestamp: nowSeconds - 10,
	})

	Process(vehicleCache, stopCache)
	// do not change stop ID
	assert.Equal(t, s1.StopID, vehicleCache[101].LastStop.Stop.StopID)
}

func TestProcessor__At_Stop(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	s1 := &Stop{
		StopID:    1,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.920733,
			Longitude: -78.902220,
		},
		Sequence: 1,
	}
	s2 := &Stop{
		StopID:    2,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.919003,
			Longitude: -78.930301,
		},
		Sequence: 2,
	}

	stopCache := make(StopCache)
	stopCache[2] = StopSlice{s1, s2}

	vehicleStop := &VehicleStop{
		Stop:      s1,
		Timestamp: nowSeconds - 200,
	}

	vehicle := &Vehicle{
		VehicleID: 101,
		PatternID: 2,
		LastStop:  vehicleStop,
		Updates:   NewUpdateQueue(10),
	}

	vehicleCache := make(VehicleCache)
	vehicleCache[vehicle.VehicleID] = vehicle

	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.919948,
			Longitude: -78.930431,
		},
		Timestamp: nowSeconds - 20,
	})
	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.9188,
			Longitude: -78.930205,
		},
		Timestamp: nowSeconds - 10,
	})

	Process(vehicleCache, stopCache)
	assert.Equal(t, s2.StopID, vehicleCache[101].LastStop.Stop.StopID)
}

func TestProcessor__Passed_Stop__Past(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	s1 := &Stop{
		StopID:    1,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.920733,
			Longitude: -78.902220,
		},
		Sequence: 1,
	}
	s2 := &Stop{
		StopID:    2,
		PatternID: 2,
		Position: &Position{
			Latitude:  35.91937,
			Longitude: -78.9310093,
		},
		Sequence: 2,
	}

	stopCache := make(StopCache)
	stopCache[2] = StopSlice{s1, s2}

	vehicleStop := &VehicleStop{
		Stop:      s1,
		Timestamp: nowSeconds - 100,
	}

	vehicle := &Vehicle{
		VehicleID: 101,
		PatternID: 2,
		LastStop:  vehicleStop,
		Updates:   NewUpdateQueue(10),
	}

	vehicleCache := make(VehicleCache)
	vehicleCache[vehicle.VehicleID] = vehicle

	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.919948,
			Longitude: -78.930431,
		},
		Timestamp: nowSeconds - 20,
	})
	vehicle.AddUpdate(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.9188,
			Longitude: -78.930205,
		},
		Timestamp: nowSeconds - 50,
	})

	Process(vehicleCache, stopCache)
	assert.Equal(t, s2.StopID, vehicleCache[101].LastStop.Stop.StopID)
}
