package sim900

const (
	CMD_AT     string = "AT"
	CMD_OK     string = "(^OK$)"
	CMD_ERROR  string = "(^ERROR$)"
	CMD_CMGF   string = "AT+CMGF=1"
	CMD_CTRL_Z string = "\x1A"
	// CMGS
	CMD_CMGS           string = "AT+CMGS=\"%s\""
	CMD_CMGS_RX_REGEXP string = "(^[+]CMGS[:] [0-9]+$)"
	// CMGR
	CMD_CMGR        string = "AT+CMGR=%s"
	CMD_CMGR_REGEXP string = "(^[+]CMGR[:] .*)"
	CMD_CMGR_RX     string = "+CMGR: "
	// CMTI - CMTI
	CMD_CMTI_REGEXP string = "(^[+]CMTI[:] \"SM\",[0-9]+$)"
	CMD_CMTI_RX     string = "+CMTI: \"SM\","
)
