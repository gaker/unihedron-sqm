package sqm

import (
	"bufio"
	"net"
)

func NewNetwork(cfg *Config) *Network {
	return &Network{
		cfg: cfg,
	}
}

// Network supports the SQM-LE
type Network struct {
	cfg  *Config
	conn *net.TCPConn
}

// Dial is needed to implement the Monitor interface
func (n *Network) Dial() error {
	var err error
	addr, _ := net.ResolveTCPAddr("tcp", n.cfg.Tcp.Addr())
	n.conn, err = net.DialTCP("tcp", nil, addr)
	return err
}

// Send is needed to implement the Monitor interface
func (n *Network) Send(command Command) error {

	issued := append([]byte(command), []byte{'\r', '\n'}...)
	_, err := n.conn.Write(issued)
	return err
}

// Read is needed to implement the Monitor interface
func (n *Network) Read(item Readable) error {

	defer n.conn.Close()

	reader := bufio.NewReader(n.conn)
	msg, err := reader.ReadBytes('\n')

	if err != nil {
		return err
	}

	return item.Parse(msg)
}
