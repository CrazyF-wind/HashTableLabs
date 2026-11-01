// 同一个函数从最啰嗦到最简的五种写法，以及一个排序比较函数的简写。
const fnV1 = function (n) {
  return function (m) {
    return m + n;
  };
};

const fnV2 = (n) => {
  return (m) => {
    return m + n;
  };
};

const fnV3 = (n) => (m) => {
  return m + n;
};

const fnV4 = (n) => (m) => m + n;

const fnV5 = n => m => m + n;

// 闭包记住字段名 p，返回一个可以直接交给 Array.prototype.sort 的比较函数。
const compare = p => (m, n) => m[p] - n[p];

const people = [
  { name: '小明', age: 12 },
  { name: '小红', age: 11 },
  { name: '小刚', age: 15 },
  { name: '小华', age: 13 },
];

function main() {
  const versions = [fnV1, fnV2, fnV3, fnV4, fnV5];
  console.log('五种写法的结果：', versions.map((fn) => fn(5)(3)));
  console.log('按 age 升序：', [...people].sort(compare('age')));
}

if (require.main === module) {
  main();
}

module.exports = { fnV1, fnV2, fnV3, fnV4, fnV5, compare, people };
