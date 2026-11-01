# HashTableLabs

技术博客 [fwind's blog](https://crazyf-wind.github.io) 中示例代码的可运行版本。

仓库按**技术领域**组织，每个示例自包含：进入对应目录，按目录内 README 运行即可。文章本身没有程序源码的，不在本仓库收录。

## 示例索引

| 领域 | 示例 | 对应文章 | 语言 | 快速运行 |
| --- | --- | --- | --- | --- |
| backend | [hash-table](backend/go/hash-table/) | [从零实现一个哈希表，理解它的每一个设计决策](https://crazyf-wind.github.io/2026/10/30/从零实现一个哈希表，理解它的每一个设计决策/) | Go | `cd backend/go/hash-table && go test ./...` |
| backend | [map-concurrency](backend/go/map-concurrency/) | [一次 map 并发写崩溃，让我重新理解了 Go runtime](https://crazyf-wind.github.io/2026/10/26/一次%20map%20并发写崩溃，让我重新理解了%20Go%20runtime/) | Go | `cd backend/go/map-concurrency && go run ./reproduce compare` |
| backend | [child-process](backend/node/child-process/) | [Node.js 下使用 gm 模块的一些思考](https://crazyf-wind.github.io/2019/11/14/Node-js-下使用-gm-模块的一些思考/) | Node.js | `cd backend/node/child-process && node spawn-ls.js` |
| frontend | [es6-shorthand](frontend/es6-shorthand/) | [ES6函数简写及思考](https://crazyf-wind.github.io/2019/09/26/ES6函数简写及思考/) | JavaScript | `cd frontend/es6-shorthand && node --test index.test.js` |
| devops | [linux-daemon](devops/linux-daemon/) | [Linux服务器守护进程方案及实现](https://crazyf-wind.github.io/2019/10/10/Linux服务器守护进程方案及实现/) | Node.js / Bash | `cd devops/linux-daemon && node server.js` |
| devops | [pm2](devops/pm2/) | [pm2学习笔记](https://crazyf-wind.github.io/2019/07/03/pm2学习笔记/) | Node.js / pm2 | `cd devops/pm2 && pm2 start ecosystem.config.js` |

## 环境要求

| 工具 | 版本 | 用在哪里 |
| --- | --- | --- |
| Go | 1.23+ | `backend/go` |
| Node.js | 18+（示例在 24.x 验证） | `backend/node`、`frontend`、`devops` |
| tmux | 任意近期版本 | `devops/linux-daemon` |
| pm2 | 任意近期版本 | `devops/pm2` |
| ImageMagick | 需提供 `convert` 命令 | `backend/node/child-process` 的 gm 示例 |

## 约定

- 每个示例目录自带构建文件（`go.mod` 或 `package.json`），依赖不跨目录共享。
- 各目录 README 中的命令默认从仓库根目录开始写。
- 文章里的代码是讲解用的片段；本仓库中的版本是能实际运行的版本，个别地方修正了中文引号、笔误等问题，并在对应 README 中注明。
- 会崩溃的示例（并发 map 复现）是刻意为之，属于演示的一部分。
