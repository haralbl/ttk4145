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
    Set_motor_direction(UP);

    for {
    	fmt.Printf("I'm running\n")
        // Change direction when we reach top/bottom floor
        if Get_floor_sensor_signal() == NFloors - 1 {
            Set_motor_direction(DOWN)
        } else if Get_floor_sensor_signal() == 0 {
            Set_motor_direction(UP)
        }

        // Stop elevator and exit program if the stop button is pressed
        if Get_stop_signal() != 0 {
            Set_motor_direction(STOP)
            break
        }
    }
}
