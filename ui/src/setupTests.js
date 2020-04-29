const {configure} = require('enzyme');
const Adapter = require('enzyme-adapter-react-16');
const initializeUnhandledRejection  = require('./hackfixes').initializeUnhandledRejection;
// const axios = require('axios');

console.log('setupTests');

initializeUnhandledRejection();


// hackfix for jsdom/axios (see https://github.com/axios/axios/issues/1754#issuecomment-572778305)
// axios.defaults.adapter = require('axios/lib/adapters/http');

configure({ adapter: new Adapter() });
// global.XMLHttpRequest = null;