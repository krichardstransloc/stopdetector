package main

import (
	"fmt"
	"sort"
	"time"
)

// Stop is the location of a stop on a pattern
type Stop struct {
	StopID    int
	PatternID int
	Position  *Position
	Sequence  int
}

// StopSlice is a slice of stops for a pattern
type StopSlice []*Stop

func (me StopSlice) Len() int {
	return len(me)
}

func (me StopSlice) Swap(i, j int) {
	me[i], me[j] = me[j], me[i]
}

func (me StopSlice) Less(i, j int) bool {
	return me[i].Sequence < me[j].Sequence
}

// StopCache maps pattern id to a slice of stops. It is refreshed by polling Operations.
type StopCache map[int]StopSlice

// PositionUpdate is a recent position update
type PositionUpdate struct {
	PatternID int
	Position  *Position
	Timestamp int
}

// VehicleStop is an actual stop
type VehicleStop struct {
	Stop       *Stop
	Timestamp  int
	IsEstimate bool
}

// Vehicle is a vehicle...
type Vehicle struct {
	VehicleID int
	PatternID int
	LastStop  *VehicleStop
	Updates   *UpdateQueue
}

// VehicleCache maps vehicle ids to vehciles. It is hydrated by Operations and enriched by rmq.
type VehicleCache map[int]*Vehicle

// AddUpdate records a new position update
func (me *Vehicle) AddUpdate(p *PositionUpdate) {
	if me.PatternID != p.PatternID {
		me.ChangePattern(p.PatternID)
	}
	me.Updates.Enqueue(p)
}

// ChangePattern resets data stores that are tied to a specific pattern
func (me *Vehicle) ChangePattern(patternID int) {
	me.PatternID = patternID
	me.Updates.Items = make([]*PositionUpdate, 0)
	me.SetLastStop(&VehicleStop{Stop: &Stop{Sequence: 0}}) // this should really be the first point in the pattern.
}

// SetLastStop records a new position update
func (me *Vehicle) SetLastStop(s *VehicleStop) {
	me.LastStop = s
}

// LastPositionUpdate returns the most recent position received
func (me *Vehicle) LastPositionUpdate() (*PositionUpdate, error) {
	return me.Updates.Front()
}

// Process loops over all vehicles and their stops, and interpolates current position based on
// average m/s of recent updates. This position is then used to detect passed stops
func Process(vehicleCache VehicleCache, stopCache StopCache) {
	for _, vehicle := range vehicleCache {
		nowSeconds := int(time.Now().Unix())
		vStops, _ := stopCache[vehicle.PatternID]

		// make sure stops are sorted sequentially
		sort.Sort(vStops)

		for _, s := range vStops {
			// grab the most recent position update
			lastPosition, err := vehicle.LastPositionUpdate()
			if err != nil {
				fmt.Println(err)
				break
			}

			// compute average speed over last updates
			avgSpeed, err := vehicle.Updates.AverageSpeed()
			if err != nil {
				fmt.Println(err)
				break
			}

			// only consider stops beyond our most recent
			if s.Sequence > vehicle.LastStop.Stop.Sequence {
				distanceToStop := lastPosition.Position.Distance(s.Position)
				// we're within the stop radius, so set this stop as our current.
				if distanceToStop < 50 {
					fmt.Println("at stop")
					vehicle.SetLastStop(&VehicleStop{
						Stop:       s,
						Timestamp:  nowSeconds,
						IsEstimate: false,
					})

					// if we're at the current stop, we've already iterated over previous stops and
					// know we aren't also at the next stop. stop iteration.
					break
				}

				// we're not within the stop radius. we need to see if we've missed a stop due to
				// sparse position updates
				timeToStop := int(distanceToStop * avgSpeed)

				// if the time at the last stop plus travel time at our average speed is less than the
				// current time, assume we've missed a stop. if last stop time plus travel time is
				// greater than now, we're still approaching the stop
				if lastPosition.Timestamp+timeToStop < nowSeconds {
					fmt.Println("passed a stop")
					vehicle.SetLastStop(&VehicleStop{
						Stop:       s,
						Timestamp:  lastPosition.Timestamp + timeToStop,
						IsEstimate: true,
					})

					// just because we passed over one stop does not mean we haven't passed another
					// as well. we need to keep iterating through stops until we hit a future stop
					continue
				} else {
					fmt.Println("approaching the next stop")

					// if we're approaching the next stop, we're also approaching all
					// subsequent stops. stop iteration
					break
				}
			}
		}
	}
}
