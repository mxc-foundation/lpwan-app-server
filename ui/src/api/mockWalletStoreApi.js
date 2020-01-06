import delay from './delay';

const getDlPrice = [
  {
    "downLinkPrice": 100,
    "userProfile": {}
  }
];

class MockWalletStoreApi {
  static getDlPrice(orgId) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign([], getDlPrice));
      }, delay);
    });
  }
}

export default MockWalletStoreApi;
