package driver

import(
	"os"
	"fmt"
	"time"
	"defines"
)

type direction int


var (
	//If number of floors is changed, you will have to add a larger number of addresses for the different buttons and lights 
	lamp_channel_matrix [defines.NumberOfFloors][defines.NumberOfButtonTypes] int = [defines.NumberOfFloors][defines.NumberOfButtonTypes]int{
	    {lIGHT_UP1, lIGHT_DOWN1, lIGHT_COMMAND1},
	    {lIGHT_UP2, lIGHT_DOWN2, lIGHT_COMMAND2},
	    {lIGHT_UP3, lIGHT_DOWN3, lIGHT_COMMAND3},
	    {lIGHT_UP4, lIGHT_DOWN4, lIGHT_COMMAND4},
	}

	//If number of floors is changed, you will have to add a larger number of addresses for the different buttons and lights
	button_channel_matrix [defines.NumberOfFloors][defines.NumberOfButtonTypes] int = [defines.NumberOfFloors][defines.NumberOfButtonTypes]int{
	    {bUTTON_UP1, bUTTON_DOWN1, bUTTON_COMMAND1},
	    {bUTTON_UP2, bUTTON_DOWN2, bUTTON_COMMAND2},
	    {bUTTON_UP3, bUTTON_DOWN3, bUTTON_COMMAND3},
	    {bUTTON_UP4, bUTTON_DOWN4, bUTTON_COMMAND4},
	}
)

func Init() int {
	fmt.Printf("I'm running driver\n")
	initStat := Io_init()
	
	if initStat == 0 {
		fmt.Printf("init failed")
		return 0
	}
	
    for i := 0; i < defines.NumberOfFloors; i++ {
        if i != 0 {
            Set_button_lamp(defines.BUTTON_CALL_DOWN, i, 0)
        }
        if i != defines.NumberOfFloors - 1 {
            Set_button_lamp(defines.BUTTON_CALL_UP, i, 0)
        }
        Set_button_lamp(defines.BUTTON_COMMAND, i, 0)
    }
	
    // Clear stop lamp, door open lamp, and set floor indicator to ground floor.
    set_stop_lamp(0)
    Set_door_open_lamp(0)
    Set_floor_indicator(0)
	
    // Return success.
    return 1
}

func UpButtonPoller(upButtonChan chan int) {
	var upButtonFlags [defines.NumberOfFloors-1]int

	for {

		for i := 0; i < defines.NumberOfFloors-1; i++ {
			if upButtonFlags[i] == 0 {
				if get_button_signal(defines.BUTTON_CALL_UP, i) == 1 {
					upButtonFlags[i] = 1
					upButtonChan <- i
				}
			} else {
				if get_button_signal(defines.BUTTON_CALL_UP, i) == 0 {
					upButtonFlags[i] = 0
				}
			}
		}
		time.Sleep(time.Millisecond*10)
	}
}

func DownButtonPoller(downButtonChan chan int) {
	var downButtonFlags [defines.NumberOfFloors-1]int
	for {
		for i := 1; i < defines.NumberOfFloors; i++ {
			if downButtonFlags[i-1] == 0 {
				if get_button_signal(defines.BUTTON_CALL_DOWN, i) == 1 {
					downButtonFlags[i-1] = 1
					downButtonChan <- i
				}
			} else {
				if get_button_signal(defines.BUTTON_CALL_DOWN, i) == 0 {
					downButtonFlags[i-1] = 0
				}
			}
		}
		time.Sleep(time.Millisecond*10)
	}
}

func CommandButtonPoller(commandButtonChan chan int) {
	var commandButtonFlags [defines.NumberOfFloors]int
	for {
		for i := 0; i < defines.NumberOfFloors; i++ {
			if commandButtonFlags[i] == 0 {
				if get_button_signal(defines.BUTTON_COMMAND, i) == 1 {
					commandButtonFlags[i] = 1
					commandButtonChan <- i
				}
			} else {
				if get_button_signal(defines.BUTTON_COMMAND, i) == 0 {
					commandButtonFlags[i] = 0
				}
			}
		}
		time.Sleep(time.Millisecond*10)
	}
}

