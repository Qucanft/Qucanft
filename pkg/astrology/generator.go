package astrology

import (
	"math"
	"math/rand"
	"time"
)

// ChartGenerator provides methods to generate astrological charts
type ChartGenerator struct {
	rand *rand.Rand
}

// NewChartGenerator creates a new chart generator
func NewChartGenerator() *ChartGenerator {
	return &ChartGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateChart creates a sample astrological chart
func (cg *ChartGenerator) GenerateChart(timestamp time.Time) *Chart {
	chart := &Chart{
		Timestamp: timestamp,
		Planets:   make([]PlanetPosition, 0, 10),
		Aspects:   make([]Aspect, 0),
	}

	// Generate house cusps (starting positions for each house)
	for i := 0; i < 12; i++ {
		chart.Houses[i] = float64(i * 30) // Simple equal house system
	}

	// Generate planet positions
	planets := []Planet{Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto}
	for _, planet := range planets {
		position := cg.generatePlanetPosition(planet)
		chart.Planets = append(chart.Planets, position)
	}

	// Generate aspects between planets
	chart.Aspects = cg.generateAspects(chart.Planets)

	return chart
}

// generatePlanetPosition creates a planet position with realistic constraints
func (cg *ChartGenerator) generatePlanetPosition(planet Planet) PlanetPosition {
	// Generate degree (0-360)
	degree := cg.rand.Float64() * 360

	// Determine zodiac sign based on degree
	sign := ZodiacSign(int(degree / 30))

	// Determine house based on degree (simplified)
	house := House(int(degree/30) + 1)
	if house > TwelfthHouse {
		house = FirstHouse
	}

	// Some planets are more likely to be retrograde
	retrograde := false
	if planet == Mercury || planet == Venus || planet == Mars {
		retrograde = cg.rand.Float64() < 0.2 // 20% chance
	} else if planet >= Jupiter {
		retrograde = cg.rand.Float64() < 0.4 // 40% chance for outer planets
	}

	return PlanetPosition{
		Planet:     planet,
		Degree:     degree,
		Sign:       sign,
		House:      house,
		Retrograde: retrograde,
	}
}

// generateAspects creates aspects between planets
func (cg *ChartGenerator) generateAspects(positions []PlanetPosition) []Aspect {
	var aspects []Aspect
	aspectTypes := []AspectType{Conjunction, Sextile, Square, Trine, Opposition}
	orbTolerance := 8.0 // degrees

	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			planet1 := positions[i]
			planet2 := positions[j]

			// Calculate angle between planets
			angle := math.Abs(planet1.Degree - planet2.Degree)
			if angle > 180 {
				angle = 360 - angle
			}

			// Check if angle forms an aspect
			for _, aspectType := range aspectTypes {
				targetAngle := aspectType.Angle()
				orb := math.Abs(angle - targetAngle)

				if orb <= orbTolerance {
					aspects = append(aspects, Aspect{
						Planet1: planet1.Planet,
						Planet2: planet2.Planet,
						Type:    aspectType,
						Angle:   angle,
						Orb:     orb,
					})
					break // Only one aspect per planet pair
				}
			}
		}
	}

	return aspects
}

// GetZodiacDegree returns the degree within the zodiac sign (0-30)
func (p *PlanetPosition) GetZodiacDegree() float64 {
	return math.Mod(p.Degree, 30)
}

// GetElementalEnergy returns the elemental energy of the planet's sign
func (p *PlanetPosition) GetElementalEnergy() string {
	elements := []string{
		"Fire", "Earth", "Air", "Water", // Aries, Taurus, Gemini, Cancer
		"Fire", "Earth", "Air", "Water", // Leo, Virgo, Libra, Scorpio
		"Fire", "Earth", "Air", "Water", // Sagittarius, Capricorn, Aquarius, Pisces
	}
	return elements[p.Sign]
}

// GetModalityEnergy returns the modality energy of the planet's sign
func (p *PlanetPosition) GetModalityEnergy() string {
	modalities := []string{
		"Cardinal", "Fixed", "Mutable", "Cardinal", // Aries, Taurus, Gemini, Cancer
		"Fixed", "Mutable", "Cardinal", "Fixed",    // Leo, Virgo, Libra, Scorpio
		"Mutable", "Cardinal", "Fixed", "Mutable",  // Sagittarius, Capricorn, Aquarius, Pisces
	}
	return modalities[p.Sign]
}

// IsHarmonicAspect returns true if the aspect is harmonious
func (a *Aspect) IsHarmonicAspect() bool {
	return a.Type == Sextile || a.Type == Trine
}

// IsChallengingAspect returns true if the aspect is challenging
func (a *Aspect) IsChallengingAspect() bool {
	return a.Type == Square || a.Type == Opposition
}

// GetIntensity returns the intensity of the aspect based on orb
func (a *Aspect) GetIntensity() float64 {
	// Closer to exact aspect = higher intensity
	maxOrb := 8.0
	return 1.0 - (a.Orb / maxOrb)
}