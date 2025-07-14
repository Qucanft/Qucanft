package time

import (
	"testing"
	"time"
)

func TestToJulianDate(t *testing.T) {
	// Test J2000.0 epoch (January 1, 2000, 12:00 UTC)
	j2000 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := ToJulianDate(j2000)
	
	expected := JulianDate(2451545.0)
	if jd != expected {
		t.Errorf("Expected J2000.0 = %f, got %f", expected, jd)
	}
}

func TestJulianDateConversion(t *testing.T) {
	// Test round-trip conversion
	originalTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := ToJulianDate(originalTime)
	convertedTime := jd.ToTime()
	
	// Allow for small rounding errors
	if convertedTime.Unix() != originalTime.Unix() {
		t.Errorf("Round-trip conversion failed: %v != %v", originalTime, convertedTime)
	}
}

func TestDaysSinceJ2000(t *testing.T) {
	// Test one day after J2000.0
	nextDay := time.Date(2000, 1, 2, 12, 0, 0, 0, time.UTC)
	jd := ToJulianDate(nextDay)
	days := jd.DaysSinceJ2000()
	
	if days != 1.0 {
		t.Errorf("Expected 1 day since J2000.0, got %f", days)
	}
}

func TestCenturiesSinceJ2000(t *testing.T) {
	// Test 100 years (1 century) after J2000.0
	century := time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := ToJulianDate(century)
	centuries := jd.CenturiesSinceJ2000()
	
	// approximately 1.0
	if centuries < 0.99 || centuries > 1.01 {
		t.Errorf("Expected ~1.0 centuries since J2000.0, got %f", centuries)
	}
}