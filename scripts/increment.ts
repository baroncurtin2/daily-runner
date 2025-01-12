import * as fs from 'fs';
import * as path from 'path';

function processLine(line: string, foundTypeScript: boolean, outputFile: fs.WriteStream): boolean {
    const trimmedLine = line.trim();

    if (foundTypeScript) {
        const number = parseInt(trimmedLine, 10);
        if (isNaN(number)) {
            throw new Error('Invalid number format under "javascript"');
        }
        outputFile.write(`${number + 1}\n`);
        return false; // Reset flag after processing
    }

    if (trimmedLine === 'typescript') {
        return true; // Set the flag when "typescript" is found
    }

    outputFile.write(`${line}\n`);
    return foundTypeScript;
}

function incrementNumber(numbersFile: string, tempFile: string): void {
    try {
        const content = fs.readFileSync(numbersFile, 'utf8');
        const lines = content.split('\n');

        const outputFile = fs.createWriteStream(tempFile);
        let foundTypeScript = false;

        for (let line of lines) {
            foundTypeScript = processLine(line, foundTypeScript, outputFile);
        }

        outputFile.end(); // Close the output file after writing
        fs.renameSync(tempFile, numbersFile); // Replace the original file with the temp file

        console.log('Successfully incremented the number for "typescript" in "numbers.txt"');
    } catch (error) {
        console.error(`Error processing file: ${error.message}`);
        if (fs.existsSync(tempFile)) {
            fs.unlinkSync(tempFile); // Clean up temp file in case of an error
        }
    }
}


function main() {
    const rootDir = __dirname;
    const numbersFile = path.join(rootDir, 'numbers.txt');
    const tempFile = path.join(rootDir, 'numbers.tmp');

    incrementNumber(numbersFile, tempFile)
}

main();
