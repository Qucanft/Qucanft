package coordinates

import (
	"testing"
	"math"
)

func TestCoordinateTransformer(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test obliquity
	expectedObliquity := J2000Obliquity
	if math.Abs(ct.GetObliquity()-expectedObliquity) > 0.001 {
		t.Errorf("Expected obliquity %.6f, got %.6f", expectedObliquity, ct.GetObliquity())
	}
	
	// Test custom obliquity
	customObliquity := 23.5
	ct.SetObliquity(customObliquity)
	if math.Abs(ct.GetObliquity()-customObliquity) > 0.001 {
		t.Errorf("Expected custom obliquity %.6f, got %.6f", customObliquity, ct.GetObliquity())
	}
}

func TestEquatorialToEcliptic(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test conversion with known values
	// Point on celestial equator
	eq := EquatorialCoordinates{
		RightAscension: 0.0,
		Declination:    0.0,
		Distance:       1.0,
	}
	
	ec := ct.EquatorialToEcliptic(eq)
	
	// Should be close to ecliptic coordinates
	if math.Abs(ec.Longitude-0.0) > 0.1 || math.Abs(ec.Latitude-0.0) > 0.1 {
		t.Errorf("Equatorial to ecliptic conversion failed: got lon=%.6f, lat=%.6f", ec.Longitude, ec.Latitude)
	}
	
	// Test distance preservation
	if math.Abs(ec.Distance-eq.Distance) > 0.001 {
		t.Errorf("Distance not preserved: expected %.6f, got %.6f", eq.Distance, ec.Distance)
	}
}

func TestEclipticToEquatorial(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test conversion with known values
	ec := EclipticCoordinates{
		Longitude: 0.0,
		Latitude:  0.0,
		Distance:  1.0,
	}
	
	eq := ct.EclipticToEquatorial(ec)
	
	// Should be close to equatorial coordinates
	if math.Abs(eq.RightAscension-0.0) > 0.1 || math.Abs(eq.Declination-0.0) > 0.1 {
		t.Errorf("Ecliptic to equatorial conversion failed: got RA=%.6f, Dec=%.6f", eq.RightAscension, eq.Declination)
	}
	
	// Test distance preservation
	if math.Abs(eq.Distance-ec.Distance) > 0.001 {
		t.Errorf("Distance not preserved: expected %.6f, got %.6f", ec.Distance, eq.Distance)
	}
}

func TestRoundTripConversion(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test round trip: equatorial -> ecliptic -> equatorial
	original := EquatorialCoordinates{
		RightAscension: 45.0,
		Declination:    30.0,
		Distance:       1.5,
	}
	
	ecliptic := ct.EquatorialToEcliptic(original)
	backToEquatorial := ct.EclipticToEquatorial(ecliptic)
	
	tolerance := 0.001
	if math.Abs(backToEquatorial.RightAscension-original.RightAscension) > tolerance ||
		math.Abs(backToEquatorial.Declination-original.Declination) > tolerance ||
		math.Abs(backToEquatorial.Distance-original.Distance) > tolerance {
		t.Errorf("Round trip conversion failed: original %v, final %v", original, backToEquatorial)
	}
}

func TestEquatorialToHorizontal(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test with object at zenith (LST = RA, altitude = latitude)
	eq := EquatorialCoordinates{
		RightAscension: 90.0,
		Declination:    50.0,
		Distance:       1.0,
	}
	
	lst := 90.0     // Local sidereal time
	latitude := 50.0 // Observer latitude
	
	hz := ct.EquatorialToHorizontal(eq, lst, latitude)
	
	// Object should be at zenith (altitude = 90°)
	expectedAltitude := 90.0
	if math.Abs(hz.Altitude-expectedAltitude) > 1.0 {
		t.Errorf("Expected altitude %.1f°, got %.6f°", expectedAltitude, hz.Altitude)
	}
}

func TestHorizontalToEquatorial(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test with object at zenith
	hz := HorizontalCoordinates{
		Azimuth:  0.0,  // North
		Altitude: 90.0, // Zenith
	}
	
	lst := 90.0     // Local sidereal time
	latitude := 50.0 // Observer latitude
	
	eq := ct.HorizontalToEquatorial(hz, lst, latitude)
	
	// Check if conversion is reasonable
	if eq.RightAscension < 0 || eq.RightAscension >= 360 {
		t.Errorf("RA out of range: %.6f", eq.RightAscension)
	}
	
	if eq.Declination < -90 || eq.Declination > 90 {
		t.Errorf("Declination out of range: %.6f", eq.Declination)
	}
}

