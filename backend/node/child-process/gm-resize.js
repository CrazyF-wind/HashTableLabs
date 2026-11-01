// gm 并没有自己实现图像处理，它通过 cross-spawn 调用 ImageMagick 的 convert 命令。
// 所以运行前需要：npm install && 安装 ImageMagick（提供 convert）。
const fs = require('fs');
const path = require('path');
const gm = require('gm').subClass({ imageMagick: true });

const input = process.argv[2];

if (!input) {
  console.error('用法: node gm-resize.js <图片路径>');
  process.exit(1);
}

if (!fs.existsSync(input)) {
  console.error(`文件不存在: ${input}`);
  process.exit(1);
}

const ext = path.extname(input);
const output = path.join(path.dirname(input), `${path.basename(input, ext)}-thumb${ext}`);

gm(input)
  .resize(160, 120)
  .write(output, (err) => {
    if (err) {
      console.error(`处理失败: ${err.message}`);
      process.exit(1);
    }
    console.log(`已生成缩略图: ${output}`);
  });
