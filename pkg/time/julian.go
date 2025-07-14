package time

import (
	"time"
)

// JulianDate represents a Julian Day Number with fractional part
type JulianDate float64

// ToJulianDate converts a time.Time to Julian Date
func ToJulianDate(t time.Time) JulianDate {
	// Convert to UTC
	utc := t.UTC()
	
	// Get components
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	hour := utc.Hour()
	minute := utc.Minute()
	second := utc.Second()
	nanosecond := utc.Nanosecond()
	
	// Convert to Julian Date using standard algorithm
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	
	jdn := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
	
	// Add fractional part for time of day
	fractionalDay := (float64(hour) + float64(minute)/60.0 + 
		(float64(second) + float64(nanosecond)/1e9)/3600.0) / 24.0
	
	return JulianDate(float64(jdn) + fractionalDay - 0.5)
}

// ToTime converts a Julian Date to time.Time
func (jd JulianDate) ToTime() time.Time {
	// Extract integer and fractional parts
	jdn := int(jd + 0.5)
	fraction := float64(jd) - float64(jdn) + 0.5
	
	// Convert Julian Day Number to calendar date
	a := jdn + 32044
	b := (4*a + 3) / 146097
	c := a - (146097*b)/4
	d := (4*c + 3) / 1461
	e := c - (1461*d)/4
	m := (5*e + 2) / 153
	
	day := e - (153*m+2)/5 + 1
	month := m + 3 - 12*(m/10)
	year := 100*b + d - 4800 + m/10
	
	// Convert fractional part to time
	totalSeconds := fraction * 24 * 3600
	hours := int(totalSeconds / 3600)
	minutes := int((totalSeconds - float64(hours)*3600) / 60)
	seconds := int(totalSeconds - float64(hours)*3600 - float64(minutes)*60)
	nanoseconds := int((totalSeconds - float64(int(totalSeconds))) * 1e9)
	
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, nanoseconds, time.UTC)
}

// J2000 returns the J2000.0 epoch (January 1, 2000, 12:00 TT)
func J2000() JulianDate {
	return JulianDate(2451545.0)
}

// DaysSinceJ2000 returns the number of days since J2000.0 epoch
func (jd JulianDate) DaysSinceJ2000() float64 {
	return float64(jd - J2000())
}

// CenturiesSinceJ2000 returns the number of centuries since J2000.0 epoch
func (jd JulianDate) CenturiesSinceJ2000() float64 {
	return jd.DaysSinceJ2000() / 36525.0
}