// 用 child_process.spawn 调用系统命令。
// 文章示例固定查看 /Users，这里改为：命令行参数 > HOME > 当前目录。
const { spawn } = require('child_process');

const target = process.argv[2] || process.env.HOME || '.';

const ls = spawn('ls', ['-lh', target]);

// 控制台正常输出
ls.stdout.on('data', (data) => {
  console.log(`stdout:\n${data}`);
});

// 控制台异常输出
ls.stderr.on('data', (data) => {
  console.error(`stderr: ${data}`);
});

// 关闭
ls.on('close', (code) => {
  console.log(`child process exited with code ${code}`);
});
