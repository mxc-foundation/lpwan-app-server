import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class DeviceQueueStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/deviceQueue.swagger.json", sessionStore.getClientOpts());
  }

  flush(devEUI, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceQueueService.Flush({
        devEUI: devEUI,
      })
        .then(checkStatus)
        .then(resp => {
          this.notify("device-queue has been flushed");
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  list(devEUI, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceQueueService.List({
        devEUI: devEUI,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
      .catch(errorHandler);
    });
  }

  enqueue(item, callbackFunc) {
    this.swagger.then(client => {
      client.apis.DeviceQueueService.Enqueue({
        "deviceQueueItem.devEUI": item.devEUI,
        body: {
          deviceQueueItem: item,
        },
      })
        .then(checkStatus)
        .then(resp => {
          this.notify("device-queue item has bee created");
          this.emit("enqueue");
          callbackFunc(resp.obj);
        })
      .catch(errorHandler);
    });
  }

  notify(msg) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: msg,
      },
    });
  }
}

const deviceQueueStore = new DeviceQueueStore();
export default deviceQueueStore;
