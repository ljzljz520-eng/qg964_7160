const fs = require('fs');
const path = require('path');
const source = fs.readFileSync(path.join(__dirname, 'index.html'), 'utf8');
const output = path.join(__dirname, 'dist');
fs.mkdirSync(output, { recursive: true });
fs.writeFileSync(path.join(output, 'index.html'), source.replace('app.js', 'app.js'));
fs.copyFileSync(path.join(__dirname, 'app.js'), path.join(output, 'app.js'));
process.stdout.write('built cold storage risk console web\n');
