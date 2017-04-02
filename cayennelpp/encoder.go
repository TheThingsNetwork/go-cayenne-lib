// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Encoder interface {
	Grow(n int)
	Bytes() []byte
	WriteTo(w io.Writer) (int64, error)
	AddDigitalInput(channel, value uint8)
	AddDigitalOutput(channel, value uint8)
	AddAnalogInput(channel uint8, value float32)
	AddAnalogOutput(channel uint8, value float32)
	AddLuminosity(channel uint8, value uint16)
	AddPresence(channel, value uint8)
	AddTemperature(channel uint8, celcius float32)
	AddRelativeHumidity(channel uint8, rh float32)
	AddAccelerometer(channel uint8, x, y, z float32)
	AddBarometricPressure(channel uint8, hpa float32)
	AddGyrometer(channel uint8, x, y, z float32)
	AddGPS(channel uint8, latitude, longitude, meters float32)
}

type encoder struct {
	buf *bytes.Buffer
}

func NewEncoder() Encoder {
	return &encoder{
		buf: new(bytes.Buffer),
	}
}

func (e *encoder) Grow(n int) {
	e.buf.Grow(n)
}

func (e *encoder) Bytes() []byte {
	return e.buf.Bytes()
}

func (e *encoder) Reset() {
	e.buf.Reset()
}

func (e *encoder) WriteTo(w io.Writer) (int64, error) {
	return e.buf.WriteTo(w)
}

func (e *encoder) AddDigitalInput(channel, value uint8) {
	e.buf.WriteByte(channel)
	e.buf.WriteByte(DigitalInput)
	e.buf.WriteByte(value)
}

func (e *encoder) AddDigitalOutput(channel, value uint8) {
	e.buf.WriteByte(channel)
	e.buf.WriteByte(DigitalOutput)
	e.buf.WriteByte(value)
}

func (e *encoder) AddAnalogInput(channel uint8, value float32) {
	val := uint16(value * 100)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(AnalogInput)
	binary.Write(e.buf, binary.BigEndian, val)
}

func (e *encoder) AddAnalogOutput(channel uint8, value float32) {
	val := uint16(value * 100)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(AnalogOutput)
	binary.Write(e.buf, binary.BigEndian, val)
}

func (e *encoder) AddLuminosity(channel uint8, value uint16) {
	e.buf.WriteByte(channel)
	e.buf.WriteByte(Luminosity)
	binary.Write(e.buf, binary.BigEndian, value)
}

func (e *encoder) AddPresence(channel, value uint8) {
	e.buf.WriteByte(channel)
	e.buf.WriteByte(Presence)
	e.buf.WriteByte(value)
}

func (e *encoder) AddTemperature(channel uint8, celcius float32) {
	val := uint16(celcius * 10)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(Temperature)
	binary.Write(e.buf, binary.BigEndian, val)
}

func (e *encoder) AddRelativeHumidity(channel uint8, rh float32) {
	e.buf.WriteByte(channel)
	e.buf.WriteByte(RelativeHumidity)
	e.buf.WriteByte(uint8(rh * 2))
}

func (e *encoder) AddAccelerometer(channel uint8, x, y, z float32) {
	valX := uint16(x * 1000)
	valY := uint16(y * 1000)
	valZ := uint16(z * 1000)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(Accelerometer)
	binary.Write(e.buf, binary.BigEndian, valX)
	binary.Write(e.buf, binary.BigEndian, valY)
	binary.Write(e.buf, binary.BigEndian, valZ)
}

func (e *encoder) AddBarometricPressure(channel uint8, hpa float32) {
	val := uint16(hpa * 10)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(BarometricPressure)
	binary.Write(e.buf, binary.BigEndian, val)
}

func (e *encoder) AddGyrometer(channel uint8, x, y, z float32) {
	valX := uint16(x * 100)
	valY := uint16(y * 100)
	valZ := uint16(z * 100)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(Gyrometer)
	binary.Write(e.buf, binary.BigEndian, valX)
	binary.Write(e.buf, binary.BigEndian, valY)
	binary.Write(e.buf, binary.BigEndian, valZ)
}

func (e *encoder) AddGPS(channel uint8, latitude, longitude, meters float32) {
	valLat := uint32(latitude * 10000)
	valLon := uint32(longitude * 10000)
	valAlt := uint32(meters * 100)
	e.buf.WriteByte(channel)
	e.buf.WriteByte(GPS)
	e.buf.WriteByte(uint8(valLat >> 16))
	e.buf.WriteByte(uint8(valLat >> 8))
	e.buf.WriteByte(uint8(valLat))
	e.buf.WriteByte(uint8(valLon >> 16))
	e.buf.WriteByte(uint8(valLon >> 8))
	e.buf.WriteByte(uint8(valLon))
	e.buf.WriteByte(uint8(valAlt >> 16))
	e.buf.WriteByte(uint8(valAlt >> 8))
	e.buf.WriteByte(uint8(valAlt))
}
