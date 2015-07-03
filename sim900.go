package sim900

import (
	"fmt"
	"https://github.com/argandas/serial"
)

const (
	CMD_ERROR_REGEXP string = "(^ERROR$)"
	CMD_OK_REGEXP string = "(^OK$)"
	CMD_AT string = "AT"
	CMD_CMGF string = "AT+CMGF=1"
	CMD_CMGS string = "AT+CMGS=\"%s\""
	CMD_CTRL_Z string = "\x1A"
	CMD_CMGS_RX_REGEXP string = "(^[+]CMGS[:] [0-9]+$)"
	CMD_CMTI_REGEXP string = "(^[+]CMTI[:] \"SM\",[0-9]+$)" 
	CMD_CMTI_RX string = "+CMTI: \"SM\","
	CMD_CMGR string = "AT+CMGR=%s"
	CMD_CMGR_REGEXP string = "(^[+]CMGR[:] .*)" 
	CMD_CMGR_RX string = "+CMGR: "
)

/*******************************************************************************************
********************************	TYPE DEFINITIONS	************************************
*******************************************************************************************/

type SIM900 struct {
	port serial.Serial
}

var(
	GSM_portOpen bool = false
	GSM_sp sp_config_t
)

/*******************************************************************************************
********************************   GSM: BASIC FUNCTIONS  ***********************************
*******************************************************************************************/

func New() SIM900 {
	return SIM900 {
		port: serial.New(),
	}
}

func (sim *SIM900) Setup(name string, br uint32) error {	
	
}

func (sim *SIM900) Teardown() error {	
	
}

// Send a SMS
func (sim *SIM900) SendSMS(number, msg string) (err error) {	
	mode := sim.mode()
	if mode == MODE_SMS {
		cmd := fmt.Sprintf(CMD_CMGS, number)
		e := sim.port.Print(cmd)
		if e !=  nil {
			err = errors.New("CMD ERROR: " + cmd + " >> " + e.Error())
		} else {
			str, e := sim.port.WaitForRegexTimeout(CMD_ERROR_REGEXP, time.Second * 1 )
			if e !=  nil {
				err = errors.New("CMD ERROR: " + cmd + " >> " + e.Error())
			} else {
				cmd := msg + CMD_CTRL_Z
				e := sim.port.Print(cmd)
				if e != nil {
					err = errors.New("CMD ERROR: " + cmd + " >> " + e.Error())
				} else {
					str, e := sim.port.WaitForRegexTimeout(CMD_CMGS_RX_REGEXP + "|" + CMD_ERROR_REGEXP, time.Second * 10 )
					if e != nil {
						err = errors.New("CMD ERROR: " + cmd + " >> " + e.Error())
					} else {
						// Check if there is an update in progress
						error, _  := regexp.Match("ERROR", []byte(str))
						if !error {
							echo = str
						} else {
							// Wait for update to be done
							err = errors.New("CMD ERROR: " + cmd + " >> " + str)
						}
					}
				}
			}
		}
	}
}

// Wait for a new SMS to come
func (sim *SIM900) WaitSMS(timeout time.Duration) error {	
	
}

// Read SMS by ID
func (sim *SIM900) ReadSMS(id string) error {	
	
}

// Check if there are unread SMS
func (sim *SIM900) UnreadSMS() bool {	
	
}

// Ping modem
func (sim *SIM900) Ping() error {	
	
}

// Check check current mode
func (sim *SIM900) getMode() Mode {	
	
}

// Set modem mode
func (sim *SIM900) setMode(mode Mode) error {	
	
}