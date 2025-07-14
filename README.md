# Qucanft - Astronomical Data Fetching and Astrological Analysis

![Python](https://img.shields.io/badge/python-3.7%2B-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

Qucanft is a Python library that provides functionality to fetch planetary positioning and celestial data from public astronomical databases such as JPL Horizons using the Astroquery library. The library processes this data into formats suitable for astrological analysis and visual representation, including zodiac signs, houses, and aspects calculations.

## Features

- **Astronomical Data Fetching**: Query JPL Horizons for planetary positions and ephemeris data
- **Zodiac Calculations**: Convert celestial coordinates to zodiac signs with symbolic representations
- **House Systems**: Calculate astrological houses using multiple systems (Equal, Placidus, Whole Sign, etc.)
- **Aspects Analysis**: Calculate planetary aspects with interpretations and orb analysis
- **Visualization Tools**: Create natal charts, planetary distribution charts, and aspect summaries
- **Custom Queries**: Support for custom date ranges, locations, and celestial objects
- **Data Export**: Export calculations and visualizations in various formats

## Installation

### Prerequisites

```bash
pip install astroquery numpy pandas matplotlib astropy pytz swisseph
```

### Install from source

```bash
git clone https://github.com/Qucanft/Qucanft.git
cd Qucanft
pip install -r requirements.txt
python setup.py install
```

## Quick Start

```python
from qucanft import AstroDataFetcher, ZodiacCalculator, HousesCalculator, AspectsCalculator, VisualizationHelper
from datetime import datetime

# Initialize components
fetcher = AstroDataFetcher()
zodiac_calc = ZodiacCalculator()
houses_calc = HousesCalculator()
aspects_calc = AspectsCalculator()
viz = VisualizationHelper()

# Fetch planetary positions
date = "2024-01-01T12:00:00"
location = "New York, NY"
planets = ["Sun", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn"]

planetary_data = fetcher.get_planet_positions(date=date, location=location, planets=planets)

# Calculate zodiac positions
zodiac_data = zodiac_calc.calculate_zodiac_positions(planetary_data)

# Calculate house positions
ascendant = 75.0  # Example ascendant in Gemini
house_cusps = houses_calc.calculate_houses(ascendant=ascendant, house_system='equal')
house_data = houses_calc.add_house_positions(zodiac_data, house_cusps)

# Calculate aspects
aspects_data = aspects_calc.calculate_all_aspects(house_data)

# Create visualizations
natal_chart = viz.create_natal_chart(house_data, house_cusps, aspects_data)
natal_chart.savefig('natal_chart.png')
```

## Detailed Usage

### 1. Astronomical Data Fetching

```python
from qucanft import AstroDataFetcher

fetcher = AstroDataFetcher()

# Fetch planetary positions for a specific date and location
planetary_data = fetcher.get_planet_positions(
    date="2024-01-01T12:00:00",
    location="New York, NY",
    planets=["Sun", "Moon", "Mercury", "Venus", "Mars"]
)

# Fetch ephemeris data over a date range
ephemeris_data = fetcher.get_ephemeris_range(
    start_date="2024-01-01",
    end_date="2024-01-31",
    step="1d",
    planet="Sun"
)

# Custom query for any celestial object
asteroid_data = fetcher.get_custom_query(
    target_id="433",  # Eros asteroid
    date="2024-01-01T12:00:00",
    quantities="1,2,3,4"
)
```

### 2. Zodiac Calculations

```python
from qucanft import ZodiacCalculator

zodiac_calc = ZodiacCalculator()

# Calculate zodiac positions from planetary data
zodiac_data = zodiac_calc.calculate_zodiac_positions(planetary_data)

# Get zodiac sign information
aries_info = zodiac_calc.get_zodiac_sign_info('Aries')
print(f"Aries: {aries_info['element']}, {aries_info['quality']}, ruled by {aries_info['ruler']}")

# Calculate zodiac compatibility
compatibility = zodiac_calc.get_zodiac_compatibility('Leo', 'Aquarius')
print(f"Leo-Aquarius compatibility: {compatibility['overall_compatibility']}")
```

### 3. House Calculations

```python
from qucanft import HousesCalculator

houses_calc = HousesCalculator()

# Calculate house cusps using different systems
equal_houses = houses_calc.calculate_houses(ascendant=75.0, house_system='equal')
placidus_houses = houses_calc.calculate_houses(
    ascendant=75.0, 
    midheaven=345.0, 
    latitude=40.7128, 
    house_system='placidus'
)

# Add house positions to planetary data
house_data = houses_calc.add_house_positions(zodiac_data, equal_houses)

# Get house information
house_info = houses_calc.get_house_info(1)
print(f"1st House: {house_info['name']} - {house_info['theme']}")
```

### 4. Aspects Analysis

```python
from qucanft import AspectsCalculator

aspects_calc = AspectsCalculator()

# Calculate all aspects between planets
aspects_data = aspects_calc.calculate_all_aspects(house_data)

# Get strongest aspects
strongest = aspects_calc.get_strongest_aspects(aspects_data, limit=5)

# Filter aspects by nature
harmonious = aspects_calc.get_aspects_by_nature(aspects_data, 'Harmonious')
challenging = aspects_calc.get_aspects_by_nature(aspects_data, 'Challenging')

# Get aspects for a specific planet
sun_aspects = aspects_calc.get_planet_aspects(aspects_data, 'Sun')
```

### 5. Visualization

```python
from qucanft import VisualizationHelper

viz = VisualizationHelper()

# Create a natal chart
natal_chart = viz.create_natal_chart(house_data, house_cusps, aspects_data)
natal_chart.savefig('natal_chart.png')

# Create planetary distribution chart
positions_chart = viz.create_planetary_positions_chart(house_data)
positions_chart.savefig('planetary_positions.png')

# Create aspects summary
aspects_chart = viz.create_aspects_summary_chart(aspects_data)
aspects_chart.savefig('aspects_summary.png')

# Export data summary
summary = viz.export_data_summary(house_data, aspects_data, house_cusps)
```

## Advanced Features

### Custom Locations

```python
# Using coordinates
location = {
    'lat': 40.7128,
    'lon': -74.0060,
    'elevation': 10
}

# Using string (common locations are recognized)
location = "Tokyo, Japan"

# Using JPL Horizons location codes
location = "500"  # Geocentric
```

### Custom Date Ranges

```python
# Single date
date = "2024-01-01T12:00:00"

# Date range for ephemeris
ephemeris = fetcher.get_ephemeris_range(
    start_date="2024-01-01",
    end_date="2024-12-31",
    step="1d",
    planet="Mars"
)
```

### House Systems

- **Equal House**: Each house is exactly 30°
- **Placidus**: Time-based house system (most common)
- **Whole Sign**: Each house occupies an entire zodiac sign
- **Koch**: Alternative time-based system
- **Campanus**: Space-based system
- **Regiomontanus**: Medieval house system

### Aspect Types

**Major Aspects:**
- Conjunction (0°) - Neutral
- Sextile (60°) - Harmonious
- Square (90°) - Challenging
- Trine (120°) - Harmonious
- Opposition (180°) - Challenging

**Minor Aspects:**
- Semisextile (30°) - Mild
- Semisquare (45°) - Challenging
- Sesquiquadrate (135°) - Challenging
- Quincunx (150°) - Adjusting
- Quintile (72°) - Creative
- Biquintile (144°) - Creative

## Examples

Run the complete example:

```bash
python examples/example_usage.py
```

This example demonstrates:
- Fetching planetary data from JPL Horizons
- Converting coordinates to zodiac positions
- Calculating house positions
- Finding planetary aspects
- Creating visualizations
- Exporting data summaries

## API Reference

### Classes

- `AstroDataFetcher`: Fetch astronomical data from JPL Horizons
- `ZodiacCalculator`: Calculate zodiac positions and compatibility
- `HousesCalculator`: Calculate astrological houses
- `AspectsCalculator`: Calculate planetary aspects
- `VisualizationHelper`: Create charts and visualizations

### Key Methods

- `get_planet_positions()`: Fetch planetary positions for a date/location
- `calculate_zodiac_positions()`: Convert coordinates to zodiac signs
- `calculate_houses()`: Calculate house cusps
- `calculate_all_aspects()`: Find all aspects between planets
- `create_natal_chart()`: Generate natal chart visualization

## Dependencies

- `astroquery>=0.4.6`: Query astronomical databases
- `numpy>=1.21.0`: Numerical computations
- `pandas>=1.3.0`: Data manipulation
- `matplotlib>=3.5.0`: Plotting and visualization
- `astropy>=5.0.0`: Astronomy calculations
- `pytz>=2021.1`: Timezone handling
- `swisseph>=2.10.0`: Swiss Ephemeris calculations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built using the [Astroquery](https://astroquery.readthedocs.io/) library
- Astronomical calculations powered by [JPL Horizons](https://ssd.jpl.nasa.gov/horizons.cgi)
- Swiss Ephemeris calculations via [pyswisseph](https://github.com/astrorigin/pyswisseph)

## Support

For questions, issues, or contributions, please visit the [GitHub repository](https://github.com/Qucanft/Qucanft).

---

*Note: This library is for educational and research purposes. Astrological interpretations are provided for cultural and historical interest and should not be used for making important life decisions.*