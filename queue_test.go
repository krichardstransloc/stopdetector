package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateQueue_Enqueue(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	// given
	queue := NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})

	// then
	assert.Equal(t, 1, len(queue.Items))

	// given
	queue = NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds,
	})

	// then
	assert.Equal(t, 2, len(queue.Items))
	assert.Equal(t, 35.996656, queue.Items[0].Position.Latitude)
	assert.Equal(t, 35.996644, queue.Items[1].Position.Latitude)

	// given
	queue = NewUpdateQueue(10)

	// when

	// then
	assert.Equal(t, 0, len(queue.Items))

	// given
	queue = NewUpdateQueue(1)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds,
	})

	// then
	assert.Equal(t, 1, len(queue.Items))
	assert.Equal(t, 35.996656, queue.Items[0].Position.Latitude)

}

func TestUpdateQueue_Front(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	// given
	queue := NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})

	// then
	front, err := queue.Front()
	assert.NoError(t, err)
	assert.Equal(t, 35.996644, front.Position.Latitude)
	assert.Equal(t, 1, len(queue.Items))

	// given
	queue = NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds,
	})

	// then
	front, err = queue.Front()
	assert.NoError(t, err)
	assert.Equal(t, 35.996656, front.Position.Latitude)
	assert.Equal(t, 2, len(queue.Items))

	// given
	queue = NewUpdateQueue(10)

	// when

	// then
	front, err = queue.Front()
	assert.Error(t, err)
	assert.Nil(t, front)

	// given
	queue = NewUpdateQueue(1)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds,
	})

	// then
	front, err = queue.Front()
	assert.NoError(t, err)
	assert.Equal(t, 35.996656, front.Position.Latitude)
	assert.Equal(t, 1, len(queue.Items))
}

func TestUpdateQueue_Dequeue(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	// given
	queue := NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})

	// then
	last, err := queue.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 35.996644, last.Position.Latitude)
	assert.Equal(t, 0, len(queue.Items))

	// given
	queue = NewUpdateQueue(10)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds,
	})

	// then
	last, err = queue.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 35.996644, last.Position.Latitude)
	assert.Equal(t, 1, len(queue.Items))

	// given
	queue = NewUpdateQueue(10)

	// when

	// then
	last, err = queue.Dequeue()
	assert.Error(t, err)
	assert.Nil(t, last)
}

func TestUpdateQueue_AverageSpeed(t *testing.T) {
	nowSeconds := int(time.Now().Unix())

	// given
	queue := NewUpdateQueue(2)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds + 100,
	})

	// then
	speed, err := queue.AverageSpeed()
	assert.NoError(t, err)
	assert.Equal(t, 1.4772431640476942, speed)

	// given
	queue = NewUpdateQueue(3)

	// when
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996644,
			Longitude: -78.902220,
		},
		Timestamp: nowSeconds,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.996656,
			Longitude: -78.903862,
		},
		Timestamp: nowSeconds + 100,
	})
	queue.Enqueue(&PositionUpdate{
		PatternID: 2,
		Position: &Position{
			Latitude:  35.995972,
			Longitude: -78.9003702,
		},
		Timestamp: nowSeconds + 300,
	})

	// then
	speed, err = queue.AverageSpeed()
	assert.NoError(t, err)
	assert.Equal(t, 1.560522607996258, speed)

	// given
	queue = NewUpdateQueue(2)

	// when

	// then
	speed, err = queue.AverageSpeed()
	assert.Error(t, err)
	assert.Equal(t, float64(0), speed)
}
