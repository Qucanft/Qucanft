package main

import (
	"fmt"
	"time"
	
	"github.com/Qucanft/Qucanft/pkg/aspects"
	"github.com/Qucanft/Qucanft/pkg/astro"
	"github.com/Qucanft/Qucanft/pkg/coordinates"
	"github.com/Qucanft/Qucanft/pkg/planets"
)

func main() {
	fmt.Println("=== Astronomical Calculations Demo ===")
	
	// Create a chart for a specific date and location
	// Example: New York City on January 1, 2024 at noon
	dateTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	location := astro.Location{
		Latitude:  40.7128,  // New York City
		Longitude: -74.0060,
		Timezone:  "America/New_York",
	}
	
	chart := astro.NewChart(dateTime, location)
	
	fmt.Printf("\nChart for: %s\n", dateTime.Format("January 2, 2006 15:04:05 MST"))
	fmt.Printf("Location: %.4f°N, %.4f°W\n", location.Latitude, -location.Longitude)
	fmt.Printf("Julian Date: %.6f\n", chart.JulianDate)
	
	// Display planetary positions
	fmt.Println("\n=== Planetary Positions ===")
	for i, position := range chart.Positions {
		sign := chart.Signs[i]
		fmt.Printf("%-8s: %s (%s)\n", 
			position.Planet.String(), 
			sign.FormatPosition(),
			sign.Element.String())
	}
	
	// Display lunar phase
	phase, angle := chart.GetLunarPhase()
	fmt.Printf("\nLunar Phase: %s (%.1f°)\n", phase, angle)
	
	// Display elemental balance
	fmt.Println("\n=== Elemental Balance ===")
	elementBalance := chart.GetElementalBalance()
	for element, count := range elementBalance {
		fmt.Printf("%-5s: %d planets\n", element.String(), count)
	}
	
	// Display quality balance
	fmt.Println("\n=== Quality Balance ===")
	qualityBalance := chart.GetQualityBalance()
	for quality, count := range qualityBalance {
		fmt.Printf("%-8s: %d planets\n", quality.String(), count)
	}
	
	// Display major aspects
	fmt.Println("\n=== Major Aspects ===")
	majorAspects := chart.GetMajorAspects()
	for _, aspect := range majorAspects {
		fmt.Printf("%s %s %s (orb: %.1f°) %s\n",
			aspect.Planet1.String(),
			aspects.GetAspectInfo(aspect.Type).Symbol,
			aspect.Planet2.String(),
			aspect.Orb,
			func() string {
				if aspect.IsApplying {
					return "applying"
				}
				return "separating"
			}())
	}
	
	// Display aspects for a specific planet (Sun)
	fmt.Println("\n=== Aspects to Sun ===")
	sunAspects := chart.GetAspectsForPlanet(planets.Sun)
	for _, aspect := range sunAspects {
		otherPlanet := aspect.Planet2
		if aspect.Planet2 == planets.Sun {
			otherPlanet = aspect.Planet1
		}
		
		aspectInfo := aspects.GetAspectInfo(aspect.Type)
		fmt.Printf("Sun %s %s: %s (orb: %.1f°)\n",
			aspectInfo.Symbol,
			otherPlanet.String(),
			aspectInfo.Meaning,
			aspect.Orb)
	}
	
	// Demonstrate coordinate transformations
	fmt.Println("\n=== Coordinate Transformations ===")
	sunPos := chart.GetPlanetPosition(planets.Sun)
	if sunPos != nil {
		// Convert to equatorial coordinates
		obliquity := coordinates.GetObliquity(float64(chart.JulianDate))
		ecliptic := coordinates.EclipticCoordinates{
			Longitude: sunPos.Longitude,
			Latitude:  sunPos.Latitude,
		}
		equatorial := ecliptic.ToEquatorial(obliquity)
		
		fmt.Printf("Sun ecliptic: %.2f° longitude, %.2f° latitude\n",
			sunPos.Longitude, sunPos.Latitude)
		fmt.Printf("Sun equatorial: %.2f° RA, %.2f° Dec\n",
			equatorial.RightAscension, equatorial.Declination)
		fmt.Printf("Obliquity: %.4f°\n", obliquity)
	}
	
	fmt.Println("\n=== Demonstration Complete ===")
}