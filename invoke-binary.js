const { spawnSync } = require('child_process');
const { chmodSync } = require('fs');
const path = require('path');

function chooseBinary() {
    const platform = process.platform;
    const arch = process.arch;

    let binaryName;

    if (platform === 'linux') {
        if (arch === 'x64') {
            binaryName = 'fm-sync-linux-amd64';
        } else if (arch === 'arm64') {
            binaryName = 'fm-sync-linux-arm64';
        } else {
            throw new Error(`Unsupported architecture for Linux: ${arch}`);
        }
    } else if (platform === 'win32') {
        if (arch === 'x64') {
            binaryName = 'fm-sync-windows-amd64.exe';
        } else if (arch === 'arm64') {
            binaryName = 'fm-sync-windows-arm64.exe';
        } else {
            throw new Error(`Unsupported architecture for Windows: ${arch}`);
        }
    } else if (platform === 'darwin') {
        if (arch === 'x64') {
            binaryName = 'fm-sync-darwin-amd64';
        } else if (arch === 'arm64') {
            binaryName = 'fm-sync-darwin-arm64';
        } else {
            throw new Error(`Unsupported architecture for macOS: ${arch}`);
        }
    } else {
        throw new Error(`Unsupported platform: ${platform}`);
    }

    return binaryName;
}

function main() {
    try {
        const binaryName = chooseBinary();
        const binaryPath = path.join(__dirname, binaryName);

        // Ensure binary is executable on Unix-like systems
        if (process.platform !== 'win32') {
            try {
                chmodSync(binaryPath, '755');
            } catch (err) {
                // Ignore errors if already executable or permissions issue
            }
        }

        console.log(`Executing binary: ${binaryName}`);

        const result = spawnSync(binaryPath, [], {
            stdio: 'inherit',
            encoding: 'utf-8'
        });

        if (result.error) {
            console.error(`Failed to execute binary: ${result.error.message}`);
            process.exit(1);
        }

        process.exit(result.status || 0);
    } catch (err) {
        console.error(`Error: ${err.message}`);
        process.exit(1);
    }
}

main();
