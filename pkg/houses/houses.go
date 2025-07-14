// Package houses provides astrological house system calculations.
package houses

import (
	"fmt"
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/planets"
	"github.com/Qucanft/Qucanft/pkg/zodiac"
)

// HouseSystem represents different house calculation systems
type HouseSystem string

const (
	// Equal houses - each house is exactly 30°
	Equal HouseSystem = "Equal"
	
	// Placidus - most common time-based system
	Placidus HouseSystem = "Placidus"
	
	// Whole Sign - each house occupies an entire zodiac sign
	WholeSign HouseSystem = "WholeSign"
	
	// Koch - alternative time-based system
	Koch HouseSystem = "Koch"
	
	// Campanus - space-based system
	Campanus HouseSystem = "Campanus"
	
	// Regiomontanus - medieval system
	Regiomontanus HouseSystem = "Regiomontanus"
)

// House represents an astrological house
type House struct {
	Number      int
	Name        string
	Theme       string
	Ruler       string
	CuspDegree  float64 // Ecliptic longitude of the house cusp
	Size        float64 // Size of the house in degrees
	Planets     []planets.Planet
	Description string
}

// HouseCalculator handles house calculations
type HouseCalculator struct {
	system HouseSystem
}

// NewHouseCalculator creates a new house calculator with the specified system
func NewHouseCalculator(system HouseSystem) *HouseCalculator {
	return &HouseCalculator{
		system: system,
	}
}

// SetSystem changes the house system
func (hc *HouseCalculator) SetSystem(system HouseSystem) {
	hc.system = system
}

// GetSystem returns the current house system
func (hc *HouseCalculator) GetSystem() HouseSystem {
	return hc.system
}

// CalculateHouseCusps calculates the house cusps for a given time and location
func (hc *HouseCalculator) CalculateHouseCusps(ascendant, midheaven, latitude float64) ([]float64, error) {
	switch hc.system {
	case Equal:
		return hc.calculateEqualHouses(ascendant), nil
	case Placidus:
		return hc.calculatePlacidusHouses(ascendant, midheaven, latitude), nil
	case WholeSign:
		return hc.calculateWholeSignHouses(ascendant), nil
	case Koch:
		return hc.calculateKochHouses(ascendant, midheaven, latitude), nil
	case Campanus:
		return hc.calculateCampanusHouses(ascendant, midheaven, latitude), nil
	case Regiomontanus:
		return hc.calculateRegiomontanusHouses(ascendant, midheaven, latitude), nil
	default:
		return nil, fmt.Errorf("unsupported house system: %s", hc.system)
	}
}

// CalculateHouses calculates complete house information
func (hc *HouseCalculator) CalculateHouses(ascendant, midheaven, latitude float64) ([]House, error) {
	cusps, err := hc.CalculateHouseCusps(ascendant, midheaven, latitude)
	if err != nil {
		return nil, err
	}
	
	houses := make([]House, 12)
	houseInfo := getHouseInformation()
	
	for i := 0; i < 12; i++ {
		nextCusp := cusps[(i+1)%12]
		size := nextCusp - cusps[i]
		if size < 0 {
			size += 360
		}
		
		houses[i] = House{
			Number:      i + 1,
			Name:        houseInfo[i].Name,
			Theme:       houseInfo[i].Theme,
			Ruler:       houseInfo[i].Ruler,
			CuspDegree:  cusps[i],
			Size:        size,
			Planets:     []planets.Planet{},
			Description: houseInfo[i].Description,
		}
	}
	
	return houses, nil
}

// AddPlanetsToHouses assigns planets to houses based on their positions
func (hc *HouseCalculator) AddPlanetsToHouses(houses []House, positions []planets.PlanetaryPosition) []House {
	// Create a copy of houses to avoid modifying the original
	result := make([]House, len(houses))
	copy(result, houses)
	
	// Clear existing planets
	for i := range result {
		result[i].Planets = []planets.Planet{}
	}
	
	// Assign planets to houses
	for _, pos := range positions {
		houseIndex := hc.findHouseForPosition(pos.Coordinates.Longitude, houses)
		if houseIndex >= 0 && houseIndex < len(result) {
			result[houseIndex].Planets = append(result[houseIndex].Planets, pos.Planet)
		}
	}
	
	return result
}

