const mockDeviceFrame = {
  "downlinkFrame": {
    "txInfo": {
      "frequency": 123,
      "modulation": "LORA",
      "loRaModulationInfo": {
        "bandwidth": 123,
        "spreadingFactor": 123,
        "codeRate": "fast",
        "polarizationInversion": true
      },
      "fskModulationInfo": {
        "bandwidth": 123,
        "bitRtate": 1
      }
    },
    "rxInfo": [
      {
        "gatewayID": "1",
        "uplinkID": "2",
        "time": "2019-12-06 15:46:12.624982",
        "timeSinceGpsEpoch": "2019-12-01 15:46:12.624982",
        "rssi": 123,
        "loraSnr": 1.2,
        "channel": 200,
        "rfChain": 10,
        "board": 33,
        "antenna": 44,
        "location": {
          "latitude": 300.1,
          "longitude": 200.5,
          "altitude": 30.6,
          "source": "GPS",
          "accuracy": 10
        },
        "fineTimestampType": "ENCRYPTED",
        "encryptedFineTimestamp": {
          "aesKeyIndex": 3,
          "encryptedNS": "123",
          "fpgaID": "456"
        },
        "plainFineTimestamp": {
          "time": "2019-1-01 15:46:12.624982"
        },
        "context": {}
      }
    ],
    "phyPayloadJSON": '{"mhdr":{"mtype":"abc"},"macPayload":{"devEUI":"70b3d5fffe1cb547","fhdr":{"devAddr":"123"}}}'
  },
  "uplinkFrame": {
    "txInfo": {
      "frequency": 123,
      "modulation": "LORA",
      "loRaModulationInfo": {
        "bandwidth": 123,
        "spreadingFactor": 123,
        "codeRate": "fast",
        "polarizationInversion": true
      },
      "fskModulationInfo": {
        "bandwidth": 123,
        "bitRtate": 1
      }
    },
    "rxInfo": [],
    "phyPayloadJSON": '{"mhdr":{"mtype":"abc"},"macPayload":{"devEUI":"70b3d5fffe1cb547","fhdr":{"devAddr":"123"}}}'
  }
};

export default mockDeviceFrame;
