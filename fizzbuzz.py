from typing import List


def generate_fizzbuzz(start: int = 1, end: int = 100) -> List[str]:
    if end < start:
        raise ValueError("end must be >= start")

    results: List[str] = []
    for number in range(start, end + 1):
        token = ""
        if number % 3 == 0:
            token += "Fizz"
        if number % 5 == 0:
            token += "Buzz"
        results.append(token or str(number))

    return results


def main() -> None:
    print("\n".join(generate_fizzbuzz()))


if __name__ == "__main__":
    main()
