package zodiac

import (
	"fmt"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
)

// Sign represents a zodiac sign
type Sign int

const (
	Aries Sign = iota
	Taurus
	Gemini
	Cancer
	Leo
	Virgo
	Libra
	Scorpio
	Sagittarius
	Capricorn
	Aquarius
	Pisces
)

// String returns the name of the zodiac sign
func (s Sign) String() string {
	names := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	if int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}

// Element represents the classical element of a zodiac sign
type Element int

const (
	Fire Element = iota
	Earth
	Air
	Water
)

// String returns the name of the element
func (e Element) String() string {
	names := []string{"Fire", "Earth", "Air", "Water"}
	if int(e) < len(names) {
		return names[e]
	}
	return "Unknown"
}

// Quality represents the quality (modality) of a zodiac sign
type Quality int

const (
	Cardinal Quality = iota
	Fixed
	Mutable
)

// String returns the name of the quality
func (q Quality) String() string {
	names := []string{"Cardinal", "Fixed", "Mutable"}
	if int(q) < len(names) {
		return names[q]
	}
	return "Unknown"
}

// SignInfo contains information about a zodiac sign
type SignInfo struct {
	Sign     Sign
	Element  Element
	Quality  Quality
	Ruler    string
	Symbol   string
	Degree   coordinates.Angle // Position within the sign (0-30 degrees)
}

// GetSignFromLongitude determines the zodiac sign from ecliptic longitude
func GetSignFromLongitude(longitude coordinates.Angle) SignInfo {
	// Normalize longitude to [0, 360)
	normalizedLongitude := longitude.Normalize()
	
	// Each sign is 30 degrees
	signIndex := int(normalizedLongitude) / 30
	degreeInSign := float64(normalizedLongitude) - float64(signIndex*30)
	
	// Ensure we're within bounds
	if signIndex >= 12 {
		signIndex = 11
	}
	
	sign := Sign(signIndex)
	
	return SignInfo{
		Sign:    sign,
		Element: getElement(sign),
		Quality: getQuality(sign),
		Ruler:   getRuler(sign),
		Symbol:  getSymbol(sign),
		Degree:  coordinates.Angle(degreeInSign),
	}
}

// getElement returns the element for a zodiac sign
func getElement(sign Sign) Element {
	elements := []Element{
		Fire,  // Aries
		Earth, // Taurus
		Air,   // Gemini
		Water, // Cancer
		Fire,  // Leo
		Earth, // Virgo
		Air,   // Libra
		Water, // Scorpio
		Fire,  // Sagittarius
		Earth, // Capricorn
		Air,   // Aquarius
		Water, // Pisces
	}
	return elements[sign]
}

// getQuality returns the quality for a zodiac sign
func getQuality(sign Sign) Quality {
	qualities := []Quality{
		Cardinal, // Aries
		Fixed,    // Taurus
		Mutable,  // Gemini
		Cardinal, // Cancer
		Fixed,    // Leo
		Mutable,  // Virgo
		Cardinal, // Libra
		Fixed,    // Scorpio
		Mutable,  // Sagittarius
		Cardinal, // Capricorn
		Fixed,    // Aquarius
		Mutable,  // Pisces
	}
	return qualities[sign]
}

// getRuler returns the traditional ruler for a zodiac sign
func getRuler(sign Sign) string {
	rulers := []string{
		"Mars",     // Aries
		"Venus",    // Taurus
		"Mercury",  // Gemini
		"Moon",     // Cancer
		"Sun",      // Leo
		"Mercury",  // Virgo
		"Venus",    // Libra
		"Mars",     // Scorpio
		"Jupiter",  // Sagittarius
		"Saturn",   // Capricorn
		"Saturn",   // Aquarius
		"Jupiter",  // Pisces
	}
	return rulers[sign]
}

// getSymbol returns the symbol for a zodiac sign
func getSymbol(sign Sign) string {
	symbols := []string{
		"♈", // Aries
		"♉", // Taurus
		"♊", // Gemini
		"♋", // Cancer
		"♌", // Leo
		"♍", // Virgo
		"♎", // Libra
		"♏", // Scorpio
		"♐", // Sagittarius
		"♑", // Capricorn
		"♒", // Aquarius
		"♓", // Pisces
	}
	return symbols[sign]
}

// IsCompatible checks basic compatibility between two signs based on element and quality
func IsCompatible(sign1, sign2 Sign) bool {
	element1 := getElement(sign1)
	element2 := getElement(sign2)
	
	// Same element signs are generally compatible
	if element1 == element2 {
		return true
	}
	
	// Fire and Air are compatible
	if (element1 == Fire && element2 == Air) || (element1 == Air && element2 == Fire) {
		return true
	}
	
	// Earth and Water are compatible
	if (element1 == Earth && element2 == Water) || (element1 == Water && element2 == Earth) {
		return true
	}
	
	return false
}

// GetOppositeSign returns the opposite sign (180 degrees away)
func GetOppositeSign(sign Sign) Sign {
	opposite := (int(sign) + 6) % 12
	return Sign(opposite)
}

// FormatPosition formats a position in a sign as a string
func (si SignInfo) FormatPosition() string {
	degrees := int(si.Degree)
	minutes := int((si.Degree - coordinates.Angle(degrees)) * 60)
	return fmt.Sprintf("%s %d°%d'", si.Sign.String(), degrees, minutes)
}