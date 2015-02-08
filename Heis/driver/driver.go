package driver

import(
	"os"
	"fmt"
)

type buttonType int
type direction int

const(
	NFloors 		int = 4
	nButtonTypes 	int = 3

	BUTTON_CALL_UP buttonType 	= 0
	BUTTON_CALL_DOWN buttonType = 1
	BUTTON_COMMAND buttonType 	= 2

    UP direction = 1
    DOWN direction = -1
    STOP direction = 0
)

var(
	lamp_channel_matrix [NFloors][nButtonTypes] int = [NFloors][nButtonTypes]int{
	    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}


	button_channel_matrix [NFloors][nButtonTypes] int = [NFloors][nButtonTypes]int{
	    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}
)



func Init() int{
	fmt.Printf("I'm running driver\n")
	initStat := Io_init()
	
	if initStat == 0 {
		fmt.Printf("init failed")
		return 0
	}
	
    for i := 0; i < NFloors; i++ {
        if i != 0 {
            set_button_lamp(BUTTON_CALL_DOWN, i, 0)
        }
        if i != NFloors - 1 {
            set_button_lamp(BUTTON_CALL_UP, i, 0)
        }
        set_button_lamp(BUTTON_COMMAND, i, 0)
    }
	
    // Clear stop lamp, door open lamp, and set floor indicator to ground floor.
    set_stop_lamp(0)
    set_door_open_lamp(0)
    set_floor_indicator(0)
	
    // Return success.
    return 1

}


func Set_motor_direction(dir direction) {
	fmt.Printf("setting motor dir\n")
    if dir == 0 {
        io_write_analog(MOTOR, 0)
    } else if dir > 0 {
        io_clear_bit(MOTORDIR)
        io_write_analog(MOTOR, 2800)
    } else if dir < 0 {
        io_set_bit(MOTORDIR)
        io_write_analog(MOTOR, 2800)
    }
}

func set_door_open_lamp(value int) {
    if value != 0 {
        io_set_bit(LIGHT_DOOR_OPEN)
    } else {
        io_clear_bit(LIGHT_DOOR_OPEN)
    }
}

 func get_obstruction_signal() int {
    return io_read_bit(OBSTRUCTION)
}

func Get_stop_signal() int {
    return io_read_bit(STOP_SENSOR)
}

func set_stop_lamp(value int) {
    if value != 0 {
        io_set_bit(LIGHT_STOP)
    } else {
        io_clear_bit(LIGHT_STOP)
    }
}

func Get_floor_sensor_signal() int{
    if io_read_bit(SENSOR_FLOOR1) != 0 {
        return 0
    } else if io_read_bit(SENSOR_FLOOR2) != 0 {
        return 1
    } else if io_read_bit(SENSOR_FLOOR3) != 0 {
        return 2
    } else if io_read_bit(SENSOR_FLOOR4) != 0{
        return 3
    } else {
        return -1
    }
}

func set_floor_indicator(floor int) {
    if floor < 0 || floor >= NFloors {
        os.Exit(1)
    }
    // Binary encoding. One light must always be on.
    if (floor & 0x02) != 0 {
        io_set_bit(LIGHT_FLOOR_IND1)
    } else {
        io_clear_bit(LIGHT_FLOOR_IND1)
    }
    
    if (floor & 0x01) != 0 {
        io_set_bit(LIGHT_FLOOR_IND2)
    } else {
        io_clear_bit(LIGHT_FLOOR_IND2)
    }
}

func get_button_signal(button buttonType,floor int) int{
    if floor < 0 || floor >= NFloors {
        fmt.Printf("bad combination of button and floor")
    }
    if (!(button == BUTTON_CALL_UP && floor == NFloors - 1)) || (!(button == BUTTON_CALL_DOWN && floor == 0)) {
        fmt.Printf("bad combination of button and floor")
    }
    if (!(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND)) {
        fmt.Printf("bad combination of button and floor")
    }

    if io_read_bit(button_channel_matrix[floor][button]) != 0 {
        return 1
    } else{
        return 0
    }
}

func set_button_lamp(button buttonType, floor int, value int) {
    if floor < 0 || floor >= NFloors {
    	//os.Exit(1)
    	fmt.Printf("bad combination of button and floor")
    }
    if (!(button == BUTTON_CALL_UP && floor == NFloors - 1)) || (!(button == BUTTON_CALL_DOWN && floor == 0)) {
    	//os.Exit(1)
    	fmt.Printf("bad combination of button and floor")
    }
    if (!(button == BUTTON_CALL_UP || button == BUTTON_CALL_DOWN || button == BUTTON_COMMAND)) {
    	//os.Exit(1)
    	fmt.Printf("bad combination of button and floor")
    }

    if value != 0 {
    	io_set_bit(lamp_channel_matrix[floor][button])
    } else {
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
}
