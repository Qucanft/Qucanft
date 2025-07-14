package zodiac

import (
	"testing"
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
)

func TestZodiacCalculator(t *testing.T) {
	zc := NewZodiacCalculator()
	
	// Test getting zodiac signs
	signs := zc.GetZodiacSigns()
	if len(signs) != 12 {
		t.Errorf("Expected 12 zodiac signs, got %d", len(signs))
	}
	
	// Test sign properties
	aries := signs[0]
	if aries.Name != "Aries" || aries.Symbol != "♈" {
		t.Errorf("First sign should be Aries (♈), got %s (%s)", aries.Name, aries.Symbol)
	}
	
	if aries.Element != "Fire" || aries.Quality != "Cardinal" {
		t.Errorf("Aries should be Fire Cardinal, got %s %s", aries.Element, aries.Quality)
	}
	
	if aries.StartDeg != 0 || aries.EndDeg != 30 {
		t.Errorf("Aries should be 0-30 degrees, got %.1f-%.1f", aries.StartDeg, aries.EndDeg)
	}
}

func TestGetSignByName(t *testing.T) {
	zc := NewZodiacCalculator()
	
	// Test existing sign
	leo, exists := zc.GetSignByName("Leo")
	if !exists {
		t.Error("Leo sign not found")
	}
	
	if leo.Name != "Leo" || leo.Symbol != "♌" {
		t.Errorf("Leo data incorrect: got %s (%s)", leo.Name, leo.Symbol)
	}
	
	// Test non-existent sign
	_, exists = zc.GetSignByName("NonExistentSign")
	if exists {
		t.Error("Non-existent sign should not be found")
	}
}

func TestEclipticToZodiac(t *testing.T) {
	zc := NewZodiacCalculator()
	
	testCases := []struct {
		longitude    float64
		expectedSign string
		expectedDeg  float64
	}{
		{0.0, "Aries", 0.0},
		{15.0, "Aries", 15.0},
		{30.0, "Taurus", 0.0},
		{45.0, "Taurus", 15.0},
		{90.0, "Cancer", 0.0},
		{120.0, "Leo", 0.0},
		{180.0, "Libra", 0.0},
		{270.0, "Capricorn", 0.0},
		{359.0, "Pisces", 29.0},
		{360.0, "Aries", 0.0}, // Should normalize to 0
	}
	
	for _, test := range testCases {
		position := zc.EclipticToZodiac(test.longitude)
		
		if position.Sign.Name != test.expectedSign {
			t.Errorf("Longitude %.1f: expected %s, got %s", test.longitude, test.expectedSign, position.Sign.Name)
		}
		
		if math.Abs(position.DegreeInSign-test.expectedDeg) > 0.1 {
			t.Errorf("Longitude %.1f: expected %.1f degrees in sign, got %.6f", test.longitude, test.expectedDeg, position.DegreeInSign)
		}
		
		if math.Abs(position.AbsoluteDeg-coordinates.NormalizeAngle(test.longitude)) > 0.1 {
			t.Errorf("Longitude %.1f: absolute degree mismatch", test.longitude)
		}
	}
}

func TestZodiacToEcliptic(t *testing.T) {
	zc := NewZodiacCalculator()
	
	// Test round trip conversion
	testLongitudes := []float64{0, 15, 30, 45, 90, 120, 180, 270, 359}
	
	for _, longitude := range testLongitudes {
		zodiacPos := zc.EclipticToZodiac(longitude)
		backToEcliptic := zc.ZodiacToEcliptic(zodiacPos)
		
		normalizedLongitude := coordinates.NormalizeAngle(longitude)
		if math.Abs(backToEcliptic-normalizedLongitude) > 0.001 {
			t.Errorf("Round trip failed: %.1f -> %.6f", longitude, backToEcliptic)
		}
	}
}

func TestGetSignCompatibility(t *testing.T) {
	zc := NewZodiacCalculator()
	
	aries, _ := zc.GetSignByName("Aries")
	leo, _ := zc.GetSignByName("Leo")
	cancer, _ := zc.GetSignByName("Cancer")
	
	// Test same element (Fire signs)
	compatibility := zc.GetSignCompatibility(aries, leo)
	if compatibility < 70 { // Fire signs should have high compatibility
		t.Errorf("Fire signs should have high compatibility, got %.1f", compatibility)
	}
	
	// Test different elements
	compatibility = zc.GetSignCompatibility(aries, cancer)
	if compatibility < 0 || compatibility > 100 {
		t.Errorf("Compatibility should be 0-100, got %.1f", compatibility)
	}
	
	// Test same sign
	compatibility = zc.GetSignCompatibility(aries, aries)
	if compatibility < 80 { // Same sign should have very high compatibility
		t.Errorf("Same sign should have very high compatibility, got %.1f", compatibility)
	}
}

func TestCalculateAspectAngle(t *testing.T) {
	zc := NewZodiacCalculator()
	
	pos1 := ZodiacPosition{
		AbsoluteDeg: 0.0,
	}
	
	pos2 := ZodiacPosition{
		AbsoluteDeg: 90.0,
	}
	
	angle := zc.CalculateAspectAngle(pos1, pos2)
	expectedAngle := 90.0
	
	if math.Abs(angle-expectedAngle) > 0.1 {
		t.Errorf("Expected angle %.1f, got %.6f", expectedAngle, angle)
	}
	
	// Test angle > 180 (should return smaller angle)
	pos2.AbsoluteDeg = 270.0
	angle = zc.CalculateAspectAngle(pos1, pos2)
	expectedAngle = 90.0 // Should be 90, not 270
	
	if math.Abs(angle-expectedAngle) > 0.1 {
		t.Errorf("Expected angle %.1f, got %.6f", expectedAngle, angle)
	}
}

