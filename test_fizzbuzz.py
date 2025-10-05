import pytest

from fizzbuzz import fizzbuzz, fizzbuzz_sequence


@pytest.mark.parametrize(
    "value, expected",
    [
        (1, "1"),
        (2, "2"),
        (3, "Fizz"),
        (4, "4"),
        (5, "Buzz"),
        (6, "Fizz"),
        (9, "Fizz"),
        (10, "Buzz"),
        (15, "FizzBuzz"),
        (30, "FizzBuzz"),
    ],
)
def test_fizzbuzz_value(value, expected):
    assert fizzbuzz(value) == expected


def test_fizzbuzz_sequence_15():
    expected = [
        "1",
        "2",
        "Fizz",
        "4",
        "Buzz",
        "Fizz",
        "7",
        "8",
        "Fizz",
        "Buzz",
        "11",
        "Fizz",
        "13",
        "14",
        "FizzBuzz",
    ]

    assert fizzbuzz_sequence(15) == expected


def test_fizzbuzz_sequence_zero_returns_empty():
    assert fizzbuzz_sequence(0) == []
