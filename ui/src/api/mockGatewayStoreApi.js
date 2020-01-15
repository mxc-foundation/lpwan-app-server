import delay from './delay';

const getGatewayList = {
  "gwProfile": [
    {
      "id": "1"
    }
  ],
  "count": "1",
  "userProfile": {}
};

class MockGatewayStoreApi {
  static getGatewayList(orgId) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, getGatewayList));
      }, delay);
    });
  }
}

export default MockGatewayStoreApi;
