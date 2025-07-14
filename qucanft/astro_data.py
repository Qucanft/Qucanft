"""
Astronomical data fetching module using Astroquery and JPL Horizons.

This module provides functionality to fetch planetary positioning and celestial data
from public astronomical databases such as JPL Horizons.
"""

from typing import Dict, List, Optional, Tuple, Union
from datetime import datetime, timezone
import pandas as pd
import numpy as np
from astroquery.jplhorizons import Horizons
from astropy.time import Time
from astropy.coordinates import EarthLocation
from astropy import units as u
import pytz


class AstroDataFetcher:
    """
    A class for fetching astronomical data from JPL Horizons using Astroquery.
    
    This class provides methods to query planetary positions, ephemeris data,
    and other celestial information for specified dates and locations.
    """
    
    # Major planets and luminaries commonly used in astrology
    PLANETS = {
        'Sun': 10,
        'Moon': 301,
        'Mercury': 199,
        'Venus': 299,
        'Mars': 499,
        'Jupiter': 599,
        'Saturn': 699,
        'Uranus': 799,
        'Neptune': 899,
        'Pluto': 999
    }
    
    def __init__(self):
        """Initialize the AstroDataFetcher."""
        self.data_cache = {}
    
    def get_planet_positions(self, 
                           date: Union[str, datetime],
                           location: Optional[Union[str, Dict[str, float]]] = None,
                           planets: Optional[List[str]] = None) -> pd.DataFrame:
        """
        Fetch planetary positions for a specific date and location.
        
        Args:
            date: Date for the query (ISO format string or datetime object)
            location: Location as string (e.g., "New York, NY") or dict with
                     'lat', 'lon', 'elevation' keys. If None, uses geocentric position.
            planets: List of planet names to query. If None, queries all major planets.
            
        Returns:
            DataFrame with planetary positions including RA, Dec, and other ephemeris data
            
        Example:
            >>> fetcher = AstroDataFetcher()
            >>> positions = fetcher.get_planet_positions(
            ...     date="2023-01-01T12:00:00",
            ...     location="New York, NY",
            ...     planets=["Sun", "Moon", "Mercury"]
            ... )
        """
        if planets is None:
            planets = list(self.PLANETS.keys())
        
        # Convert date to proper format
        if isinstance(date, str):
            query_date = Time(date)
        elif isinstance(date, datetime):
            query_date = Time(date)
        else:
            query_date = Time(date)
        
        # Setup location
        if location is None:
            location_code = '500'  # Geocentric
        elif isinstance(location, str):
            location_code = location
        elif isinstance(location, dict):
            # Convert coordinates to location string
            lat = location.get('lat', 0)
            lon = location.get('lon', 0)
            elevation = location.get('elevation', 0)
            location_code = f"{lon:+.6f},{lat:+.6f},{elevation:.0f}"
        else:
            location_code = '500'  # Default to geocentric
        
        results = []
        
        for planet_name in planets:
            if planet_name not in self.PLANETS:
                print(f"Warning: Unknown planet '{planet_name}', skipping...")
                continue
                
            planet_id = self.PLANETS[planet_name]
            
            try:
                # Query JPL Horizons
                obj = Horizons(
                    id=planet_id,
                    location=location_code,
                    epochs=query_date.jd
                )
                
                # Get ephemeris data
                ephemeris = obj.ephemerides(
                    quantities='1,2,3,4,8,9,19,20,23,24'  # RA, Dec, distance, magnitude, etc.
                )
                
                # Extract relevant data
                row_data = {
                    'Planet': planet_name,
                    'Date': query_date.iso,
                    'RA': ephemeris['RA'][0],  # Right Ascension
                    'Dec': ephemeris['DEC'][0],  # Declination
                    'Distance_AU': ephemeris['delta'][0],  # Distance in AU
                    'Magnitude': ephemeris['V'][0] if 'V' in ephemeris.colnames else np.nan,
                    'Elongation': ephemeris['elong'][0] if 'elong' in ephemeris.colnames else np.nan,
                    'Phase': ephemeris['alpha'][0] if 'alpha' in ephemeris.colnames else np.nan,
                }
                
                results.append(row_data)
                
            except Exception as e:
                print(f"Error fetching data for {planet_name}: {e}")
                continue
        
        return pd.DataFrame(results)
    
    def get_ephemeris_range(self,
                           start_date: Union[str, datetime],
                           end_date: Union[str, datetime],
                           step: str = '1d',
                           planet: str = 'Sun',
                           location: Optional[Union[str, Dict[str, float]]] = None) -> pd.DataFrame:
        """
        Fetch ephemeris data for a planet over a date range.
        
        Args:
            start_date: Start date for the range
            end_date: End date for the range
            step: Time step (e.g., '1d' for daily, '1h' for hourly)
            planet: Planet name to query
            location: Location specification
            
        Returns:
            DataFrame with ephemeris data over the specified range
            
        Example:
            >>> fetcher = AstroDataFetcher()
            >>> sun_data = fetcher.get_ephemeris_range(
            ...     start_date="2023-01-01",
            ...     end_date="2023-01-31",
            ...     step="1d",
            ...     planet="Sun"
            ... )
        """
        if planet not in self.PLANETS:
            raise ValueError(f"Unknown planet: {planet}")
        
        # Convert dates
        start_time = Time(start_date)
        end_time = Time(end_date)
        
        # Setup location
        if location is None:
            location_code = '500'  # Geocentric
        elif isinstance(location, str):
            location_code = location
        elif isinstance(location, dict):
            lat = location.get('lat', 0)
            lon = location.get('lon', 0)
            elevation = location.get('elevation', 0)
            location_code = f"{lon:+.6f},{lat:+.6f},{elevation:.0f}"
        else:
            location_code = '500'
        
        planet_id = self.PLANETS[planet]
        
        try:
            # Query JPL Horizons
            obj = Horizons(
                id=planet_id,
                location=location_code,
                epochs={
                    'start': start_time.iso,
                    'stop': end_time.iso,
                    'step': step
                }
            )
            
            # Get ephemeris data
            ephemeris = obj.ephemerides(
                quantities='1,2,3,4,8,9,19,20,23,24'
            )
            
            # Convert to DataFrame
            df = ephemeris.to_pandas()
            df['Planet'] = planet
            
            return df
            
        except Exception as e:
            raise RuntimeError(f"Error fetching ephemeris data for {planet}: {e}")
    
    def get_custom_query(self,
                        target_id: Union[str, int],
                        date: Union[str, datetime],
                        location: Optional[Union[str, Dict[str, float]]] = None,
                        quantities: str = '1,2,3,4') -> pd.DataFrame:
        """
        Perform a custom query to JPL Horizons for any celestial object.
        
        Args:
            target_id: JPL Horizons object ID (int) or name (str)
            date: Date for the query
            location: Location specification
            quantities: Quantities to retrieve (JPL Horizons format)
            
        Returns:
            DataFrame with query results
            
        Example:
            >>> fetcher = AstroDataFetcher()
            >>> asteroid_data = fetcher.get_custom_query(
            ...     target_id="433",  # Eros asteroid
            ...     date="2023-01-01T12:00:00",
            ...     quantities="1,2,3,4,8,9"
            ... )
        """
        # Convert date
        if isinstance(date, str):
            query_date = Time(date)
        elif isinstance(date, datetime):
            query_date = Time(date)
        else:
            query_date = Time(date)
        
        # Setup location
        if location is None:
            location_code = '500'  # Geocentric
        elif isinstance(location, str):
            location_code = location
        elif isinstance(location, dict):
            lat = location.get('lat', 0)
            lon = location.get('lon', 0)
            elevation = location.get('elevation', 0)
            location_code = f"{lon:+.6f},{lat:+.6f},{elevation:.0f}"
        else:
            location_code = '500'
        
        try:
            # Query JPL Horizons
            obj = Horizons(
                id=target_id,
                location=location_code,
                epochs=query_date.jd
            )
            
            # Get ephemeris data
            ephemeris = obj.ephemerides(quantities=quantities)
            
            # Convert to DataFrame
            df = ephemeris.to_pandas()
            df['Target_ID'] = target_id
            df['Query_Date'] = query_date.iso
            
            return df
            
        except Exception as e:
            raise RuntimeError(f"Error in custom query for target {target_id}: {e}")
    
    def get_location_from_string(self, location_string: str) -> Dict[str, float]:
        """
        Convert a location string to coordinates (placeholder for geocoding).
        
        Args:
            location_string: Location as string (e.g., "New York, NY")
            
        Returns:
            Dictionary with 'lat', 'lon', 'elevation' keys
            
        Note:
            This is a simplified implementation. In a production environment,
            you would use a geocoding service like Google Maps API or similar.
        """
        # This is a placeholder implementation
        # In a real application, you would use a geocoding service
        common_locations = {
            'New York, NY': {'lat': 40.7128, 'lon': -74.0060, 'elevation': 10},
            'Los Angeles, CA': {'lat': 34.0522, 'lon': -118.2437, 'elevation': 71},
            'London, UK': {'lat': 51.5074, 'lon': -0.1278, 'elevation': 11},
            'Tokyo, Japan': {'lat': 35.6762, 'lon': 139.6503, 'elevation': 6},
            'Sydney, Australia': {'lat': -33.8688, 'lon': 151.2093, 'elevation': 3}
        }
        
        if location_string in common_locations:
            return common_locations[location_string]
        else:
            # Default to Greenwich, UK if location not found
            print(f"Warning: Location '{location_string}' not found, using Greenwich, UK")
            return {'lat': 51.4769, 'lon': 0.0, 'elevation': 0}
    
    def clear_cache(self):
        """Clear the internal data cache."""
        self.data_cache.clear()