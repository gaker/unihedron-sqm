package sqm

type Conn interface {
	SendCommand(cmd Command) error
	ReadResponse(in Readable) error
}
