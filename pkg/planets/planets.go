// Package planets provides planetary position calculations and celestial body definitions.
package planets

import (
	"fmt"
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	timeutil "github.com/Qucanft/Qucanft/pkg/time"
)

// Planet represents a celestial body with its orbital characteristics
type Planet struct {
	Name   string
	Symbol string
	
	// Orbital elements (simplified for demonstration)
	SemimajorAxis      float64 // AU
	Eccentricity       float64
	Inclination        float64 // degrees
	LongitudeOfNode    float64 // degrees
	ArgumentOfPeriapsis float64 // degrees
	MeanAnomalyAtEpoch float64 // degrees
	MeanMotion         float64 // degrees per day
}

// PlanetaryPosition represents a planet's position at a specific time
type PlanetaryPosition struct {
	Planet      Planet
	Time        timeutil.JulianDay
	Coordinates coordinates.EclipticCoordinates
}

// PlanetaryCalculator handles planetary position calculations
type PlanetaryCalculator struct {
	planets map[string]Planet
}

// NewPlanetaryCalculator creates a new planetary calculator with default planet definitions
func NewPlanetaryCalculator() *PlanetaryCalculator {
	return &PlanetaryCalculator{
		planets: getDefaultPlanets(),
	}
}

// GetPlanet returns a planet by name
func (pc *PlanetaryCalculator) GetPlanet(name string) (Planet, bool) {
	planet, exists := pc.planets[name]
	return planet, exists
}

// GetAllPlanets returns all available planets
func (pc *PlanetaryCalculator) GetAllPlanets() map[string]Planet {
	return pc.planets
}

// CalculatePosition calculates the position of a planet at a given time
func (pc *PlanetaryCalculator) CalculatePosition(planetName string, jd timeutil.JulianDay) (PlanetaryPosition, error) {
	planet, exists := pc.planets[planetName]
	if !exists {
		return PlanetaryPosition{}, fmt.Errorf("planet %s not found", planetName)
	}
	
	// Calculate time since J2000.0 epoch
	tc := timeutil.NewTimeConverter()
	_ = tc.JulianCenturies(jd)
	
	// Calculate mean anomaly
	meanAnomaly := planet.MeanAnomalyAtEpoch + planet.MeanMotion*float64(jd-timeutil.J2000)
	meanAnomaly = coordinates.NormalizeAngle(meanAnomaly)
	
	// Solve Kepler's equation for eccentric anomaly (simplified)
	eccentricAnomaly := solveKeplerEquation(meanAnomaly*coordinates.DegreesToRadians, planet.Eccentricity)
	
	// Calculate true anomaly
	trueAnomaly := 2 * math.Atan2(
		math.Sqrt(1+planet.Eccentricity)*math.Sin(eccentricAnomaly/2),
		math.Sqrt(1-planet.Eccentricity)*math.Cos(eccentricAnomaly/2),
	)
	
	// Calculate heliocentric distance
	distance := planet.SemimajorAxis * (1 - planet.Eccentricity*math.Cos(eccentricAnomaly))
	
	// Calculate position in orbital plane
	argumentOfPeriapsis := planet.ArgumentOfPeriapsis * coordinates.DegreesToRadians
	inclination := planet.Inclination * coordinates.DegreesToRadians
	longitudeOfNode := planet.LongitudeOfNode * coordinates.DegreesToRadians
	
	// Position in orbital plane
	x := distance * math.Cos(trueAnomaly)
	y := distance * math.Sin(trueAnomaly)
	z := 0.0
	
	// Rotate to ecliptic coordinates
	// First rotation: argument of periapsis
	x1 := x*math.Cos(argumentOfPeriapsis) - y*math.Sin(argumentOfPeriapsis)
	y1 := x*math.Sin(argumentOfPeriapsis) + y*math.Cos(argumentOfPeriapsis)
	z1 := z
	
	// Second rotation: inclination
	x2 := x1
	y2 := y1*math.Cos(inclination) - z1*math.Sin(inclination)
	z2 := y1*math.Sin(inclination) + z1*math.Cos(inclination)
	
	// Third rotation: longitude of ascending node
	x3 := x2*math.Cos(longitudeOfNode) - y2*math.Sin(longitudeOfNode)
	y3 := x2*math.Sin(longitudeOfNode) + y2*math.Cos(longitudeOfNode)
	z3 := z2
	
	// Convert to ecliptic longitude and latitude
	longitude := math.Atan2(y3, x3) * coordinates.RadiansToDegrees
	latitude := math.Atan2(z3, math.Sqrt(x3*x3+y3*y3)) * coordinates.RadiansToDegrees
	
	longitude = coordinates.NormalizeAngle(longitude)
	
	return PlanetaryPosition{
		Planet: planet,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: longitude,
			Latitude:  latitude,
			Distance:  distance,
		},
	}, nil
}

