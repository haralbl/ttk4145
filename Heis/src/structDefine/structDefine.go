package structDefine

const (
	numberOfElevators	= 3
	numberOfFloors		= 4
)

type ElevatorStatus_t struct {
	ActiveElevators		[numberOfElevators]string // IP addresses
 
	PreviousFloors		[numberOfElevators]int
	InFloor				[numberOfElevators]int
	Directions			[numberOfElevators]int

	OrdersUp			[numberOfElevators][numberOfFloors]int
	OrdersDown			[numberOfElevators][numberOfFloors]int
	OrdersOut			[numberOfElevators][numberOfFloors]int

	State				int
	
	MessageType			string
	OrderedButtonType	int
	OrderedElevator		string
	OrderedFloor		int	
}











