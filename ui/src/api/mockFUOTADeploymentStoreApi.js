import delay from './delay';

const get = {
  "fuotaDeployment": {
    "id": "1361C787-1EC4-48C2-86A7-2A7BA39AD4F3",
    "name": "Test Deployment 1",
    "createdAt": "2019-12-06 15:46:12.624982",
    "updatedAt": "2019-12-06 15:46:12.624982",
    "groupType": "CLASS_C",
    "dr": 1,
    "frequency": 100,
    "redundancy": 100,
    "multicastTimeout": 100,
    "unicastTimeout": "100",
    "state": "1",
    "nextStepAfter": "2019-12-06 15:46:12.624982",
  }
};

const list = {
  result: [{
    "id": "1361C787-1EC4-48C2-86A7-2A7BA39AD4F3",
    "createdAt": "2019-12-06 15:46:12.624982",
    "updatedAt": "2019-12-06 15:46:12.624982",
    "name": "Test Deployment 1",
    "state": "1",
    "nextStepAfter": "2019-12-06 15:46:12.624982"
  }],
  totalCount: 1
}

const listDeploymentDevices = {
  result: [{
    "devEUI": "70b3d5fffe1cb547",
    "deviceName": "Test Deployment 1",
    "state": "1",
    "errorMessage": "",
    "createdAt": "2019-12-06 15:46:12.624982",
    "updatedAt": "2019-12-06 15:46:12.624982",
  }]
}

const getDeploymentDevice = {
  deploymentDevice: {
    "devEUI": "70b3d5fffe1cb547",
    "deviceName": "Test Deployment 1",
    "state": "1",
    "errorMessage": "",
    "createdAt": "2019-12-06 15:46:12.624982",
    "updatedAt": "2019-12-06 15:46:12.624982",
  }
}

class MockFUOTADeploymentStoreApi {
  static get() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, get));
      }, delay);
    });
  }

  static list() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, list));
      }, delay);
    });
  }

  static listDeploymentDevices() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, listDeploymentDevices));
      }, delay);
    });
  }

  static getDeploymentDevice() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, getDeploymentDevice));
      }, delay);
    });
  }
}

export default MockFUOTADeploymentStoreApi;
