package sqm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqm "github.com/gaker/unihedron-sqm"
)

func PointIt[T any](input T) *T {
	return &input
}

var _ = Describe("SQM", func() {
	Context("New", func() {
		It("should require a valid config", func() {
			cfg := &sqm.Config{}

			mon, err := sqm.New(cfg)
			Expect(mon).To(BeNil())
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("[SQM] one of http or serial should be set in the config"))
		})

		It("should create a network interface", func() {
			cfg := &sqm.Config{
				Tcp: &sqm.Tcp{
					Host: PointIt("0.0.0.0"),
					Port: PointIt("10001"),
				},
			}

			mon, err := sqm.New(cfg)
			Expect(err).To(BeNil())
			_, ok := mon.(*sqm.Network)

			Expect(ok).To(BeTrue())
		})

		It("should create a serial interface", func() {
			cfg := &sqm.Config{
				Serial: &sqm.Serial{
					Port: PointIt("/dev/tty/usb1.1"),
				},
			}

			mon, err := sqm.New(cfg)
			Expect(err).To(BeNil())
			_, ok := mon.(*sqm.Usb)

			Expect(ok).To(BeTrue())
		})
	})

	Context("Network", func() {
		It("should do something", func() {
			Expect(true).To(BeTrue())
		})
	})
})
