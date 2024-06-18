package sqm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqm "github.com/gaker/unihedron-sqm"
)

var _ = Describe("Commands", func() {
	When("Trying to parse a reading", func() {
		Context("when the input is invalid", func() {
			It("should require the right length", func() {
				reading := &sqm.Reading{}
				err := reading.Parse([]byte("x,-09.42m,0000005915Hz,0000000000c,0000000.000s\r\n"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid length"))
			})

			It("should require a valid first character", func() {
				reading := &sqm.Reading{}
				err := reading.Parse([]byte("x,-09.42m,0000005915Hz,0000000000c,0000000.000s, 027.0C\r\n"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid reading. 'x' start position not valid"))
			})

			It("should return zero values for total nonsense", func() {
				reading := &sqm.Reading{}
				err := reading.Parse([]byte("r,-0S.42m,00XXX05915Hz,000000SSSSc,000SSSS.000s, ABC.0C\r\n"))

				Expect(err).To(BeNil())
				Expect(reading).ToNot(BeNil())

				Expect(reading.Averaged).To(BeTrue())
				Expect(reading.Reading).To(Equal(0.0))
				Expect(reading.Frequency).To(Equal(0))
				Expect(reading.Counts).To(Equal(0))
				Expect(reading.Millis).To(Equal(0.0))
				Expect(reading.Temp).To(Equal(0.0))
			})
		})

		Context("the input is valid", func() {
			It("should parse the averaged reading", func() {
				reading := &sqm.Reading{}
				err := reading.Parse([]byte("r,-09.42m,0000005915Hz,1000100000c,0006000.001s, 027.1C\r\n"))

				Expect(err).To(BeNil())

				Expect(reading.Averaged).To(BeTrue())
				Expect(reading.Reading).To(Equal(-9.4))
				Expect(reading.Frequency).To(Equal(5915))
				Expect(reading.Counts).To(Equal(1000100000))
				Expect(reading.Millis).To(Equal(6000.001))
				Expect(reading.Temp).To(Equal(27.1))
			})

			It("should parse the unaveraged reading", func() {
				reading := &sqm.Reading{}
				err := reading.Parse([]byte("u,-09.42m,0000005915Hz,1000100000c,0006000.001s, 027.1C\r\n"))

				Expect(err).To(BeNil())

				Expect(reading.Averaged).To(BeFalse())
				Expect(reading.Reading).To(Equal(-9.4))
				Expect(reading.Frequency).To(Equal(5915))
				Expect(reading.Counts).To(Equal(1000100000))
				Expect(reading.Millis).To(Equal(6000.001))
				Expect(reading.Temp).To(Equal(27.1))
			})
		})
	})

	When("parsing UnitInfo", func() {
		Context("it should validate the incoming packet", func() {
			It("should be of the right length", func() {
				info := &sqm.UnitInfo{}

				err := info.Parse([]byte("foobar"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid length"))
			})

			It("should be of the correct start", func() {
				info := &sqm.UnitInfo{}

				err := info.Parse([]byte("r,00000002,00000003,00000001,00000413\r\n"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid reading. 'r' start position not valid"))
			})
		})

		Context("when everything is good", func() {
			It("should parse a response", func() {
				info := &sqm.UnitInfo{}

				err := info.Parse([]byte("i,00000002,00000003,00000001,00000413\r\n"))
				Expect(err).To(BeNil())

				Expect(info.Protocol).To(Equal(2))
				Expect(info.Model).To(Equal(3))
				Expect(info.Feature).To(Equal(1))
				Expect(info.SerialNumber).To(Equal(413))
			})
		})
	})

	When("parsing CalibrationInfo", func() {
		Context("it should validate the incoming packet", func() {
			It("should be of the right length", func() {
				info := &sqm.CalibrationInfo{}

				err := info.Parse([]byte("foobar"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid length"))
			})

			It("should be of the correct start", func() {
				info := &sqm.CalibrationInfo{}

				err := info.Parse([]byte("r,00000017.60m,0000000.000s, 039.4C,00000008.71m, 039.4C\r\n"))
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid reading. 'r' start position not valid"))
			})
		})

		Context("When things are good", func() {
			It("should parse the result", func() {
				info := &sqm.CalibrationInfo{}
				err := info.Parse([]byte("c,00000017.60m,0000000.000s, 039.4C,00000008.71m, 039.4C\r\n"))
				Expect(err).To(BeNil())

				Expect(info.LightCalibrationOffset).To(Equal(17.60))
				Expect(info.LightCalibrationTemp).To(Equal(0.0))
				Expect(info.DarkCalibrationTimePeriod).To(Equal(39.4))
				Expect(info.DarkCalibrationTemp).To(Equal(8.71))
				Expect(info.Offset).To(Equal(39.4))
			})

			It("should handle negative temp values", func() {
				info := &sqm.CalibrationInfo{}
				err := info.Parse([]byte("c,00000017.60m,0000000.000s,-039.4C,00000008.71m,-039.4C\r\n"))
				Expect(err).To(BeNil())

				Expect(info.LightCalibrationOffset).To(Equal(17.60))
				Expect(info.LightCalibrationTemp).To(Equal(0.0))
				Expect(info.DarkCalibrationTimePeriod).To(Equal(-39.4))
				Expect(info.DarkCalibrationTemp).To(Equal(8.71))
				Expect(info.Offset).To(Equal(-39.4))
			})
		})
	})
})
