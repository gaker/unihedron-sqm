package sqm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqm "github.com/gaker/unihedron-sqm"
)

var _ = Describe("Config", func() {
	Context("Validate", func() {
		It("should require a valid config", func() {
			c := &sqm.Config{}
			err := c.Validate()
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).
				To(Equal("[SQM] one of http or serial should be set in the config"))
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

		It("should require fields for Serial connection", func() {
			c := sqm.Config{
				Serial: &sqm.Serial{
					// Port: PointIt("1200"),
				},
			}
			err := c.Validate()
			Expect(err).ToNot(BeNil())

			c.Serial.Port = PointIt("1200")

			err = c.Validate()
			Expect(err).To(BeNil())
		})
	})
})
