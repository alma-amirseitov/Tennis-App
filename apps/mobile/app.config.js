const { expo } = require('./app.json');
module.exports = {
  expo: {
    ...expo,
    extra: {
      apiUrl: process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080',
    },
  },
};
