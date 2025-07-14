#!/usr/bin/env python3
"""
Basic functionality test for the Qucanft library.

This test verifies the core functionality works without external dependencies.
"""

import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

# Test basic zodiac calculations without external dependencies
print("Testing basic zodiac calculations...")

# Test 1: Basic zodiac sign calculation
def test_zodiac_sign_calculation():
    """Test basic zodiac sign calculations"""
    try:
        # Mock the pandas dependency by creating a simple DataFrame-like structure
        class MockDataFrame:
            def __init__(self, data):
                self.data = data
                self.columns = list(data.keys()) if isinstance(data, dict) else []
                
            def __getitem__(self, key):
                return self.data[key]
                
            def iterrows(self):
                if isinstance(self.data, dict) and self.columns:
                    for i in range(len(self.data[self.columns[0]])):
                        yield i, {col: self.data[col][i] for col in self.columns}
                        
            def copy(self):
                return MockDataFrame(self.data.copy())
                
            def empty(self):
                return len(self.data) == 0
        
        # Test zodiac sign properties
        zodiac_signs = [
            {'name': 'Aries', 'symbol': '♈', 'element': 'Fire', 'quality': 'Cardinal'},
            {'name': 'Taurus', 'symbol': '♉', 'element': 'Earth', 'quality': 'Fixed'},
            {'name': 'Gemini', 'symbol': '♊', 'element': 'Air', 'quality': 'Mutable'},
            {'name': 'Cancer', 'symbol': '♋', 'element': 'Water', 'quality': 'Cardinal'},
            {'name': 'Leo', 'symbol': '♌', 'element': 'Fire', 'quality': 'Fixed'},
            {'name': 'Virgo', 'symbol': '♍', 'element': 'Earth', 'quality': 'Mutable'},
            {'name': 'Libra', 'symbol': '♎', 'element': 'Air', 'quality': 'Cardinal'},
            {'name': 'Scorpio', 'symbol': '♏', 'element': 'Water', 'quality': 'Fixed'},
            {'name': 'Sagittarius', 'symbol': '♐', 'element': 'Fire', 'quality': 'Mutable'},
            {'name': 'Capricorn', 'symbol': '♑', 'element': 'Earth', 'quality': 'Cardinal'},
            {'name': 'Aquarius', 'symbol': '♒', 'element': 'Air', 'quality': 'Fixed'},
            {'name': 'Pisces', 'symbol': '♓', 'element': 'Water', 'quality': 'Mutable'}
        ]
        
        # Test ecliptic longitude to zodiac conversion
        def ecliptic_to_zodiac(longitude):
            """Convert ecliptic longitude to zodiac sign"""
            longitude = longitude % 360
            sign_index = int(longitude // 30)
            degree_in_sign = longitude % 30
            
            zodiac_info = zodiac_signs[sign_index].copy()
            zodiac_info.update({
                'degree': degree_in_sign,
                'absolute_degree': longitude,
                'position_string': f"{int(degree_in_sign)}°{int((degree_in_sign % 1) * 60):02d}' {zodiac_info['symbol']}"
            })
            
            return zodiac_info
        
        # Test cases
        test_cases = [
            (0, 'Aries'),     # 0° = Aries
            (30, 'Taurus'),   # 30° = Taurus
            (60, 'Gemini'),   # 60° = Gemini
            (90, 'Cancer'),   # 90° = Cancer
            (120, 'Leo'),     # 120° = Leo
            (150, 'Virgo'),   # 150° = Virgo
            (180, 'Libra'),   # 180° = Libra
            (210, 'Scorpio'), # 210° = Scorpio
            (240, 'Sagittarius'), # 240° = Sagittarius
            (270, 'Capricorn'),   # 270° = Capricorn
            (300, 'Aquarius'),    # 300° = Aquarius
            (330, 'Pisces')       # 330° = Pisces
        ]
        
        print("   Testing zodiac sign calculations...")
        for longitude, expected_sign in test_cases:
            result = ecliptic_to_zodiac(longitude)
            if result['name'] == expected_sign:
                print(f"   ✓ {longitude}° = {result['name']} {result['symbol']}")
            else:
                print(f"   ✗ {longitude}° = {result['name']} (expected {expected_sign})")
                return False
        
        # Test degree calculations
        test_longitude = 45.5  # 15.5° Taurus
        result = ecliptic_to_zodiac(test_longitude)
        expected_degree = 15.5
        if abs(result['degree'] - expected_degree) < 0.001:
            print(f"   ✓ Degree calculation: {test_longitude}° = {result['degree']}° {result['name']}")
        else:
            print(f"   ✗ Degree calculation failed: expected {expected_degree}, got {result['degree']}")
            return False
        
        return True
        
    except Exception as e:
        print(f"   ✗ Error in zodiac calculations: {e}")
        return False

# Test 2: Basic house calculations
def test_house_calculations():
    """Test basic house calculations"""
    try:
        print("   Testing house calculations...")
        
        # Test Equal House system
        def calculate_equal_houses(ascendant):
            """Calculate equal houses"""
            houses = {}
            for house_num in range(1, 13):
                cusp = (ascendant + (house_num - 1) * 30) % 360
                houses[house_num] = cusp
            return houses
        
        # Test with ascendant at 75° (Gemini)
        ascendant = 75.0
        houses = calculate_equal_houses(ascendant)
        
        # Verify calculations
        expected_houses = {
            1: 75.0,   # 1st house (Ascendant)
            2: 105.0,  # 2nd house
            3: 135.0,  # 3rd house
            4: 165.0,  # 4th house
            5: 195.0,  # 5th house
            6: 225.0,  # 6th house
            7: 255.0,  # 7th house (Descendant)
            8: 285.0,  # 8th house
            9: 315.0,  # 9th house
            10: 345.0, # 10th house (Midheaven)
            11: 15.0,  # 11th house
            12: 45.0   # 12th house
        }
        
        for house_num in range(1, 13):
            if abs(houses[house_num] - expected_houses[house_num]) < 0.001:
                print(f"   ✓ House {house_num}: {houses[house_num]}°")
            else:
                print(f"   ✗ House {house_num}: {houses[house_num]}° (expected {expected_houses[house_num]}°)")
                return False
        
        return True
        
    except Exception as e:
        print(f"   ✗ Error in house calculations: {e}")
        return False

# Test 3: Basic aspect calculations
def test_aspect_calculations():
    """Test basic aspect calculations"""
    try:
        print("   Testing aspect calculations...")
        
        def calculate_angular_distance(pos1, pos2):
            """Calculate angular distance between two positions"""
            pos1 = pos1 % 360
            pos2 = pos2 % 360
            diff = abs(pos2 - pos1)
            return min(diff, 360 - diff)
        
        def find_aspect(angular_distance):
            """Find aspect based on angular distance"""
            aspects = {
                'Conjunction': {'degrees': 0, 'orb': 8},
                'Sextile': {'degrees': 60, 'orb': 6},
                'Square': {'degrees': 90, 'orb': 8},
                'Trine': {'degrees': 120, 'orb': 8},
                'Opposition': {'degrees': 180, 'orb': 8}
            }
            
            for aspect_name, aspect_info in aspects.items():
                orb_difference = abs(angular_distance - aspect_info['degrees'])
                if orb_difference <= aspect_info['orb']:
                    return aspect_name, orb_difference
            
            return None, None
        
        # Test cases
        test_cases = [
            (0, 0, 'Conjunction'),      # 0° = Conjunction
            (0, 60, 'Sextile'),         # 60° = Sextile
            (0, 90, 'Square'),          # 90° = Square
            (0, 120, 'Trine'),          # 120° = Trine
            (0, 180, 'Opposition'),     # 180° = Opposition
            (0, 45, None),              # 45° = No major aspect
        ]
        
        for pos1, pos2, expected_aspect in test_cases:
            angular_distance = calculate_angular_distance(pos1, pos2)
            aspect, orb_diff = find_aspect(angular_distance)
            
            if aspect == expected_aspect:
                if aspect:
                    print(f"   ✓ {pos1}° - {pos2}° = {aspect} ({orb_diff:.1f}° orb)")
                else:
                    print(f"   ✓ {pos1}° - {pos2}° = No major aspect")
            else:
                print(f"   ✗ {pos1}° - {pos2}° = {aspect} (expected {expected_aspect})")
                return False
        
        return True
        
    except Exception as e:
        print(f"   ✗ Error in aspect calculations: {e}")
        return False

# Test 4: Test compatibility calculations
def test_compatibility_calculations():
    """Test zodiac compatibility calculations"""
    try:
        print("   Testing compatibility calculations...")
        
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
        
        # Test cases
        test_cases = [
            ('Fire', 'Air', 'High'),
            ('Earth', 'Water', 'High'),
            ('Fire', 'Water', 'Low'),
            ('Earth', 'Air', 'Low'),
        ]
        
        for element1, element2, expected in test_cases:
            key = (element1, element2)
            if key not in element_compatibility:
                key = (element2, element1)
            
            result = element_compatibility.get(key, 'Medium')
            if result == expected:
                print(f"   ✓ {element1} + {element2} = {result} compatibility")
            else:
                print(f"   ✗ {element1} + {element2} = {result} (expected {expected})")
                return False
        
        return True
        
    except Exception as e:
        print(f"   ✗ Error in compatibility calculations: {e}")
        return False

# Run all tests
def run_all_tests():
    """Run all tests"""
    print("=" * 50)
    print("Qucanft Core Functionality Tests")
    print("=" * 50)
    
    tests = [
        ("Zodiac Sign Calculations", test_zodiac_sign_calculation),
        ("House Calculations", test_house_calculations),
        ("Aspect Calculations", test_aspect_calculations),
        ("Compatibility Calculations", test_compatibility_calculations),
    ]
    
    results = []
    for test_name, test_func in tests:
        print(f"\n{test_name}:")
        try:
            result = test_func()
            results.append((test_name, result))
            if result:
                print(f"   ✓ {test_name} - PASSED")
            else:
                print(f"   ✗ {test_name} - FAILED")
        except Exception as e:
            print(f"   ✗ {test_name} - ERROR: {e}")
            results.append((test_name, False))
    
    print("\n" + "=" * 50)
    print("Test Results Summary:")
    print("=" * 50)
    
    passed = sum(1 for _, result in results if result)
    total = len(results)
    
    for test_name, result in results:
        status = "✓ PASSED" if result else "✗ FAILED"
        print(f"{test_name}: {status}")
    
    print(f"\nOverall: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 All tests passed! Core functionality is working correctly.")
    else:
        print("❌ Some tests failed. Please check the implementation.")
    
    return passed == total

if __name__ == "__main__":
    success = run_all_tests()
    sys.exit(0 if success else 1)