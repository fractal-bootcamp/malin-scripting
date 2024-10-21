#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

function getExecutableName() {
  const platform = os.platform();
  const arch = os.arch();

  if (platform === 'win32') {
    return 'go-fullstack-cli.exe';
  } else if (platform === 'darwin') {
    return 'go-fullstack-cli-macos';
  } else if (platform === 'linux') {
    return 'go-fullstack-cli-linux';
  }

  throw new Error(`Unsupported platform: ${platform}`);
}

const executablePath = path.join(__dirname, 'bin', getExecutableName());

const child = spawn(executablePath, process.argv.slice(2), { stdio: 'inherit' });

child.on('error', (error) => {
  console.error(`Error: ${error.message}`);
  process.exit(1);
});

child.on('close', (code) => {
  process.exit(code);
});