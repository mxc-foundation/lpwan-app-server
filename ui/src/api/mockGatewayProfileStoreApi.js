import delay from './delay';

const get = {
  "gatewayProfile": {
    "apiGatewayProfile": {
      "id": "0b2767d2-6315-4a14-a014-082b83954508",
      "name": "My Gateway Profile",
      "networkServerID": "1",
      "channels": [10, 20],
      "extraChannels": [
        {
          "modulation": "LORA",
          "frequency": 30,
          "bandwidth": 1000,
          "bitrate": 50,
          "spreadingFactors": [12,13]
        }
      ]
    }
  },
  "createdAt": "2019-12-06 15:46:12.624982",
  "updatedAt": "2019-12-06 15:46:12.624982"
};

class MockGatewayProfileStoreApi {
  static get(id) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, get));
      }, delay);
    });
  }
}

export default MockGatewayProfileStoreApi;
