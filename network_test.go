package sqm_test

import (
	"bufio"
	"io"
	"net"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqm "github.com/gaker/unihedron-sqm"
)

// StartServer starts the TCP server on the given port.
func startTestServer() (net.Listener, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	Expect(err).To(BeNil())

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			handleConnection(conn)
		}
	}()

	return listener, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			continue
		}
		if err != nil {
			continue
		}
		message = strings.TrimRight(message, "\r\n")

		var output []byte

		switch message {
		case "rx":
			output = []byte("r,-09.42m,0000005915Hz,0000000000c,0000000.000s, 027.0C\r\n")
		}
		conn.Write(output)
	}
}

var _ = Describe("Network", func() {
	var (
		lc  net.Listener
		cfg *sqm.Config
	)

	var _ = BeforeEach(func() {
		var err error
		lc, err = startTestServer()
		Expect(err).To(BeNil())

		parts := strings.Split(lc.Addr().String(), ":")

		cfg = &sqm.Config{
			Tcp: &sqm.Tcp{
				Host: &parts[0],
				Port: &parts[1],
			},
		}
	})

	var _ = AfterEach(func() {
		lc.Close()
	})

	When("it is dialing", func() {
		It("should handle a failure to connect", func() {

			cfg.Tcp.Port = PointIt("0")
			n := sqm.NewNetwork(cfg)

			err := n.Dial()
			Expect(err).ToNot(BeNil())
		})

		It("should create a connection", func() {
			n := sqm.NewNetwork(cfg)
			Expect(n.Dial()).To(BeNil())
		})
	})

	When("sending commands", func() {
		It("should get a response", func() {
			n := sqm.NewNetwork(cfg)
			Expect(n.Dial()).To(BeNil())
			Expect(n.Send(sqm.CmdAveragedReading)).To(BeNil())

			reading := &sqm.Reading{}
			Expect(n.Read(reading)).To(BeNil())
		})
	})
})