func FloorPoller(floorReachedChan chan int, floorLeftChan chan string) {
	floorReachedFlag 	:= 0
	floorLeftFlag 		:= 0 
	currFloor 			:= -1 
	for { 
		currFloor = Get_floor_sensor_signal() 
		if floorReachedFlag == 0 { 
			if currFloor != -1 { 
				floorReachedFlag = 1
				floorReachedChan <- currFloor
				floorLeftFlag = 0
			}
		} else { 
			if currFloor == -1 { 
				floorReachedFlag = 0
				if floorLeftFlag == 0 {
					floorLeftFlag = 1
					floorLeftChan <- "floor left"
				}
			} 
		} 
		time.Sleep(time.Millisecond*10)
	} 
}

func Set_motor_direction(dir direction) {
    if dir == 0 {
        io_write_analog(mOTOR, 0)
    } else if dir > 0 {
        io_clear_bit(mOTORDIR)
        io_write_analog(mOTOR, 2800)
    } else if dir < 0 {
        io_set_bit(mOTORDIR)
        io_write_analog(mOTOR, 2800)
    }
}

func Set_door_open_lamp(value int) {
    if value != 0 {
        io_set_bit(lIGHT_DOOR_OPEN)
    } else {
        io_clear_bit(lIGHT_DOOR_OPEN)
    }
}

func get_obstruction_signal() int {
	return io_read_bit(oBSTRUCTION)
}

func Get_stop_signal() int {
    return io_read_bit(sTOP_SENSOR)
}

func set_stop_lamp(value int) {
    if value != 0 {
        io_set_bit(lIGHT_STOP)
    } else {
        io_clear_bit(lIGHT_STOP)
    }
}

func Get_floor_sensor_signal() int {
    if io_read_bit(sENSOR_FLOOR1) != 0 {
        return 0
    } else if io_read_bit(sENSOR_FLOOR2) != 0 {
        return 1
    } else if io_read_bit(sENSOR_FLOOR3) != 0 {
        return 2
    } else if io_read_bit(sENSOR_FLOOR4) != 0{
        return 3
    } else {
        return -1
    }
}

func Set_floor_indicator(floor int) {
    if floor < 0 || floor >= defines.NumberOfFloors {
        os.Exit(1)
    }
    // Binary encoding. One light must always be on.
    if (floor & 0x02) != 0 {
        io_set_bit(lIGHT_FLOOR_IND1)
    } else {
        io_clear_bit(lIGHT_FLOOR_IND1)
    }
    
    if (floor & 0x01) != 0 {
        io_set_bit(lIGHT_FLOOR_IND2)
    } else {
        io_clear_bit(lIGHT_FLOOR_IND2)
    }
}

func get_button_signal(button int,floor int) int {
    if floor < 0 || floor >= defines.NumberOfFloors {
        fmt.Printf("reading button on invalid floor")
    }
    if (!(button == defines.BUTTON_CALL_UP && floor == defines.NumberOfFloors - 1)) || (!(button == defines.BUTTON_CALL_DOWN && floor == 0)) {
        //fmt.Printf("bad combination of button and floor")
    }
    if (!(button == defines.BUTTON_CALL_UP || button == defines.BUTTON_CALL_DOWN || button == defines.BUTTON_COMMAND)) {
        fmt.Printf("nonvalid button type")
    }

    if io_read_bit(button_channel_matrix[floor][button]) != 0 {
        return 1
    } else{
        return 0
    }
}

func Set_button_lamp(button int, floor int, value int) {
    if floor < 0 || floor >= defines.NumberOfFloors {
    	//os.Exit(1)
    	fmt.Printf("setting lamp in nonvalid floor")
    }
    if (!(button == defines.BUTTON_CALL_UP && floor == defines.NumberOfFloors - 1)) || (!(button == defines.BUTTON_CALL_DOWN && floor == 0)) {
    	//os.Exit(1)
    	//fmt.Printf("bad combination of button and floor")
    }
    if (!(button == defines.BUTTON_CALL_UP || button == defines.BUTTON_CALL_DOWN || button == defines.BUTTON_COMMAND)) {
    	//os.Exit(1)
    	fmt.Printf("nonvalid button type")
    }

    if value != 0 {
    	io_set_bit(lamp_channel_matrix[floor][button])
    } else {
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
}













