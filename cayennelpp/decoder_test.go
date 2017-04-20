// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

import (
	"bytes"
	"testing"

	"io"

	. "github.com/smartystreets/assertions"
)

type target struct {
	values map[uint8]interface{}
}

func (t *target) Port(channel uint8, value float32) {
	t.values[channel] = value
}

func (t *target) DigitalInput(channel, value uint8) {
	t.values[channel] = value
}

func (t *target) DigitalOutput(channel, value uint8) {
	t.values[channel] = value
}

func (t *target) AnalogInput(channel uint8, value float32) {
	t.values[channel] = value
}

func (t *target) AnalogOutput(channel uint8, value float32) {
	t.values[channel] = value
}

func (t *target) Luminosity(channel uint8, value uint16) {
	t.values[channel] = value
}

func (t *target) Presence(channel, value uint8) {
	t.values[channel] = value
}

func (t *target) Temperature(channel uint8, celcius float32) {
	t.values[channel] = celcius
}

func (t *target) RelativeHumidity(channel uint8, rh float32) {
	t.values[channel] = rh
}

func (t *target) Accelerometer(channel uint8, x, y, z float32) {
	t.values[channel] = []float32{x, y, z}
}

func (t *target) BarometricPressure(channel uint8, hpa float32) {
	t.values[channel] = hpa
}

func (t *target) Gyrometer(channel uint8, x, y, z float32) {
	t.values[channel] = []float32{x, y, z}
}

func (t *target) GPS(channel uint8, latitude, longitude, altitude float32) {
	t.values[channel] = []float32{latitude, longitude, altitude}
}

func TestDecode(t *testing.T) {
	a := New(t)

	// Happy flow: uplink
	{
		buf := []byte{
			1, DigitalInput, 255,
			2, DigitalOutput, 100,
			3, AnalogInput, 21, 74,
			4, AnalogOutput, 234, 182,
			5, Luminosity, 1, 244,
			6, Presence, 50,
			7, Temperature, 255, 100,
			8, RelativeHumidity, 99,
			9, Accelerometer, 254, 88, 0, 15, 6, 130,
			10, BarometricPressure, 41, 239,
			11, Gyrometer, 1, 99, 2, 49, 254, 102,
			12, GPS, 7, 253, 135, 0, 190, 245, 0, 8, 106,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))
		target := &target{make(map[uint8]interface{})}

		err := decoder.DecodeUplink(target)
		a.So(err, ShouldBeNil)
		a.So(target.values[1], ShouldEqual, 255)
		a.So(target.values[2], ShouldEqual, 100)
		a.So(target.values[3], ShouldEqual, 54.5)
		a.So(target.values[4], ShouldEqual, -54.5)
		a.So(target.values[5], ShouldEqual, 500)
		a.So(target.values[6], ShouldEqual, 50)
		a.So(target.values[7], ShouldEqual, -15.6)
		a.So(target.values[8], ShouldEqual, 49.5)
		a.So(target.values[9], ShouldResemble, []float32{-0.424, 0.015, 1.666})
		a.So(target.values[10], ShouldEqual, 1073.5)
		a.So(target.values[11], ShouldResemble, []float32{3.55, 5.61, -4.10})
		a.So(target.values[12], ShouldResemble, []float32{52.3655, 4.8885, 21.54})
	}

	// Happy flow: downlink
	{
		buf := []byte{
			1, 0, 100,
			2, 234, 182,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))
		target := &target{make(map[uint8]interface{})}

		err := decoder.DecodeDownlink(target)
		a.So(err, ShouldBeNil)
		a.So(target.values[1], ShouldEqual, 1)
		a.So(target.values[2], ShouldEqual, -54.5)
	}

	// Invalid data type
	{
		buf := []byte{
			1, 255, 255,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))
		target := &target{make(map[uint8]interface{})}

		err := decoder.DecodeUplink(target)
		a.So(err, ShouldEqual, ErrInvalidChannel)
	}

	// Not enough data: uplink
	{
		buf := []byte{
			12, GPS, 7, 253, 135, 0, 190,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))
		target := &target{make(map[uint8]interface{})}

		err := decoder.DecodeUplink(target)
		a.So(err, ShouldEqual, io.ErrUnexpectedEOF)
	}

	// Not enough data: downlink
	{
		buf := []byte{
			12, 1,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))
		target := &target{make(map[uint8]interface{})}

		err := decoder.DecodeDownlink(target)
		a.So(err, ShouldEqual, io.ErrUnexpectedEOF)
	}
}
