package planets

import (
	"testing"
	"math"
	
	timeutil "github.com/Qucanft/Qucanft/pkg/time"
)

func TestPlanetaryCalculator(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	// Test getting planets
	planets := pc.GetAllPlanets()
	expectedPlanets := []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto"}
	
	for _, planetName := range expectedPlanets {
		if _, exists := planets[planetName]; !exists {
			t.Errorf("Planet %s not found", planetName)
		}
	}
	
	// Test getting specific planet
	mars, exists := pc.GetPlanet("Mars")
	if !exists {
		t.Error("Mars planet not found")
	}
	
	if mars.Name != "Mars" || mars.Symbol != "♂" {
		t.Errorf("Mars planet data incorrect: got %s (%s)", mars.Name, mars.Symbol)
	}
	
	// Test getting non-existent planet
	_, exists = pc.GetPlanet("NonExistentPlanet")
	if exists {
		t.Error("Non-existent planet should not be found")
	}
}

func TestPlanetaryPositionCalculation(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	// Test calculation for J2000.0 epoch
	jd := timeutil.J2000
	
	// Test Mars position calculation
	position, err := pc.CalculatePosition("Mars", jd)
	if err != nil {
		t.Errorf("Error calculating Mars position: %v", err)
	}
	
	// Check that position is valid
	if position.Planet.Name != "Mars" {
		t.Errorf("Expected Mars, got %s", position.Planet.Name)
	}
	
	if position.Time != jd {
		t.Errorf("Expected JD %.6f, got %.6f", jd, position.Time)
	}
	
	// Check longitude is in valid range
	if position.Coordinates.Longitude < 0 || position.Coordinates.Longitude >= 360 {
		t.Errorf("Longitude out of range: %.6f", position.Coordinates.Longitude)
	}
	
	// Check latitude is in valid range
	if position.Coordinates.Latitude < -90 || position.Coordinates.Latitude > 90 {
		t.Errorf("Latitude out of range: %.6f", position.Coordinates.Latitude)
	}
	
	// Check distance is positive
	if position.Coordinates.Distance <= 0 {
		t.Errorf("Distance should be positive: %.6f", position.Coordinates.Distance)
	}
}

func TestSunPositionCalculation(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	// Test Sun position at J2000.0
	jd := timeutil.J2000
	position, err := pc.CalculateSunPosition(jd)
	if err != nil {
		t.Errorf("Error calculating Sun position: %v", err)
	}
	
	// Check that position is valid
	if position.Planet.Name != "Sun" {
		t.Errorf("Expected Sun, got %s", position.Planet.Name)
	}
	
	// Sun's ecliptic latitude should be 0
	if math.Abs(position.Coordinates.Latitude) > 0.001 {
		t.Errorf("Sun's ecliptic latitude should be 0, got %.6f", position.Coordinates.Latitude)
	}
	
	// Check longitude is in valid range
	if position.Coordinates.Longitude < 0 || position.Coordinates.Longitude >= 360 {
		t.Errorf("Longitude out of range: %.6f", position.Coordinates.Longitude)
	}
	
	// Distance should be approximately 1 AU
	if math.Abs(position.Coordinates.Distance-1.0) > 0.1 {
		t.Errorf("Sun distance should be ~1 AU, got %.6f", position.Coordinates.Distance)
	}
}

func TestMultiplePlanetCalculation(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	jd := timeutil.J2000
	planetNames := []string{"Sun", "Moon", "Mercury", "Venus", "Mars"}
	
	positions, err := pc.CalculateMultiplePositions(planetNames, jd)
	if err != nil {
		t.Errorf("Error calculating multiple positions: %v", err)
	}
	
	if len(positions) != len(planetNames) {
		t.Errorf("Expected %d positions, got %d", len(planetNames), len(positions))
	}
	
	// Check that all planets are included
	for i, position := range positions {
		if position.Planet.Name != planetNames[i] {
			t.Errorf("Expected planet %s, got %s", planetNames[i], position.Planet.Name)
		}
	}
}

func TestMultiplePlanetCalculationError(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	jd := timeutil.J2000
	planetNames := []string{"Sun", "NonExistentPlanet", "Mars"}
	
	_, err := pc.CalculateMultiplePositions(planetNames, jd)
	if err == nil {
		t.Error("Expected error for non-existent planet, got nil")
	}
}

func TestPlanetPositionCalculationError(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	jd := timeutil.J2000
	
	_, err := pc.CalculatePosition("NonExistentPlanet", jd)
	if err == nil {
		t.Error("Expected error for non-existent planet, got nil")
	}
}

