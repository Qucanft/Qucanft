// Package aspects provides planetary aspect calculations and interpretations.
package aspects

import (
	"fmt"
	"math"
	"sort"
	
	"github.com/Qucanft/Qucanft/pkg/planets"
	"github.com/Qucanft/Qucanft/pkg/zodiac"
)

// AspectType represents the type of aspect
type AspectType struct {
	Name        string
	Symbol      string
	Angle       float64
	Orb         float64
	Nature      string // "Harmonious", "Challenging", "Neutral", "Minor"
	Description string
}

// Aspect represents a planetary aspect
type Aspect struct {
	Planet1     planets.Planet
	Planet2     planets.Planet
	Type        AspectType
	Angle       float64
	Orb         float64
	IsApplying  bool
	Strength    float64 // 0-100
	Description string
}

// AspectCalculator handles aspect calculations
type AspectCalculator struct {
	aspectTypes []AspectType
}

// NewAspectCalculator creates a new aspect calculator
func NewAspectCalculator() *AspectCalculator {
	return &AspectCalculator{
		aspectTypes: getAspectTypes(),
	}
}

// GetAspectTypes returns all available aspect types
func (ac *AspectCalculator) GetAspectTypes() []AspectType {
	return ac.aspectTypes
}

// GetAspectTypeByName returns an aspect type by name
func (ac *AspectCalculator) GetAspectTypeByName(name string) (AspectType, bool) {
	for _, aspectType := range ac.aspectTypes {
		if aspectType.Name == name {
			return aspectType, true
		}
	}
	return AspectType{}, false
}

// CalculateAspect calculates the aspect between two planetary positions
func (ac *AspectCalculator) CalculateAspect(pos1, pos2 planets.PlanetaryPosition) *Aspect {
	// Calculate angular separation
	angle := math.Abs(pos1.Coordinates.Longitude - pos2.Coordinates.Longitude)
	if angle > 180 {
		angle = 360 - angle
	}
	
	// Find the closest aspect type
	var closestAspect *AspectType
	var smallestDiff float64 = 999
	
	for _, aspectType := range ac.aspectTypes {
		diff := math.Abs(angle - aspectType.Angle)
		if diff <= aspectType.Orb && diff < smallestDiff {
			smallestDiff = diff
			closestAspect = &aspectType
		}
	}
	
	if closestAspect == nil {
		return nil // No aspect found within orb
	}
	
	// Calculate strength based on orb
	strength := ((closestAspect.Orb - smallestDiff) / closestAspect.Orb) * 100
	
	// Determine if aspect is applying or separating (simplified)
	isApplying := ac.isApplying(pos1, pos2, *closestAspect)
	
	return &Aspect{
		Planet1:     pos1.Planet,
		Planet2:     pos2.Planet,
		Type:        *closestAspect,
		Angle:       angle,
		Orb:         smallestDiff,
		IsApplying:  isApplying,
		Strength:    strength,
		Description: ac.generateDescription(*closestAspect, pos1.Planet, pos2.Planet),
	}
}

// CalculateAllAspects calculates all aspects between a set of planetary positions
func (ac *AspectCalculator) CalculateAllAspects(positions []planets.PlanetaryPosition) []Aspect {
	var aspects []Aspect
	
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			aspect := ac.CalculateAspect(positions[i], positions[j])
			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}
	
	// Sort aspects by strength (strongest first)
	sort.Slice(aspects, func(i, j int) bool {
		return aspects[i].Strength > aspects[j].Strength
	})
	
	return aspects
}

// GetAspectsByPlanet returns all aspects involving a specific planet
func (ac *AspectCalculator) GetAspectsByPlanet(aspects []Aspect, planetName string) []Aspect {
	var planetAspects []Aspect
	
	for _, aspect := range aspects {
		if aspect.Planet1.Name == planetName || aspect.Planet2.Name == planetName {
			planetAspects = append(planetAspects, aspect)
		}
	}
	
	return planetAspects
}

// GetAspectsByType returns all aspects of a specific type
func (ac *AspectCalculator) GetAspectsByType(aspects []Aspect, aspectType string) []Aspect {
	var typeAspects []Aspect
	
	for _, aspect := range aspects {
		if aspect.Type.Name == aspectType {
			typeAspects = append(typeAspects, aspect)
		}
	}
	
	return typeAspects
}

