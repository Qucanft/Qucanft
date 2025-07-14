package zodiac

import (
	"testing"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
)

func TestGetSignFromLongitude(t *testing.T) {
	tests := []struct {
		longitude coordinates.Angle
		sign      Sign
		element   Element
		quality   Quality
	}{
		{coordinates.Angle(0), Aries, Fire, Cardinal},
		{coordinates.Angle(30), Taurus, Earth, Fixed},
		{coordinates.Angle(60), Gemini, Air, Mutable},
		{coordinates.Angle(90), Cancer, Water, Cardinal},
		{coordinates.Angle(120), Leo, Fire, Fixed},
		{coordinates.Angle(150), Virgo, Earth, Mutable},
		{coordinates.Angle(180), Libra, Air, Cardinal},
		{coordinates.Angle(210), Scorpio, Water, Fixed},
		{coordinates.Angle(240), Sagittarius, Fire, Mutable},
		{coordinates.Angle(270), Capricorn, Earth, Cardinal},
		{coordinates.Angle(300), Aquarius, Air, Fixed},
		{coordinates.Angle(330), Pisces, Water, Mutable},
	}
	
	for _, test := range tests {
		result := GetSignFromLongitude(test.longitude)
		if result.Sign != test.sign {
			t.Errorf("Longitude %f should be %s, got %s", test.longitude, test.sign, result.Sign)
		}
		if result.Element != test.element {
			t.Errorf("Longitude %f should have element %s, got %s", test.longitude, test.element, result.Element)
		}
		if result.Quality != test.quality {
			t.Errorf("Longitude %f should have quality %s, got %s", test.longitude, test.quality, result.Quality)
		}
	}
}

func TestSignStrings(t *testing.T) {
	expectedNames := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	
	for i, expected := range expectedNames {
		sign := Sign(i)
		if sign.String() != expected {
			t.Errorf("Sign(%d).String() = %s, expected %s", i, sign.String(), expected)
		}
	}
}

func TestElementStrings(t *testing.T) {
	expectedNames := []string{"Fire", "Earth", "Air", "Water"}
	
	for i, expected := range expectedNames {
		element := Element(i)
		if element.String() != expected {
			t.Errorf("Element(%d).String() = %s, expected %s", i, element.String(), expected)
		}
	}
}

func TestQualityStrings(t *testing.T) {
	expectedNames := []string{"Cardinal", "Fixed", "Mutable"}
	
	for i, expected := range expectedNames {
		quality := Quality(i)
		if quality.String() != expected {
			t.Errorf("Quality(%d).String() = %s, expected %s", i, quality.String(), expected)
		}
	}
}

func TestIsCompatible(t *testing.T) {
	// Test same element compatibility
	if !IsCompatible(Aries, Leo) {
		t.Error("Fire signs should be compatible")
	}
	
	if !IsCompatible(Taurus, Virgo) {
		t.Error("Earth signs should be compatible")
	}
	
	// Test Fire-Air compatibility
	if !IsCompatible(Aries, Gemini) {
		t.Error("Fire and Air signs should be compatible")
	}
	
	// Test Earth-Water compatibility
	if !IsCompatible(Taurus, Cancer) {
		t.Error("Earth and Water signs should be compatible")
	}
	
	// Test incompatible elements
	if IsCompatible(Aries, Taurus) {
		t.Error("Fire and Earth signs should not be compatible")
	}
}

func TestGetOppositeSign(t *testing.T) {
	tests := []struct {
		sign     Sign
		opposite Sign
	}{
		{Aries, Libra},
		{Taurus, Scorpio},
		{Gemini, Sagittarius},
		{Cancer, Capricorn},
		{Leo, Aquarius},
		{Virgo, Pisces},
	}
	
	for _, test := range tests {
		result := GetOppositeSign(test.sign)
		if result != test.opposite {
			t.Errorf("Opposite of %s should be %s, got %s", test.sign, test.opposite, result)
		}
	}
}

func TestSignDegreeCalculation(t *testing.T) {
	// Test a position at 15 degrees Aries (15 degrees total longitude)
	longitude := coordinates.Angle(15)
	result := GetSignFromLongitude(longitude)
	
	if result.Sign != Aries {
		t.Errorf("Expected Aries, got %s", result.Sign)
	}
	
	if result.Degree != coordinates.Angle(15) {
		t.Errorf("Expected 15° in sign, got %f°", result.Degree)
	}
	
	// Test a position at 15 degrees Taurus (45 degrees total longitude)
	longitude2 := coordinates.Angle(45)
	result2 := GetSignFromLongitude(longitude2)
	
	if result2.Sign != Taurus {
		t.Errorf("Expected Taurus, got %s", result2.Sign)
	}
	
	if result2.Degree != coordinates.Angle(15) {
		t.Errorf("Expected 15° in sign, got %f°", result2.Degree)
	}
}