import fs from 'fs';
import path from 'path';

function processLine(line, foundJS, outputFile) {
    const trimmedLine = line.trim();

    if (foundJS) {
        // increment the number if 'javascript' was found
        const number = parseInt(trimmedLine, 10);
        outputFile.write(`${number + 1}\n`);
        return false;
    }

    if (trimmedLine === 'javascript') {
        return true
    }

    outputFile.write(`${line}\n`);
    return foundJS;
}

function incrementNumber(numbersFile, tempFile) {
    try {
        const inputFile = fs.readFileSync(numbersFile, 'utf8');
        const lines = inputFile.split('\n');

        const outputFile = fs.createWriteStream(tempFile);
        let foundJS = false;

        lines.forEach(line => {
            foundJS = processLine(line, foundJS, outputFile);
        });

        outputFile.end();
        fs.renameSync(tempFile, numbersFile);
        console.log("Successfully incremented the number for 'javascript' in numbers.txt)
    } catch (error) {
        console.error(`Error processing file: ${error.message}`);
        if (fs.existsSync(tempFile)) {
            fs.unlinkSync(tempFile);
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