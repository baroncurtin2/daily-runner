import sys
from pathlib import Path
from typing import TextIO


def process_line(line: str, found_python: bool, output_file: TextIO) -> bool:
    """Porcess a single line from the numbers file"""
    line = line.strip()

    if found_python:
        output_file.write(f"{int(line) + 1}\n")
        return False

    if line == "python":
        found_python = True

    output_file.write(f"{line}\n")
    return found_python


def increment_number(numbers_file: Path, temp_file: Path) -> None:
    """Increment the Python number in the file"""
    try:
        with numbers_file.open("r") as input_file, temp_file.open("w") as output_file:
            found_python = False

            for line in input_file:
                found_python = process_line(line, found_python, output_file)

        temp_file.replace(numbers_file)
    except (IOError, ValueError) as e:
        sys.exit(f"Error processing file: {e}")
    except Exception as e:
        sys.exit(f"Unexpected error: {e}")
    finally:
        if temp_file.exists():
            temp_file.unlink(missing_ok=True)


def main() -> None:
    """Main entry point for the script"""
    root_dir = Path(__file__).parent.parent
    numbers_file = root_dir / "numbers.txt"
    temp_file = root_dir / "numbers.tmp"

    increment_number(numbers_file, temp_file)


if __name__ == "__main__":
    main()
