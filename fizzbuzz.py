from typing import List


def fizzbuzz(value: int) -> str:
    """Return Fizz/Buzz/FizzBuzz or the number as a string.

    - Multiples of 3 => "Fizz"
    - Multiples of 5 => "Buzz"
    - Multiples of 3 and 5 => "FizzBuzz"
    - Otherwise => the decimal string of the value
    """
    if value % 15 == 0:
        return "FizzBuzz"
    if value % 3 == 0:
        return "Fizz"
    if value % 5 == 0:
        return "Buzz"
    return str(value)


def fizzbuzz_sequence(count: int) -> List[str]:
    """Generate the FizzBuzz sequence from 1 to `count` (inclusive).

    Returns an empty list when `count` is less than 1.
    """
    if count < 1:
        return []
    return [fizzbuzz(i) for i in range(1, count + 1)]
