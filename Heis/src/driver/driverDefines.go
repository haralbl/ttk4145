package driver

//in port 4
const pORT4  			=   3
const oBSTRUCTION 		=   (0x300+23)
const sTOP_SENSOR		=   (0x300+22)
const bUTTON_COMMAND1	=  	(0x300+21)
const bUTTON_COMMAND2   =   (0x300+20)
const bUTTON_COMMAND3   =   (0x300+19)
const bUTTON_COMMAND4   =   (0x300+18)
const bUTTON_UP1        =   (0x300+17)
const bUTTON_UP2        =   (0x300+16)

//in port 1
const pORT1             =  2
const bUTTON_DOWN2      =  (0x200+0)
const bUTTON_UP3        =  (0x200+1)
const bUTTON_DOWN3      =  (0x200+2)
const bUTTON_DOWN4      =  (0x200+3)
const sENSOR_FLOOR1     =  (0x200+4)
const sENSOR_FLOOR2     =  (0x200+5)
const sENSOR_FLOOR3     =  (0x200+6)
const sENSOR_FLOOR4     =  (0x200+7)

//out port 3
const pORT3             =  3
const mOTORDIR          =  (0x300+15)
const lIGHT_STOP        =  (0x300+14)
const lIGHT_COMMAND1    =  (0x300+13)
const lIGHT_COMMAND2    =  (0x300+12)
const lIGHT_COMMAND3    =  (0x300+11)
const lIGHT_COMMAND4    =  (0x300+10)
const lIGHT_UP1         =  (0x300+9)
const lIGHT_UP2         =  (0x300+8)

//out port 2
const pORT2             =  3
const lIGHT_DOWN2       =  (0x300+7)
const lIGHT_UP3         =  (0x300+6)
const lIGHT_DOWN3       =  (0x300+5)
const lIGHT_DOWN4       =  (0x300+4)
const lIGHT_DOOR_OPEN   =  (0x300+3)
const lIGHT_FLOOR_IND2  =  (0x300+1)
const lIGHT_FLOOR_IND1  =  (0x300+0)

//out port 0
const pORT0              = 1
const mOTOR              = (0x100+0)

//non-existing ports (for alignment)
const bUTTON_DOWN1       = -1
const bUTTON_UP4         = -1
const lIGHT_DOWN1        = -1
const lIGHT_UP4          = -1