// GetAspectsByNature returns all aspects of a specific nature
func (ac *AspectCalculator) GetAspectsByNature(aspects []Aspect, nature string) []Aspect {
	var natureAspects []Aspect
	
	for _, aspect := range aspects {
		if aspect.Type.Nature == nature {
			natureAspects = append(natureAspects, aspect)
		}
	}
	
	return natureAspects
}

// GetStrongestAspects returns the strongest aspects up to a limit
func (ac *AspectCalculator) GetStrongestAspects(aspects []Aspect, limit int) []Aspect {
	if len(aspects) <= limit {
		return aspects
	}
	
	// Sort by strength (already sorted in CalculateAllAspects)
	return aspects[:limit]
}

// CalculateAspectPattern detects aspect patterns like Grand Trine, T-Square, etc.
func (ac *AspectCalculator) CalculateAspectPattern(positions []planets.PlanetaryPosition) []AspectPattern {
	aspects := ac.CalculateAllAspects(positions)
	var patterns []AspectPattern
	
	// Check for Grand Trine (3 planets in trine aspect)
	patterns = append(patterns, ac.findGrandTrines(aspects, positions)...)
	
	// Check for T-Square (2 squares and 1 opposition)
	patterns = append(patterns, ac.findTSquares(aspects, positions)...)
	
	// Check for Grand Cross (4 planets in square/opposition)
	patterns = append(patterns, ac.findGrandCrosses(aspects, positions)...)
	
	// Check for Stellium (3+ planets in same sign)
	patterns = append(patterns, ac.findStelliums(positions)...)
	
	return patterns
}

// AspectPattern represents a configuration of multiple aspects
type AspectPattern struct {
	Name        string
	Planets     []planets.Planet
	Aspects     []Aspect
	Description string
	Strength    float64
}

// isApplying determines if an aspect is applying (getting closer) or separating
func (ac *AspectCalculator) isApplying(pos1, pos2 planets.PlanetaryPosition, aspectType AspectType) bool {
	// This is a simplified calculation
	// In reality, you'd need to consider orbital velocities and directions
	
	// For now, assume the faster planet is applying to the slower one
	fasterPlanet := ac.getFasterPlanet(pos1.Planet, pos2.Planet)
	
	if fasterPlanet == pos1.Planet {
		return pos1.Coordinates.Longitude < pos2.Coordinates.Longitude
	}
	
	return pos2.Coordinates.Longitude < pos1.Coordinates.Longitude
}

