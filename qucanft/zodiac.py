"""
Zodiac calculations module for converting celestial coordinates to zodiac signs.

This module provides functionality to calculate zodiac signs, degrees, and 
symbolic representations from astronomical coordinates.
"""

from typing import Dict, List, Optional, Tuple
import numpy as np
import pandas as pd
from datetime import datetime
from astropy.coordinates import SkyCoord
from astropy import units as u


class ZodiacCalculator:
    """
    A class for calculating zodiac signs and positions from astronomical coordinates.
    
    This class converts Right Ascension and Declination coordinates to zodiac
    positions and provides symbolic representations.
    """
    
    # Zodiac signs with their celestial longitude ranges (in degrees)
    ZODIAC_SIGNS = [
        {'name': 'Aries', 'symbol': '♈', 'element': 'Fire', 'quality': 'Cardinal', 'ruler': 'Mars'},
        {'name': 'Taurus', 'symbol': '♉', 'element': 'Earth', 'quality': 'Fixed', 'ruler': 'Venus'},
        {'name': 'Gemini', 'symbol': '♊', 'element': 'Air', 'quality': 'Mutable', 'ruler': 'Mercury'},
        {'name': 'Cancer', 'symbol': '♋', 'element': 'Water', 'quality': 'Cardinal', 'ruler': 'Moon'},
        {'name': 'Leo', 'symbol': '♌', 'element': 'Fire', 'quality': 'Fixed', 'ruler': 'Sun'},
        {'name': 'Virgo', 'symbol': '♍', 'element': 'Earth', 'quality': 'Mutable', 'ruler': 'Mercury'},
        {'name': 'Libra', 'symbol': '♎', 'element': 'Air', 'quality': 'Cardinal', 'ruler': 'Venus'},
        {'name': 'Scorpio', 'symbol': '♏', 'element': 'Water', 'quality': 'Fixed', 'ruler': 'Mars'},
        {'name': 'Sagittarius', 'symbol': '♐', 'element': 'Fire', 'quality': 'Mutable', 'ruler': 'Jupiter'},
        {'name': 'Capricorn', 'symbol': '♑', 'element': 'Earth', 'quality': 'Cardinal', 'ruler': 'Saturn'},
        {'name': 'Aquarius', 'symbol': '♒', 'element': 'Air', 'quality': 'Fixed', 'ruler': 'Saturn'},
        {'name': 'Pisces', 'symbol': '♓', 'element': 'Water', 'quality': 'Mutable', 'ruler': 'Jupiter'}
    ]
    
    def __init__(self):
        """Initialize the ZodiacCalculator."""
        pass
    
    def ra_dec_to_ecliptic(self, ra: float, dec: float, epoch: str = 'J2000') -> Tuple[float, float]:
        """
        Convert Right Ascension and Declination to ecliptic coordinates.
        
        Args:
            ra: Right Ascension in degrees
            dec: Declination in degrees
            epoch: Epoch for the coordinate system (default: J2000)
            
        Returns:
            Tuple of (ecliptic_longitude, ecliptic_latitude) in degrees
        """
        # Create SkyCoord object
        coord = SkyCoord(ra=ra*u.degree, dec=dec*u.degree, frame='icrs')
        
        # Convert to ecliptic coordinates
        ecliptic_coord = coord.geocentrictrueecliptic
        
        return ecliptic_coord.lon.degree, ecliptic_coord.lat.degree
    
    def ecliptic_to_zodiac(self, ecliptic_longitude: float) -> Dict[str, any]:
        """
        Convert ecliptic longitude to zodiac sign and degree.
        
        Args:
            ecliptic_longitude: Ecliptic longitude in degrees
            
        Returns:
            Dictionary with zodiac information including sign, degree, symbol, etc.
        """
        # Normalize longitude to 0-360 range
        longitude = ecliptic_longitude % 360
        
        # Calculate zodiac sign (each sign is 30 degrees)
        sign_index = int(longitude // 30)
        degree_in_sign = longitude % 30
        
        # Get zodiac sign information
        zodiac_info = self.ZODIAC_SIGNS[sign_index].copy()
        zodiac_info.update({
            'degree': degree_in_sign,
            'absolute_degree': longitude,
            'degree_minutes': int((degree_in_sign % 1) * 60),
            'degree_seconds': int(((degree_in_sign % 1) * 60 % 1) * 60),
            'position_string': f"{int(degree_in_sign)}°{int((degree_in_sign % 1) * 60):02d}' {zodiac_info['symbol']}"
        })
        
        return zodiac_info
    
    def calculate_zodiac_positions(self, planetary_data: pd.DataFrame) -> pd.DataFrame:
        """
        Calculate zodiac positions for planetary data.
        
        Args:
            planetary_data: DataFrame with planetary positions (must have 'RA' and 'Dec' columns)
            
        Returns:
            DataFrame with added zodiac position columns
        """
        if 'RA' not in planetary_data.columns or 'Dec' not in planetary_data.columns:
            raise ValueError("DataFrame must contain 'RA' and 'Dec' columns")
        
        result_df = planetary_data.copy()
        
        # Calculate zodiac positions for each planet
        zodiac_data = []
        
        for _, row in planetary_data.iterrows():
            # Convert RA/Dec to ecliptic coordinates
            ecl_lon, ecl_lat = self.ra_dec_to_ecliptic(row['RA'], row['Dec'])
            
            # Get zodiac information
            zodiac_info = self.ecliptic_to_zodiac(ecl_lon)
            
            zodiac_data.append({
                'Ecliptic_Longitude': ecl_lon,
                'Ecliptic_Latitude': ecl_lat,
                'Zodiac_Sign': zodiac_info['name'],
                'Zodiac_Symbol': zodiac_info['symbol'],
                'Zodiac_Degree': zodiac_info['degree'],
                'Zodiac_Element': zodiac_info['element'],
                'Zodiac_Quality': zodiac_info['quality'],
                'Zodiac_Ruler': zodiac_info['ruler'],
                'Position_String': zodiac_info['position_string']
            })
        
        # Add zodiac columns to the result DataFrame
        zodiac_df = pd.DataFrame(zodiac_data)
        result_df = pd.concat([result_df, zodiac_df], axis=1)
        
        return result_df
    
    def get_zodiac_sign_info(self, sign_name: str) -> Optional[Dict[str, any]]:
        """
        Get detailed information about a zodiac sign.
        
        Args:
            sign_name: Name of the zodiac sign
            
        Returns:
            Dictionary with zodiac sign information, or None if not found
        """
        for sign in self.ZODIAC_SIGNS:
            if sign['name'].lower() == sign_name.lower():
                return sign.copy()
        return None
    
    def get_zodiac_compatibility(self, sign1: str, sign2: str) -> Dict[str, any]:
        """
        Calculate basic zodiac compatibility based on elements and qualities.
        
        Args:
            sign1: First zodiac sign name
            sign2: Second zodiac sign name
            
        Returns:
            Dictionary with compatibility information
        """
        info1 = self.get_zodiac_sign_info(sign1)
        info2 = self.get_zodiac_sign_info(sign2)
        
        if not info1 or not info2:
            return {'compatibility': 'Unknown', 'reason': 'Invalid zodiac sign'}
        
        # Element compatibility
        element_compatibility = {
            ('Fire', 'Fire'): 'High',
            ('Fire', 'Air'): 'High',
            ('Fire', 'Water'): 'Low',
            ('Fire', 'Earth'): 'Medium',
            ('Earth', 'Earth'): 'High',
            ('Earth', 'Water'): 'High',
            ('Earth', 'Air'): 'Low',
            ('Air', 'Air'): 'High',
            ('Air', 'Water'): 'Low',
            ('Water', 'Water'): 'High'
        }
        
        element_pair = (info1['element'], info2['element'])
        if element_pair not in element_compatibility:
            element_pair = (info2['element'], info1['element'])
        
        element_comp = element_compatibility.get(element_pair, 'Medium')
        
        # Quality compatibility
        quality_compatibility = {
            ('Cardinal', 'Cardinal'): 'Medium',
            ('Cardinal', 'Fixed'): 'Low',
            ('Cardinal', 'Mutable'): 'High',
            ('Fixed', 'Fixed'): 'Low',
            ('Fixed', 'Mutable'): 'Medium',
            ('Mutable', 'Mutable'): 'High'
        }
        
        quality_pair = (info1['quality'], info2['quality'])
        if quality_pair not in quality_compatibility:
            quality_pair = (info2['quality'], info1['quality'])
        
        quality_comp = quality_compatibility.get(quality_pair, 'Medium')
        
        # Overall compatibility
        comp_scores = {'High': 3, 'Medium': 2, 'Low': 1}
        avg_score = (comp_scores[element_comp] + comp_scores[quality_comp]) / 2
        
        if avg_score >= 2.5:
            overall = 'High'
        elif avg_score >= 1.5:
            overall = 'Medium'
        else:
            overall = 'Low'
        
        return {
            'sign1': sign1,
            'sign2': sign2,
            'element_compatibility': element_comp,
            'quality_compatibility': quality_comp,
            'overall_compatibility': overall,
            'description': f"{sign1} ({info1['element']}, {info1['quality']}) and {sign2} ({info2['element']}, {info2['quality']})"
        }
    
    def calculate_midpoint(self, pos1: float, pos2: float) -> float:
        """
        Calculate the midpoint between two zodiac positions.
        
        Args:
            pos1: First position in degrees (0-360)
            pos2: Second position in degrees (0-360)
            
        Returns:
            Midpoint in degrees (0-360)
        """
        # Normalize positions
        pos1 = pos1 % 360
        pos2 = pos2 % 360
        
        # Calculate shortest arc between positions
        diff = abs(pos2 - pos1)
        if diff > 180:
            diff = 360 - diff
            # Use the longer arc
            midpoint = (pos1 + pos2 + 360) / 2
        else:
            midpoint = (pos1 + pos2) / 2
        
        return midpoint % 360
    
    def degrees_to_dms(self, degrees: float) -> Tuple[int, int, int]:
        """
        Convert decimal degrees to degrees, minutes, seconds.
        
        Args:
            degrees: Decimal degrees
            
        Returns:
            Tuple of (degrees, minutes, seconds)
        """
        deg = int(degrees)
        min_float = (degrees - deg) * 60
        minutes = int(min_float)
        seconds = int((min_float - minutes) * 60)
        
        return deg, minutes, seconds
    
    def format_zodiac_position(self, longitude: float, include_symbol: bool = True) -> str:
        """
        Format a zodiac position as a readable string.
        
        Args:
            longitude: Ecliptic longitude in degrees
            include_symbol: Whether to include the zodiac symbol
            
        Returns:
            Formatted position string
        """
        zodiac_info = self.ecliptic_to_zodiac(longitude)
        
        if include_symbol:
            return zodiac_info['position_string']
        else:
            deg, min_val, sec = self.degrees_to_dms(zodiac_info['degree'])
            return f"{deg}°{min_val:02d}'{sec:02d}\" {zodiac_info['name']}"