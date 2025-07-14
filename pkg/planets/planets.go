package planets

import (
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/time"
)

// Planet represents a celestial body
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
	if int(p) < len(names) {
		return names[p]
	}
	return "Unknown"
}

// PlanetPosition represents the position of a planet
type PlanetPosition struct {
	Planet      Planet
	Longitude   coordinates.Angle
	Latitude    coordinates.Angle
	Distance    float64 // in AU
	JulianDate  time.JulianDate
}

// orbitalElements represents the orbital elements of a planet
type orbitalElements struct {
	a  float64 // semi-major axis (AU)
	e  float64 // eccentricity
	i  float64 // inclination (degrees)
	L  float64 // mean longitude (degrees)
	w  float64 // longitude of perihelion (degrees)
	O  float64 // longitude of ascending node (degrees)
}

// getOrbitalElements returns simplified orbital elements for a planet at J2000.0
func getOrbitalElements(planet Planet, t float64) orbitalElements {
	// Simplified orbital elements (J2000.0 epoch)
	// These are approximate values for demonstration
	switch planet {
	case Mercury:
		return orbitalElements{
			a: 0.38709927 + 0.00000037*t,
			e: 0.20563593 + 0.00001906*t,
			i: 7.00497902 - 0.00594749*t,
			L: 252.25032350 + 149472.67411175*t,
			w: 77.45779628 + 0.16047689*t,
			O: 48.33076593 - 0.12534081*t,
		}
	case Venus:
		return orbitalElements{
			a: 0.72333566 + 0.00000390*t,
			e: 0.00677672 - 0.00004107*t,
			i: 3.39467605 - 0.00078890*t,
			L: 181.97909950 + 58517.81538729*t,
			w: 131.60246718 + 0.00268329*t,
			O: 76.67984255 - 0.27769418*t,
		}
	case Mars:
		return orbitalElements{
			a: 1.52371034 + 0.00001847*t,
			e: 0.09339410 + 0.00007882*t,
			i: 1.84969142 - 0.00813131*t,
			L: -4.55343205 + 19140.30268499*t,
			w: -23.94362959 + 0.44441088*t,
			O: 49.55953891 - 0.29257343*t,
		}
	case Jupiter:
		return orbitalElements{
			a: 5.20288700 - 0.00011607*t,
			e: 0.04838624 - 0.00013253*t,
			i: 1.30439695 - 0.00183714*t,
			L: 34.39644051 + 3034.74612775*t,
			w: 14.72847983 + 0.21252668*t,
			O: 100.47390909 + 0.20469106*t,
		}
	case Saturn:
		return orbitalElements{
			a: 9.53667594 - 0.00125060*t,
			e: 0.05386179 - 0.00050991*t,
			i: 2.48599187 + 0.00193609*t,
			L: 49.95424423 + 1222.49362201*t,
			w: 92.59887831 - 0.41897216*t,
			O: 113.66242448 - 0.28867794*t,
		}
	default:
		// Default to Earth-like orbit for unknown planets
		return orbitalElements{
			a: 1.0,
			e: 0.0167,
			i: 0.0,
			L: 0.0,
			w: 0.0,
			O: 0.0,
		}
	}
}

// CalculatePosition calculates the position of a planet at a given Julian Date
func CalculatePosition(planet Planet, jd time.JulianDate) PlanetPosition {
	// Special cases for Sun and Moon
	if planet == Sun {
		return calculateSunPosition(jd)
	}
	if planet == Moon {
		return calculateMoonPosition(jd)
	}
	
	// Time in centuries from J2000.0
	t := jd.CenturiesSinceJ2000()
	
	// Get orbital elements
	elements := getOrbitalElements(planet, t)
	
	// Calculate mean anomaly
	M := elements.L - elements.w
	M = math.Mod(M, 360.0)
	if M < 0 {
		M += 360.0
	}
	
	// Convert to radians
	M_rad := M * math.Pi / 180.0
	
	// Solve Kepler's equation (simplified)
	E := M_rad + elements.e*math.Sin(M_rad)
	
	// Calculate true anomaly
	nu := 2.0 * math.Atan2(
		math.Sqrt(1.0+elements.e)*math.Sin(E/2.0),
		math.Sqrt(1.0-elements.e)*math.Cos(E/2.0),
	)
	
	// Calculate distance
	r := elements.a * (1.0 - elements.e*math.Cos(E))
	
	// Calculate longitude
	longitude := nu*180.0/math.Pi + elements.w
	longitude = math.Mod(longitude, 360.0)
	if longitude < 0 {
		longitude += 360.0
	}
	
	return PlanetPosition{
		Planet:     planet,
		Longitude:  coordinates.Angle(longitude),
		Latitude:   coordinates.Angle(0.0), // Simplified - assuming orbital plane
		Distance:   r,
		JulianDate: jd,
	}
}

// calculateSunPosition calculates the Sun's apparent position
func calculateSunPosition(jd time.JulianDate) PlanetPosition {
	// Simplified solar position calculation
	n := jd.DaysSinceJ2000()
	
	// Mean longitude of the Sun
	L := 280.460 + 0.9856474*n
	L = math.Mod(L, 360.0)
	if L < 0 {
		L += 360.0
	}
	
	// Mean anomaly of the Sun
	g := 357.528 + 0.9856003*n
	g = math.Mod(g, 360.0)
	if g < 0 {
		g += 360.0
	}
	
	// Ecliptic longitude
	lambda := L + 1.915*math.Sin(g*math.Pi/180.0) + 0.020*math.Sin(2.0*g*math.Pi/180.0)
	lambda = math.Mod(lambda, 360.0)
	if lambda < 0 {
		lambda += 360.0
	}
	
	return PlanetPosition{
		Planet:     Sun,
		Longitude:  coordinates.Angle(lambda),
		Latitude:   coordinates.Angle(0.0),
		Distance:   1.0, // 1 AU by definition
		JulianDate: jd,
	}
}

// calculateMoonPosition calculates the Moon's position (simplified)
func calculateMoonPosition(jd time.JulianDate) PlanetPosition {
	// Simplified lunar position calculation
	n := jd.DaysSinceJ2000()
	
	// Mean longitude of the Moon
	L := 218.316 + 13.176396*n
	L = math.Mod(L, 360.0)
	if L < 0 {
		L += 360.0
	}
	
	// Mean anomaly of the Moon
	M := 134.963 + 13.064993*n
	M = math.Mod(M, 360.0)
	if M < 0 {
		M += 360.0
	}
	
	// Mean anomaly of the Sun
	M_sun := 357.528 + 0.9856003*n
	M_sun = math.Mod(M_sun, 360.0)
	if M_sun < 0 {
		M_sun += 360.0
	}
	
	// Ecliptic longitude (simplified)
	lambda := L + 6.289*math.Sin(M*math.Pi/180.0) - 1.274*math.Sin((M-2.0*134.963)*math.Pi/180.0)
	lambda = math.Mod(lambda, 360.0)
	if lambda < 0 {
		lambda += 360.0
	}
	
	// Ecliptic latitude (simplified)
	beta := 5.128 * math.Sin((L-218.316)*math.Pi/180.0)
	
	return PlanetPosition{
		Planet:     Moon,
		Longitude:  coordinates.Angle(lambda),
		Latitude:   coordinates.Angle(beta),
		Distance:   0.002570, // Approximate distance in AU
		JulianDate: jd,
	}
}