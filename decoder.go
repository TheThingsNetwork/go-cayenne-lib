// Copyright © 2021 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

import (
	"encoding/binary"
	"errors"
	"io"
)

// ErrInvalidChannelType indicates that the channel type is invalid.
var ErrInvalidChannelType = errors.New("cayennelpp: unknown type")

// UplinkTarget represents a target that processes decoded uplink values.
type UplinkTarget interface {
	DigitalInput(channel, value uint8)
	DigitalOutput(channel, value uint8)
	AnalogInput(channel uint8, value float64)
	AnalogOutput(channel uint8, value float64)
	Luminosity(channel uint8, value uint16)
	Presence(channel, value uint8)
	Temperature(channel uint8, celcius float64)
	RelativeHumidity(channel uint8, rh float64)
	Accelerometer(channel uint8, x, y, z float64)
	BarometricPressure(channel uint8, hpa float64)
	Gyrometer(channel uint8, x, y, z float64)
	GPS(channel uint8, latitude, longitude, altitude float64)
}

// DownlinkTarget represents a target that processes decoded downlink values.
type DownlinkTarget interface {
	Port(channel uint8, value float64)
}

// Decoder decodes CayenneLPP encoded values.
type Decoder interface {
	DecodeUplink(target UplinkTarget) error
	DecodeDownlink(target DownlinkTarget) error
}

type decoder struct {
	r io.Reader
}

// NewDecoder instantiates a CayenneLPP decoder.
func NewDecoder(r io.Reader) Decoder {
	return &decoder{r}
}

func (d *decoder) DecodeUplink(target UplinkTarget) error {
	buf := make([]byte, 2)
	for {
		_, err := io.ReadFull(d.r, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch buf[1] {
		case DigitalInput:
			err = d.decodeDigitalInput(buf[0], target)
		case DigitalOutput:
			err = d.decodeDigitalOutput(buf[0], target)
		case AnalogInput:
			err = d.decodeAnalogInput(buf[0], target)
		case AnalogOutput:
			err = d.decodeAnalogOutput(buf[0], target)
		case Luminosity:
			err = d.decodeLuminosity(buf[0], target)
		case Presence:
			err = d.decodePresence(buf[0], target)
		case Temperature:
			err = d.decodeTemperature(buf[0], target)
		case RelativeHumidity:
			err = d.decodeRelativeHumidity(buf[0], target)
		case Accelerometer:
			err = d.decodeAccelerometer(buf[0], target)
		case BarometricPressure:
			err = d.decodeBarometricPressure(buf[0], target)
		case Gyrometer:
			err = d.decodeGyrometer(buf[0], target)
		case GPS:
			err = d.decodeGPS(buf[0], target)
		default:
			err = ErrInvalidChannelType
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) DecodeDownlink(target DownlinkTarget) error {
	buf := make([]byte, 1)
	for {
		_, err := io.ReadFull(d.r, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if buf[0] == 0xFF {
			break
		}
		var val int16
		if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
			return err
		}
		target.Port(buf[0], float64(val)/100)
	}
	return nil
}

func (d *decoder) decodeDigitalInput(channel uint8, target UplinkTarget) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.DigitalInput(channel, val)
	return nil
}

func (d *decoder) decodeDigitalOutput(channel uint8, target UplinkTarget) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.DigitalOutput(channel, val)
	return nil
}

func (d *decoder) decodeAnalogInput(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.AnalogInput(channel, float64(val)/100)
	return nil
}

func (d *decoder) decodeAnalogOutput(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.AnalogOutput(channel, float64(val)/100)
	return nil
}

func (d *decoder) decodeLuminosity(channel uint8, target UplinkTarget) error {
	var val uint16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.Luminosity(channel, val)
	return nil
}

func (d *decoder) decodePresence(channel uint8, target UplinkTarget) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.Presence(channel, val)
	return nil
}

func (d *decoder) decodeTemperature(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.Temperature(channel, float64(val)/10)
	return nil
}

func (d *decoder) decodeRelativeHumidity(channel uint8, target UplinkTarget) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.RelativeHumidity(channel, float64(val)/2)
	return nil
}

func (d *decoder) decodeAccelerometer(channel uint8, target UplinkTarget) error {
	var valX, valY, valZ int16
	if err := binary.Read(d.r, binary.BigEndian, &valX); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valY); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valZ); err != nil {
		return err
	}
	target.Accelerometer(channel, float64(valX)/1000, float64(valY)/1000, float64(valZ)/1000)
	return nil
}

func (d *decoder) decodeBarometricPressure(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.BarometricPressure(channel, float64(val)/10)
	return nil
}

func (d *decoder) decodeGyrometer(channel uint8, target UplinkTarget) error {
	var valX, valY, valZ int16
	if err := binary.Read(d.r, binary.BigEndian, &valX); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valY); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valZ); err != nil {
		return err
	}
	target.Gyrometer(channel, float64(valX)/100, float64(valY)/100, float64(valZ)/100)
	return nil
}

func (d *decoder) decodeGPS(channel uint8, target UplinkTarget) error {
	buf := make([]byte, 9)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return err
	}
	latitude := make([]byte, 4)
	copy(latitude, buf[0:3])
	longitude := make([]byte, 4)
	copy(longitude, buf[3:6])
	altitude := make([]byte, 4)
	copy(altitude, buf[6:9])
	target.GPS(channel,
		float64(int32(binary.BigEndian.Uint32(latitude))>>8)/10000,
		float64(int32(binary.BigEndian.Uint32(longitude))>>8)/10000,
		float64(int32(binary.BigEndian.Uint32(altitude))>>8)/100,
	)
	return nil
}
