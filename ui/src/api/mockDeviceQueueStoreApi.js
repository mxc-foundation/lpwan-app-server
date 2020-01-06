import delay from './delay';

const getDeviceQueueList = {
  "deviceQueueItems": [
    {
      "devEUI": "70b3d5fffe1cb777",
      "confirmed": true,
      "fCnt": 10,
      "fPort": 1,
      "data": "",
      "jsonObject": ""
    },
    {
      "devEUI": "70b3d5fffe1cb778",
      "confirmed": false,
      "fCnt": 15,
      "fPort": 2,
      "data": "",
      "jsonObject": ""
    }
  ]
};

class MockDeviceQueueStoreApi {
  static getDeviceQueueList() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, getDeviceQueueList));
      }, delay);
    });
  }
}

export default MockDeviceQueueStoreApi;