func TestAngularSeparation(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	// Test with identical coordinates
	coord1 := EquatorialCoordinates{
		RightAscension: 45.0,
		Declination:    30.0,
		Distance:       1.0,
	}
	
	coord2 := coord1
	
	separation := ct.AngularSeparation(coord1, coord2)
	if math.Abs(separation) > 0.001 {
		t.Errorf("Expected separation of 0°, got %.6f°", separation)
	}
	
	// Test with coordinates 90° apart
	coord2.RightAscension = 135.0
	coord2.Declination = 30.0
	
	separation = ct.AngularSeparation(coord1, coord2)
	// Angular separation calculation is more complex than simple RA difference
	if separation < 60 || separation > 120 {
		t.Errorf("Expected separation between 60-120°, got %.6f°", separation)
	}
}

func TestPositionAngle(t *testing.T) {
	ct := NewCoordinateTransformer()
	
	coord1 := EquatorialCoordinates{
		RightAscension: 0.0,
		Declination:    0.0,
		Distance:       1.0,
	}
	
	coord2 := EquatorialCoordinates{
		RightAscension: 90.0,
		Declination:    0.0,
		Distance:       1.0,
	}
	
	pa := ct.PositionAngle(coord1, coord2)
	
	// Position angle should be in range [0, 360)
	if pa < 0 || pa >= 360 {
		t.Errorf("Position angle out of range: %.6f", pa)
	}
}

func TestNormalizeAngle(t *testing.T) {
	testCases := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{360.0, 0.0},
		{-90.0, 270.0},
		{450.0, 90.0},
		{-450.0, 270.0},
		{720.0, 0.0},
	}
	
	for _, test := range testCases {
		result := NormalizeAngle(test.input)
		if math.Abs(result-test.expected) > 0.001 {
			t.Errorf("NormalizeAngle(%.1f): expected %.1f, got %.6f", test.input, test.expected, result)
		}
	}
}

func TestAngleDifference(t *testing.T) {
	testCases := []struct {
		angle1   float64
		angle2   float64
		expected float64
	}{
		{0.0, 90.0, 90.0},
		{90.0, 0.0, -90.0},
		{350.0, 10.0, 20.0},
		{10.0, 350.0, -20.0},
		{0.0, 180.0, 180.0},
		{180.0, 0.0, -180.0},
	}
	
	for _, test := range testCases {
		result := AngleDifference(test.angle1, test.angle2)
		if math.Abs(result-test.expected) > 0.001 {
			t.Errorf("AngleDifference(%.1f, %.1f): expected %.1f, got %.6f", test.angle1, test.angle2, test.expected, result)
		}
	}
}

func TestCoordinateStringMethods(t *testing.T) {
	eq := EquatorialCoordinates{
		RightAscension: 45.123456,
		Declination:    30.654321,
		Distance:       1.5,
	}
	
	str := eq.String()
	if str == "" {
		t.Error("EquatorialCoordinates String() returned empty string")
	}
	
	ec := EclipticCoordinates{
		Longitude: 120.123456,
		Latitude:  -15.654321,
		Distance:  2.0,
	}
	
	str = ec.String()
	if str == "" {
		t.Error("EclipticCoordinates String() returned empty string")
	}
	
	hz := HorizontalCoordinates{
		Azimuth:  270.123456,
		Altitude: 45.654321,
	}
	
	str = hz.String()
	if str == "" {
		t.Error("HorizontalCoordinates String() returned empty string")
	}
}

func BenchmarkEquatorialToEcliptic(b *testing.B) {
	ct := NewCoordinateTransformer()
	eq := EquatorialCoordinates{
		RightAscension: 45.0,
		Declination:    30.0,
		Distance:       1.0,
	}
	
	for i := 0; i < b.N; i++ {
		ct.EquatorialToEcliptic(eq)
	}
}

func BenchmarkEclipticToEquatorial(b *testing.B) {
	ct := NewCoordinateTransformer()
	ec := EclipticCoordinates{
		Longitude: 45.0,
		Latitude:  30.0,
		Distance:  1.0,
	}
	
	for i := 0; i < b.N; i++ {
		ct.EclipticToEquatorial(ec)
	}
}

func BenchmarkAngularSeparation(b *testing.B) {
	ct := NewCoordinateTransformer()
	coord1 := EquatorialCoordinates{
		RightAscension: 45.0,
		Declination:    30.0,
		Distance:       1.0,
	}
	coord2 := EquatorialCoordinates{
		RightAscension: 135.0,
		Declination:    -30.0,
		Distance:       1.0,
	}
	
	for i := 0; i < b.N; i++ {
		ct.AngularSeparation(coord1, coord2)
	}
}