// findHouseForPosition finds which house a given ecliptic longitude belongs to
func (hc *HouseCalculator) findHouseForPosition(longitude float64, houses []House) int {
	longitude = coordinates.NormalizeAngle(longitude)
	
	for i, house := range houses {
		nextHouseIndex := (i + 1) % len(houses)
		nextHouseCusp := houses[nextHouseIndex].CuspDegree
		
		// Handle wrap-around at 360°/0°
		if house.CuspDegree <= nextHouseCusp {
			if longitude >= house.CuspDegree && longitude < nextHouseCusp {
				return i
			}
		} else {
			// House spans across 0° point
			if longitude >= house.CuspDegree || longitude < nextHouseCusp {
				return i
			}
		}
	}
	
	return -1 // Not found (shouldn't happen)
}

// calculateEqualHouses calculates equal house cusps
func (hc *HouseCalculator) calculateEqualHouses(ascendant float64) []float64 {
	cusps := make([]float64, 12)
	
	for i := 0; i < 12; i++ {
		cusps[i] = coordinates.NormalizeAngle(ascendant + float64(i)*30)
	}
	
	return cusps
}

// calculatePlacidusHouses calculates Placidus house cusps
func (hc *HouseCalculator) calculatePlacidusHouses(ascendant, midheaven, latitude float64) []float64 {
	cusps := make([]float64, 12)
	
	// Set the main angles
	cusps[0] = ascendant                                      // 1st house (Ascendant)
	cusps[3] = coordinates.NormalizeAngle(ascendant + 180)    // 4th house (IC)
	cusps[6] = coordinates.NormalizeAngle(ascendant + 180)    // 7th house (Descendant)
	cusps[9] = midheaven                                      // 10th house (MC)
	
	// Calculate intermediate houses using Placidus method
	latRad := latitude * coordinates.DegreesToRadians
	
	// Calculate 2nd and 3rd houses
	for i := 1; i <= 2; i++ {
		t := float64(i) / 3.0
		cusps[i] = hc.calculatePlacidusHouse(ascendant, midheaven, latRad, t)
	}
	
	// Calculate 5th and 6th houses
	for i := 4; i <= 5; i++ {
		t := float64(i-3) / 3.0
		cusps[i] = hc.calculatePlacidusHouse(cusps[3], cusps[6], latRad, t)
	}
	
	// Calculate 8th and 9th houses
	for i := 7; i <= 8; i++ {
		t := float64(i-6) / 3.0
		cusps[i] = hc.calculatePlacidusHouse(cusps[6], cusps[9], latRad, t)
	}
	
	// Calculate 11th and 12th houses
	for i := 10; i <= 11; i++ {
		t := float64(i-9) / 3.0
		cusps[i] = hc.calculatePlacidusHouse(cusps[9], cusps[0], latRad, t)
	}
	
	return cusps
}

// calculatePlacidusHouse calculates a single Placidus house cusp
func (hc *HouseCalculator) calculatePlacidusHouse(start, end, latitude, t float64) float64 {
	// This is a simplified Placidus calculation
	// Real implementation would involve more complex spherical trigonometry
	
	diff := end - start
	if diff < 0 {
		diff += 360
	}
	
	// Apply time-based adjustment
	adjustment := math.Sin(t*math.Pi/2) * math.Tan(latitude) * 5 // Simplified
	
	result := start + diff*t + adjustment
	return coordinates.NormalizeAngle(result)
}

// calculateWholeSignHouses calculates Whole Sign house cusps
func (hc *HouseCalculator) calculateWholeSignHouses(ascendant float64) []float64 {
	cusps := make([]float64, 12)
	
	// Find the zodiac sign of the ascendant
	zc := zodiac.NewZodiacCalculator()
	ascendantZodiac := zc.EclipticToZodiac(ascendant)
	
	// Each house starts at the beginning of a zodiac sign
	startSign := int(ascendantZodiac.Sign.StartDeg / 30)
	
	for i := 0; i < 12; i++ {
		signIndex := (startSign + i) % 12
		cusps[i] = float64(signIndex * 30)
	}
	
	return cusps
}

// calculateKochHouses calculates Koch house cusps
func (hc *HouseCalculator) calculateKochHouses(ascendant, midheaven, latitude float64) []float64 {
	// Koch system is similar to Placidus but with different calculation method
	// This is a simplified implementation
	return hc.calculatePlacidusHouses(ascendant, midheaven, latitude)
}

// calculateCampanusHouses calculates Campanus house cusps
func (hc *HouseCalculator) calculateCampanusHouses(ascendant, midheaven, latitude float64) []float64 {
	// Campanus is a space-based system
	// This is a simplified implementation
	cusps := make([]float64, 12)
	
	// Use equal division as a base and apply spatial adjustments
	base := hc.calculateEqualHouses(ascendant)
	
	for i := 0; i < 12; i++ {
		// Apply spatial adjustment based on latitude
		adjustment := math.Sin(float64(i)*math.Pi/6) * latitude * 0.1
		cusps[i] = coordinates.NormalizeAngle(base[i] + adjustment)
	}
	
	return cusps
}