// getFasterPlanet returns the planet with faster orbital motion
func (ac *AspectCalculator) getFasterPlanet(planet1, planet2 planets.Planet) planets.Planet {
	// Order from fastest to slowest
	order := []string{"Moon", "Sun", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto"}
	
	index1 := ac.getPlanetIndex(planet1.Name, order)
	index2 := ac.getPlanetIndex(planet2.Name, order)
	
	if index1 < index2 {
		return planet1
	}
	
	return planet2
}

// getPlanetIndex returns the index of a planet in the order array
func (ac *AspectCalculator) getPlanetIndex(planetName string, order []string) int {
	for i, name := range order {
		if name == planetName {
			return i
		}
	}
	return 999 // Unknown planet, treat as slowest
}

// generateDescription generates a description for an aspect
func (ac *AspectCalculator) generateDescription(aspectType AspectType, planet1, planet2 planets.Planet) string {
	return fmt.Sprintf("%s %s %s: %s", planet1.Name, aspectType.Name, planet2.Name, aspectType.Description)
}

// findGrandTrines finds Grand Trine patterns
func (ac *AspectCalculator) findGrandTrines(aspects []Aspect, positions []planets.PlanetaryPosition) []AspectPattern {
	var patterns []AspectPattern
	
	// Find all trine aspects
	trines := ac.GetAspectsByType(aspects, "Trine")
	
	// Check for three planets forming a Grand Trine
	for i := 0; i < len(trines); i++ {
		for j := i + 1; j < len(trines); j++ {
			for k := j + 1; k < len(trines); k++ {
				if ac.formsGrandTrine(trines[i], trines[j], trines[k]) {
					planets := []planets.Planet{trines[i].Planet1, trines[i].Planet2, trines[j].Planet2}
					aspectsInPattern := []Aspect{trines[i], trines[j], trines[k]}
					
					patterns = append(patterns, AspectPattern{
						Name:        "Grand Trine",
						Planets:     planets,
						Aspects:     aspectsInPattern,
						Description: "A harmonious triangle of energy flow between three planets",
						Strength:    ac.calculatePatternStrength(aspectsInPattern),
					})
				}
			}
		}
	}
	
	return patterns
}

// findTSquares finds T-Square patterns
func (ac *AspectCalculator) findTSquares(aspects []Aspect, positions []planets.PlanetaryPosition) []AspectPattern {
	var patterns []AspectPattern
	
	// Find all square and opposition aspects
	squares := ac.GetAspectsByType(aspects, "Square")
	oppositions := ac.GetAspectsByType(aspects, "Opposition")
	
	// Check for T-Square pattern (2 squares + 1 opposition)
	for _, opp := range oppositions {
		for _, sq1 := range squares {
			for _, sq2 := range squares {
				if ac.formsTSquare(opp, sq1, sq2) {
					planets := []planets.Planet{opp.Planet1, opp.Planet2, sq1.Planet2}
					aspectsInPattern := []Aspect{opp, sq1, sq2}
					
					patterns = append(patterns, AspectPattern{
						Name:        "T-Square",
						Planets:     planets,
						Aspects:     aspectsInPattern,
						Description: "A challenging configuration creating tension and drive",
						Strength:    ac.calculatePatternStrength(aspectsInPattern),
					})
				}
			}
		}
	}
	
	return patterns
}

// findGrandCrosses finds Grand Cross patterns
func (ac *AspectCalculator) findGrandCrosses(aspects []Aspect, positions []planets.PlanetaryPosition) []AspectPattern {
	var patterns []AspectPattern
	
	// Find all square and opposition aspects
	squares := ac.GetAspectsByType(aspects, "Square")
	oppositions := ac.GetAspectsByType(aspects, "Opposition")
	
	// Check for Grand Cross pattern (4 squares + 2 oppositions)
	if len(squares) >= 4 && len(oppositions) >= 2 {
		// This is a simplified check - a full implementation would be more complex
		for _, opp1 := range oppositions {
			for _, opp2 := range oppositions {
				if ac.formsGrandCross(opp1, opp2, squares) {
					planets := []planets.Planet{opp1.Planet1, opp1.Planet2, opp2.Planet1, opp2.Planet2}
					aspectsInPattern := append([]Aspect{opp1, opp2}, squares[:4]...)
					
					patterns = append(patterns, AspectPattern{
						Name:        "Grand Cross",
						Planets:     planets,
						Aspects:     aspectsInPattern,
						Description: "A powerful cross configuration creating maximum tension and potential",
						Strength:    ac.calculatePatternStrength(aspectsInPattern),
					})
				}
			}
		}
	}
	
	return patterns
}

// findStelliums finds Stellium patterns
func (ac *AspectCalculator) findStelliums(positions []planets.PlanetaryPosition) []AspectPattern {
	var patterns []AspectPattern
	
	// Group planets by zodiac sign
	zc := zodiac.NewZodiacCalculator()
	signGroups := make(map[string][]planets.Planet)
	
	for _, pos := range positions {
		zodiacPos := zc.EclipticToZodiac(pos.Coordinates.Longitude)
		signGroups[zodiacPos.Sign.Name] = append(signGroups[zodiacPos.Sign.Name], pos.Planet)
	}
	
	// Find signs with 3+ planets (Stellium)
	for signName, planetsInSign := range signGroups {
		if len(planetsInSign) >= 3 {
			patterns = append(patterns, AspectPattern{
				Name:        "Stellium",
				Planets:     planetsInSign,
				Aspects:     []Aspect{}, // No specific aspects, just proximity
				Description: fmt.Sprintf("A concentration of %d planets in %s", len(planetsInSign), signName),
				Strength:    float64(len(planetsInSign)) * 20, // Strength based on number of planets
			})
		}
	}
	
	return patterns
}

// Helper functions for pattern detection
func (ac *AspectCalculator) formsGrandTrine(aspect1, aspect2, aspect3 Aspect) bool {
	// Check if three aspects form a closed triangle of trines
	// This is a simplified check
	return true // Placeholder
}

func (ac *AspectCalculator) formsTSquare(opposition, square1, square2 Aspect) bool {
	// Check if aspects form a T-Square pattern
	// This is a simplified check
	return true // Placeholder
}

func (ac *AspectCalculator) formsGrandCross(opp1, opp2 Aspect, squares []Aspect) bool {
	// Check if aspects form a Grand Cross pattern
	// This is a simplified check
	return len(squares) >= 4 // Placeholder
}

func (ac *AspectCalculator) calculatePatternStrength(aspects []Aspect) float64 {
	totalStrength := 0.0
	for _, aspect := range aspects {
		totalStrength += aspect.Strength
	}
	return totalStrength / float64(len(aspects))
}

// getAspectTypes returns the standard astrological aspects
func getAspectTypes() []AspectType {
	return []AspectType{
		{
			Name:        "Conjunction",
			Symbol:      "☌",
			Angle:       0,
			Orb:         8,
			Nature:      "Neutral",
			Description: "Union of energies, intensity, new beginnings",
		},
		{
			Name:        "Sextile",
			Symbol:      "⚹",
			Angle:       60,
			Orb:         6,
			Nature:      "Harmonious",
			Description: "Opportunity, cooperation, creative potential",
		},
		{
			Name:        "Square",
			Symbol:      "□",
			Angle:       90,
			Orb:         8,
			Nature:      "Challenging",
			Description: "Tension, conflict, catalyst for growth",
		},
		{
			Name:        "Trine",
			Symbol:      "△",
			Angle:       120,
			Orb:         8,
			Nature:      "Harmonious",
			Description: "Flow, ease, natural talent, harmony",
		},
		{
			Name:        "Opposition",
			Symbol:      "☍",
			Angle:       180,
			Orb:         8,
			Nature:      "Challenging",
			Description: "Polarity, awareness, balance needed",
		},
		{
			Name:        "Semisextile",
			Symbol:      "⚺",
			Angle:       30,
			Orb:         2,
			Nature:      "Minor",
			Description: "Mild connection, subtle influence",
		},
		{
			Name:        "Semisquare",
			Symbol:      "∠",
			Angle:       45,
			Orb:         2,
			Nature:      "Minor",
			Description: "Mild friction, minor irritation",
		},
		{
			Name:        "Sesquiquadrate",
			Symbol:      "⚼",
			Angle:       135,
			Orb:         2,
			Nature:      "Minor",
			Description: "Adjustment needed, minor challenge",
		},
		{
			Name:        "Quincunx",
			Symbol:      "⚻",
			Angle:       150,
			Orb:         2,
			Nature:      "Minor",
			Description: "Adjustment, adaptation, awkward energy",
		},
		{
			Name:        "Quintile",
			Symbol:      "Q",
			Angle:       72,
			Orb:         1,
			Nature:      "Minor",
			Description: "Creativity, special talent, artistic expression",
		},
		{
			Name:        "Biquintile",
			Symbol:      "bQ",
			Angle:       144,
			Orb:         1,
			Nature:      "Minor",
			Description: "Enhanced creativity, artistic mastery",
		},
	}
}

// String methods
func (at AspectType) String() string {
	return fmt.Sprintf("%s (%s) - %s at %g° (orb: %g°)", at.Name, at.Symbol, at.Nature, at.Angle, at.Orb)
}

func (a Aspect) String() string {
	return fmt.Sprintf("%s %s %s (%.1f°, orb: %.1f°, strength: %.1f%%)", 
		a.Planet1.Name, a.Type.Name, a.Planet2.Name, a.Angle, a.Orb, a.Strength)
}

func (ap AspectPattern) String() string {
	planetNames := make([]string, len(ap.Planets))
	for i, planet := range ap.Planets {
		planetNames[i] = planet.Name
	}
	return fmt.Sprintf("%s: %v (strength: %.1f%%)", ap.Name, planetNames, ap.Strength)
}