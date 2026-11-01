const http = require('http');

const port = Number(process.env.PORT) || 3000;

const server = http.createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ pid: process.pid }));
});

server.listen(port, () => {
  console.log(`app listening on ${port}, pid ${process.pid}`);
});

// pm2 gracefulReload 会向进程发送 'shutdown' 消息。
process.on('shutdown', () => {
  console.log('收到 shutdown，停止接收新连接');
  server.close(() => {
    console.log('连接已关闭，退出');
    process.exit(0);
  });
});

// pm2 stop / restart 发送 SIGINT，这里复用同一套优雅退出逻辑。
process.on('SIGINT', () => {
  process.emit('shutdown');
});
