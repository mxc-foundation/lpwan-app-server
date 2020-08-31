const express = require('express');
const bodyParser = require('body-parser');
const path = require('path');
const {createProxyMiddleware} = require('http-proxy-middleware');

const app = express();
app.use(express.static(path.join(__dirname, 'build')));
// real service is at https://appserver:8080
app.use('/', createProxyMiddleware({target: 'http://appserver:8080', changeOrigin: true}));

// forward all request to react-app
app.get('*', function (req, res) {
    res.sendFile(path.join(__dirname, 'build', 'index.html'));
});

app.listen(8080, () =>
    console.log('Express server is running on :8080'));