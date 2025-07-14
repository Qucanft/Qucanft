// Package main demonstrates basic usage of the Qucanft library.
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/Qucanft/Qucanft/pkg/planets"
	timeutil "github.com/Qucanft/Qucanft/pkg/time"
	"github.com/Qucanft/Qucanft/pkg/zodiac"
)

func main() {
	fmt.Println("=== Basic Astrological Calculations ===")
	
	// Create a calculator
	pc := planets.NewPlanetaryCalculator()
	zc := zodiac.NewZodiacCalculator()
	tc := timeutil.NewTimeConverter()
	
	// Calculate planet positions for a specific date
	birthTime := time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC)
	jd := tc.ToJulianDay(birthTime)
	
	fmt.Printf("Birth Date: %s\n", birthTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Julian Day: %s\n\n", jd.String())
	
	// Calculate positions for inner planets
	planetNames := []string{"Sun", "Moon", "Mercury", "Venus", "Mars"}
	
	fmt.Println("Planetary Positions:")
	for _, planetName := range planetNames {
		pos, err := pc.CalculatePosition(planetName, jd)
		if err != nil {
			log.Printf("Error calculating %s: %v", planetName, err)
			continue
		}
		
		zodiacPos := zc.EclipticToZodiac(pos.Coordinates.Longitude)
		formatted := zc.FormatZodiacPosition(zodiacPos)
		
		fmt.Printf("  %s: %s\n", planetName, formatted)
	}
	
	fmt.Println("\n=== Example Complete ===")
}