module.exports = {
  '**/*.{ts,tsx}': 'eslint',
  '**/*.ts?(x)': () => 'tsc --noEmit',
};
