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
			Expect(err.Error()).To(Equal("tcp configuration required"))
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

	})
})
