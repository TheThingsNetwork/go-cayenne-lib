// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrInvalidChannel = errors.New("cayennelpp: unknown type")

type UplinkTarget interface {
	DigitalInput(channel, value uint8)
	DigitalOutput(channel, value uint8)
	AnalogInput(channel uint8, value float32)
	AnalogOutput(channel uint8, value float32)
	Luminosity(channel uint8, value uint16)
	Presence(channel, value uint8)
	Temperature(channel uint8, celcius float32)
	RelativeHumidity(channel uint8, rh float32)
	Accelerometer(channel uint8, x, y, z float32)
	BarometricPressure(channel uint8, hpa float32)
	Gyrometer(channel uint8, x, y, z float32)
	GPS(channel uint8, latitude, longitude, altitude float32)
}

type DownlinkTarget interface {
	Port(channel uint8, value float32)
}

type Decoder interface {
	DecodeUplink(target UplinkTarget) error
	DecodeDownlink(target DownlinkTarget) error
}

type decoder struct {
	r io.Reader
}

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
			err = ErrInvalidChannel
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
		var val int16
		if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
			return err
		}
		target.Port(buf[0], float32(val)/100)
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
	target.AnalogInput(channel, float32(val)/100)
	return nil
}

func (d *decoder) decodeAnalogOutput(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.AnalogOutput(channel, float32(val)/100)
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
	target.Temperature(channel, float32(val)/10)
	return nil
}

func (d *decoder) decodeRelativeHumidity(channel uint8, target UplinkTarget) error {
	var val int8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.RelativeHumidity(channel, float32(val)/2)
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
	target.Accelerometer(channel, float32(valX)/1000, float32(valY)/1000, float32(valZ)/1000)
	return nil
}

func (d *decoder) decodeBarometricPressure(channel uint8, target UplinkTarget) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.BarometricPressure(channel, float32(val)/10)
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
	target.Gyrometer(channel, float32(valX)/100, float32(valY)/100, float32(valZ)/100)
	return nil
}

func (d *decoder) decodeGPS(channel uint8, target UplinkTarget) error {
	buf := make([]byte, 9)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return err
	}
	var latitude, longitude, altitude int32
	latitude = int32(buf[0])<<16 | int32(buf[1])<<8 | int32(buf[2])
	longitude = int32(buf[3])<<16 | int32(buf[4])<<8 | int32(buf[5])
	altitude = int32(buf[6])<<16 | int32(buf[7])<<8 | int32(buf[8])
	target.GPS(channel, float32(latitude)/10000, float32(longitude)/10000, float32(altitude)/100)
	return nil
}
