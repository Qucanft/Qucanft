package aspects

import (
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/planets"
)

// AspectType represents different types of aspects
type AspectType int

const (
	Conjunction AspectType = iota
	Sextile
	Square
	Trine
	Opposition
	Quincunx
	Semisextile
	Semisquare
	Sesquiquadrate
)

// String returns the name of the aspect
func (a AspectType) String() string {
	names := []string{
		"Conjunction", "Sextile", "Square", "Trine", "Opposition",
		"Quincunx", "Semisextile", "Semisquare", "Sesquiquadrate",
	}
	if int(a) < len(names) {
		return names[a]
	}
	return "Unknown"
}

// AspectInfo contains information about an aspect
type AspectInfo struct {
	Type         AspectType
	Angle        coordinates.Angle
	Orb          coordinates.Angle
	Symbol       string
	Meaning      string
	IsHarmonious bool
}

// Aspect represents an aspect between two planets
type Aspect struct {
	Planet1      planets.Planet
	Planet2      planets.Planet
	Position1    coordinates.Angle
	Position2    coordinates.Angle
	Type         AspectType
	Orb          coordinates.Angle
	ExactAngle   coordinates.Angle
	IsApplying   bool
	IsHarmonious bool
}

// GetAspectInfo returns information about an aspect type
func GetAspectInfo(aspectType AspectType) AspectInfo {
	aspects := []AspectInfo{
		{
			Type:         Conjunction,
			Angle:        0,
			Orb:          8,
			Symbol:       "☌",
			Meaning:      "Unity, blending of energies",
			IsHarmonious: true,
		},
		{
			Type:         Sextile,
			Angle:        60,
			Orb:          6,
			Symbol:       "⚹",
			Meaning:      "Opportunity, cooperation",
			IsHarmonious: true,
		},
		{
			Type:         Square,
			Angle:        90,
			Orb:          8,
			Symbol:       "□",
			Meaning:      "Tension, challenge, action",
			IsHarmonious: false,
		},
		{
			Type:         Trine,
			Angle:        120,
			Orb:          8,
			Symbol:       "△",
			Meaning:      "Harmony, ease, flow",
			IsHarmonious: true,
		},
		{
			Type:         Opposition,
			Angle:        180,
			Orb:          8,
			Symbol:       "☍",
			Meaning:      "Polarity, awareness, balance",
			IsHarmonious: false,
		},
		{
			Type:         Quincunx,
			Angle:        150,
			Orb:          3,
			Symbol:       "⚻",
			Meaning:      "Adjustment, adaptation",
			IsHarmonious: false,
		},
		{
			Type:         Semisextile,
			Angle:        30,
			Orb:          2,
			Symbol:       "⚺",
			Meaning:      "Subtle connection, growth",
			IsHarmonious: true,
		},
		{
			Type:         Semisquare,
			Angle:        45,
			Orb:          2,
			Symbol:       "∠",
			Meaning:      "Irritation, minor stress",
			IsHarmonious: false,
		},
		{
			Type:         Sesquiquadrate,
			Angle:        135,
			Orb:          2,
			Symbol:       "⚼",
			Meaning:      "Crisis, breakthrough",
			IsHarmonious: false,
		},
	}
	
	if int(aspectType) < len(aspects) {
		return aspects[aspectType]
	}
	
	return AspectInfo{Type: aspectType, Angle: 0, Orb: 0, Symbol: "?", Meaning: "Unknown"}
}

// CalculateAspects calculates all aspects between a list of planet positions
func CalculateAspects(positions []planets.PlanetPosition) []Aspect {
	var aspects []Aspect
	
	// Check all pairs of planets
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			aspect := calculateAspectBetween(positions[i], positions[j])
			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}
	
	return aspects
}

