package mc3

import (
	"github.com/mitchellh/mapstructure"
	"encoding/json"
	"regexp"
	"errors"
	"time"
	"fmt"
)

const(
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
****************************	 GSM: TYPE DEFINITIONS	************************************
*******************************************************************************************/

var(
	GSM_portOpen bool = false
	GSM_sp sp_config_t
)

/*******************************************************************************************
********************************   GSM: BASIC FUNCTIONS  ***********************************
*******************************************************************************************/

// This function opens the selected COMPort (name), at a configurable Baud Rate (br).
// The common used baud rate is 115200 but this can change.
//   mc3 := NewMC3()
//   err := mc3.GSM_OpenPort("COM1", 115200)
//   if err != nil {
//   	// Port Open NOK. Maybe log that error?. Do things	
//   } else {
//   	// Port Open OK. Enjoy!	
//   }
func (mc3 *MC3) GSM_OpenPort(name string, br uint32) (err error) {	
	//Marshal Serial port configuration before calling Serial Port PIWI
	var portConfig = sp_config_t { Port: name, Baudrate: br	}
	json_2_send, e := json.Marshal(portConfig)
	if e != nil {
		err = e
	} else {
		//Call "Open" Method from Serial Port PIWI, return port status (0.-Closed, 1.-Open*, 2.-Error )
		var payload_rx piwiError
		err = mapstructure.Decode(mc3.model.Call("serialport", "Open", json_2_send), &payload_rx)
		if err != nil {
			// Do nothing
		} else {
			err = payload_rx.GetError()
			if err != nil {
				// Do nothing
			} else {
				GSM_portOpen = true
				GSM_sp = portConfig
				mc3.Log("INF","GSM","Serial port open.")
			}
		}
	}
	mc3.gsm_log_error(err)
	return
}

// This function close the GSM COM Port ( previously defined by GSM_OpenPort).
//   mc3 := NewMC3()
//   err := mc3.GSM_OpenPort("COM1", 115200)
//   if err != nil {
//   	// Port Open NOK. Maybe log that error?. Do things	
//   } else {
//   	// Port Open OK. Enjoy!	
//   	// Do things
//   	mc3.GSM_ClosePort()	
//   }
// NOTE: If you have Opened the GSM Serial Port you must use the GSM_ClosePort function to prevent future errors.
func (mc3 *MC3) GSM_ClosePort() (err error) {
	if GSM_portOpen {
		//Call "Open" Method from Serial Port PIWI, return port status (0.-Closed, 1.-Open*, 2.-Error )
		var payload_rx piwiError
		err = mapstructure.Decode(mc3.model.Call("serialport", "Close", uint8(0)), &payload_rx)
		if err != nil {
			// Do nothing
		} else {
			err = payload_rx.GetError()
			if err != nil {
				// Do nothing
			} else {
				GSM_portOpen = false
				GSM_sp.Port = ""
				GSM_sp.Baudrate = 0
				mc3.Log("INF","GSM", "Serial port closed.")
			}
		}
	} else {
		err = errors.New("Unable to Close Port, Serial Port is not open.")
	}
	mc3.gsm_log_error(err)
	return
}

/*
This function call the "Write" method from "serialport" PIWI. 
Returns nil if the serial port write was succesful, otherwise returns an error.
*/
func (mc3 *MC3) GSM_print(data string) (err error) {
	var payload_rx piwiError
	mapstructure.Decode(mc3.model.Call("serialport", "Write", data), &payload_rx)
	err = payload_rx.GetError()
	if err != nil {
		// Do nothing
	} else {
		mc3.Log("INF","GSM","Tx >> " + data)
	}
	mc3.gsm_log_error(err)
	return
}

/*
This function call the "WriteLine" method from "serialport" PIWI. 
Returns nil if the serial port write was succesful, otherwise returns an error.
*/
func (mc3 *MC3) GSM_println(data string) (err error) {
	var payload_rx piwiError
	mapstructure.Decode(mc3.model.Call("serialport", "WriteLine", data), &payload_rx)
	err = payload_rx.GetError()
	if err != nil {
		// Do nothing
	} else {
		mc3.Log("INF","GSM","Tx >> " + data)
	}
	mc3.gsm_log_error(err)
	return
}

/*
This function writes a value to the specified EID (both string format).
Example: To write 20099 (integer value) to the EID 0x41B should be written like
	GSM_WriteEID("0x41B","20099")
Example: For string values they must be written as raw string literals (between back quotes ``):
	GSM_WriteEID("0x41A",`"thunderfish.moreycorp.com"`)
*/
func (mc3 *MC3) GSM_SendSMS(number, msg string) (echo string, err error) {
	_, err = mc3.GSM_CheckSMSMode()
	if err != nil {
		// Do nothing
	} else {
		cmd := fmt.Sprintf(CMD_CMGS + "\r\n", number)
		err = mc3.GSM_print(cmd)
		if err !=  nil {
			// Do nothing
		} else {
			str, error_found := mc3.GSM_WaitForRegexTimeout(CMD_ERROR_REGEXP, time.Second * 1 )
			if !error_found {
				cmd := msg + CMD_CTRL_Z
				err = mc3.GSM_print(cmd)
				if err != nil {
					// Do nothing
				} else {
					str, valid := mc3.GSM_WaitForRegexTimeout(CMD_CMGS_RX_REGEXP + "|" + CMD_ERROR_REGEXP, time.Second * 10 )
					if valid {
						mc3.Log("INF","GSM","Rx << " + str)
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
			} else {
				err = errors.New("CMD ERROR: " + cmd + " >> " + str)
			}
		}
		mc3.gsm_log_error(err)
	}
	return
}

func (mc3 *MC3) GSM_WaitSMS(timeout time.Duration) (id string, err error) {
	_, err = mc3.GSM_CheckSMSMode()
	if err != nil {
		// Do nothing
	} else {
		str, found := mc3.GSM_WaitForRegexTimeout(CMD_CMTI_REGEXP, timeout)
		if found {
			if len(str) > len(CMD_CMTI_RX) {
				id = str[len(CMD_CMTI_RX):]
			}
		} else {
			err = errors.New("Timeout expired")
		}
		mc3.gsm_log_error(err)
	}
	return
}

func (mc3 *MC3) GSM_ReadSMS(id string) (msg string, err error) {
	_, err = mc3.GSM_CheckSMSMode()
	if err != nil {
		// Do nothing
	} else {
		cmd := fmt.Sprintf(CMD_CMGR + "\r\n", id)
		err = mc3.GSM_print(cmd)
		if err !=  nil {
			// Do nothing
		} else {
			_, found := mc3.GSM_WaitForRegexTimeout(CMD_CMGR_REGEXP, time.Second * 10)
			if found {
				// if len(str) > len(CMD_CMGR_RX) {
				// 	msg = str[len(CMD_CMGR_RX):]
				// }
				msg, _ = mc3.GSM_WaitForRegexTimeout(".*", time.Second * 1)
			} else {
				err = errors.New("Timeout expired")
			}
			mc3.gsm_log_error(err)
		}
	}
	return
}

func (mc3 *MC3) GSM_CheckSMSMode() (echo string, err error) {
	_, err = mc3.GSM_Ping()
	if err != nil {
		// Do nothing
	} else {
		cmd := CMD_CMGF + "\r\n"
		err = mc3.GSM_print(cmd)
		if err !=  nil {
			// Do nothing
		} else {
			str, valid := mc3.GSM_WaitForRegexTimeout(CMD_OK_REGEXP + "|" + CMD_ERROR_REGEXP, time.Second * 1 )
			if valid {
				mc3.Log("INF","GSM","Rx << " + str)
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
		mc3.gsm_log_error(err)
	}
	return
}

func (mc3 *MC3) GSM_Ping() (echo string, err error) {
	cmd := CMD_AT
	err = mc3.GSM_print(cmd + "\r\n")
	if err !=  nil {
		// Do nothing
	} else {
		str, valid := mc3.GSM_WaitForRegexTimeout(CMD_OK_REGEXP + "|" + CMD_ERROR_REGEXP, time.Second * 1 )
		if valid {
			echo = str
			mc3.Log("INF","GSM","Rx << " + str)
			// Check if there is an update in progress
			error, _  := regexp.Match("ERROR", []byte(str))
			if error {
				err = errors.New("CMD ERROR: " + cmd + " >> " + str)
			}
		}
	}
	mc3.gsm_log_error(err)
	return
}

/*
This function wait for a regular expression to be received by the "serialport" PIWI including a timeout
Is recommended to use this function instead of "CLI_WaitForRegex"
Returns 2 parameters:
	string: The string in wich the RegExp was found
	bool: Will return true if the RegExp was found
*/
func (mc3 *MC3) GSM_WaitForRegexTimeout(regexp string, timeout time.Duration) (string, bool) {
	//Encode data to send
	var data_tx = Regexp_Timeout_t {
	 	regexp,
		timeout,
	}
	payload_tx, _ := json.Marshal(data_tx)

	//Wait for RegEx from Serialport
	payload_rx := mc3.model.Call("serialport", "WaitForRegexTimeout", payload_tx).(interface{})

	//Decode received data
	var data_rx Regexp_Response_t
	json.Unmarshal(payload_rx.([]byte), &data_rx)
	return data_rx.Regexp, data_rx.Found
}

func (mc3 *MC3) gsm_log_error(e error) {
	if e != nil {
		mc3.Log("ERR","GSM", e.Error())
	}
}