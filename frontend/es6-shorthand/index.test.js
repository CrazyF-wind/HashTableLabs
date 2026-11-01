const test = require('node:test');
const assert = require('node:assert');

const { fnV1, fnV2, fnV3, fnV4, fnV5, compare, people } = require('./index.js');

test('五种写法的结果一致', () => {
  for (const fn of [fnV1, fnV2, fnV3, fnV4, fnV5]) {
    assert.strictEqual(fn(5)(3), 8);
  }
});

test('compare 按指定字段升序比较', () => {
  const sorted = [...people].sort(compare('age'));
  assert.deepStrictEqual(sorted.map((p) => p.age), [11, 12, 13, 15]);
});

test('compare 只负责比较，不修改传入对象', () => {
  assert.strictEqual(compare('age')({ age: 3 }, { age: 5 }), -2);
});