func TestSolveKeplerEquation(t *testing.T) {
	testCases := []struct {
		meanAnomaly  float64
		eccentricity float64
		tolerance    float64
	}{
		{0.0, 0.0, 0.001},        // Circular orbit
		{math.Pi/2, 0.1, 0.001},  // Low eccentricity
		{math.Pi, 0.5, 0.001},    // High eccentricity
		{3*math.Pi/2, 0.9, 0.001}, // Very high eccentricity
	}
	
	for _, test := range testCases {
		E := solveKeplerEquation(test.meanAnomaly, test.eccentricity)
		
		// Verify Kepler's equation: E - e*sin(E) = M
		calculated := E - test.eccentricity*math.Sin(E)
		if math.Abs(calculated-test.meanAnomaly) > test.tolerance {
			t.Errorf("Kepler equation solution failed: M=%.6f, e=%.6f, E=%.6f, verification=%.6f", 
				test.meanAnomaly, test.eccentricity, E, calculated)
		}
	}
}

func TestPlanetOrbitalElements(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	// Test that all planets have valid orbital elements
	planets := pc.GetAllPlanets()
	
	for name, planet := range planets {
		// Skip Sun as it has special orbital elements
		if name == "Sun" {
			continue
		}
		
		// Check semi-major axis is positive
		if planet.SemimajorAxis <= 0 {
			t.Errorf("Planet %s has invalid semi-major axis: %.6f", name, planet.SemimajorAxis)
		}
		
		// Check eccentricity is in valid range [0, 1)
		if planet.Eccentricity < 0 || planet.Eccentricity >= 1 {
			t.Errorf("Planet %s has invalid eccentricity: %.6f", name, planet.Eccentricity)
		}
		
		// Check inclination is in reasonable range
		if planet.Inclination < 0 || planet.Inclination > 180 {
			t.Errorf("Planet %s has invalid inclination: %.6f", name, planet.Inclination)
		}
		
		// Check mean motion is positive
		if planet.MeanMotion <= 0 {
			t.Errorf("Planet %s has invalid mean motion: %.6f", name, planet.MeanMotion)
		}
	}
}

func TestPlanetPositionConsistency(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	// Test position consistency over time
	jd1 := timeutil.J2000
	jd2 := jd1.Add(1.0) // One day later
	
	position1, err := pc.CalculatePosition("Mars", jd1)
	if err != nil {
		t.Errorf("Error calculating Mars position at JD1: %v", err)
	}
	
	position2, err := pc.CalculatePosition("Mars", jd2)
	if err != nil {
		t.Errorf("Error calculating Mars position at JD2: %v", err)
	}
	
	// Position should change over time
	if position1.Coordinates.Longitude == position2.Coordinates.Longitude {
		t.Error("Planet position should change over time")
	}
	
	// Change should be reasonable (not too large)
	longitudeDiff := math.Abs(position2.Coordinates.Longitude - position1.Coordinates.Longitude)
	if longitudeDiff > 180 {
		longitudeDiff = 360 - longitudeDiff
	}
	
	if longitudeDiff > 10 { // More than 10 degrees per day is unrealistic
		t.Errorf("Planet position change too large: %.6f degrees per day", longitudeDiff)
	}
}

func TestPlanetStringMethods(t *testing.T) {
	pc := NewPlanetaryCalculator()
	
	mars, _ := pc.GetPlanet("Mars")
	str := mars.String()
	if str == "" {
		t.Error("Planet String() returned empty string")
	}
	
	// Test PlanetaryPosition string method
	jd := timeutil.J2000
	position, _ := pc.CalculatePosition("Mars", jd)
	str = position.String()
	if str == "" {
		t.Error("PlanetaryPosition String() returned empty string")
	}
}

func BenchmarkCalculatePosition(b *testing.B) {
	pc := NewPlanetaryCalculator()
	jd := timeutil.J2000
	
	for i := 0; i < b.N; i++ {
		pc.CalculatePosition("Mars", jd)
	}
}

func BenchmarkCalculateSunPosition(b *testing.B) {
	pc := NewPlanetaryCalculator()
	jd := timeutil.J2000
	
	for i := 0; i < b.N; i++ {
		pc.CalculateSunPosition(jd)
	}
}

func BenchmarkSolveKeplerEquation(b *testing.B) {
	meanAnomaly := math.Pi / 2
	eccentricity := 0.1
	
	for i := 0; i < b.N; i++ {
		solveKeplerEquation(meanAnomaly, eccentricity)
	}
}