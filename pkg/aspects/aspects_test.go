package aspects

import (
	"testing"
	"math"
	
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/planets"
	timeutil "github.com/Qucanft/Qucanft/pkg/time"
)

func TestAspectCalculator(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Test getting aspect types
	aspectTypes := ac.GetAspectTypes()
	if len(aspectTypes) == 0 {
		t.Error("No aspect types found")
	}
	
	// Test major aspects are present
	majorAspects := []string{"Conjunction", "Sextile", "Square", "Trine", "Opposition"}
	for _, aspectName := range majorAspects {
		aspectType, exists := ac.GetAspectTypeByName(aspectName)
		if !exists {
			t.Errorf("Major aspect %s not found", aspectName)
		}
		
		if aspectType.Name != aspectName {
			t.Errorf("Aspect name mismatch: expected %s, got %s", aspectName, aspectType.Name)
		}
	}
	
	// Test non-existent aspect
	_, exists := ac.GetAspectTypeByName("NonExistentAspect")
	if exists {
		t.Error("Non-existent aspect should not be found")
	}
}

func TestAspectTypeProperties(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Test conjunction
	conjunction, _ := ac.GetAspectTypeByName("Conjunction")
	if conjunction.Angle != 0 || conjunction.Symbol != "☌" {
		t.Errorf("Conjunction properties incorrect: angle=%.1f, symbol=%s", conjunction.Angle, conjunction.Symbol)
	}
	
	// Test square
	square, _ := ac.GetAspectTypeByName("Square")
	if square.Angle != 90 || square.Nature != "Challenging" {
		t.Errorf("Square properties incorrect: angle=%.1f, nature=%s", square.Angle, square.Nature)
	}
	
	// Test trine
	trine, _ := ac.GetAspectTypeByName("Trine")
	if trine.Angle != 120 || trine.Nature != "Harmonious" {
		t.Errorf("Trine properties incorrect: angle=%.1f, nature=%s", trine.Angle, trine.Nature)
	}
	
	// Test opposition
	opposition, _ := ac.GetAspectTypeByName("Opposition")
	if opposition.Angle != 180 || opposition.Nature != "Challenging" {
		t.Errorf("Opposition properties incorrect: angle=%.1f, nature=%s", opposition.Angle, opposition.Nature)
	}
}

