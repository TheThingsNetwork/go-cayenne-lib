// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

const (
	// DigitalInput (1 byte)
	DigitalInput = 0
	// DigitalOutput (1 byte)
	DigitalOutput = 1
	// AnalogInput (2 bytes, 0.01 signed)
	AnalogInput = 2
	// AnalogOutput (2 bytes, 0.01 signed)
	AnalogOutput = 3
	// Luminosity (2 bytes, 1 lux unsigned)
	Luminosity = 101
	// Presence (1 byte, 1)
	Presence = 102
	// Temperature (2 bytes, 0.1°C signed)
	Temperature = 103
	// RelativeHumidity (1 byte, 0.5% unsigned )
	RelativeHumidity = 104
	// Accelerometer 2 bytes per axis, 0.001G
	Accelerometer = 113
	// BarometricPressure (2 bytes 0.1 hPa Unsigned)
	BarometricPressure = 115
	// Gyrometer (2 bytes per axis, 0.01 °/s)
	Gyrometer = 134
	// GPS (3 byte lon/lat 0.0001 °, 3 bytes alt 0.01 meter)
	GPS = 136
)
