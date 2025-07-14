// Package main demonstrates the usage of the Qucanft astrological calculations library.
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/Qucanft/Qucanft/pkg/aspects"
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/houses"
	"github.com/Qucanft/Qucanft/pkg/planets"
	timeutil "github.com/Qucanft/Qucanft/pkg/time"
	"github.com/Qucanft/Qucanft/pkg/zodiac"
)

func main() {
	fmt.Println("=== Qucanft Astrological Calculations Demo ===")
	fmt.Println()
	
	// Example birth data
	birthTime := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	latitude := 40.7128  // New York City
	longitude := -74.0060
	
	// Convert to Julian Day
	tc := timeutil.NewTimeConverter()
	jd := tc.ToJulianDay(birthTime)
	
	fmt.Printf("Birth Time: %s\n", birthTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Julian Day: %s\n", jd.String())
	fmt.Printf("Location: %.4f°N, %.4f°W\n", latitude, longitude)
	fmt.Println()
	
	// 1. Calculate planetary positions
	fmt.Println("=== PLANETARY POSITIONS ===")
	pc := planets.NewPlanetaryCalculator()
	
	planetNames := []string{"Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn"}
	positions, err := pc.CalculateMultiplePositions(planetNames, jd)
	if err != nil {
		log.Fatal(err)
	}
	
	// 2. Convert to zodiac positions
	fmt.Println("Planet Positions:")
	zc := zodiac.NewZodiacCalculator()
	
	for _, pos := range positions {
		zodiacPos := zc.EclipticToZodiac(pos.Coordinates.Longitude)
		formatted := zc.FormatZodiacPosition(zodiacPos)
		fmt.Printf("  %s: %s\n", pos.Planet.Name, formatted)
	}
	fmt.Println()
	
	// 3. Calculate house positions
	fmt.Println("=== HOUSE SYSTEM ===")
	
	// Calculate ascendant (simplified - using Sun position + 90°)
	ascendant := coordinates.NormalizeAngle(positions[0].Coordinates.Longitude + 90)
	midheaven := coordinates.NormalizeAngle(ascendant + 90)
	
	fmt.Printf("Ascendant: %.2f° (%s)\n", ascendant, zc.FormatZodiacPosition(zc.EclipticToZodiac(ascendant)))
	fmt.Printf("Midheaven: %.2f° (%s)\n", midheaven, zc.FormatZodiacPosition(zc.EclipticToZodiac(midheaven)))
	fmt.Println()
	
	// Calculate houses using Equal House system
	hc := houses.NewHouseCalculator(houses.Equal)
	houseList, err := hc.CalculateHouses(ascendant, midheaven, latitude)
	if err != nil {
		log.Fatal(err)
	}
	
	// Add planets to houses
	housesWithPlanets := hc.AddPlanetsToHouses(houseList, positions)
	
	fmt.Println("Houses with Planets:")
	for _, house := range housesWithPlanets {
		if len(house.Planets) > 0 {
			planetNames := make([]string, len(house.Planets))
			for i, planet := range house.Planets {
				planetNames[i] = planet.Name
			}
			fmt.Printf("  %s: %v\n", house.Name, planetNames)
		}
	}
	fmt.Println()
	
	// 4. Calculate aspects
	fmt.Println("=== PLANETARY ASPECTS ===")
	ac := aspects.NewAspectCalculator()
	
	aspectList := ac.CalculateAllAspects(positions)
	
	fmt.Printf("Found %d aspects:\n", len(aspectList))
	for i, aspect := range aspectList {
		if i >= 10 { // Show only top 10
			break
		}
		fmt.Printf("  %s (%.1f%% strength)\n", aspect.String(), aspect.Strength)
	}
	fmt.Println()
	
	// 5. Analyze aspect patterns
	fmt.Println("=== ASPECT PATTERNS ===")
	patterns := ac.CalculateAspectPattern(positions)
	
	if len(patterns) > 0 {
		fmt.Printf("Found %d aspect patterns:\n", len(patterns))
		for _, pattern := range patterns {
			fmt.Printf("  %s\n", pattern.String())
		}
	} else {
		fmt.Println("No major aspect patterns found.")
	}
	fmt.Println()
	
	// 6. Coordinate system transformations
	fmt.Println("=== COORDINATE TRANSFORMATIONS ===")
	ct := coordinates.NewCoordinateTransformer()
	
	// Example: Convert Sun position to equatorial coordinates
	sunEcliptic := positions[0].Coordinates
	sunEquatorial := ct.EclipticToEquatorial(sunEcliptic)
	
	fmt.Printf("Sun Position:\n")
	fmt.Printf("  Ecliptic: %s\n", sunEcliptic.String())
	fmt.Printf("  Equatorial: %s\n", sunEquatorial.String())
	fmt.Println()
	
	// 7. Calculate Local Sidereal Time
	fmt.Println("=== TIME CALCULATIONS ===")
	
	lst := tc.LocalSiderealTime(jd, longitude)
	fmt.Printf("Local Sidereal Time: %.2f° (%.2f hours)\n", lst, lst/15.0)
	
	// Convert to horizontal coordinates for the Sun
	sunHorizontal := ct.EquatorialToHorizontal(sunEquatorial, lst, latitude)
	fmt.Printf("Sun Horizontal Position: %s\n", sunHorizontal.String())
	fmt.Println()
	
	// 8. Zodiac sign compatibility
	fmt.Println("=== ZODIAC COMPATIBILITY ===")
	
	sunSign := zc.EclipticToZodiac(positions[0].Coordinates.Longitude).Sign
	moonSign := zc.EclipticToZodiac(positions[1].Coordinates.Longitude).Sign
	
	compatibility := zc.GetSignCompatibility(sunSign, moonSign)
	fmt.Printf("Sun-Moon Compatibility (%s & %s): %.1f%%\n", sunSign.Name, moonSign.Name, compatibility)
	fmt.Println()
	
	// 9. Show all zodiac signs
	fmt.Println("=== ZODIAC SIGNS ===")
	signs := zc.GetZodiacSigns()
	
	fmt.Println("The 12 Zodiac Signs:")
	for _, sign := range signs {
		fmt.Printf("  %s\n", sign.String())
	}
	fmt.Println()
	
	fmt.Println("=== Demo Complete ===")
}