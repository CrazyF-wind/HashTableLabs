# Node.js 子进程与 gm

对应文章：[Node.js 下使用 gm 模块的一些思考](https://crazyf-wind.github.io/2019/11/14/Node-js-下使用-gm-模块的一些思考/)

文章从 `child_process` 讲起，说明 `gm` 本质上也是通过子进程调用 ImageMagick。

## 文件

- `spawn-ls.js`：零依赖，用 `spawn` 调用系统 `ls`
- `gm-resize.js`：用 `gm` 生成缩略图，依赖 `gm` 与 ImageMagick

## 运行

```bash
cd backend/node/child-process

# 零依赖，直接运行
node spawn-ls.js
node spawn-ls.js /tmp

# gm 示例
npm install
convert -size 400x300 gradient:blue-white sample.png
node gm-resize.js sample.png
# 生成 sample-thumb.png
```

若本机没有 ImageMagick，`gm-resize.js` 会报找不到 `convert` 命令——这正是文章里遇到的错误。

> 说明：`gm` 已于 2025-02 停止维护，这里保留它只是为了复现文章中的结论；新项目应直接调用 ImageMagick 或使用仍在维护的图像库。
