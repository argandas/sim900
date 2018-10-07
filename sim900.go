package sim900

import (
	"errors"
	"fmt"
	"github.com/argandas/serial"
	"log"
	"os"
	"strings"
	"time"
)

/*******************************************************************************************
********************************	TYPE DEFINITIONS	************************************
*******************************************************************************************/

// A SIM900 is the representation of a SIM900 GSM modem with several utility features.
type SIM900 struct {
	port   *serial.SerialPort
	logger *log.Logger
}

/*******************************************************************************************
********************************   GSM: BASIC FUNCTIONS  ***********************************
*******************************************************************************************/

// New creates and initializes a new SIM900 device.
func New() *SIM900 {
	return &SIM900{
		port:   serial.New(),
		logger: log.New(os.Stdout, "[sim900] ", log.LstdFlags),
	}
}

// Connect creates a connection with the SIM900 modem via serial port and test communications.
func (s *SIM900) Connect(port string, baud int) error {
	// Open device serial port
	if err := s.port.Open(port, baud, time.Millisecond*100); err != nil {
		return err
	}
	// Ping to Modem
	return s.Ping()
}

func (sim *SIM900) Disconnect() error {
	// Close device serial port
	return sim.port.Close()
}

// Ignore response
func (sim *SIM900) ignoreReponse(cmd string, timeout time.Duration) (error) {
	// Send command
	if err := sim.port.Println(cmd); err != nil {
		return err
	}
	// ignore
	_, err := sim.port.ReadLine()
	if err != nil {
		return err
	}
	return nil
}

func (sim *SIM900) wait4response(cmd, expected string, timeout time.Duration) (string, error) {
	// Send command
	if err := sim.port.Println(cmd); err != nil {
		return "", err
	}
	// Wait for command response
	regexp := expected + "|" + CMD_ERROR
	response, err := sim.port.WaitForRegexTimeout(regexp, timeout)
	if err != nil {
		return "", err
	}
	// Check if response is an error
	if strings.Contains(response, "ERROR") {
		return response, errors.New("Errors found on command response")
	}
	// Response received succesfully
	return response, nil
}

// Send a SMS
func (s *SIM900) SendSMS(number, msg string) error {
	// Set message format
	if err := s.SetSMSMode(TEXT_MODE); err != nil {
		return err
	}
	// Send command
	cmd := fmt.Sprintf(CMD_CMGS, number)
	if err := s.port.Println(cmd); err != nil {
		return err
	}
	// Wait modem to be ready
	time.Sleep(time.Second * 1)
	// Send message
	_, err := s.wait4response(msg+CMD_CTRL_Z, CMD_OK, time.Second*5)
	if err != nil {
		return err
	}
	// Message sent succesfully
	return nil
}

// WaitSMS will return when either a new SMS is recived or the timeout has expired.
// The returned value is the memory ID of the received SMS, use ReadSMS to read SMS content.
func (s *SIM900) WaitSMS(timeout time.Duration) (id string, err error) {
	id, err = s.wait4response("", CMD_CMTI_REGEXP, timeout)
	if err != nil {
		return
	}
	if len(id) >= len(CMD_CMTI_RX) {
		id = id[len(CMD_CMTI_RX):]
	}
	return
}

// ReadAllSMS retrieves all SMS text
func (s *SIM900) ReadAllSMS() (msg string, err error) {
	// Set message format
	if err := s.SetSMSMode(TEXT_MODE); err != nil {
		return "", err
	}
	// Send command
	cmd := fmt.Sprintf(CMD_CMGL_ALL)
	if _, err := s.wait4response(cmd, CMD_CMGR_REGEXP, time.Second*5); err != nil {
		return "", err
	}
	// Reading succesful get message data
	return s.port.ReadLine()
}

// ReadSMS retrieves SMS text from inbox memory by ID.
func (s *SIM900) ReadSMS(id string) (msg string, err error) {
	// Set message format
	if err := s.SetSMSMode(TEXT_MODE); err != nil {
		return "", err
	}
	// Send command
	cmd := fmt.Sprintf(CMD_CMGR, id)
	if _, err := s.wait4response(cmd, CMD_CMGR_REGEXP, time.Second*5); err != nil {
		return "", err
	}
	// Reading succesful get message data
	return s.port.ReadLine()
}

// ReadSMS deletes SMS from inbox memory by ID.
func (s *SIM900) DeleteSMS(id string) error {
	// Send command
	cmd := fmt.Sprintf(CMD_CMGD, id)
	_, err := s.wait4response(cmd, CMD_OK, time.Second*1)
	return err
}

// Ping modem
func (s *SIM900) Ping() error {
	_, err := s.wait4response(CMD_AT, CMD_OK, time.Second*1)
	return err
}


// Modem echo
func (s *SIM900) Echo(enable bool) error {
	var cmd string
	if enable {
		cmd = fmt.Sprintf(CMD_ATE, 1)
	} else {
		cmd = fmt.Sprintf(CMD_ATE, 0)
	}
	_, err := s.wait4response(cmd, "OK", time.Second*1)
	return err
}

// SetSMSMode selects SMS Message Format ("0" = PDU mode, "1" = Text mode)
func (s *SIM900) SetSMSMode(mode string) error {
	cmd := fmt.Sprintf(CMD_CMGF_SET, mode)
	_, err := s.wait4response(cmd, CMD_OK, time.Second*1)
	return err
}

// SetSMSMode reads SMS Message Format (0 = PDU mode, 1 = Text mode)
func (s *SIM900) SMSMode() (mode string, err error) {
	mode, err = s.wait4response(CMD_CMGF, CMD_CMGF_REGEXP, time.Second*1)
	if err != nil {
		return
	}
	if len(mode) >= len(CMD_CMGF_RX) {
		mode = mode[len(CMD_CMGF_RX):]
	}
	return
}

// SetSMSMode selects SMS Message Format (0 = PDU mode, 1 = Text mode)
func (s *SIM900) CheckSMSTextMode(mode int) error {
	cmd := fmt.Sprintf(CMD_CMGF, mode)
	_, err := s.wait4response(cmd, CMD_OK, time.Second*1)
	return err
}
