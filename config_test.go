package sqm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqm "github.com/gaker/unihedron-sqm"
)

var _ = Describe("Config", func() {

	It("should return an empty string when host and port are not set", func() {
		c := sqm.Tcp{}
		Expect(c.Addr()).To(Equal(""))
	})

	It("should return an empty string when host and port are not set", func() {
		c := sqm.Tcp{
			Host: PointIt("localhost"),
			Port: PointIt("8080"),
		}
		Expect(c.Addr()).To(Equal("localhost:8080"))
	})

	Context("Validate", func() {
		It("should require a valid config", func() {
			c := &sqm.Config{}
			err := c.Validate()
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).
				To(Equal("tcp configuration required"))
		})

		It("should require fields for TCP connection", func() {
			configs := map[sqm.Config]bool{
				{
					Tcp: &sqm.Tcp{},
				}: false,
				{
					Tcp: &sqm.Tcp{
						Host: PointIt("localhost"),
					},
				}: false,
				{
					Tcp: &sqm.Tcp{
						Port: PointIt("12001"),
					},
				}: false,
				{
					Tcp: &sqm.Tcp{
						Host: PointIt("localhost"),
						Port: PointIt("12001"),
					},
				}: true,
			}

			for c, expected := range configs {
				err := c.Validate()
				if expected {
					Expect(err).To(BeNil())
				} else {
					Expect(err).ToNot(BeNil())
				}
			}
		})
	})
})
