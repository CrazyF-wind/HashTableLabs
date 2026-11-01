# Linux 服务器守护进程

对应文章：[Linux服务器守护进程方案及实现](https://crazyf-wind.github.io/2019/10/10/Linux服务器守护进程方案及实现/)

方案：无登录账号 + tmux 保活 + crontab 每分钟检查。

## 文件

- `server.js`：被守护的 Node 服务
- `server.sh`：检查进程，不存在就用 tmux 拉起

## 本地验证 server.js

```bash
cd devops/linux-daemon
node server.js
curl http://localhost:8080
# {"server":"nodejs"}
```

## 部署成守护进程

以下步骤需要 root 权限，并确认机器上有 tmux 与 crontab。

1. 建无登录账号：

   ```bash
   useradd -M -s /sbin/nologin nodejsUser
   ```

2. 把 `server.js` 放到 `/opt/nodejsUser/` 并授权：

   ```bash
   chown nodejsUser /opt/nodejsUser/server.js
   ```

3. 把 `server.sh` 也放过去，手动执行一次：

   ```bash
   SERVER_USER=nodejsUser SERVER_SCRIPT=/opt/nodejsUser/server.js ./server.sh
   tmux attach -t server
   ```

4. 加入 cron：

   ```bash
   crontab -e
   ```

   ```cron
   */1 * * * * /bin/bash /opt/nodejsUser/server.sh >> /opt/nodejsUser/run.log 2>&1
   ```

## 相对文章的两处修正

- 文章代码里的中文弯引号已改成标准引号，否则 JS 和 Shell 都无法执行。
- 原脚本在会话已存在时会直接 `tmux attach`；cron 没有终端，attach 会失败。这里改为：会话存在就只发送启动命令，需要查看时手动 `tmux attach`。
