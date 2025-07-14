package time

import (
	"testing"
	"time"
	"math"
)

func TestTimeConverter(t *testing.T) {
	tc := NewTimeConverter()
	
	// Test J2000.0 epoch
	j2000 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := tc.ToJulianDay(j2000)
	
	if math.Abs(float64(jd-J2000)) > 0.001 {
		t.Errorf("Expected JD %.6f, got %.6f", J2000, jd)
	}
	
	// Test round trip conversion
	backToTime := tc.FromJulianDay(jd)
	if math.Abs(float64(j2000.Unix()-backToTime.Unix())) > 1 {
		t.Errorf("Round trip conversion failed: original %v, converted %v", j2000, backToTime)
	}
}

func TestJulianDayOperations(t *testing.T) {
	jd := JulianDay(2451545.0) // J2000.0
	
	// Test Add
	jdPlusOne := jd.Add(1.0)
	if math.Abs(float64(jdPlusOne-jd)-1.0) > 0.001 {
		t.Errorf("Add operation failed: expected difference of 1.0, got %.6f", jdPlusOne-jd)
	}
	
	// Test Sub
	diff := jdPlusOne.Sub(jd)
	if math.Abs(diff-1.0) > 0.001 {
		t.Errorf("Sub operation failed: expected 1.0, got %.6f", diff)
	}
	
	// Test String
	str := jd.String()
	expected := "JD 2451545.000000"
	if str != expected {
		t.Errorf("String() failed: expected %s, got %s", expected, str)
	}
}

func TestJulianCenturies(t *testing.T) {
	tc := NewTimeConverter()
	
	// Test J2000.0 epoch (should be 0 centuries)
	centuries := tc.JulianCenturies(J2000)
	if math.Abs(centuries) > 0.001 {
		t.Errorf("Expected 0 centuries for J2000.0, got %.6f", centuries)
	}
	
	// Test one century later
	jdOneCentury := J2000.Add(DaysPerCentury)
	centuries = tc.JulianCenturies(jdOneCentury)
	if math.Abs(centuries-1.0) > 0.001 {
		t.Errorf("Expected 1.0 centuries, got %.6f", centuries)
	}
}

func TestSiderealTime(t *testing.T) {
	tc := NewTimeConverter()
	
	// Test known values (approximate)
	jd := J2000
	gmst := tc.SiderealTime(jd)
	
	// GMST should be in range [0, 360)
	if gmst < 0 || gmst >= 360 {
		t.Errorf("GMST out of range: %.6f", gmst)
	}
	
	// Test Local Sidereal Time
	longitude := 0.0 // Greenwich
	lst := tc.LocalSiderealTime(jd, longitude)
	
	// For Greenwich, LST should equal GMST
	if math.Abs(lst-gmst) > 0.001 {
		t.Errorf("LST should equal GMST at Greenwich: GMST=%.6f, LST=%.6f", gmst, lst)
	}
	
	// Test with different longitude
	longitude = 90.0 // 90° East
	lst = tc.LocalSiderealTime(jd, longitude)
	expectedLst := gmst + longitude
	if expectedLst >= 360 {
		expectedLst -= 360
	}
	
	if math.Abs(lst-expectedLst) > 0.001 {
		t.Errorf("LST calculation failed: expected %.6f, got %.6f", expectedLst, lst)
	}
}

func TestDeltaT(t *testing.T) {
	tc := NewTimeConverter()
	
	// Test known approximate values
	testCases := []struct {
		year     int
		expected float64
		tolerance float64
	}{
		{2000, 64, 10},    // Around 64 seconds in 2000
		{2020, 69, 10},    // Around 69 seconds in 2020
	}
	
	for _, test := range testCases {
		deltaT := tc.DeltaT(test.year)
		if math.Abs(deltaT-test.expected) > test.tolerance {
			t.Errorf("DeltaT for year %d: expected ~%.1f, got %.6f", test.year, test.expected, deltaT)
		}
	}
}

func TestTimeConversions(t *testing.T) {
	tc := NewTimeConverter()
	
	// Test various dates
	testDates := []time.Time{
		time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 15, 18, 30, 0, 0, time.UTC),
		time.Date(1950, 6, 1, 6, 0, 0, 0, time.UTC),
		time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC),
	}
	
	for _, testDate := range testDates {
		jd := tc.ToJulianDay(testDate)
		backToTime := tc.FromJulianDay(jd)
		
		// Check if the conversion is accurate within 1 second
		diff := math.Abs(float64(testDate.Unix() - backToTime.Unix()))
		if diff > 1 {
			t.Errorf("Time conversion failed for %v: got %v, difference: %.1f seconds", testDate, backToTime, diff)
		}
	}
}

func TestJulianDayToTime(t *testing.T) {
	// Test the method on JulianDay type
	jd := JulianDay(2451545.0) // J2000.0
	convertedTime := jd.ToTime()
	
	expected := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	
	// Check if the conversion is accurate within 1 second
	diff := math.Abs(float64(expected.Unix() - convertedTime.Unix()))
	if diff > 1 {
		t.Errorf("JulianDay.ToTime() failed: expected %v, got %v", expected, convertedTime)
	}
}

func BenchmarkToJulianDay(b *testing.B) {
	tc := NewTimeConverter()
	testTime := time.Date(2024, 3, 15, 18, 30, 0, 0, time.UTC)
	
	for i := 0; i < b.N; i++ {
		tc.ToJulianDay(testTime)
	}
}

func BenchmarkFromJulianDay(b *testing.B) {
	tc := NewTimeConverter()
	testJD := JulianDay(2451545.0)
	
	for i := 0; i < b.N; i++ {
		tc.FromJulianDay(testJD)
	}
}

func BenchmarkSiderealTime(b *testing.B) {
	tc := NewTimeConverter()
	testJD := J2000
	
	for i := 0; i < b.N; i++ {
		tc.SiderealTime(testJD)
	}
}