package defines

const (
	NumberOfElevators	= 10
	NumberOfFloors		= 4
	NumberOfButtonTypes = 3
	
	IDLE				= 0
	DOOR_OPEN			= 1
	MOVING				= 2
	
	BUTTON_CALL_UP 		int = 0
	BUTTON_CALL_DOWN 	int = 1
	BUTTON_COMMAND 		int = 2
	
	UP 					= 1
	DOWN 				= -1
	STOP				= 0
)

type ElevatorStatus_t struct {
	ActiveElevators		[NumberOfElevators]string // IP addresses
 
	PreviousFloors		[NumberOfElevators]int
	InFloor				[NumberOfElevators]int
	Directions			[NumberOfElevators]int

	OrdersUp			[NumberOfElevators][NumberOfFloors]int
	OrdersDown			[NumberOfElevators][NumberOfFloors]int
	OrdersOut			[NumberOfElevators][NumberOfFloors]int

	State				int
	
	MessageType			string
	OrderedButtonType	int
	OrderedElevator		string
	OrderedFloor		int	
}












