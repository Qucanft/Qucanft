package astrology

import (
	"time"
)

// ZodiacSign represents the 12 zodiac signs
type ZodiacSign int

const (
	Aries ZodiacSign = iota
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
func (z ZodiacSign) String() string {
	names := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	if z < 0 || int(z) >= len(names) {
		return "Unknown"
	}
	return names[z]
}

// Planet represents celestial bodies in astrology
type Planet int

const (
	Sun Planet = iota
	Moon
	Mercury
	Venus
	Mars
	Jupiter
	Saturn
	Uranus
	Neptune
	Pluto
)

// String returns the name of the planet
func (p Planet) String() string {
	names := []string{
		"Sun", "Moon", "Mercury", "Venus", "Mars",
		"Jupiter", "Saturn", "Uranus", "Neptune", "Pluto",
	}
	if p < 0 || int(p) >= len(names) {
		return "Unknown"
	}
	return names[p]
}

// House represents the 12 astrological houses
type House int

const (
	FirstHouse House = iota + 1
	SecondHouse
	ThirdHouse
	FourthHouse
	FifthHouse
	SixthHouse
	SeventhHouse
	EighthHouse
	NinthHouse
	TenthHouse
	EleventhHouse
	TwelfthHouse
)

// String returns the name of the house
func (h House) String() string {
	names := []string{
		"First", "Second", "Third", "Fourth", "Fifth", "Sixth",
		"Seventh", "Eighth", "Ninth", "Tenth", "Eleventh", "Twelfth",
	}
	if h < 1 || int(h) > len(names) {
		return "Unknown"
	}
	return names[h-1] + " House"
}

// AspectType represents the angular relationships between planets
type AspectType int

const (
	Conjunction AspectType = iota // 0°
	Sextile                       // 60°
	Square                        // 90°
	Trine                         // 120°
	Opposition                    // 180°
)

// String returns the name of the aspect
func (a AspectType) String() string {
	names := []string{"Conjunction", "Sextile", "Square", "Trine", "Opposition"}
	if a < 0 || int(a) >= len(names) {
		return "Unknown"
	}
	return names[a]
}

// Angle returns the angle in degrees for the aspect
func (a AspectType) Angle() float64 {
	angles := []float64{0, 60, 90, 120, 180}
	if a < 0 || int(a) >= len(angles) {
		return 0
	}
	return angles[a]
}

// PlanetPosition represents a planet's position in the zodiac
type PlanetPosition struct {
	Planet   Planet
	Degree   float64 // 0-360 degrees
	Sign     ZodiacSign
	House    House
	Retrograde bool
}

// Aspect represents an angular relationship between two planets
type Aspect struct {
	Planet1 Planet
	Planet2 Planet
	Type    AspectType
	Angle   float64 // Exact angle between planets
	Orb     float64 // Deviation from perfect aspect
}

// Chart represents a complete astrological chart
type Chart struct {
	Timestamp time.Time
	Planets   []PlanetPosition
	Aspects   []Aspect
	Houses    [12]float64 // House cusps in degrees
}

// GetPlanetPosition returns the position of a specific planet
func (c *Chart) GetPlanetPosition(planet Planet) (*PlanetPosition, bool) {
	for _, pos := range c.Planets {
		if pos.Planet == planet {
			return &pos, true
		}
	}
	return nil, false
}

// GetAspects returns all aspects involving a specific planet
func (c *Chart) GetAspects(planet Planet) []Aspect {
	var aspects []Aspect
	for _, aspect := range c.Aspects {
		if aspect.Planet1 == planet || aspect.Planet2 == planet {
			aspects = append(aspects, aspect)
		}
	}
	return aspects
}