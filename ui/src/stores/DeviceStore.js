import { EventEmitter } from "events";
import RobustWebSocket from "robust-websocket";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import { checkStatus, errorHandler, errorHandlerIgnoreNotFoundWithCallback } from "./helpers";
import sessionStore from "./SessionStore";




class DeviceStore extends EventEmitter {
  constructor() {
    super();
    this.wsDataStatus = null;
    this.wsFramesStatus = null;
    this.swagger = new Swagger("/swagger/device.swagger.json", sessionStore.getClientOpts());
    this.swaggerM2M = new Swagger("/swagger/m2mserver_device.swagger.json", sessionStore.getClientOpts());
  }

  getWSDataStatus() {
    return this.wsDataStatus;
  }

  getWSFramesStatus() {
    return this.wsFramesStatus;
  }

  getDeviceList(orgId, offset, limit, callbackFunc, errorCallbackFunc) {
    this.swaggerM2M.then(client => {
      client.apis.DSDeviceService.GetDeviceList({
        orgId,
        offset,
        limit
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.body);
      })
      .catch((error) => {
        errorHandler(error);
        if (errorCallbackFunc)
          errorCallbackFunc(error);
      });
    });
  }

  getDeviceHistory(orgId, devId, offset, limit, callbackFunc) {    
    this.swaggerM2M.then(client => {
      client.apis.DSDeviceService.GetDeviceHistory({
        orgId,
        devId,
        offset,
        limit
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.body);
      })
      .catch(errorHandler);
    });
  }

  setDeviceMode(orgId, devId, devMode, callbackFunc) {
    this.swaggerM2M.then(client => {
    client.apis.DSDeviceService.SetDeviceMode({
      "orgId": orgId,
      "devId": devId,
      body: {
        orgId,
        devId,
        devMode
      },
    })
    .then(checkStatus)
    .then(resp => {
      this.emit("update");
      this.notify("updated");
      callbackFunc(resp.obj);
    })
    .catch(errorHandler);
    });
  }

  create(device, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.Create({
        body: {
          device: device,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("created");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  get(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.Get({
        devEUI: id,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  update(device, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.Update({
        "device.devEui": device.devEUI,
        body: {
          device: device,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.emit("update");
        this.notify("updated");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });

  }
//please fix this
  delete(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.Delete({
        devEUI: id,
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("deleted");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  list(filters, callbackFunc) {
    console.log('device list filters:', filters);
    this.swagger.then(client => {
      client.apis.DeviceService.List(filters)
      .then(checkStatus)
      .then(resp => {
        console.log('device list:', resp);
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getKeys(devEUI, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.GetKeys({
        devEUI,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandlerIgnoreNotFoundWithCallback(callbackFunc));
    });
  }

  createKeys(deviceKeys, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.CreateKeys({
        "deviceKeys.devEui": deviceKeys.devEUI,
        body: {
          deviceKeys: deviceKeys,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notifyKeys("created");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  updateKeys(deviceKeys, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.UpdateKeys({
        "deviceKeys.devEui": deviceKeys.devEUI,
        body: {
          deviceKeys: deviceKeys,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notifyKeys("updated");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getActivation(devEUI, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.GetActivation({
        devEUI,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandlerIgnoreNotFoundWithCallback(callbackFunc));
    });
  }

  activate(deviceActivation, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.Activate({
        "deviceActivation.devEui": deviceActivation.devEUI,
        body: {
          deviceActivation: deviceActivation,
        },
      })
      .then(checkStatus)
      .then(resp => {
        dispatcher.dispatch({
          type: "CREATE_NOTIFICATION",
          notification: {
            type: "success",
            message: "device has been (re)activated",
          },
        });
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getRandomDevAddr(devEUI, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceService.GetRandomDevAddr({
        devEUI,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getDataLogsConnection(devEUI, onData) {
    const loc = window.location;
    const wsURL = (() => {
      if (loc.host === "localhost:3000" || loc.host === "localhost:3001") {
        return `wss://localhost:8080/api/devices/${devEUI}/events`;
      }

      const wsProtocol = loc.protocol === "https:" ? "wss:" : "ws:";
      return `${wsProtocol}//${loc.host}/api/devices/${devEUI}/events`;
    })();

    const conn = new RobustWebSocket(wsURL, ["Bearer", sessionStore.getToken()], {});

    conn.addEventListener("open", () => {
      this.wsDataStatus = "CONNECTED";
      this.emit("ws.status.change");
    });

    conn.addEventListener("message", (e) => {
      const msg = JSON.parse(e.data);
      if (msg.error !== undefined) {
        dispatcher.dispatch({
          type: "CREATE_NOTIFICATION",
          notification: {
            type: "error",
            message: msg.error.message,
          },
        });
      } else if (msg.result !== undefined) {
        onData(msg.result);
      }
    });

    conn.addEventListener("close", () => {
      this.wsDataStatus = null;
      this.emit("ws.status.change");
    });

    conn.addEventListener("error", () => {
      this.wsDataStatus = "ERROR";
      this.emit("ws.status.change");
    });

    return conn;
  }

  getFrameLogsConnection(devEUI, onData) {
    const loc = window.location;
    const wsURL = (() => {
      if (loc.host === "localhost:3000" || loc.host === "localhost:3001") {
        return `wss://localhost:8080/api/devices/${devEUI}/frames`;
      }

      const wsProtocol = loc.protocol === "https:" ? "wss:" : "ws:";
      return `${wsProtocol}//${loc.host}/api/devices/${devEUI}/frames`;
    })();

    const conn = new RobustWebSocket(wsURL, ["Bearer", sessionStore.getToken()], {});

    conn.addEventListener("open", () => {
      this.wsFramesStatus = "CONNECTED";
      this.emit("ws.status.change");
    });

    conn.addEventListener("message", (e) => {
      const msg = JSON.parse(e.data);
      if (msg.error !== undefined) {
        dispatcher.dispatch({
          type: "CREATE_NOTIFICATION",
          notification: {
            type: "error",
            message: msg.error.message,
          },
        });
      } else if (msg.result !== undefined) {
        onData(msg.result);
      }
    });

    conn.addEventListener("close", () => {
      this.wsFramesStatus = null;
      this.emit("ws.status.change");
    });

    conn.addEventListener("error", (e) => {
      this.wsFramesStatus = "ERROR";
      this.emit("ws.status.change");
    });

    return conn;
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "device has been " + action,
      },
    });
  }

  notifyKeys(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "device-keys have been " + action,
      },
    });
  }
}

const deviceStore = new DeviceStore();
export default deviceStore;
