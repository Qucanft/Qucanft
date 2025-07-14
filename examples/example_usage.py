#!/usr/bin/env python3
"""
Example script demonstrating the Qucanft library functionality.

This script shows how to use the library to fetch planetary data,
calculate zodiac positions, houses, and aspects, and create visualizations.
"""

import sys
import os
import pandas as pd
from datetime import datetime
import matplotlib.pyplot as plt

# Add the qucanft package to the path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

try:
    from qucanft import (
        AstroDataFetcher, 
        ZodiacCalculator, 
        HousesCalculator, 
        AspectsCalculator,
        VisualizationHelper
    )
except ImportError as e:
    print(f"Error importing qucanft modules: {e}")
    print("Please ensure all dependencies are installed: pip install -r requirements.txt")
    sys.exit(1)


def main():
    """Main function demonstrating the library functionality."""
    print("=" * 60)
    print("Qucanft - Astronomical Data Fetching and Astrological Analysis")
    print("=" * 60)
    
    # Initialize the components
    print("\n1. Initializing components...")
    fetcher = AstroDataFetcher()
    zodiac_calc = ZodiacCalculator()
    houses_calc = HousesCalculator()
    aspects_calc = AspectsCalculator()
    viz_helper = VisualizationHelper()
    
    # Example 1: Fetch planetary positions for a specific date and location
    print("\n2. Fetching planetary positions...")
    try:
        # Example: January 1, 2024, noon UTC, New York
        date = "2024-01-01T12:00:00"
        location = "New York, NY"
        
        print(f"   Date: {date}")
        print(f"   Location: {location}")
        
        # Fetch planetary data
        planetary_data = fetcher.get_planet_positions(
            date=date,
            location=location,
            planets=["Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn"]
        )
        
        if not planetary_data.empty:
            print(f"   Successfully fetched data for {len(planetary_data)} planets")
            print("\n   Raw planetary positions:")
            print(planetary_data[['Planet', 'RA', 'Dec', 'Distance_AU']].to_string(index=False))
        else:
            print("   No planetary data fetched. Using sample data for demonstration.")
            # Create sample data for demonstration
            planetary_data = pd.DataFrame({
                'Planet': ['Sun', 'Moon', 'Mercury', 'Venus', 'Mars', 'Jupiter', 'Saturn'],
                'RA': [280.5, 45.2, 290.1, 315.8, 120.3, 30.7, 350.2],
                'Dec': [-23.1, 15.6, -20.8, -18.4, 25.7, 12.3, -22.5],
                'Distance_AU': [0.983, 0.002, 0.652, 0.287, 1.456, 5.234, 9.876]
            })
            
    except Exception as e:
        print(f"   Error fetching data: {e}")
        print("   Using sample data for demonstration.")
        # Create sample data for demonstration
        planetary_data = pd.DataFrame({
            'Planet': ['Sun', 'Moon', 'Mercury', 'Venus', 'Mars', 'Jupiter', 'Saturn'],
            'RA': [280.5, 45.2, 290.1, 315.8, 120.3, 30.7, 350.2],
            'Dec': [-23.1, 15.6, -20.8, -18.4, 25.7, 12.3, -22.5],
            'Distance_AU': [0.983, 0.002, 0.652, 0.287, 1.456, 5.234, 9.876]
        })
    
    # Example 2: Calculate zodiac positions
    print("\n3. Calculating zodiac positions...")
    try:
        zodiac_data = zodiac_calc.calculate_zodiac_positions(planetary_data)
        print("   Successfully calculated zodiac positions:")
        print(zodiac_data[['Planet', 'Zodiac_Sign', 'Position_String', 'Zodiac_Element']].to_string(index=False))
    except Exception as e:
        print(f"   Error calculating zodiac positions: {e}")
        # Add mock zodiac data for demonstration
        zodiac_data = planetary_data.copy()
        zodiac_data['Zodiac_Sign'] = ['Capricorn', 'Taurus', 'Capricorn', 'Aquarius', 'Cancer', 'Taurus', 'Pisces']
        zodiac_data['Position_String'] = ['10°30\' ♑', '15°12\' ♉', '20°06\' ♑', '15°48\' ♒', '00°18\' ♋', '00°42\' ♉', '20°12\' ♓']
        zodiac_data['Ecliptic_Longitude'] = [280.5, 45.2, 290.1, 315.8, 120.3, 30.7, 350.2]
    
    # Example 3: Calculate house positions
    print("\n4. Calculating house positions...")
    try:
        # Calculate ascendant and midheaven (simplified)
        ascendant = 75.0  # Example ascendant in Gemini
        midheaven = 345.0  # Example midheaven in Pisces
        
        # Calculate house cusps using Equal House system
        house_cusps = houses_calc.calculate_houses(
            ascendant=ascendant,
            house_system='equal'
        )
        
        # Add house positions to planetary data
        if 'Ecliptic_Longitude' in zodiac_data.columns:
            house_data = houses_calc.add_house_positions(zodiac_data, house_cusps)
            print("   Successfully calculated house positions:")
            print(house_data[['Planet', 'Zodiac_Sign', 'House', 'House_Meaning']].to_string(index=False))
        else:
            print("   Missing ecliptic longitude data for house calculations")
            house_data = zodiac_data.copy()
            house_data['House'] = [1, 2, 3, 4, 5, 6, 7]  # Mock data
            
    except Exception as e:
        print(f"   Error calculating house positions: {e}")
        house_data = zodiac_data.copy()
        house_data['House'] = [1, 2, 3, 4, 5, 6, 7]  # Mock data
    
    # Example 4: Calculate aspects
    print("\n5. Calculating planetary aspects...")
    try:
        if 'Ecliptic_Longitude' in house_data.columns:
            aspects_data = aspects_calc.calculate_all_aspects(house_data)
            if not aspects_data.empty:
                print(f"   Found {len(aspects_data)} aspects:")
                # Show top 5 strongest aspects
                strongest_aspects = aspects_calc.get_strongest_aspects(aspects_data, limit=5)
                print(strongest_aspects[['planet1', 'planet2', 'aspect', 'orb_difference', 'nature']].to_string(index=False))
            else:
                print("   No aspects found within orb limits")
        else:
            print("   Missing ecliptic longitude data for aspect calculations")
            aspects_data = pd.DataFrame()
    except Exception as e:
        print(f"   Error calculating aspects: {e}")
        aspects_data = pd.DataFrame()
    
    # Example 5: Create visualizations
    print("\n6. Creating visualizations...")
    try:
        # Create planetary positions chart
        if 'Zodiac_Sign' in house_data.columns:
            positions_fig = viz_helper.create_planetary_positions_chart(house_data)
            positions_fig.savefig('/tmp/planetary_positions.png', dpi=150, bbox_inches='tight')
            print("   ✓ Planetary positions chart saved to /tmp/planetary_positions.png")
        
        # Create natal chart (if we have all required data)
        if 'Ecliptic_Longitude' in house_data.columns and 'house_cusps' in locals():
            natal_fig = viz_helper.create_natal_chart(
                house_data, 
                house_cusps, 
                aspects_data if not aspects_data.empty else None
            )
            natal_fig.savefig('/tmp/natal_chart.png', dpi=150, bbox_inches='tight')
            print("   ✓ Natal chart saved to /tmp/natal_chart.png")
        
        # Create aspects summary chart
        if not aspects_data.empty:
            aspects_fig = viz_helper.create_aspects_summary_chart(aspects_data)
            aspects_fig.savefig('/tmp/aspects_summary.png', dpi=150, bbox_inches='tight')
            print("   ✓ Aspects summary chart saved to /tmp/aspects_summary.png")
        
        plt.close('all')  # Close all figures to free memory
        
    except Exception as e:
        print(f"   Error creating visualizations: {e}")
    
    # Example 6: Export data summary
    print("\n7. Exporting data summary...")
    try:
        summary = viz_helper.export_data_summary(
            house_data,
            aspects_data if not aspects_data.empty else None,
            house_cusps if 'house_cusps' in locals() else None
        )
        
        print("   Data summary:")
        print(f"   - Planets analyzed: {len(summary['planetary_positions'])}")
        print(f"   - Zodiac distribution: {summary['zodiac_distribution']}")
        if summary['aspects_summary']:
            print(f"   - Total aspects: {summary['aspects_summary']['total_aspects']}")
        
        # Save summary to file
        import json
        with open('/tmp/astrological_summary.json', 'w') as f:
            json.dump(summary, f, indent=2, default=str)
        print("   ✓ Summary saved to /tmp/astrological_summary.json")
        
    except Exception as e:
        print(f"   Error exporting summary: {e}")
    
    # Example 7: Demonstrate custom queries
    print("\n8. Demonstrating custom queries...")
    try:
        # Example: Query for a specific asteroid
        print("   Attempting custom query for asteroid Ceres (1)...")
        custom_data = fetcher.get_custom_query(
            target_id="1",  # Ceres
            date="2024-01-01T12:00:00",
            quantities="1,2,3"  # RA, Dec, Distance
        )
        
        if not custom_data.empty:
            print(f"   ✓ Custom query successful: {custom_data.shape[0]} records")
        else:
            print("   Custom query returned no data")
            
    except Exception as e:
        print(f"   Custom query failed: {e}")
        print("   This is expected if astroquery is not properly installed")
    
    # Example 8: Zodiac compatibility example
    print("\n9. Zodiac compatibility example...")
    try:
        compatibility = zodiac_calc.get_zodiac_compatibility('Leo', 'Aquarius')
        print(f"   Leo ♌ and Aquarius ♒ compatibility:")
        print(f"   Overall: {compatibility['overall_compatibility']}")
        print(f"   Element: {compatibility['element_compatibility']}")
        print(f"   Quality: {compatibility['quality_compatibility']}")
        
    except Exception as e:
        print(f"   Error calculating compatibility: {e}")
    
    print("\n" + "=" * 60)
    print("Example completed successfully!")
    print("Check the /tmp/ directory for generated files:")
    print("- planetary_positions.png")
    print("- natal_chart.png")
    print("- aspects_summary.png")
    print("- astrological_summary.json")
    print("=" * 60)


if __name__ == "__main__":
    main()