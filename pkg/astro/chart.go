package astro

import (
	"time"
	
	"github.com/Qucanft/Qucanft/pkg/aspects"
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/planets"
	astrotime "github.com/Qucanft/Qucanft/pkg/time"
	"github.com/Qucanft/Qucanft/pkg/zodiac"
)

// Chart represents an astrological chart
type Chart struct {
	DateTime      time.Time
	JulianDate    astrotime.JulianDate
	Positions     []planets.PlanetPosition
	Signs         []zodiac.SignInfo
	Aspects       []aspects.Aspect
	Location      Location
}

// Location represents a geographical location
type Location struct {
	Latitude  float64
	Longitude float64
	Timezone  string
}

// NewChart creates a new astrological chart for the given time and location
func NewChart(dateTime time.Time, location Location) *Chart {
	jd := astrotime.ToJulianDate(dateTime)
	
	// Calculate positions for all planets
	allPlanets := []planets.Planet{
		planets.Sun, planets.Moon, planets.Mercury, planets.Venus,
		planets.Mars, planets.Jupiter, planets.Saturn, planets.Uranus,
		planets.Neptune, planets.Pluto,
	}
	
	var positions []planets.PlanetPosition
	var signs []zodiac.SignInfo
	
	for _, planet := range allPlanets {
		position := planets.CalculatePosition(planet, jd)
		positions = append(positions, position)
		
		sign := zodiac.GetSignFromLongitude(position.Longitude)
		signs = append(signs, sign)
	}
	
	// Calculate aspects
	chartAspects := aspects.CalculateAspects(positions)
	
	return &Chart{
		DateTime:   dateTime,
		JulianDate: jd,
		Positions:  positions,
		Signs:      signs,
		Aspects:    chartAspects,
		Location:   location,
	}
}

// GetPlanetPosition returns the position of a specific planet
func (c *Chart) GetPlanetPosition(planet planets.Planet) *planets.PlanetPosition {
	for _, position := range c.Positions {
		if position.Planet == planet {
			return &position
		}
	}
	return nil
}

// GetPlanetSign returns the zodiac sign of a specific planet
func (c *Chart) GetPlanetSign(planet planets.Planet) *zodiac.SignInfo {
	for i, position := range c.Positions {
		if position.Planet == planet {
			return &c.Signs[i]
		}
	}
	return nil
}

// GetAspectsBetween returns all aspects between two planets
func (c *Chart) GetAspectsBetween(planet1, planet2 planets.Planet) []aspects.Aspect {
	var planetAspects []aspects.Aspect
	
	for _, aspect := range c.Aspects {
		if (aspect.Planet1 == planet1 && aspect.Planet2 == planet2) ||
			(aspect.Planet1 == planet2 && aspect.Planet2 == planet1) {
			planetAspects = append(planetAspects, aspect)
		}
	}
	
	return planetAspects
}

// GetAspectsForPlanet returns all aspects involving a specific planet
func (c *Chart) GetAspectsForPlanet(planet planets.Planet) []aspects.Aspect {
	var planetAspects []aspects.Aspect
	
	for _, aspect := range c.Aspects {
		if aspect.Planet1 == planet || aspect.Planet2 == planet {
			planetAspects = append(planetAspects, aspect)
		}
	}
	
	return planetAspects
}

// GetMajorAspects returns only the major aspects in the chart
func (c *Chart) GetMajorAspects() []aspects.Aspect {
	return aspects.GetMajorAspects(c.Aspects)
}

// GetHarmoniousAspects returns only the harmonious aspects in the chart
func (c *Chart) GetHarmoniousAspects() []aspects.Aspect {
	return aspects.GetHarmoniousAspects(c.Aspects)
}

// GetChallengingAspects returns only the challenging aspects in the chart
func (c *Chart) GetChallengingAspects() []aspects.Aspect {
	return aspects.GetChallengingAspects(c.Aspects)
}

// GetElementalBalance returns the distribution of planets across elements
func (c *Chart) GetElementalBalance() map[zodiac.Element]int {
	balance := make(map[zodiac.Element]int)
	
	for _, sign := range c.Signs {
		balance[sign.Element]++
	}
	
	return balance
}

// GetQualityBalance returns the distribution of planets across qualities
func (c *Chart) GetQualityBalance() map[zodiac.Quality]int {
	balance := make(map[zodiac.Quality]int)
	
	for _, sign := range c.Signs {
		balance[sign.Quality]++
	}
	
	return balance
}

// CalculateHousePositions calculates house positions (placeholder for house system)
func (c *Chart) CalculateHousePositions() []coordinates.Angle {
	// This is a simplified placeholder
	// In a real implementation, you would use a house system like Placidus or Equal House
	var houseCusps []coordinates.Angle
	
	// For now, just return equal houses starting from 0 degrees
	for i := 0; i < 12; i++ {
		cusp := coordinates.Angle(float64(i * 30)).Normalize()
		houseCusps = append(houseCusps, cusp)
	}
	
	return houseCusps
}

// GetAngles returns the four angles of the chart (ASC, MC, DESC, IC)
func (c *Chart) GetAngles() map[string]coordinates.Angle {
	// This is a simplified placeholder
	// In a real implementation, you would calculate these based on location and time
	angles := make(map[string]coordinates.Angle)
	
	// Placeholder values
	angles["ASC"] = coordinates.Angle(0)   // Ascendant
	angles["MC"] = coordinates.Angle(90)   // Midheaven
	angles["DESC"] = coordinates.Angle(180) // Descendant
	angles["IC"] = coordinates.Angle(270)   // Imum Coeli
	
	return angles
}

// CompareCharts performs basic chart comparison (synastry)
func CompareCharts(chart1, chart2 *Chart) []aspects.Aspect {
	var comparisonAspects []aspects.Aspect
	
	// Compare each planet in chart1 with each planet in chart2
	for _, pos1 := range chart1.Positions {
		for _, pos2 := range chart2.Positions {
			// Create temporary positions for aspect calculation
			tempPos1 := pos1
			tempPos2 := pos2
			
			// Calculate aspects between the two positions
			positions := []planets.PlanetPosition{tempPos1, tempPos2}
			chartAspects := aspects.CalculateAspects(positions)
			
			comparisonAspects = append(comparisonAspects, chartAspects...)
		}
	}
	
	return comparisonAspects
}

// GetLunarPhase calculates the lunar phase based on Sun and Moon positions
func (c *Chart) GetLunarPhase() (string, coordinates.Angle) {
	sunPos := c.GetPlanetPosition(planets.Sun)
	moonPos := c.GetPlanetPosition(planets.Moon)
	
	if sunPos == nil || moonPos == nil {
		return "Unknown", 0
	}
	
	// Calculate the angle between Sun and Moon
	angle := coordinates.AngularDistance(
		coordinates.EclipticCoordinates{Longitude: sunPos.Longitude, Latitude: sunPos.Latitude},
		coordinates.EclipticCoordinates{Longitude: moonPos.Longitude, Latitude: moonPos.Latitude},
	)
	
	// Determine phase based on angle
	switch {
	case angle < 45:
		return "New Moon", angle
	case angle < 135:
		return "Waxing Moon", angle
	case angle < 225:
		return "Full Moon", angle
	default:
		return "Waning Moon", angle
	}
}