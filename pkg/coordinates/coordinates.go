// Package coordinates provides coordinate system transformations for astrological calculations.
package coordinates

import (
	"fmt"
	"math"
)

// Constants for coordinate calculations
const (
	// DegreesToRadians converts degrees to radians
	DegreesToRadians = math.Pi / 180.0
	
	// RadiansToDegrees converts radians to degrees
	RadiansToDegrees = 180.0 / math.Pi
	
	// J2000Obliquity is the obliquity of the ecliptic at J2000.0 epoch in degrees
	J2000Obliquity = 23.4392911
)

// EquatorialCoordinates represents coordinates in the equatorial system
type EquatorialCoordinates struct {
	// RightAscension in degrees (0-360)
	RightAscension float64
	
	// Declination in degrees (-90 to +90)
	Declination float64
	
	// Distance in AU (optional, default 1.0)
	Distance float64
}

// EclipticCoordinates represents coordinates in the ecliptic system
type EclipticCoordinates struct {
	// Longitude in degrees (0-360)
	Longitude float64
	
	// Latitude in degrees (-90 to +90)
	Latitude float64
	
	// Distance in AU (optional, default 1.0)
	Distance float64
}

// HorizontalCoordinates represents coordinates in the horizontal system
type HorizontalCoordinates struct {
	// Azimuth in degrees (0-360, measured from North)
	Azimuth float64
	
	// Altitude in degrees (-90 to +90)
	Altitude float64
}

// GalacticCoordinates represents coordinates in the galactic system
type GalacticCoordinates struct {
	// Longitude in degrees (0-360)
	Longitude float64
	
	// Latitude in degrees (-90 to +90)
	Latitude float64
}

// CoordinateTransformer handles transformations between coordinate systems
type CoordinateTransformer struct {
	// Obliquity of the ecliptic in degrees
	obliquity float64
}

// NewCoordinateTransformer creates a new CoordinateTransformer with default obliquity
func NewCoordinateTransformer() *CoordinateTransformer {
	return &CoordinateTransformer{
		obliquity: J2000Obliquity,
	}
}

// NewCoordinateTransformerWithObliquity creates a CoordinateTransformer with custom obliquity
func NewCoordinateTransformerWithObliquity(obliquity float64) *CoordinateTransformer {
	return &CoordinateTransformer{
		obliquity: obliquity,
	}
}

// SetObliquity sets the obliquity of the ecliptic
func (ct *CoordinateTransformer) SetObliquity(obliquity float64) {
	ct.obliquity = obliquity
}

// GetObliquity returns the current obliquity of the ecliptic
func (ct *CoordinateTransformer) GetObliquity() float64 {
	return ct.obliquity
}

// EquatorialToEcliptic converts equatorial coordinates to ecliptic coordinates
func (ct *CoordinateTransformer) EquatorialToEcliptic(eq EquatorialCoordinates) EclipticCoordinates {
	// Convert to radians
	ra := eq.RightAscension * DegreesToRadians
	dec := eq.Declination * DegreesToRadians
	eps := ct.obliquity * DegreesToRadians
	
	// Calculate ecliptic longitude
	y := math.Sin(ra) * math.Cos(eps) + math.Tan(dec) * math.Sin(eps)
	x := math.Cos(ra)
	longitude := math.Atan2(y, x) * RadiansToDegrees
	
	// Calculate ecliptic latitude
	latitude := math.Asin(math.Sin(dec)*math.Cos(eps) - math.Cos(dec)*math.Sin(eps)*math.Sin(ra)) * RadiansToDegrees
	
	// Normalize longitude to 0-360 degrees
	longitude = normalizeAngle(longitude)
	
	return EclipticCoordinates{
		Longitude: longitude,
		Latitude:  latitude,
		Distance:  eq.Distance,
	}
}

// EclipticToEquatorial converts ecliptic coordinates to equatorial coordinates
func (ct *CoordinateTransformer) EclipticToEquatorial(ec EclipticCoordinates) EquatorialCoordinates {
	// Convert to radians
	lon := ec.Longitude * DegreesToRadians
	lat := ec.Latitude * DegreesToRadians
	eps := ct.obliquity * DegreesToRadians
	
	// Calculate right ascension
	y := math.Sin(lon) * math.Cos(eps) - math.Tan(lat) * math.Sin(eps)
	x := math.Cos(lon)
	ra := math.Atan2(y, x) * RadiansToDegrees
	
	// Calculate declination
	dec := math.Asin(math.Sin(lat)*math.Cos(eps) + math.Cos(lat)*math.Sin(eps)*math.Sin(lon)) * RadiansToDegrees
	
	// Normalize right ascension to 0-360 degrees
	ra = normalizeAngle(ra)
	
	return EquatorialCoordinates{
		RightAscension: ra,
		Declination:    dec,
		Distance:       ec.Distance,
	}
}

