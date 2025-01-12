use std::fs::{self, File};
use std::io::{self, BufRead, BufReader, Write};
use std::path::Path;

fn process_line(line: &str, found_rust: &mut bool, output_file: &mut File) -> io::Result<()> {
    let trimmed_line = line.trim();

    if *found_rust {
        let number: i32 = trimmed_line.parse().map_err(|_| io::Error::new(io::ErrorKind::InvalidData, "Expected a valid number"))?;
        writeln!(output_file, "{}", number + 1)?;
        *found_rust = false;
    } else {
        if trimmed_line == "rust" {
            *found_rust = true;
        }
        writeln!(output_file, "{}", trimmed_line)?;
    }
    Ok(())
}

fn increment_number(numbers_file: &Path, temp_file: &Path) -> io::Result<()> {
    let input_file = File::open(numbers_file)?;
    let mut output_file = File::create(temp_file)?;

    let reader = BufReader::new(input_file);
    let mut found_rust = false;

    for line in reader.lines() {
        process_line(&line?, &mut found_rust, &mut output_file)?;
    }

    drop(output_file); // Ensure the output file is closed before renaming
    fs::rename(temp_file, numbers_file)?;
    println!("Successfully incremented the number for 'rust' in 'numbers.txt'");
    Ok(())
}

fn main() -> io::Result<()> {
    let root_dir = Path::new(env!("CARGO_MANIFEST_DIR"));
    let numbers_file = root_dir.join("numbers.txt");
    let temp_file = root_dir.join("numbers.tmp");

    if let Err(e) = increment_number(&numbers_file, &temp_file) {
        eprintln!("Error processing file: {}", e);
        if temp_file.exists() {
            fs::remove_file(&temp_file)?;
        }
    }

    Ok(())
}
