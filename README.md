# SIM900 Go's package
This package uses a serialport to communicate with the SIM900 GSM Modem.

## How to install

- You'll need Golang v1.3+
- SIM900 Package uses the [serial](https://github.com/argandas/serial) package in order to communicate with the modem via AT commands, you will need to install both SIM900 and serial packages.

```bash
go get github.com/argandas/serial  # installs the serial package
go get github.com/argandas/sim900  # installs the SIM900 package
```

## How to use

- You'll need an available serial port, SIM900 boards usually works with 5V TTL signals so you can get a USB-to-Serial TTL converter, I recommend you to use the [FTDI Cable](https://www.sparkfun.com/products/9718) for this, but you can use any USB-to-Serial adapters there are plenty of them. 
![SIM900: FTDI Cable](TBD)

- Connect carefuly your serialport to your SIM900 board.
![SIM900: Connection diagram](TBD)

## Example code

```go
package main
import "github.com/argandas/sim900"

func main() {
	gsm := sim900.New()
	err := gsm.Connect("COM1", 9600)
	if err != nil {
		panic(err)
	}
	defer gsm.Disconnect()
	phoneNumber := "XXXXXXXXXX" // The number to send the SMS
	gsm.SendSMS(phoneNumber, "Hello World!")
}
```

## Reference

- List of available SIM900 commands can be found [here](http://wm.sim.com/upfile/2013424141114f.pdf).
- For more information about available SIM900 methods please check godoc for this package.

Go explore!

## License

SIM900 package is MIT-Licensed