func TestFormatZodiacPosition(t *testing.T) {
	zc := NewZodiacCalculator()
	
	position := ZodiacPosition{
		DegreeInSign: 15.5,
		AbsoluteDeg:  15.5,
	}
	position.Sign, _ = zc.GetSignByName("Aries")
	
	formatted := zc.FormatZodiacPosition(position)
	if formatted == "" {
		t.Error("FormatZodiacPosition returned empty string")
	}
	
	// Should contain degrees, minutes, seconds, and sign name
	if !contains(formatted, "15°") || !contains(formatted, "Aries") {
		t.Errorf("Formatted position missing expected elements: %s", formatted)
	}
}

func TestIsRetrograde(t *testing.T) {
	zc := NewZodiacCalculator()
	
	// Test forward motion (not retrograde)
	retrograde := zc.IsRetrograde("Mercury", 0.0, 1.0, 1.0)
	if retrograde {
		t.Error("Forward motion should not be retrograde")
	}
	
	// Test backward motion (retrograde) - longitude decreases over time
	retrograde = zc.IsRetrograde("Mercury", 1.0, 0.0, 1.0)
	if !retrograde {
		t.Error("Backward motion should be retrograde")
	}
	
	// Test that function works with other planets
	retrograde = zc.IsRetrograde("Jupiter", 1.0, 0.0, 1.0)
	if !retrograde {
		t.Error("Backward motion should be retrograde for Jupiter")
	}
}

func TestZodiacSignElements(t *testing.T) {
	zc := NewZodiacCalculator()
	
	elementCounts := map[string]int{
		"Fire":  0,
		"Earth": 0,
		"Air":   0,
		"Water": 0,
	}
	
	qualityCounts := map[string]int{
		"Cardinal": 0,
		"Fixed":    0,
		"Mutable":  0,
	}
	
	signs := zc.GetZodiacSigns()
	for _, sign := range signs {
		elementCounts[sign.Element]++
		qualityCounts[sign.Quality]++
	}
	
	// Each element should have 3 signs
	for element, count := range elementCounts {
		if count != 3 {
			t.Errorf("Element %s should have 3 signs, got %d", element, count)
		}
	}
	
	// Each quality should have 4 signs
	for quality, count := range qualityCounts {
		if count != 4 {
			t.Errorf("Quality %s should have 4 signs, got %d", quality, count)
		}
	}
}

func TestZodiacSignDegrees(t *testing.T) {
	zc := NewZodiacCalculator()
	
	signs := zc.GetZodiacSigns()
	
	for i, sign := range signs {
		expectedStart := float64(i * 30)
		expectedEnd := float64((i + 1) * 30)
		
		if sign.StartDeg != expectedStart {
			t.Errorf("Sign %s start degree: expected %.1f, got %.1f", sign.Name, expectedStart, sign.StartDeg)
		}
		
		if sign.EndDeg != expectedEnd {
			t.Errorf("Sign %s end degree: expected %.1f, got %.1f", sign.Name, expectedEnd, sign.EndDeg)
		}
		
		// Check sign size is 30 degrees
		if (sign.EndDeg - sign.StartDeg) != 30.0 {
			t.Errorf("Sign %s should be 30 degrees, got %.1f", sign.Name, sign.EndDeg-sign.StartDeg)
		}
	}
}

func TestZodiacStringMethods(t *testing.T) {
	zc := NewZodiacCalculator()
	
	aries, _ := zc.GetSignByName("Aries")
	str := aries.String()
	if str == "" {
		t.Error("ZodiacSign String() returned empty string")
	}
	
	position := zc.EclipticToZodiac(15.5)
	str = position.String()
	if str == "" {
		t.Error("ZodiacPosition String() returned empty string")
	}
}

func TestCompatibilityFunctions(t *testing.T) {
	// Test element compatibility
	if !isCompatibleElement("Fire", "Air") {
		t.Error("Fire and Air should be compatible")
	}
	
	if !isCompatibleElement("Earth", "Water") {
		t.Error("Earth and Water should be compatible")
	}
	
	if isCompatibleElement("Fire", "Water") {
		t.Error("Fire and Water should not be compatible")
	}
	
	// Test quality compatibility
	if !isCompatibleQuality("Cardinal", "Mutable") {
		t.Error("Cardinal and Mutable should be compatible")
	}
	
	if !isCompatibleQuality("Fixed", "Cardinal") {
		t.Error("Fixed and Cardinal should be compatible")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   len(s) > len(substr) && s[:len(substr)] == substr ||
		   (len(s) > len(substr) && func() bool {
			   for i := 0; i <= len(s)-len(substr); i++ {
				   if s[i:i+len(substr)] == substr {
					   return true
				   }
			   }
			   return false
		   }())
}

func BenchmarkEclipticToZodiac(b *testing.B) {
	zc := NewZodiacCalculator()
	
	for i := 0; i < b.N; i++ {
		zc.EclipticToZodiac(123.456)
	}
}

func BenchmarkGetSignCompatibility(b *testing.B) {
	zc := NewZodiacCalculator()
	aries, _ := zc.GetSignByName("Aries")
	leo, _ := zc.GetSignByName("Leo")
	
	for i := 0; i < b.N; i++ {
		zc.GetSignCompatibility(aries, leo)
	}
}

func BenchmarkFormatZodiacPosition(b *testing.B) {
	zc := NewZodiacCalculator()
	position := zc.EclipticToZodiac(123.456)
	
	for i := 0; i < b.N; i++ {
		zc.FormatZodiacPosition(position)
	}
}