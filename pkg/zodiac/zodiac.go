// Package zodiac provides zodiac sign calculations and interpretations.
package zodiac

import (
	"fmt"
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
)

// ZodiacSign represents a zodiac sign
type ZodiacSign struct {
	Name     string
	Symbol   string
	Element  string
	Quality  string
	Ruler    string
	StartDeg float64 // Starting degree in the ecliptic
	EndDeg   float64 // Ending degree in the ecliptic
}

// ZodiacPosition represents a position within a zodiac sign
type ZodiacPosition struct {
	Sign        ZodiacSign
	DegreeInSign float64 // 0-30 degrees within the sign
	AbsoluteDeg  float64 // 0-360 degrees absolute position
}

// ZodiacCalculator handles zodiac-related calculations
type ZodiacCalculator struct {
	signs []ZodiacSign
}

// NewZodiacCalculator creates a new zodiac calculator
func NewZodiacCalculator() *ZodiacCalculator {
	return &ZodiacCalculator{
		signs: getZodiacSigns(),
	}
}

// GetZodiacSigns returns all zodiac signs
func (zc *ZodiacCalculator) GetZodiacSigns() []ZodiacSign {
	return zc.signs
}

// GetSignByName returns a zodiac sign by name
func (zc *ZodiacCalculator) GetSignByName(name string) (ZodiacSign, bool) {
	for _, sign := range zc.signs {
		if sign.Name == name {
			return sign, true
		}
	}
	return ZodiacSign{}, false
}

// EclipticToZodiac converts ecliptic longitude to zodiac position
func (zc *ZodiacCalculator) EclipticToZodiac(longitude float64) ZodiacPosition {
	// Normalize longitude to 0-360 range
	longitude = coordinates.NormalizeAngle(longitude)
	
	// Find the zodiac sign
	signIndex := int(longitude / 30.0)
	if signIndex >= len(zc.signs) {
		signIndex = len(zc.signs) - 1
	}
	
	sign := zc.signs[signIndex]
	degreeInSign := longitude - sign.StartDeg
	
	return ZodiacPosition{
		Sign:        sign,
		DegreeInSign: degreeInSign,
		AbsoluteDeg:  longitude,
	}
}

// ZodiacToEcliptic converts zodiac position back to ecliptic longitude
func (zc *ZodiacCalculator) ZodiacToEcliptic(position ZodiacPosition) float64 {
	return position.Sign.StartDeg + position.DegreeInSign
}

// GetSignCompatibility returns compatibility score between two signs (0-100)
func (zc *ZodiacCalculator) GetSignCompatibility(sign1, sign2 ZodiacSign) float64 {
	// Simple compatibility based on elements and qualities
	elementScore := 0.0
	qualityScore := 0.0
	
	// Element compatibility
	if sign1.Element == sign2.Element {
		elementScore = 40.0 // Same element
	} else if isCompatibleElement(sign1.Element, sign2.Element) {
		elementScore = 30.0 // Compatible elements
	} else {
		elementScore = 10.0 // Less compatible elements
	}
	
	// Quality compatibility
	if sign1.Quality == sign2.Quality {
		qualityScore = 20.0 // Same quality
	} else if isCompatibleQuality(sign1.Quality, sign2.Quality) {
		qualityScore = 25.0 // Compatible qualities
	} else {
		qualityScore = 15.0 // Less compatible qualities
	}
	
	// Angular relationship bonus
	angle := math.Abs(sign1.StartDeg - sign2.StartDeg)
	if angle > 180 {
		angle = 360 - angle
	}
	
	angleScore := 0.0
	switch {
	case angle == 0:
		angleScore = 35.0 // Same sign
	case angle == 60 || angle == 120:
		angleScore = 35.0 // Trine or sextile
	case angle == 90 || angle == 180:
		angleScore = 15.0 // Square or opposition
	default:
		angleScore = 25.0 // Other angles
	}
	
	return elementScore + qualityScore + angleScore
}

// CalculateAspectAngle calculates the aspect angle between two zodiac positions
func (zc *ZodiacCalculator) CalculateAspectAngle(pos1, pos2 ZodiacPosition) float64 {
	angle := math.Abs(pos1.AbsoluteDeg - pos2.AbsoluteDeg)
	if angle > 180 {
		angle = 360 - angle
	}
	return angle
}

// FormatZodiacPosition formats a zodiac position as a string
func (zc *ZodiacCalculator) FormatZodiacPosition(position ZodiacPosition) string {
	degrees := int(position.DegreeInSign)
	minutes := int((position.DegreeInSign - float64(degrees)) * 60)
	seconds := int(((position.DegreeInSign - float64(degrees)) * 60 - float64(minutes)) * 60)
	
	return fmt.Sprintf("%d°%d'%d\" %s", degrees, minutes, seconds, position.Sign.Name)
}

// IsRetrograde determines if a planet appears to be in retrograde motion
// This is a simplified calculation based on orbital mechanics
func (zc *ZodiacCalculator) IsRetrograde(planetName string, longitude1, longitude2 float64, timeDiff float64) bool {
	// Calculate apparent motion
	motion := (longitude2 - longitude1) / timeDiff
	
	// Normalize motion
	if motion > 180 {
		motion -= 360
	} else if motion < -180 {
		motion += 360
	}
	
	// Different planets have different retrograde thresholds
	threshold := getRetrogradeThreshold(planetName)
	
	return motion < threshold
}

