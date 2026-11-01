const http = require('http');

const port = Number(process.env.PORT) || 8080;

const server = http.createServer((req, res) => {
  res.writeHead(200, {
    'Content-Type': 'application/json',
  });
  res.end('{"server":"nodejs"}');
});

server.listen(port, () => {
  console.log(`port = ${port}`);
});
