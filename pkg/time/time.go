// Package time provides utilities for time and date calculations in astrological contexts.
package time

import (
	"fmt"
	"math"
	"time"
)

// JulianDay represents a Julian Day Number (JDN)
type JulianDay float64

// Constants for time calculations
const (
	// J2000 is the Julian Day Number for J2000.0 epoch (January 1, 2000 12:00:00 UTC)
	J2000 JulianDay = 2451545.0
	
	// SecondsPerDay is the number of seconds in a day
	SecondsPerDay = 86400.0
	
	// DaysPerCentury is the number of days in a Julian century
	DaysPerCentury = 36525.0
)

// TimeConverter handles conversions between different time systems
type TimeConverter struct{}

// NewTimeConverter creates a new TimeConverter instance
func NewTimeConverter() *TimeConverter {
	return &TimeConverter{}
}

// ToJulianDay converts a standard time.Time to Julian Day Number
func (tc *TimeConverter) ToJulianDay(t time.Time) JulianDay {
	// Convert to UTC
	utc := t.UTC()
	
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	hour := utc.Hour()
	minute := utc.Minute()
	second := utc.Second()
	
	// Decimal day
	decimalDay := float64(day) + float64(hour)/24.0 + float64(minute)/1440.0 + float64(second)/86400.0
	
	// Adjust for January and February
	if month <= 2 {
		year--
		month += 12
	}
	
	// Calculate Julian Day Number
	a := year / 100
	b := 2 - a + a/4
	
	jd := math.Floor(365.25*float64(year+4716)) + math.Floor(30.6001*float64(month+1)) + decimalDay + float64(b) - 1524.5
	
	return JulianDay(jd)
}

// FromJulianDay converts a Julian Day Number to time.Time
func (tc *TimeConverter) FromJulianDay(jd JulianDay) time.Time {
	// Algorithm from Meeus, "Astronomical Algorithms"
	jd += 0.5
	z := math.Floor(float64(jd))
	f := float64(jd) - z
	
	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}
	
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)
	
	day := b - d - math.Floor(30.6001*e) + f
	
	var month float64
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	
	var year float64
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}
	
	// Extract time components
	dayInt := int(day)
	dayFrac := day - float64(dayInt)
	
	hours := dayFrac * 24
	hoursInt := int(hours)
	hoursFrac := hours - float64(hoursInt)
	
	minutes := hoursFrac * 60
	minutesInt := int(minutes)
	minutesFrac := minutes - float64(minutesInt)
	
	seconds := minutesFrac * 60
	secondsInt := int(seconds)
	nanoseconds := int((seconds - float64(secondsInt)) * 1e9)
	
	return time.Date(int(year), time.Month(month), dayInt, hoursInt, minutesInt, secondsInt, nanoseconds, time.UTC)
}

// JulianCenturies returns the number of Julian centuries since J2000.0
func (tc *TimeConverter) JulianCenturies(jd JulianDay) float64 {
	return float64(jd-J2000) / DaysPerCentury
}

// SiderealTime calculates the Greenwich Mean Sidereal Time (GMST)
func (tc *TimeConverter) SiderealTime(jd JulianDay) float64 {
	t := tc.JulianCenturies(jd)
	
	// GMST at 0h UT
	gmst0 := 280.46061837 + 360.98564736629*float64(jd-J2000) + 0.000387933*t*t - t*t*t/38710000.0
	
	// Normalize to 0-360 degrees
	gmst0 = math.Mod(gmst0, 360.0)
	if gmst0 < 0 {
		gmst0 += 360.0
	}
	
	return gmst0
}

// LocalSiderealTime calculates the Local Sidereal Time for a given longitude
func (tc *TimeConverter) LocalSiderealTime(jd JulianDay, longitude float64) float64 {
	gmst := tc.SiderealTime(jd)
	lst := gmst + longitude
	
	// Normalize to 0-360 degrees
	lst = math.Mod(lst, 360.0)
	if lst < 0 {
		lst += 360.0
	}
	
	return lst
}

// DeltaT returns the difference between Terrestrial Time and Universal Time
// This is a simplified approximation for the period 1620-2100
func (tc *TimeConverter) DeltaT(year int) float64 {
	y := float64(year)
	
	switch {
	case year < 1620:
		t := (y - 1600) / 100
		return 120 - 0.9808*t - 2.532*t*t + 0.1427*t*t*t - 0.0288*t*t*t*t
	case year < 1900:
		t := (y - 1900) / 100
		return -2.79 + 149.4119*t - 598.939*t*t + 6196.6*t*t*t - 19700*t*t*t*t
	case year < 2000:
		t := y - 2000
		return 63.86 + 0.3345*t - 0.060374*t*t + 0.0017275*t*t*t + 0.000651814*t*t*t*t + 0.00002373599*t*t*t*t*t
	case year <= 2100:
		t := y - 2000
		return 62.92 + 0.32217*t + 0.005589*t*t
	default:
		// Extrapolation for years > 2100
		t := (y - 2000) / 100
		return -20 + 32*t*t
	}
}

// String implements the Stringer interface for JulianDay
func (jd JulianDay) String() string {
	return fmt.Sprintf("JD %.6f", float64(jd))
}

// ToTime converts JulianDay to time.Time using the default TimeConverter
func (jd JulianDay) ToTime() time.Time {
	tc := NewTimeConverter()
	return tc.FromJulianDay(jd)
}

// Add adds days to a JulianDay
func (jd JulianDay) Add(days float64) JulianDay {
	return JulianDay(float64(jd) + days)
}

// Sub subtracts another JulianDay and returns the difference in days
func (jd JulianDay) Sub(other JulianDay) float64 {
	return float64(jd) - float64(other)
}