// EquatorialToHorizontal converts equatorial coordinates to horizontal coordinates
func (ct *CoordinateTransformer) EquatorialToHorizontal(eq EquatorialCoordinates, lst, latitude float64) HorizontalCoordinates {
	// Convert to radians
	ra := eq.RightAscension * DegreesToRadians
	dec := eq.Declination * DegreesToRadians
	lstRad := lst * DegreesToRadians
	latRad := latitude * DegreesToRadians
	
	// Calculate hour angle
	hourAngle := lstRad - ra
	
	// Calculate altitude
	altitude := math.Asin(math.Sin(dec)*math.Sin(latRad) + math.Cos(dec)*math.Cos(latRad)*math.Cos(hourAngle))
	
	// Calculate azimuth
	y := -math.Sin(hourAngle)
	x := math.Tan(dec)*math.Cos(latRad) - math.Sin(latRad)*math.Cos(hourAngle)
	azimuth := math.Atan2(y, x)
	
	// Convert to degrees and normalize
	altitudeDeg := altitude * RadiansToDegrees
	azimuthDeg := normalizeAngle(azimuth * RadiansToDegrees)
	
	return HorizontalCoordinates{
		Azimuth:  azimuthDeg,
		Altitude: altitudeDeg,
	}
}

// HorizontalToEquatorial converts horizontal coordinates to equatorial coordinates
func (ct *CoordinateTransformer) HorizontalToEquatorial(hz HorizontalCoordinates, lst, latitude float64) EquatorialCoordinates {
	// Convert to radians
	az := hz.Azimuth * DegreesToRadians
	alt := hz.Altitude * DegreesToRadians
	lstRad := lst * DegreesToRadians
	latRad := latitude * DegreesToRadians
	
	// Calculate declination
	dec := math.Asin(math.Sin(alt)*math.Sin(latRad) + math.Cos(alt)*math.Cos(latRad)*math.Cos(az))
	
	// Calculate hour angle
	y := -math.Sin(az)
	x := math.Tan(alt)*math.Cos(latRad) - math.Sin(latRad)*math.Cos(az)
	hourAngle := math.Atan2(y, x)
	
	// Calculate right ascension
	ra := lstRad - hourAngle
	
	// Convert to degrees and normalize
	decDeg := dec * RadiansToDegrees
	raDeg := normalizeAngle(ra * RadiansToDegrees)
	
	return EquatorialCoordinates{
		RightAscension: raDeg,
		Declination:    decDeg,
		Distance:       1.0,
	}
}

// AngularSeparation calculates the angular separation between two points
func (ct *CoordinateTransformer) AngularSeparation(coord1, coord2 EquatorialCoordinates) float64 {
	// Convert to radians
	ra1 := coord1.RightAscension * DegreesToRadians
	dec1 := coord1.Declination * DegreesToRadians
	ra2 := coord2.RightAscension * DegreesToRadians
	dec2 := coord2.Declination * DegreesToRadians
	
	// Use the haversine formula
	deltaRA := ra2 - ra1
	
	a := math.Sin(dec1)*math.Sin(dec2) + math.Cos(dec1)*math.Cos(dec2)*math.Cos(deltaRA)
	
	// Clamp to prevent numerical errors
	if a > 1.0 {
		a = 1.0
	} else if a < -1.0 {
		a = -1.0
	}
	
	separation := math.Acos(a) * RadiansToDegrees
	
	return separation
}

// PositionAngle calculates the position angle from coord1 to coord2
func (ct *CoordinateTransformer) PositionAngle(coord1, coord2 EquatorialCoordinates) float64 {
	// Convert to radians
	ra1 := coord1.RightAscension * DegreesToRadians
	dec1 := coord1.Declination * DegreesToRadians
	ra2 := coord2.RightAscension * DegreesToRadians
	dec2 := coord2.Declination * DegreesToRadians
	
	deltaRA := ra2 - ra1
	
	y := math.Sin(deltaRA)
	x := math.Cos(dec1)*math.Tan(dec2) - math.Sin(dec1)*math.Cos(deltaRA)
	
	pa := math.Atan2(y, x) * RadiansToDegrees
	
	return normalizeAngle(pa)
}

// normalizeAngle normalizes an angle to the range [0, 360) degrees
func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360.0)
	if angle < 0 {
		angle += 360.0
	}
	return angle
}

// NormalizeAngle is a public function to normalize angles
func NormalizeAngle(angle float64) float64 {
	return normalizeAngle(angle)
}

// AngleDifference calculates the shortest angular difference between two angles
func AngleDifference(angle1, angle2 float64) float64 {
	diff := angle2 - angle1
	if diff > 180 {
		diff -= 360
	} else if diff < -180 {
		diff += 360
	}
	return diff
}

// String methods for coordinate types
func (eq EquatorialCoordinates) String() string {
	return fmt.Sprintf("RA: %.6f°, Dec: %.6f°, Dist: %.6f AU", eq.RightAscension, eq.Declination, eq.Distance)
}

func (ec EclipticCoordinates) String() string {
	return fmt.Sprintf("Lon: %.6f°, Lat: %.6f°, Dist: %.6f AU", ec.Longitude, ec.Latitude, ec.Distance)
}

func (hz HorizontalCoordinates) String() string {
	return fmt.Sprintf("Az: %.6f°, Alt: %.6f°", hz.Azimuth, hz.Altitude)
}

func (gc GalacticCoordinates) String() string {
	return fmt.Sprintf("Lon: %.6f°, Lat: %.6f°", gc.Longitude, gc.Latitude)
}