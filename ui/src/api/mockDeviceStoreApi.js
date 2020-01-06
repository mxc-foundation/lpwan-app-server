import delay from './delay';

const get = {
  // device object
  "device": {
    "devEUI": "70b3d5fffe1cb547",
    "name": "My Device",
    "applicationID": "2",
    "description": "My Device Description",
    "deviceProfileID": "1361C787-1EC4-48C2-86A7-2A7BA39AD4F3",
    "skipFCntCheck": false,
    "referenceAltitude": 2.1,
    "variables": {
      my_variable_key_1: "my variable value 1",
      my_variable_key_2: "my variable value 2"
    },
    "tags": {
      my_tag_key_1: "my tag value 1",
      my_tag_key_2: "my tag value 2"
    },
  },
  "lastSeenAt": "2019-12-06 15:46:12.624982",
  "deviceStatusBattery": "254",
  "deviceStatusMargin": "32",
  // commonLocation
  "location": {
    "latitude": 150.1,
    "longitude": 200.5,
    "altitude": 4000,
    "source": "GEO_RESOLVER",
    "accuracy": 10
  }
};

const getDeviceList = {
  "devProfile": [
    {
      "id": "1",
      "devEui": "70b3d5fffe1cb777",
      "fkWallet": "123",
      "mode": "DV_WHOLE_NETWORK",
      "createdAt": "2019-12-06 15:44:58.767722",
      "lastSeenAt": "2019-12-06 15:44:58.767722",
      // "applicationId": "2",
      "name": "My Device M2M - WITHOUT Application"
    },
    {
      "id": "2",
      "devEui": "70b3d5fffe1cb778",
      "fkWallet": "123",
      "mode": "DV_WHOLE_NETWORK",
      "createdAt": "2019-12-06 15:44:58.767722",
      "lastSeenAt": "2019-12-06 15:44:58.767722",
      "applicationId": "2",
      "name": "My Device M2M - WITH Application"
    }
  ],
  "count": 2
};

const getDeviceActivation = {
  "deviceActivation": {
    "id": "1",
    "devEUI": "70b3d5fffe1cb777",
    "devAddr": "abc12345",
    "appSKey": "abc45677778888999900009999888877",
    "nwkSEncKey": "abc78977778888999900009999888877",
    "sNwkSIntKey": "abc12377778888999900009999888877",
    "fNwkSIntKey": "abc34577778888999900009999888877",
    "fCntUp": 10,
    "nFCntDown": 20,
    "aFCntDown": 30
  }
};

class MockDeviceStoreApi {
  static get() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, get));
      }, delay);
    });
  }

  static getDeviceList() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, getDeviceList));
      }, delay);
    });
  }

  static getDeviceActivation() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, getDeviceActivation));
      }, delay);
    });
  }
}

export default MockDeviceStoreApi;
