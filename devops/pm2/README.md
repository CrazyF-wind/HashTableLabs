# pm2 学习笔记

对应文章：[pm2学习笔记](https://crazyf-wind.github.io/2019/07/03/pm2学习笔记/)

## 文件

- `app.js`：一个可优雅退出的 HTTP 服务
- `ecosystem.config.js`：启动两个实例的配置

## 运行

```bash
npm install -g pm2
cd devops/pm2

pm2 start ecosystem.config.js
pm2 status
pm2 logs app1

curl http://localhost:3000
curl http://localhost:3001

pm2 reload app1     # 更新：先起新进程，再停旧进程
pm2 stop app1
pm2 delete all
```

## 优雅退出

`app.js` 监听了 `shutdown` 消息：

```js
process.on('shutdown', () => {
  server.close(() => process.exit(0));
});
```

文章中的 `pm2 gracefulReload` 会向旧进程发送该消息。不同 pm2 版本对这个子命令的支持有差异，若本机不支持，用 `pm2 reload` 代替，并以 `pm2 --help` 为准。
