#!/bin/bash
# 守护脚本：进程不存在时，在 tmux 会话里重新拉起。
# 用法：SERVER_USER=nodejsUser SERVER_SCRIPT=/opt/nodejsUser/server.js ./server.sh

export DISABLE_AUTO_TITLE="true"

SESSION="${SESSION:-server}"
SERVER_USER="${SERVER_USER:-nodejsUser}"
SERVER_SCRIPT="${SERVER_SCRIPT:-/opt/nodejsUser/server.js}"

if ps -ef | grep "node ${SERVER_SCRIPT}" | grep -v grep > /dev/null; then
    echo "running....."
    exit 0
fi

echo "start process....."

if ! tmux has-session -t "${SESSION}" 2>/dev/null; then
    tmux new-session -d -s "${SESSION}"
fi

tmux send-keys -t "${SESSION}" "sudo -u ${SERVER_USER} node ${SERVER_SCRIPT}" C-m
