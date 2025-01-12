package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrNoGolangLine = errors.New("golang line not found in file")
	ErrInvalidNumber = errors.New("invalid number format after golang line")
)

type File struct {
	path string
	temp string
}

// New creates a new File instances with proper path initiatlization
func New() (*File, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("finding project root: %w", err)
	}

	return &File {
		path: filepath.Join(root, "numbers.txt"),
		temp: filepath.Join(root, "numbers.txt.tmp"),
	}, nil
}

// findProjectRoot locates the root directory containing numbers.txt
func findProjectRoot() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "",  fmt.Errorf("getting executable path: %w, err")
	}

	return filepath.Dir(filepath.Dir(exe)), nil
}

func processLine(line string, foundGolang bool, w *bufio.Writer) (bool, error) {
	line = strings.TrimSpace(line)

	if foundGolang {
		num, err := strconv.Atoi(line)
		if err != nil {
			return false, fmt.Errorf("%w: %v", ErrInvalidNumber, err)
		}

		if err := writeLine(w, strconv.Itoa(num+1)); err != nil {
			return false, err
		}

		return false, nil
	}

	if line == "golang" {
		foundGolang = true
	}

	if err := writeLine(w, line); err != nil {
		return false, err
	}
}

// writeLine writes a line to the writer with proper error handling
func writeLine(w *bufio.Writer, line string) error {
	_, err := fmt.Fprintln(w, line)
	return err
}

// process handles the file processing logic
func (f *File) process() error {
	input, err := os.Open(f.path)
	if err != nil {
		return fmt.Errorf("opening input file: %w", err)
	}
	defer input.Close()

	output, err := os.OpenFile(f.temp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer output.Close()

	if err := f.processFiles(input, output); err != nil {
		return err
	}

	// close files before rename
	input.Close()
	output.Close()

	return f.finalize()
}

// processFiles handles the core file processing logic
func (f *File) processFiles(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	var foundGolang bool
	var foundNumber bool

	for scanner.Scan() {
		var err error
		foundGolang, err = processLine(scanner.Text(), foundGolang, writer)
		if err != nil {
			return err
		}

		if foundGolang {
			foundNumber = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning file: %w", err)
	}

	if !foundNumber {
		return ErrNoGolangLine
	}

	return writer.Flush()
}

// finalize performaces the atomic file replacement
func (f *File) finalize() error {
	if err := os.Rename(f.temp, f.path); err != nil {
		return fmt.Errorf("replacing original file: %w", err)
	}

	return nil
}

// cleanup removes the temporary file if it exists
func (f *File) cleanup() {
	os.Remove(f.temp)
}

func run() error {
	f, err := New()
	if err != nil {
		return err
	}

	defer f.cleanup()

	if err := f.process(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}