// calculateAspectBetween calculates the aspect between two planet positions
func calculateAspectBetween(pos1, pos2 planets.PlanetPosition) *Aspect {
	// Calculate angular separation
	angle1 := pos1.Longitude
	angle2 := pos2.Longitude
	
	separation := calculateAngularSeparation(angle1, angle2)
	
	// Check each aspect type
	for aspectType := Conjunction; aspectType <= Sesquiquadrate; aspectType++ {
		aspectInfo := GetAspectInfo(aspectType)
		
		// Calculate difference from exact aspect
		diff := math.Abs(float64(separation) - float64(aspectInfo.Angle))
		
		// Check if within orb
		if diff <= float64(aspectInfo.Orb) {
			return &Aspect{
				Planet1:      pos1.Planet,
				Planet2:      pos2.Planet,
				Position1:    angle1,
				Position2:    angle2,
				Type:         aspectType,
				Orb:          coordinates.Angle(diff),
				ExactAngle:   aspectInfo.Angle,
				IsApplying:   isApplying(pos1, pos2, aspectInfo.Angle),
				IsHarmonious: aspectInfo.IsHarmonious,
			}
		}
	}
	
	return nil
}

// calculateAngularSeparation calculates the angular separation between two angles
func calculateAngularSeparation(angle1, angle2 coordinates.Angle) coordinates.Angle {
	diff := math.Abs(float64(angle1) - float64(angle2))
	
	// Use the smaller angle
	if diff > 180 {
		diff = 360 - diff
	}
	
	return coordinates.Angle(diff)
}

// isApplying determines if an aspect is applying (getting closer) or separating
func isApplying(pos1, pos2 planets.PlanetPosition, exactAngle coordinates.Angle) bool {
	// This is a simplified calculation
	// In reality, you'd need to consider the planets' daily motions
	
	// For now, assume faster planets are applying to slower ones
	fasterPlanet := getFasterPlanet(pos1.Planet, pos2.Planet)
	
	var fasterPos, slowerPos coordinates.Angle
	if fasterPlanet == pos1.Planet {
		fasterPos = pos1.Longitude
		slowerPos = pos2.Longitude
	} else {
		fasterPos = pos2.Longitude
		slowerPos = pos1.Longitude
	}
	
	// Calculate if the faster planet is approaching the aspect
	targetPos := slowerPos + exactAngle
	if targetPos >= 360 {
		targetPos -= 360
	}
	
	// Simplified logic: if faster planet is within 180 degrees behind the target, it's applying
	diff := float64(targetPos) - float64(fasterPos)
	if diff < 0 {
		diff += 360
	}
	
	return diff < 180
}

// getFasterPlanet returns the planet with faster motion
func getFasterPlanet(p1, p2 planets.Planet) planets.Planet {
	// Order from fastest to slowest motion
	order := []planets.Planet{
		planets.Moon, planets.Mercury, planets.Venus, planets.Sun,
		planets.Mars, planets.Jupiter, planets.Saturn, planets.Uranus,
		planets.Neptune, planets.Pluto,
	}
	
	for _, planet := range order {
		if planet == p1 {
			return p1
		}
		if planet == p2 {
			return p2
		}
	}
	
	return p1 // Default
}

// GetMajorAspects returns only the major aspects (conjunction, sextile, square, trine, opposition)
func GetMajorAspects(aspects []Aspect) []Aspect {
	var majorAspects []Aspect
	
	for _, aspect := range aspects {
		if aspect.Type <= Opposition {
			majorAspects = append(majorAspects, aspect)
		}
	}
	
	return majorAspects
}

// GetHarmoniousAspects returns only the harmonious aspects
func GetHarmoniousAspects(aspects []Aspect) []Aspect {
	var harmoniousAspects []Aspect
	
	for _, aspect := range aspects {
		if aspect.IsHarmonious {
			harmoniousAspects = append(harmoniousAspects, aspect)
		}
	}
	
	return harmoniousAspects
}

// GetChallengingAspects returns only the challenging aspects
func GetChallengingAspects(aspects []Aspect) []Aspect {
	var challengingAspects []Aspect
	
	for _, aspect := range aspects {
		if !aspect.IsHarmonious {
			challengingAspects = append(challengingAspects, aspect)
		}
	}
	
	return challengingAspects
}