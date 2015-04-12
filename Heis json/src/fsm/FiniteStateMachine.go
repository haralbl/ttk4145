package fsm

import (
	."fmt"
)

const (
	IDLE			= 0
	DOOR_OPEN		= 1
	MOVING			= 2
)

var (
	state		int	= 0
)

func floorReached() {
	// set floor lights
	switch state {
	case IDLE:
		// do nothing? (shouldnt happen)
	case DOOR_OPEN:
		// do nothing? (shouldnt happen)
	case MOVING:
		// maybe stop
	}
}

func doorTimerOut() {
	switch state {
	case IDLE:
		// close door? (shouldnt happen)
	case DOOR_OPEN:
		// close door
	case MOVING:
		// close door? (shouldnt happen)
	}
}

func newOrder() {
	// send/add the order
	switch state {
	case IDLE:
	case DOOR_OPEN:
	case MOVING:
	}
	
	// sette igang heisen...
}

func outOrder() {
	switch state {
	case IDLE:
		if /*in floor*/ {
			// open door
			// reset door timer
		} else {
			// update status without sending? (takes this oneself anyways)
		}
	case DOOR_OPEN:
		if /*in floor*/ {
			// reset door timer
		} else {
			// update status without sending? (takes this oneself anyways)
		}
	case MOVING:
		// update status without sending? (takes this oneself anyways)
	}
}








