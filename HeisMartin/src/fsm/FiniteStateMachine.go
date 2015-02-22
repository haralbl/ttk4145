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
	// TODO set floor lights
	switch state {
	case IDLE:
	case DOOR_OPEN:
	case MOVING:
	}
}

func doorTimerOut() {
	switch state {
	case IDLE:
	case DOOR_OPEN:
	case MOVING:
	}
}

func newOrder() {
	switch state {
	case IDLE:
	case DOOR_OPEN:
	case MOVING:
	}
}

func outOrder() {
	switch state {
	case IDLE:
	case DOOR_OPEN:
	case MOVING:
	}
}