// getRetrogradeThreshold returns the retrograde motion threshold for a planet
func getRetrogradeThreshold(planetName string) float64 {
	thresholds := map[string]float64{
		"Mercury": -0.5,
		"Venus":   -0.3,
		"Mars":    -0.2,
		"Jupiter": -0.1,
		"Saturn":  -0.05,
		"Uranus":  -0.02,
		"Neptune": -0.01,
		"Pluto":   -0.008,
	}
	
	if threshold, exists := thresholds[planetName]; exists {
		return threshold
	}
	
	return -0.1 // Default threshold
}

// isCompatibleElement checks if two elements are compatible
func isCompatibleElement(element1, element2 string) bool {
	compatible := map[string][]string{
		"Fire":  {"Air"},
		"Earth": {"Water"},
		"Air":   {"Fire"},
		"Water": {"Earth"},
	}
	
	if compatibleElements, exists := compatible[element1]; exists {
		for _, elem := range compatibleElements {
			if elem == element2 {
				return true
			}
		}
	}
	
	return false
}

// isCompatibleQuality checks if two qualities are compatible
func isCompatibleQuality(quality1, quality2 string) bool {
	// Cardinal and mutable are generally compatible
	// Fixed signs can be compatible with both but with more challenge
	compatible := map[string][]string{
		"Cardinal": {"Mutable"},
		"Fixed":    {"Cardinal", "Mutable"},
		"Mutable":  {"Cardinal"},
	}
	
	if compatibleQualities, exists := compatible[quality1]; exists {
		for _, qual := range compatibleQualities {
			if qual == quality2 {
				return true
			}
		}
	}
	
	return false
}

// getZodiacSigns returns the 12 zodiac signs with their properties
func getZodiacSigns() []ZodiacSign {
	return []ZodiacSign{
		{
			Name:     "Aries",
			Symbol:   "♈",
			Element:  "Fire",
			Quality:  "Cardinal",
			Ruler:    "Mars",
			StartDeg: 0,
			EndDeg:   30,
		},
		{
			Name:     "Taurus",
			Symbol:   "♉",
			Element:  "Earth",
			Quality:  "Fixed",
			Ruler:    "Venus",
			StartDeg: 30,
			EndDeg:   60,
		},
		{
			Name:     "Gemini",
			Symbol:   "♊",
			Element:  "Air",
			Quality:  "Mutable",
			Ruler:    "Mercury",
			StartDeg: 60,
			EndDeg:   90,
		},
		{
			Name:     "Cancer",
			Symbol:   "♋",
			Element:  "Water",
			Quality:  "Cardinal",
			Ruler:    "Moon",
			StartDeg: 90,
			EndDeg:   120,
		},
		{
			Name:     "Leo",
			Symbol:   "♌",
			Element:  "Fire",
			Quality:  "Fixed",
			Ruler:    "Sun",
			StartDeg: 120,
			EndDeg:   150,
		},
		{
			Name:     "Virgo",
			Symbol:   "♍",
			Element:  "Earth",
			Quality:  "Mutable",
			Ruler:    "Mercury",
			StartDeg: 150,
			EndDeg:   180,
		},
		{
			Name:     "Libra",
			Symbol:   "♎",
			Element:  "Air",
			Quality:  "Cardinal",
			Ruler:    "Venus",
			StartDeg: 180,
			EndDeg:   210,
		},
		{
			Name:     "Scorpio",
			Symbol:   "♏",
			Element:  "Water",
			Quality:  "Fixed",
			Ruler:    "Mars",
			StartDeg: 210,
			EndDeg:   240,
		},
		{
			Name:     "Sagittarius",
			Symbol:   "♐",
			Element:  "Fire",
			Quality:  "Mutable",
			Ruler:    "Jupiter",
			StartDeg: 240,
			EndDeg:   270,
		},
		{
			Name:     "Capricorn",
			Symbol:   "♑",
			Element:  "Earth",
			Quality:  "Cardinal",
			Ruler:    "Saturn",
			StartDeg: 270,
			EndDeg:   300,
		},
		{
			Name:     "Aquarius",
			Symbol:   "♒",
			Element:  "Air",
			Quality:  "Fixed",
			Ruler:    "Uranus",
			StartDeg: 300,
			EndDeg:   330,
		},
		{
			Name:     "Pisces",
			Symbol:   "♓",
			Element:  "Water",
			Quality:  "Mutable",
			Ruler:    "Neptune",
			StartDeg: 330,
			EndDeg:   360,
		},
	}
}

// String methods
func (zs ZodiacSign) String() string {
	return fmt.Sprintf("%s (%s) - %s %s, ruled by %s", zs.Name, zs.Symbol, zs.Element, zs.Quality, zs.Ruler)
}

func (zp ZodiacPosition) String() string {
	return fmt.Sprintf("%.2f° %s (%.2f° absolute)", zp.DegreeInSign, zp.Sign.Name, zp.AbsoluteDeg)
}