// CalculateMultiplePositions calculates positions for multiple planets at once
func (pc *PlanetaryCalculator) CalculateMultiplePositions(planetNames []string, jd timeutil.JulianDay) ([]PlanetaryPosition, error) {
	positions := make([]PlanetaryPosition, 0, len(planetNames))
	
	for _, name := range planetNames {
		pos, err := pc.CalculatePosition(name, jd)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate position for %s: %w", name, err)
		}
		positions = append(positions, pos)
	}
	
	return positions, nil
}

// solveKeplerEquation solves Kepler's equation using Newton's method
func solveKeplerEquation(meanAnomaly, eccentricity float64) float64 {
	// Initial guess
	E := meanAnomaly
	
	// Newton's method iteration
	for i := 0; i < 10; i++ {
		deltaE := (E - eccentricity*math.Sin(E) - meanAnomaly) / (1 - eccentricity*math.Cos(E))
		E -= deltaE
		
		// Check for convergence
		if math.Abs(deltaE) < 1e-10 {
			break
		}
	}
	
	return E
}

// getDefaultPlanets returns the default planet definitions
func getDefaultPlanets() map[string]Planet {
	return map[string]Planet{
		"Sun": {
			Name:               "Sun",
			Symbol:             "☉",
			SemimajorAxis:      0.0, // Special case for Sun
			Eccentricity:       0.0,
			Inclination:        0.0,
			LongitudeOfNode:    0.0,
			ArgumentOfPeriapsis: 0.0,
			MeanAnomalyAtEpoch: 0.0,
			MeanMotion:         0.9856, // degrees per day (approximate)
		},
		"Moon": {
			Name:               "Moon",
			Symbol:             "☽",
			SemimajorAxis:      0.00257, // AU (Earth-Moon distance)
			Eccentricity:       0.0549,
			Inclination:        5.145,
			LongitudeOfNode:    125.1228,
			ArgumentOfPeriapsis: 318.0634,
			MeanAnomalyAtEpoch: 115.3654,
			MeanMotion:         13.1764, // degrees per day
		},
		"Mercury": {
			Name:               "Mercury",
			Symbol:             "☿",
			SemimajorAxis:      0.3871,
			Eccentricity:       0.2056,
			Inclination:        7.005,
			LongitudeOfNode:    48.331,
			ArgumentOfPeriapsis: 29.124,
			MeanAnomalyAtEpoch: 174.796,
			MeanMotion:         4.0923,
		},
		"Venus": {
			Name:               "Venus",
			Symbol:             "♀",
			SemimajorAxis:      0.7233,
			Eccentricity:       0.0067,
			Inclination:        3.395,
			LongitudeOfNode:    76.680,
			ArgumentOfPeriapsis: 54.884,
			MeanAnomalyAtEpoch: 50.115,
			MeanMotion:         1.6021,
		},
		"Mars": {
			Name:               "Mars",
			Symbol:             "♂",
			SemimajorAxis:      1.5237,
			Eccentricity:       0.0934,
			Inclination:        1.850,
			LongitudeOfNode:    49.558,
			ArgumentOfPeriapsis: 286.502,
			MeanAnomalyAtEpoch: 19.373,
			MeanMotion:         0.5240,
		},
		"Jupiter": {
			Name:               "Jupiter",
			Symbol:             "♃",
			SemimajorAxis:      5.2026,
			Eccentricity:       0.0484,
			Inclination:        1.303,
			LongitudeOfNode:    100.464,
			ArgumentOfPeriapsis: 273.867,
			MeanAnomalyAtEpoch: 20.020,
			MeanMotion:         0.0831,
		},
		"Saturn": {
			Name:               "Saturn",
			Symbol:             "♄",
			SemimajorAxis:      9.5549,
			Eccentricity:       0.0555,
			Inclination:        2.485,
			LongitudeOfNode:    113.665,
			ArgumentOfPeriapsis: 339.392,
			MeanAnomalyAtEpoch: 317.020,
			MeanMotion:         0.0334,
		},
		"Uranus": {
			Name:               "Uranus",
			Symbol:             "♅",
			SemimajorAxis:      19.2184,
			Eccentricity:       0.0463,
			Inclination:        0.773,
			LongitudeOfNode:    74.006,
			ArgumentOfPeriapsis: 96.998,
			MeanAnomalyAtEpoch: 142.238,
			MeanMotion:         0.0117,
		},
		"Neptune": {
			Name:               "Neptune",
			Symbol:             "♆",
			SemimajorAxis:      30.1104,
			Eccentricity:       0.0095,
			Inclination:        1.770,
			LongitudeOfNode:    131.784,
			ArgumentOfPeriapsis: 276.336,
			MeanAnomalyAtEpoch: 256.228,
			MeanMotion:         0.0060,
		},
		"Pluto": {
			Name:               "Pluto",
			Symbol:             "♇",
			SemimajorAxis:      39.4821,
			Eccentricity:       0.2488,
			Inclination:        17.16,
			LongitudeOfNode:    110.299,
			ArgumentOfPeriapsis: 113.834,
			MeanAnomalyAtEpoch: 14.882,
			MeanMotion:         0.0040,
		},
	}
}