func TestCalculateAspect(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Create test planetary positions
	jd := timeutil.J2000
	
	// Mars at 0°
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	marsPos := planets.PlanetaryPosition{
		Planet: mars,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 0.0,
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	// Jupiter at 90° (square aspect)
	jupiter := planets.Planet{Name: "Jupiter", Symbol: "♃"}
	jupiterPos := planets.PlanetaryPosition{
		Planet: jupiter,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 90.0,
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	// Calculate aspect
	aspect := ac.CalculateAspect(marsPos, jupiterPos)
	if aspect == nil {
		t.Error("Expected aspect to be found")
	}
	
	if aspect.Type.Name != "Square" {
		t.Errorf("Expected Square aspect, got %s", aspect.Type.Name)
	}
	
	if math.Abs(aspect.Angle-90.0) > 1.0 {
		t.Errorf("Expected angle ~90°, got %.1f°", aspect.Angle)
	}
	
	if aspect.Strength < 80 { // Should be strong since it's exact
		t.Errorf("Expected strong aspect, got strength %.1f", aspect.Strength)
	}
}

func TestCalculateAspectNoAspect(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Create test planetary positions with no valid aspect
	jd := timeutil.J2000
	
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	marsPos := planets.PlanetaryPosition{
		Planet: mars,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 0.0,
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	jupiter := planets.Planet{Name: "Jupiter", Symbol: "♃"}
	jupiterPos := planets.PlanetaryPosition{
		Planet: jupiter,
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 50.0, // No major aspect at 50°
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	// Calculate aspect
	aspect := ac.CalculateAspect(marsPos, jupiterPos)
	if aspect != nil {
		t.Errorf("Expected no aspect, got %s", aspect.Type.Name)
	}
}

func TestCalculateAllAspects(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Create test planetary positions
	jd := timeutil.J2000
	
	positions := []planets.PlanetaryPosition{
		{
			Planet: planets.Planet{Name: "Sun", Symbol: "☉"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 0.0,
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
		{
			Planet: planets.Planet{Name: "Moon", Symbol: "☽"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 90.0, // Square to Sun
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
		{
			Planet: planets.Planet{Name: "Mars", Symbol: "♂"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 180.0, // Opposition to Sun
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
	}
	
	aspects := ac.CalculateAllAspects(positions)
	
	// Should find 2 aspects: Sun-Moon square, Sun-Mars opposition
	if len(aspects) < 2 {
		t.Errorf("Expected at least 2 aspects, got %d", len(aspects))
	}
	
	// Aspects should be sorted by strength (strongest first)
	if len(aspects) > 1 {
		if aspects[0].Strength < aspects[1].Strength {
			t.Error("Aspects should be sorted by strength (strongest first)")
		}
	}
}

func TestGetAspectsByPlanet(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Create test aspects
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	
	square, _ := ac.GetAspectTypeByName("Square")
	trine, _ := ac.GetAspectTypeByName("Trine")
	
	aspects := []Aspect{
		{
			Planet1: sun,
			Planet2: moon,
			Type:    square,
			Angle:   90.0,
			Strength: 85.0,
		},
		{
			Planet1: sun,
			Planet2: mars,
			Type:    trine,
			Angle:   120.0,
			Strength: 80.0,
		},
		{
			Planet1: moon,
			Planet2: mars,
			Type:    square,
			Angle:   90.0,
			Strength: 75.0,
		},
	}
	
	// Get aspects involving Sun
	sunAspects := ac.GetAspectsByPlanet(aspects, "Sun")
	if len(sunAspects) != 2 {
		t.Errorf("Expected 2 Sun aspects, got %d", len(sunAspects))
	}
	
	// Get aspects involving Moon
	moonAspects := ac.GetAspectsByPlanet(aspects, "Moon")
	if len(moonAspects) != 2 {
		t.Errorf("Expected 2 Moon aspects, got %d", len(moonAspects))
	}
	
	// Get aspects involving non-existent planet
	unknownAspects := ac.GetAspectsByPlanet(aspects, "Unknown")
	if len(unknownAspects) != 0 {
		t.Errorf("Expected 0 unknown planet aspects, got %d", len(unknownAspects))
	}
}

func TestGetAspectsByType(t *testing.T) {
	ac := NewAspectCalculator()
	
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	
	square, _ := ac.GetAspectTypeByName("Square")
	trine, _ := ac.GetAspectTypeByName("Trine")
	
	aspects := []Aspect{
		{
			Planet1: sun,
			Planet2: moon,
			Type:    square,
			Angle:   90.0,
		},
		{
			Planet1: sun,
			Planet2: mars,
			Type:    trine,
			Angle:   120.0,
		},
		{
			Planet1: moon,
			Planet2: mars,
			Type:    square,
			Angle:   90.0,
		},
	}
	
	// Get square aspects
	squareAspects := ac.GetAspectsByType(aspects, "Square")
	if len(squareAspects) != 2 {
		t.Errorf("Expected 2 square aspects, got %d", len(squareAspects))
	}
	
	// Get trine aspects
	trineAspects := ac.GetAspectsByType(aspects, "Trine")
	if len(trineAspects) != 1 {
		t.Errorf("Expected 1 trine aspect, got %d", len(trineAspects))
	}
}

func TestGetAspectsByNature(t *testing.T) {
	ac := NewAspectCalculator()
	
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	
	square, _ := ac.GetAspectTypeByName("Square")
	trine, _ := ac.GetAspectTypeByName("Trine")
	
	aspects := []Aspect{
		{
			Planet1: sun,
			Planet2: moon,
			Type:    square,
			Angle:   90.0,
		},
		{
			Planet1: sun,
			Planet2: mars,
			Type:    trine,
			Angle:   120.0,
		},
	}
	
	// Get challenging aspects
	challengingAspects := ac.GetAspectsByNature(aspects, "Challenging")
	if len(challengingAspects) != 1 {
		t.Errorf("Expected 1 challenging aspect, got %d", len(challengingAspects))
	}
	
	// Get harmonious aspects
	harmoniousAspects := ac.GetAspectsByNature(aspects, "Harmonious")
	if len(harmoniousAspects) != 1 {
		t.Errorf("Expected 1 harmonious aspect, got %d", len(harmoniousAspects))
	}
}

func TestGetStrongestAspects(t *testing.T) {
	ac := NewAspectCalculator()
	
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	
	square, _ := ac.GetAspectTypeByName("Square")
	
	aspects := []Aspect{
		{
			Planet1:  sun,
			Planet2:  mars,
			Type:     square,
			Strength: 95.0,
		},
		{
			Planet1:  sun,
			Planet2:  moon,
			Type:     square,
			Strength: 85.0,
		},
		{
			Planet1:  moon,
			Planet2:  mars,
			Type:     square,
			Strength: 75.0,
		},
	}
	
	// Get top 2 strongest aspects
	strongest := ac.GetStrongestAspects(aspects, 2)
	if len(strongest) != 2 {
		t.Errorf("Expected 2 strongest aspects, got %d", len(strongest))
	}
	
	// Should be sorted by strength
	if strongest[0].Strength < strongest[1].Strength {
		t.Error("Strongest aspects should be sorted by strength")
	}
	
	// Test with limit larger than available aspects
	allStrongest := ac.GetStrongestAspects(aspects, 10)
	if len(allStrongest) != len(aspects) {
		t.Errorf("Expected %d aspects, got %d", len(aspects), len(allStrongest))
	}
}

func TestGetFasterPlanet(t *testing.T) {
	ac := NewAspectCalculator()
	
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	mars := planets.Planet{Name: "Mars", Symbol: "♂"}
	jupiter := planets.Planet{Name: "Jupiter", Symbol: "♃"}
	
	// Sun should be faster than Mars
	faster := ac.getFasterPlanet(sun, mars)
	if faster.Name != "Sun" {
		t.Errorf("Expected Sun to be faster than Mars, got %s", faster.Name)
	}
	
	// Mars should be faster than Jupiter
	faster = ac.getFasterPlanet(mars, jupiter)
	if faster.Name != "Mars" {
		t.Errorf("Expected Mars to be faster than Jupiter, got %s", faster.Name)
	}
	
	// Test with Moon (fastest)
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	faster = ac.getFasterPlanet(moon, jupiter)
	if faster.Name != "Moon" {
		t.Errorf("Expected Moon to be faster than Jupiter, got %s", faster.Name)
	}
}

func TestAspectStringMethods(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Test AspectType string method
	square, _ := ac.GetAspectTypeByName("Square")
	str := square.String()
	if str == "" {
		t.Error("AspectType String() returned empty string")
	}
	
	// Test Aspect string method
	sun := planets.Planet{Name: "Sun", Symbol: "☉"}
	moon := planets.Planet{Name: "Moon", Symbol: "☽"}
	
	aspect := Aspect{
		Planet1:  sun,
		Planet2:  moon,
		Type:     square,
		Angle:    90.0,
		Orb:      2.0,
		Strength: 85.0,
	}
	
	str = aspect.String()
	if str == "" {
		t.Error("Aspect String() returned empty string")
	}
}

func TestAspectPatternDetection(t *testing.T) {
	ac := NewAspectCalculator()
	
	// Create test positions for pattern detection
	jd := timeutil.J2000
	
	positions := []planets.PlanetaryPosition{
		{
			Planet: planets.Planet{Name: "Sun", Symbol: "☉"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 0.0,
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
		{
			Planet: planets.Planet{Name: "Moon", Symbol: "☽"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 120.0, // Trine to Sun
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
		{
			Planet: planets.Planet{Name: "Mars", Symbol: "♂"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{
				Longitude: 240.0, // Trine to both Sun and Moon
				Latitude:  0.0,
				Distance:  1.0,
			},
		},
	}
	
	patterns := ac.CalculateAspectPattern(positions)
	
	// Should find at least one pattern
	if len(patterns) == 0 {
		t.Error("Expected to find at least one aspect pattern")
	}
	
	// Check pattern properties
	for _, pattern := range patterns {
		if pattern.Name == "" {
			t.Error("Pattern name should not be empty")
		}
		
		if len(pattern.Planets) == 0 {
			t.Error("Pattern should have planets")
		}
		
		if pattern.Strength < 0 || pattern.Strength > 100 {
			t.Errorf("Pattern strength should be 0-100, got %.1f", pattern.Strength)
		}
	}
}

func BenchmarkCalculateAspect(b *testing.B) {
	ac := NewAspectCalculator()
	
	jd := timeutil.J2000
	
	marsPos := planets.PlanetaryPosition{
		Planet: planets.Planet{Name: "Mars", Symbol: "♂"},
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 0.0,
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	jupiterPos := planets.PlanetaryPosition{
		Planet: planets.Planet{Name: "Jupiter", Symbol: "♃"},
		Time:   jd,
		Coordinates: coordinates.EclipticCoordinates{
			Longitude: 90.0,
			Latitude:  0.0,
			Distance:  1.0,
		},
	}
	
	for i := 0; i < b.N; i++ {
		ac.CalculateAspect(marsPos, jupiterPos)
	}
}

func BenchmarkCalculateAllAspects(b *testing.B) {
	ac := NewAspectCalculator()
	
	jd := timeutil.J2000
	
	positions := []planets.PlanetaryPosition{
		{
			Planet: planets.Planet{Name: "Sun", Symbol: "☉"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{Longitude: 0.0, Latitude: 0.0, Distance: 1.0},
		},
		{
			Planet: planets.Planet{Name: "Moon", Symbol: "☽"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{Longitude: 90.0, Latitude: 0.0, Distance: 1.0},
		},
		{
			Planet: planets.Planet{Name: "Mars", Symbol: "♂"},
			Time:   jd,
			Coordinates: coordinates.EclipticCoordinates{Longitude: 180.0, Latitude: 0.0, Distance: 1.0},
		},
	}
	
	for i := 0; i < b.N; i++ {
		ac.CalculateAllAspects(positions)
	}
}