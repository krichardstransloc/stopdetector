package main

import (
	"errors"
)

// Queue is an interface for FIFO structures
type Queue interface {
	Enqueue(elem interface{})
	Dequeue() (interface{}, error)
	Front() (interface{}, error)
}

// UpdateQueue implements Queue
type UpdateQueue struct {
	Items    []*PositionUpdate
	Capacity int
}

// Enqueue adds an element to the front of the queue
func (me *UpdateQueue) Enqueue(update *PositionUpdate) {
	me.Items = append([]*PositionUpdate{update}, me.Items...)
	if len(me.Items) > me.Capacity {
		me.Items = me.Items[:me.Capacity]
	}
}

// Dequeue removes the last element from the queue
func (me *UpdateQueue) Dequeue() (*PositionUpdate, error) {
	nItems := len(me.Items)
	if nItems == 0 {
		return nil, errors.New("nothing to dequeue")
	} else if nItems == 1 {
		elem := me.Items[0]
		me.Items = make([]*PositionUpdate, 0)
		return elem, nil
	}

	elem, queue := me.Items[nItems-1], me.Items[:nItems-1]
	me.Items = queue
	return elem, nil
}

// Front returns the head of the queue
func (me *UpdateQueue) Front() (*PositionUpdate, error) {
	if len(me.Items) == 0 {
		return nil, errors.New("queue is empty")
	}

	return me.Items[0], nil
}

// AverageSpeed returns the weighted average of vehicle speeds, weighing recent updates more heavily
func (me *UpdateQueue) AverageSpeed() (float64, error) {
	nUpdates := len(me.Items)
	if nUpdates < 2 {
		return 0, errors.New("too few stops")
	}

	var numerator float64
	var denominator float64

	lowIdx := 0
	highIdx := 1
	weight := float64(nUpdates)
	for {
		if highIdx >= nUpdates {
			break
		}
		lowStop := me.Items[lowIdx]
		highStop := me.Items[highIdx]
		deltaDistance := lowStop.Position.Distance(highStop.Position)
		deltaTime := float64(lowStop.Timestamp - highStop.Timestamp)

		if deltaTime > 0 {
			numerator += deltaDistance / deltaTime * weight
		} else {
			return 0, errors.New("invalid time delta")
		}

		denominator += weight

		lowIdx++
		highIdx++
		weight--
	}

	if denominator == 0 {
		return 0, errors.New("divide by zero")
	}

	return numerator / denominator, nil
}

// NewUpdateQueue returns a pointer to an UpdateQueue
func NewUpdateQueue(capacity int) *UpdateQueue {
	items := make([]*PositionUpdate, 0)
	return &UpdateQueue{
		Items:    items,
		Capacity: capacity,
	}
}