// CalculateSunPosition calculates the Sun's position (geocentric)
func (pc *PlanetaryCalculator) CalculateSunPosition(jd timeutil.JulianDay) (PlanetaryPosition, error) {
	// Simplified solar position calculation
	tc := timeutil.NewTimeConverter()
	t := tc.JulianCenturies(jd)
	
	// Mean longitude of the Sun
	L0 := 280.46646 + 36000.76983*t + 0.0003032*t*t
	L0 = coordinates.NormalizeAngle(L0)
	
	// Mean anomaly of the Sun
	M := 357.52911 + 35999.05029*t - 0.0001537*t*t
	M = coordinates.NormalizeAngle(M)
	
	// Equation of center
	C := (1.914602 - 0.004817*t - 0.000014*t*t) * math.Sin(M*coordinates.DegreesToRadians)
	C += (0.019993 - 0.000101*t) * math.Sin(2*M*coordinates.DegreesToRadians)
	C += 0.000289 * math.Sin(3*M*coordinates.DegreesToRadians)
	
	// True longitude
	trueLongitude := L0 + C
	trueLongitude = coordinates.NormalizeAngle(trueLongitude)
	
	// Distance (AU)
	distance := 1.000001018 * (1 - 0.01671123*math.Cos(M*coordinates.DegreesToRadians) - 0.00014*math.Cos(2*M*coordinates.DegreesToRadians))
	
	sun := pc.planets["Sun"]
	
	return PlanetaryPosition{
		Planet: sun,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: trueLongitude,
			Latitude:  0.0, // Sun's latitude is always 0 in ecliptic coordinates
			Distance:  distance,
		},
	}, nil
}

// String method for Planet
func (p Planet) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Symbol)
}

// String method for PlanetaryPosition
func (pp PlanetaryPosition) String() string {
	return fmt.Sprintf("%s at %s: %s", pp.Planet.Name, pp.Time.String(), pp.Coordinates.String())
}