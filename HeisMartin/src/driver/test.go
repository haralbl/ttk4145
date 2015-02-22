package driver
import(
    "os"
    "fmt"
)

func Test() {
	fmt.Printf("I'm running test\n")
    // Initialize hardware
    temp := Init()
    if temp == 0 {
    	fmt.Printf("Exiting\n")
        os.Exit(1)
    }
	fmt.Printf("I'm running\n")
    //set_motor_direction(UP);
	
	//prev_floor := 0
	set_door_open_lamp(1)
	/*
    for {
    	if Get_floor_sensor_signal() != -1  && 0 == 1 {
    		fmt.Println(Get_floor_sensor_signal())
    	}
    */	
    	/*
    	//fmt.Printf("I'm running\n")
        // Change direction when we reach top/bottom floor
        if Get_floor_sensor_signal() == NFloors - 1 {
        	//prev_floor = NFloors - 1
            set_motor_direction(DOWN)
			set_floor_indicator(NFloors-1)
        } else if Get_floor_sensor_signal() == 0 {
            set_motor_direction(UP)
            //prev_floor = 0
			set_floor_indicator(0)
        } else if Get_floor_sensor_signal() == 1 {
        	//prev_floor = 1
        	set_floor_indicator(1)
        } else if Get_floor_sensor_signal() == 2 {
        	//prev_floor = 2
        	set_floor_indicator(2)
        }
        
        for flr := 0; flr < NFloors; flr++ {
        	for btntype := 0; btntype < NButtonTypes; btntype++ {
        		set_button_lamp(btntype, flr, get_button_signal(btntype, flr))
        	}
        }
        */
        
        /*} else {
        	for i := 0; i < 4; i++
        }*/
	/*	
		if get_obstruction_signal() == 1 {
			set_stop_lamp(1)
		}
		
        // Stop elevator and exit program if the stop button is pressed
        if Get_stop_signal() != 0 {
            set_motor_direction(STOP)
            set_stop_lamp(1)
            break
        }
    }
    */ 
}
