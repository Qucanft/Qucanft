package coordinates

import (
	"math"
)

// Angle represents an angle in degrees
type Angle float64

// ToRadians converts degrees to radians
func (a Angle) ToRadians() float64 {
	return float64(a) * math.Pi / 180.0
}

// Normalize normalizes an angle to the range [0, 360)
func (a Angle) Normalize() Angle {
	angle := math.Mod(float64(a), 360.0)
	if angle < 0 {
		angle += 360.0
	}
	return Angle(angle)
}

// EquatorialCoordinates represents right ascension and declination
type EquatorialCoordinates struct {
	RightAscension Angle // in degrees
	Declination    Angle // in degrees
}

// EclipticCoordinates represents longitude and latitude on the ecliptic
type EclipticCoordinates struct {
	Longitude Angle // in degrees
	Latitude  Angle // in degrees
}

// ToEcliptic converts equatorial coordinates to ecliptic coordinates
func (eq EquatorialCoordinates) ToEcliptic(obliquity Angle) EclipticCoordinates {
	// Convert to radians
	ra := eq.RightAscension.ToRadians()
	dec := eq.Declination.ToRadians()
	obl := obliquity.ToRadians()
	
	// Convert using standard formulas
	sinLon := math.Sin(ra)*math.Cos(obl) + math.Tan(dec)*math.Sin(obl)
	cosLon := math.Cos(ra)
	longitude := math.Atan2(sinLon, cosLon)
	
	sinLat := math.Sin(dec)*math.Cos(obl) - math.Cos(dec)*math.Sin(obl)*math.Sin(ra)
	latitude := math.Asin(sinLat)
	
	return EclipticCoordinates{
		Longitude: Angle(longitude * 180.0 / math.Pi).Normalize(),
		Latitude:  Angle(latitude * 180.0 / math.Pi),
	}
}

// ToEquatorial converts ecliptic coordinates to equatorial coordinates
func (ec EclipticCoordinates) ToEquatorial(obliquity Angle) EquatorialCoordinates {
	// Convert to radians
	lon := ec.Longitude.ToRadians()
	lat := ec.Latitude.ToRadians()
	obl := obliquity.ToRadians()
	
	// Convert using standard formulas
	sinRA := math.Sin(lon)*math.Cos(obl) - math.Tan(lat)*math.Sin(obl)
	cosRA := math.Cos(lon)
	ra := math.Atan2(sinRA, cosRA)
	
	sinDec := math.Sin(lat)*math.Cos(obl) + math.Cos(lat)*math.Sin(obl)*math.Sin(lon)
	dec := math.Asin(sinDec)
	
	return EquatorialCoordinates{
		RightAscension: Angle(ra * 180.0 / math.Pi).Normalize(),
		Declination:    Angle(dec * 180.0 / math.Pi),
	}
}

// GetObliquity calculates the obliquity of the ecliptic for a given Julian Date
func GetObliquity(jd float64) Angle {
	// Time in centuries from J2000.0
	t := (jd - 2451545.0) / 36525.0
	
	// Mean obliquity in arcseconds
	obliquity := 23.43929111 - 0.013004166*t - 0.0000001639*t*t + 0.0000005036*t*t*t
	
	return Angle(obliquity)
}

// AngularDistance calculates the angular distance between two points
func AngularDistance(coord1, coord2 EclipticCoordinates) Angle {
	// Convert to radians
	lon1 := coord1.Longitude.ToRadians()
	lat1 := coord1.Latitude.ToRadians()
	lon2 := coord2.Longitude.ToRadians()
	lat2 := coord2.Latitude.ToRadians()
	
	// Use spherical law of cosines
	cosDistance := math.Sin(lat1)*math.Sin(lat2) + 
		math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)
	
	// Clamp to avoid numerical errors
	if cosDistance > 1.0 {
		cosDistance = 1.0
	} else if cosDistance < -1.0 {
		cosDistance = -1.0
	}
	
	distance := math.Acos(cosDistance)
	return Angle(distance * 180.0 / math.Pi)
}