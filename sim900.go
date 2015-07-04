package sim900

import (
	"errors"
	"fmt"
	"github.com/argandas/serial"
	"time"
)

/*******************************************************************************************
********************************	TYPE DEFINITIONS	************************************
*******************************************************************************************/

type SIM900 struct {
	port serial.SerialPort
}

/*******************************************************************************************
********************************   GSM: BASIC FUNCTIONS  ***********************************
*******************************************************************************************/

func New() SIM900 {
	return SIM900{
		port: serial.New(),
	}
}

func (sim *SIM900) Setup(name string, baud int) error {
	return sim.port.Open(name, baud, time.Millisecond*100)
}

func (sim *SIM900) Teardown() error {
	return sim.port.Close()
}

func (sim *SIM900) echo(data, echo string) error {
	err := sim.Ping()
	if err != nil {
		return err
	} else {
		err := sim.port.Println(data)
		if err != nil {
			return err
		} else {
			_, err := sim.port.WaitForRegexTimeout(echo, time.Second*1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Send a SMS
func (sim *SIM900) SendSMS(number, msg string) error {
	err := sim.echo(CMD_CMGF, CMD_OK)
	if err != nil {
		return err
	} else {
		cmd := fmt.Sprintf(CMD_CMGS, number)
		err := sim.echo(cmd, ">")
		if err != nil {
			return errors.New("CMD ERROR: " + cmd + " >> " + err.Error())
		} else {
			cmd := msg + CMD_CTRL_Z
			err := sim.echo(cmd, CMD_OK)
			if err != nil {
				return errors.New("CMD ERROR: " + cmd + " >> " + err.Error())
			} else {
				_, err := sim.port.WaitForRegexTimeout(CMD_CMGS_RX_REGEXP, time.Second*5)
				if err != nil {
					return errors.New("CMD ERROR: " + err.Error())
				}
			}
		}
	}
	return nil
}

// Wait for a new SMS to come
func (sim *SIM900) WaitSMS(timeout time.Duration) (string, error) {
	err := sim.echo(CMD_CMGF, CMD_OK)
	if err != nil {
		return "", err
	} else {
		data, err := sim.port.WaitForRegexTimeout(CMD_CMTI_REGEXP, timeout)
		if err != nil {
			return "", errors.New("CMD ERROR: " + err.Error())
		} else {
			if len(data) > len(CMD_CMTI_RX) {
				return data[len(CMD_CMTI_RX):], nil
			}
		}
	}
	return "", nil
}

// Read SMS by ID
func (sim *SIM900) ReadSMS(id string) (string, error) {
	err := sim.echo(CMD_CMGF, CMD_OK)
	if err != nil {
		return "", err
	} else {
		cmd := fmt.Sprintf(CMD_CMGR, id)
		err := sim.echo(cmd, CMD_CMGR_REGEXP)
		if err != nil {
			return "", err
		} else {
			return sim.port.ReadLine()
		}
	}
	return "", nil
}

// Check if there are unread SMS
func (sim *SIM900) UnreadSMS() bool {
	return false
}

// Ping modem
func (sim *SIM900) Ping() error {
	return sim.echo(CMD_AT, CMD_OK)
}
