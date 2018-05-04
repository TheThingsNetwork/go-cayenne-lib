// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayennelpp

import (
	"testing"

	. "github.com/smartystreets/assertions"
)

func TestEncode(t *testing.T) {
	a := New(t)

	// Uplink encoding
	{
		e := NewEncoder()

		e.AddDigitalInput(1, 255)
		a.So(e.Bytes()[0:3], ShouldResemble, []byte{1, DigitalInput, 255})

		e.AddDigitalOutput(2, 100)
		a.So(e.Bytes()[3:6], ShouldResemble, []byte{2, DigitalOutput, 100})

		e.AddAnalogInput(3, 54.5)
		a.So(e.Bytes()[6:10], ShouldResemble, []byte{3, AnalogInput, 21, 74})

		e.AddAnalogOutput(4, -54.5)
		a.So(e.Bytes()[10:14], ShouldResemble, []byte{4, AnalogOutput, 234, 182})

		e.AddLuminosity(5, 500)
		a.So(e.Bytes()[14:18], ShouldResemble, []byte{5, Luminosity, 1, 244})

		e.AddPresence(6, 50)
		a.So(e.Bytes()[18:21], ShouldResemble, []byte{6, Presence, 50})

		e.AddTemperature(7, -15.65)
		a.So(e.Bytes()[21:25], ShouldResemble, []byte{7, Temperature, 255, 100})

		e.AddRelativeHumidity(8, 49.65)
		a.So(e.Bytes()[25:28], ShouldResemble, []byte{8, RelativeHumidity, 99})

		e.AddAccelerometer(9, -0.424, 0.015, 1.666)
		a.So(e.Bytes()[28:36], ShouldResemble, []byte{9, Accelerometer, 254, 88, 0, 15, 6, 130})

		e.AddBarometricPressure(10, 1073.5)
		a.So(e.Bytes()[36:40], ShouldResemble, []byte{10, BarometricPressure, 41, 239})

		e.AddGyrometer(11, 3.55, 5.61, -4.10)
		a.So(e.Bytes()[40:48], ShouldResemble, []byte{11, Gyrometer, 1, 99, 2, 49, 254, 102})

		e.AddGPS(12, 52.3655, 4.8885, 21.54)
		a.So(e.Bytes()[48:59], ShouldResemble, []byte{12, GPS, 7, 253, 135, 0, 190, 245, 0, 8, 106})

		e.AddVoltage(13, 54.5)
		a.So(e.Bytes()[59:63], ShouldResemble, []byte{13, Voltage, 21, 74})

		e.AddCurrent(14, 2.4)
		a.So(e.Bytes()[63:67], ShouldResemble, []byte{14, Current, 0, 240})

		e.AddFrequency(15, 55)
		a.So(e.Bytes()[67:71], ShouldResemble, []byte{15, Frequency, 21, 124})

		e.AddEnergy(16, 789.1234)
		a.So(e.Bytes()[71:76], ShouldResemble, []byte{16, Energy, 1, 52, 64})

		e.AddEnergy(17, 1675.887)
		a.So(e.Bytes()[76:81], ShouldResemble, []byte{17, Energy, 2, 142, 164})
	}

	// Downlink encoding
	{
		e := NewEncoder()

		e.AddPort(1, 54.5)
		a.So(e.Bytes()[0:3], ShouldResemble, []byte{1, 21, 74})

		e.AddPort(2, -54.5)
		a.So(e.Bytes()[3:6], ShouldResemble, []byte{2, 234, 182})
	}
}