// calculateRegiomontanusHouses calculates Regiomontanus house cusps
func (hc *HouseCalculator) calculateRegiomontanusHouses(ascendant, midheaven, latitude float64) []float64 {
	// Regiomontanus is a medieval system
	// This is a simplified implementation
	return hc.calculatePlacidusHouses(ascendant, midheaven, latitude)
}

// GetHousePosition returns the house position for a given ecliptic longitude
func (hc *HouseCalculator) GetHousePosition(longitude float64, houses []House) (int, float64, error) {
	houseIndex := hc.findHouseForPosition(longitude, houses)
	if houseIndex < 0 {
		return -1, 0, fmt.Errorf("could not determine house for longitude %.2f", longitude)
	}
	
	// Calculate position within the house (0-1)
	house := houses[houseIndex]
	nextHouseIndex := (houseIndex + 1) % len(houses)
	nextHouseCusp := houses[nextHouseIndex].CuspDegree
	
	var positionInHouse float64
	if house.CuspDegree <= nextHouseCusp {
		positionInHouse = (longitude - house.CuspDegree) / (nextHouseCusp - house.CuspDegree)
	} else {
		// Handle wrap-around
		if longitude >= house.CuspDegree {
			positionInHouse = (longitude - house.CuspDegree) / (360 - house.CuspDegree + nextHouseCusp)
		} else {
			positionInHouse = (360 - house.CuspDegree + longitude) / (360 - house.CuspDegree + nextHouseCusp)
		}
	}
	
	return houseIndex + 1, positionInHouse, nil
}

// HouseInfo represents basic information about a house
type HouseInfo struct {
	Name        string
	Theme       string
	Ruler       string
	Description string
}

// getHouseInformation returns basic information about the 12 houses
func getHouseInformation() []HouseInfo {
	return []HouseInfo{
		{
			Name:        "1st House",
			Theme:       "Self, Identity, Appearance",
			Ruler:       "Mars",
			Description: "The house of self, personality, and how you appear to others",
		},
		{
			Name:        "2nd House", 
			Theme:       "Money, Possessions, Values",
			Ruler:       "Venus",
			Description: "The house of personal resources, money, and material possessions",
		},
		{
			Name:        "3rd House",
			Theme:       "Communication, Siblings, Short Trips",
			Ruler:       "Mercury",
			Description: "The house of communication, learning, and immediate environment",
		},
		{
			Name:        "4th House",
			Theme:       "Home, Family, Roots",
			Ruler:       "Moon",
			Description: "The house of home, family, and emotional foundation",
		},
		{
			Name:        "5th House",
			Theme:       "Creativity, Romance, Children",
			Ruler:       "Sun",
			Description: "The house of creativity, romance, and self-expression",
		},
		{
			Name:        "6th House",
			Theme:       "Work, Health, Daily Routine",
			Ruler:       "Mercury",
			Description: "The house of work, health, and daily responsibilities",
		},
		{
			Name:        "7th House",
			Theme:       "Partnerships, Marriage, Others",
			Ruler:       "Venus",
			Description: "The house of partnerships, marriage, and open enemies",
		},
		{
			Name:        "8th House",
			Theme:       "Transformation, Shared Resources, Death",
			Ruler:       "Mars",
			Description: "The house of transformation, shared resources, and hidden things",
		},
		{
			Name:        "9th House",
			Theme:       "Philosophy, Higher Learning, Travel",
			Ruler:       "Jupiter",
			Description: "The house of higher learning, philosophy, and long-distance travel",
		},
		{
			Name:        "10th House",
			Theme:       "Career, Reputation, Authority",
			Ruler:       "Saturn",
			Description: "The house of career, reputation, and public standing",
		},
		{
			Name:        "11th House",
			Theme:       "Friends, Groups, Hopes",
			Ruler:       "Uranus",
			Description: "The house of friends, groups, and hopes and wishes",
		},
		{
			Name:        "12th House",
			Theme:       "Spirituality, Subconscious, Hidden",
			Ruler:       "Neptune",
			Description: "The house of spirituality, subconscious, and hidden enemies",
		},
	}
}

// String methods
func (h House) String() string {
	planetNames := make([]string, len(h.Planets))
	for i, planet := range h.Planets {
		planetNames[i] = planet.Name
	}
	return fmt.Sprintf("%s (%.1f°): %s - Planets: %v", h.Name, h.CuspDegree, h.Theme, planetNames)
}

func (hs HouseSystem) String() string {
	return string(hs)
}