import delay from './delay';

const get = {
  "deviceProfile": {
    "id": "1361C787-1EC4-48C2-86A7-2A7BA39AD4F3",
    "name": "My Device Profile",
    "organizationID": "2",
    "networkServerID": "1",
    "supportsClassB": true,
    "classBTimeout": 5,
    "pingSlotPeriod": 50,
    "pingSlotDR": 50,
    "pingSlotFreq": 50,
    "supportsClassC": true,
    "classCTimeout": 5,
    "macVersion": "1.0.1",
    "regParamsRevision": "B",
    "rxDelay1": 5,
    "rxDROffset1": 20,
    "rxDataRate2": 150,
    "rxFreq2": 40,
    "factoryPresetFreqs": [10, 20],
    "maxEIRP": 50,
    "maxDutyCycle": 100,
    "supportsJoin": false,
    "rfRegion": "region1",
    "supports32BitFCnt": true,
    "payloadCodec": "CAYENNE_LPP",
    "payloadEncoderScript": "",
    "payloadDecoderScript": "",
    "geolocBufferTTL": 3,
    "geolocMinBufferSize": 10
  },
  "createdAt": "2019-12-18 12:53:17.306048",
  "updatedAt": "2019-12-18 12:53:17.306048"
};

const list = {
  "result": [
    {
      "id": "1361C787-1EC4-48C2-86A7-2A7BA39AD4F3",
      "name": "My Device Profile",
      "organizationID": "2",
      "networkServerID": "1",
      "createdAt": "2019-12-18 12:53:17.306048",
      "updatedAt": "2019-12-18 12:53:17.306048"
    }
  ],
  "totalCount": "1"
}

class MockDeviceProfileStoreApi {
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
}

export default MockDeviceProfileStoreApi;
