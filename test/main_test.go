package main

import (
	"testing"
	"time"

	"github.com/Qucanft/Qucanft/pkg/astrology"
	"github.com/Qucanft/Qucanft/pkg/visualization"
)

func TestAstrologyTypes(t *testing.T) {
	// Test zodiac signs
	if astrology.Aries.String() != "Aries" {
		t.Errorf("Expected Aries, got %s", astrology.Aries.String())
	}
	
	// Test planets
	if astrology.Sun.String() != "Sun" {
		t.Errorf("Expected Sun, got %s", astrology.Sun.String())
	}
	
	// Test aspects
	if astrology.Trine.Angle() != 120.0 {
		t.Errorf("Expected 120.0 degrees for Trine, got %f", astrology.Trine.Angle())
	}
}

func TestChartGeneration(t *testing.T) {
	generator := astrology.NewChartGenerator()
	chart := generator.GenerateChart(time.Now())
	
	// Check that we have the expected number of planets
	if len(chart.Planets) != 10 {
		t.Errorf("Expected 10 planets, got %d", len(chart.Planets))
	}
	
	// Check that all planets are present
	expectedPlanets := []astrology.Planet{
		astrology.Sun, astrology.Moon, astrology.Mercury, astrology.Venus, astrology.Mars,
		astrology.Jupiter, astrology.Saturn, astrology.Uranus, astrology.Neptune, astrology.Pluto,
	}
	
	for _, expectedPlanet := range expectedPlanets {
		found := false
		for _, planetPos := range chart.Planets {
			if planetPos.Planet == expectedPlanet {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected planet %s not found in chart", expectedPlanet.String())
		}
	}
}

func TestVisualizationGeneration(t *testing.T) {
	generator := astrology.NewChartGenerator()
	chart := generator.GenerateChart(time.Now())
	
	config := visualization.GetDefaultConfig()
	config.Width = 400
	config.Height = 300
	
	artGen := visualization.NewArtGenerator(config)
	img, err := artGen.GenerateVisualization(chart)
	
	if err != nil {
		t.Errorf("Error generating visualization: %v", err)
	}
	
	if img == nil {
		t.Error("Generated image is nil")
	}
	
	bounds := img.Bounds()
	if bounds.Max.X != 400 || bounds.Max.Y != 300 {
		t.Errorf("Expected image size 400x300, got %dx%d", bounds.Max.X, bounds.Max.Y)
	}
}

func TestPlanetPositionMethods(t *testing.T) {
	position := astrology.PlanetPosition{
		Planet:     astrology.Sun,
		Degree:     125.5,
		Sign:       astrology.Leo,
		House:      astrology.FifthHouse,
		Retrograde: false,
	}
	
	// Test GetZodiacDegree
	expected := 5.5 // 125.5 % 30 = 5.5
	if position.GetZodiacDegree() != expected {
		t.Errorf("Expected zodiac degree %f, got %f", expected, position.GetZodiacDegree())
	}
	
	// Test GetElementalEnergy
	if position.GetElementalEnergy() != "Fire" {
		t.Errorf("Expected Fire element for Leo, got %s", position.GetElementalEnergy())
	}
	
	// Test GetModalityEnergy
	if position.GetModalityEnergy() != "Fixed" {
		t.Errorf("Expected Fixed modality for Leo, got %s", position.GetModalityEnergy())
	}
}

func TestAspectMethods(t *testing.T) {
	aspect := astrology.Aspect{
		Planet1: astrology.Sun,
		Planet2: astrology.Moon,
		Type:    astrology.Trine,
		Angle:   120.0,
		Orb:     2.0,
	}
	
	// Test IsHarmonicAspect
	if !aspect.IsHarmonicAspect() {
		t.Error("Expected Trine to be harmonic")
	}
	
	// Test IsChallengingAspect
	if aspect.IsChallengingAspect() {
		t.Error("Expected Trine to not be challenging")
	}
	
	// Test GetIntensity
	intensity := aspect.GetIntensity()
	if intensity < 0 || intensity > 1 {
		t.Errorf("Expected intensity between 0 and 1, got %f", intensity)
	}
}

func TestColorSchemes(t *testing.T) {
	config := visualization.GetDefaultConfig()
	
	// Test different color schemes
	schemes := []visualization.ColorScheme{
		visualization.Cosmic,
		visualization.Earthy,
		visualization.Oceanic,
		visualization.Sunset,
	}
	
	for _, scheme := range schemes {
		config.ColorScheme = scheme
		artGen := visualization.NewArtGenerator(config)
		
		// Test that we can get colors without panicking
		generator := astrology.NewChartGenerator()
		chart := generator.GenerateChart(time.Now())
		
		_, err := artGen.GenerateVisualization(chart)
		if err != nil {
			t.Errorf("Error generating visualization with scheme %d: %v", scheme, err)
		}
	}
}

func TestArtStyles(t *testing.T) {
	config := visualization.GetDefaultConfig()
	
	// Test different art styles
	styles := []visualization.ArtStyle{
		visualization.Mandala,
		visualization.Geometric,
		visualization.Organic,
		visualization.Minimalist,
	}
	
	for _, style := range styles {
		config.Style = style
		artGen := visualization.NewArtGenerator(config)
		
		// Test that we can generate art without panicking
		generator := astrology.NewChartGenerator()
		chart := generator.GenerateChart(time.Now())
		
		_, err := artGen.GenerateVisualization(chart)
		if err != nil {
			t.Errorf("Error generating visualization with style %d: %v", style, err)
		}
	}
}