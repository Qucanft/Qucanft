"""
Astrological houses calculation module.

This module provides functionality to calculate astrological houses using
different house systems (Placidus, Koch, Equal, etc.) and determine
planetary house positions.
"""

from typing import Dict, List, Optional, Tuple
import numpy as np
import pandas as pd
from datetime import datetime
import math


class HousesCalculator:
    """
    A class for calculating astrological houses and planetary house positions.
    
    This class supports multiple house systems and provides methods to
    calculate house cusps and determine which house planets are in.
    """
    
    # House systems available
    HOUSE_SYSTEMS = {
        'placidus': 'Placidus',
        'koch': 'Koch',
        'equal': 'Equal House',
        'whole': 'Whole Sign',
        'campanus': 'Campanus',
        'regiomontanus': 'Regiomontanus'
    }
    
    # House meanings and themes
    HOUSE_MEANINGS = {
        1: {'name': 'Ascendant/Self', 'theme': 'Identity, appearance, first impressions'},
        2: {'name': 'Resources', 'theme': 'Money, possessions, values, self-worth'},
        3: {'name': 'Communication', 'theme': 'Communication, siblings, short trips, learning'},
        4: {'name': 'Home/Family', 'theme': 'Home, family, roots, emotional foundation'},
        5: {'name': 'Creativity', 'theme': 'Creativity, romance, children, self-expression'},
        6: {'name': 'Work/Health', 'theme': 'Work, health, daily routines, service'},
        7: {'name': 'Partnerships', 'theme': 'Marriage, partnerships, open enemies'},
        8: {'name': 'Transformation', 'theme': 'Death, rebirth, shared resources, occult'},
        9: {'name': 'Philosophy', 'theme': 'Higher learning, philosophy, long journeys'},
        10: {'name': 'Career/Status', 'theme': 'Career, reputation, public standing'},
        11: {'name': 'Friends/Goals', 'theme': 'Friends, groups, hopes, wishes'},
        12: {'name': 'Subconscious', 'theme': 'Subconscious, karma, hidden enemies, sacrifice'}
    }
    
    def __init__(self):
        """Initialize the HousesCalculator."""
        pass
    
    def calculate_ascendant(self, 
                          local_sidereal_time: float,
                          latitude: float,
                          obliquity: float = 23.4367) -> float:
        """
        Calculate the Ascendant (1st house cusp) for given time and location.
        
        Args:
            local_sidereal_time: Local sidereal time in hours
            latitude: Geographic latitude in degrees
            obliquity: Obliquity of the ecliptic in degrees
            
        Returns:
            Ascendant position in degrees of ecliptic longitude
        """
        # Convert to radians
        lst_rad = math.radians(local_sidereal_time * 15)  # LST in degrees
        lat_rad = math.radians(latitude)
        obl_rad = math.radians(obliquity)
        
        # Calculate ascendant using spherical trigonometry
        # This is a simplified calculation - in practice, you'd use more precise algorithms
        y = -math.cos(lst_rad)
        x = math.sin(lst_rad) * math.cos(obl_rad) + math.tan(lat_rad) * math.sin(obl_rad)
        
        ascendant_rad = math.atan2(y, x)
        ascendant_deg = math.degrees(ascendant_rad)
        
        # Normalize to 0-360 range
        return ascendant_deg % 360
    
    def calculate_midheaven(self, 
                           local_sidereal_time: float,
                           obliquity: float = 23.4367) -> float:
        """
        Calculate the Midheaven (10th house cusp) for given time.
        
        Args:
            local_sidereal_time: Local sidereal time in hours
            obliquity: Obliquity of the ecliptic in degrees
            
        Returns:
            Midheaven position in degrees of ecliptic longitude
        """
        # Convert LST to degrees
        lst_deg = local_sidereal_time * 15
        
        # Calculate MC using spherical trigonometry
        # This is a simplified calculation
        mc_deg = lst_deg % 360
        
        return mc_deg
    
    def calculate_equal_houses(self, ascendant: float) -> Dict[int, float]:
        """
        Calculate house cusps using the Equal House system.
        
        Args:
            ascendant: Ascendant position in degrees
            
        Returns:
            Dictionary with house numbers as keys and cusp positions as values
        """
        houses = {}
        
        for house_num in range(1, 13):
            # In Equal House system, each house is exactly 30 degrees
            cusp = (ascendant + (house_num - 1) * 30) % 360
            houses[house_num] = cusp
        
        return houses
    
    def calculate_whole_sign_houses(self, ascendant: float) -> Dict[int, float]:
        """
        Calculate house cusps using the Whole Sign system.
        
        Args:
            ascendant: Ascendant position in degrees
            
        Returns:
            Dictionary with house numbers as keys and cusp positions as values
        """
        # In Whole Sign system, the 1st house starts at 0° of the ascendant sign
        ascendant_sign_start = (int(ascendant // 30)) * 30
        
        houses = {}
        for house_num in range(1, 13):
            cusp = (ascendant_sign_start + (house_num - 1) * 30) % 360
            houses[house_num] = cusp
        
        return houses
    
    def calculate_placidus_houses(self, 
                                 ascendant: float,
                                 midheaven: float,
                                 latitude: float,
                                 obliquity: float = 23.4367) -> Dict[int, float]:
        """
        Calculate house cusps using the Placidus system (simplified).
        
        Args:
            ascendant: Ascendant position in degrees
            midheaven: Midheaven position in degrees
            latitude: Geographic latitude in degrees
            obliquity: Obliquity of the ecliptic in degrees
            
        Returns:
            Dictionary with house numbers as keys and cusp positions as values
            
        Note:
            This is a simplified implementation. A full Placidus calculation
            requires complex iterative algorithms and is typically done using
            specialized libraries like Swiss Ephemeris.
        """
        houses = {}
        
        # Set the angular houses (1st, 4th, 7th, 10th)
        houses[1] = ascendant
        houses[4] = (midheaven + 180) % 360  # IC (Imum Coeli)
        houses[7] = (ascendant + 180) % 360  # Descendant
        houses[10] = midheaven  # MC (Medium Coeli)
        
        # Simplified calculation for intermediate houses
        # In practice, this would use complex trigonometry
        for house_num in [2, 3, 5, 6, 8, 9, 11, 12]:
            if house_num in [2, 3]:
                # Houses 2 and 3 between ASC and MC
                factor = (house_num - 1) / 3
                angle_diff = (midheaven - ascendant) % 360
                if angle_diff > 180:
                    angle_diff = 360 - angle_diff
                houses[house_num] = (ascendant + factor * angle_diff) % 360
            elif house_num in [5, 6]:
                # Houses 5 and 6 between MC and DSC
                factor = (house_num - 4) / 3
                angle_diff = (houses[7] - midheaven) % 360
                if angle_diff > 180:
                    angle_diff = 360 - angle_diff
                houses[house_num] = (midheaven + factor * angle_diff) % 360
            elif house_num in [8, 9]:
                # Houses 8 and 9 between DSC and IC
                factor = (house_num - 7) / 3
                angle_diff = (houses[4] - houses[7]) % 360
                if angle_diff > 180:
                    angle_diff = 360 - angle_diff
                houses[house_num] = (houses[7] + factor * angle_diff) % 360
            elif house_num in [11, 12]:
                # Houses 11 and 12 between IC and ASC
                factor = (house_num - 10) / 3
                angle_diff = (ascendant - houses[4]) % 360
                if angle_diff > 180:
                    angle_diff = 360 - angle_diff
                houses[house_num] = (houses[4] + factor * angle_diff) % 360
        
        return houses
    
    def calculate_houses(self, 
                        ascendant: float,
                        midheaven: Optional[float] = None,
                        latitude: Optional[float] = None,
                        house_system: str = 'equal') -> Dict[int, float]:
        """
        Calculate house cusps using the specified house system.
        
        Args:
            ascendant: Ascendant position in degrees
            midheaven: Midheaven position in degrees (required for some systems)
            latitude: Geographic latitude in degrees (required for some systems)
            house_system: House system to use ('equal', 'whole', 'placidus', etc.)
            
        Returns:
            Dictionary with house numbers as keys and cusp positions as values
        """
        house_system = house_system.lower()
        
        if house_system == 'equal':
            return self.calculate_equal_houses(ascendant)
        elif house_system == 'whole':
            return self.calculate_whole_sign_houses(ascendant)
        elif house_system == 'placidus':
            if midheaven is None or latitude is None:
                raise ValueError("Placidus system requires midheaven and latitude")
            return self.calculate_placidus_houses(ascendant, midheaven, latitude)
        else:
            # Default to equal house system
            print(f"Warning: House system '{house_system}' not implemented, using Equal House")
            return self.calculate_equal_houses(ascendant)
    
    def determine_planet_house(self, 
                             planet_longitude: float,
                             house_cusps: Dict[int, float]) -> int:
        """
        Determine which house a planet is in based on its longitude.
        
        Args:
            planet_longitude: Planet's ecliptic longitude in degrees
            house_cusps: Dictionary of house cusps
            
        Returns:
            House number (1-12) that the planet is in
        """
        # Normalize planet longitude to 0-360
        planet_longitude = planet_longitude % 360
        
        for house_num in range(1, 13):
            cusp_current = house_cusps[house_num]
            cusp_next = house_cusps[house_num + 1] if house_num < 12 else house_cusps[1]
            
            # Handle the case where the house crosses 0°
            if cusp_current > cusp_next:
                # House crosses 0°
                if planet_longitude >= cusp_current or planet_longitude < cusp_next:
                    return house_num
            else:
                # Normal case
                if cusp_current <= planet_longitude < cusp_next:
                    return house_num
        
        # Fallback (shouldn't happen with correct calculations)
        return 1
    
    def add_house_positions(self, 
                           planetary_data: pd.DataFrame,
                           house_cusps: Dict[int, float]) -> pd.DataFrame:
        """
        Add house positions to planetary data.
        
        Args:
            planetary_data: DataFrame with planetary positions
            house_cusps: Dictionary of house cusps
            
        Returns:
            DataFrame with added house position columns
        """
        if 'Ecliptic_Longitude' not in planetary_data.columns:
            raise ValueError("DataFrame must contain 'Ecliptic_Longitude' column")
        
        result_df = planetary_data.copy()
        
        # Calculate house positions
        house_positions = []
        house_meanings = []
        
        for _, row in planetary_data.iterrows():
            house_num = self.determine_planet_house(row['Ecliptic_Longitude'], house_cusps)
            house_positions.append(house_num)
            house_meanings.append(self.HOUSE_MEANINGS[house_num]['theme'])
        
        # Add house columns
        result_df['House'] = house_positions
        result_df['House_Meaning'] = house_meanings
        
        return result_df
    
    def get_house_info(self, house_number: int) -> Optional[Dict[str, str]]:
        """
        Get information about a specific house.
        
        Args:
            house_number: House number (1-12)
            
        Returns:
            Dictionary with house information, or None if invalid house number
        """
        if house_number in self.HOUSE_MEANINGS:
            return self.HOUSE_MEANINGS[house_number].copy()
        return None
    
    def calculate_house_strengths(self, 
                                planetary_data: pd.DataFrame,
                                house_cusps: Dict[int, float]) -> Dict[int, int]:
        """
        Calculate the number of planets in each house.
        
        Args:
            planetary_data: DataFrame with planetary positions
            house_cusps: Dictionary of house cusps
            
        Returns:
            Dictionary with house numbers as keys and planet counts as values
        """
        if 'Ecliptic_Longitude' not in planetary_data.columns:
            raise ValueError("DataFrame must contain 'Ecliptic_Longitude' column")
        
        house_counts = {i: 0 for i in range(1, 13)}
        
        for _, row in planetary_data.iterrows():
            house_num = self.determine_planet_house(row['Ecliptic_Longitude'], house_cusps)
            house_counts[house_num] += 1
        
        return house_counts
    
    def format_house_cusp(self, house_num: int, cusp_longitude: float) -> str:
        """
        Format a house cusp position as a readable string.
        
        Args:
            house_num: House number (1-12)
            cusp_longitude: Cusp longitude in degrees
            
        Returns:
            Formatted string with house number and position
        """
        # Convert longitude to zodiac position
        sign_index = int(cusp_longitude // 30)
        degree_in_sign = cusp_longitude % 30
        
        zodiac_signs = ['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo',
                       'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces']
        
        sign_name = zodiac_signs[sign_index]
        
        return f"House {house_num}: {int(degree_in_sign)}°{int((degree_in_sign % 1) * 60):02d}' {sign_name}"