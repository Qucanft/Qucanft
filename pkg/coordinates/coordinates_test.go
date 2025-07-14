package coordinates

import (
	"math"
	"testing"
)

func TestAngleNormalization(t *testing.T) {
	tests := []struct {
		input    Angle
		expected Angle
	}{
		{Angle(0), Angle(0)},
		{Angle(180), Angle(180)},
		{Angle(360), Angle(0)},
		{Angle(450), Angle(90)},
		{Angle(-90), Angle(270)},
		{Angle(-180), Angle(180)},
	}
	
	for _, test := range tests {
		result := test.input.Normalize()
		if result != test.expected {
			t.Errorf("Normalize(%f) = %f, expected %f", test.input, result, test.expected)
		}
	}
}

func TestAngleToRadians(t *testing.T) {
	tests := []struct {
		degrees Angle
		radians float64
	}{
		{Angle(0), 0},
		{Angle(90), math.Pi / 2},
		{Angle(180), math.Pi},
		{Angle(270), 3 * math.Pi / 2},
		{Angle(360), 2 * math.Pi},
	}
	
	for _, test := range tests {
		result := test.degrees.ToRadians()
		if math.Abs(result-test.radians) > 1e-10 {
			t.Errorf("ToRadians(%f) = %f, expected %f", test.degrees, result, test.radians)
		}
	}
}

func TestCoordinateTransformations(t *testing.T) {
	obliquity := Angle(23.4367) // Approximate obliquity for J2000.0
	
	// Test a point on the ecliptic
	ecliptic := EclipticCoordinates{
		Longitude: Angle(0),   // Vernal equinox
		Latitude:  Angle(0),
	}
	
	equatorial := ecliptic.ToEquatorial(obliquity)
	
	// At vernal equinox, RA should be 0 and Dec should be 0
	if math.Abs(float64(equatorial.RightAscension)) > 1e-10 {
		t.Errorf("Expected RA = 0 at vernal equinox, got %f", equatorial.RightAscension)
	}
	if math.Abs(float64(equatorial.Declination)) > 1e-10 {
		t.Errorf("Expected Dec = 0 at vernal equinox, got %f", equatorial.Declination)
	}
	
	// Test round-trip conversion
	backToEcliptic := equatorial.ToEcliptic(obliquity)
	if math.Abs(float64(backToEcliptic.Longitude-ecliptic.Longitude)) > 1e-10 {
		t.Errorf("Round-trip longitude conversion failed")
	}
	if math.Abs(float64(backToEcliptic.Latitude-ecliptic.Latitude)) > 1e-10 {
		t.Errorf("Round-trip latitude conversion failed")
	}
}

func TestAngularDistance(t *testing.T) {
	// Test distance between two points on the ecliptic
	coord1 := EclipticCoordinates{Longitude: Angle(0), Latitude: Angle(0)}
	coord2 := EclipticCoordinates{Longitude: Angle(90), Latitude: Angle(0)}
	
	distance := AngularDistance(coord1, coord2)
	expected := Angle(90)
	
	if math.Abs(float64(distance-expected)) > 1e-10 {
		t.Errorf("Expected distance = %f, got %f", expected, distance)
	}
	
	// Test distance between opposite points
	coord3 := EclipticCoordinates{Longitude: Angle(180), Latitude: Angle(0)}
	distance2 := AngularDistance(coord1, coord3)
	expected2 := Angle(180)
	
	if math.Abs(float64(distance2-expected2)) > 1e-10 {
		t.Errorf("Expected distance = %f, got %f", expected2, distance2)
	}
}