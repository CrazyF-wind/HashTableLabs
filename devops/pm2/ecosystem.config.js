module.exports = {
  apps: [
    {
      name: 'app1',
      script: './app.js',
      watch: true,
      env: { PORT: 3000 },
    },
    {
      name: 'app2',
      script: './app.js',
      watch: false,
      env: { PORT: 3001 },
    },
  ],
};
