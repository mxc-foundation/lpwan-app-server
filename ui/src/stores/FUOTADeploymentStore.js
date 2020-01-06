import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";
import MockFUOTADeploymentStoreApi from '../api/mockFUOTADeploymentStoreApi';
import isDev from '../util/isDev';


class FUOTADeploymentStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/fuotaDeployment.swagger.json", sessionStore.getClientOpts());
  }

  createForDevice(devEUI, fuotaDeployment, callbackFunc) {
    this.swagger.then(client => {
      client.apis.FUOTADeploymentService.CreateForDevice({
        devEui: devEUI,
        body: {
          devEUI: devEUI,
          fuotaDeployment: fuotaDeployment,
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
    // Run the following in development environment and early exit from function
    if (isDev) {
      (async () => callbackFunc(await MockFUOTADeploymentStoreApi.get()))();
      return;
    }

    this.swagger.then(client => {
      client.apis.FUOTADeploymentService.Get({
        id: id,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  list(filters, callbackFunc) {
    // Run the following in development environment and early exit from function
    if (isDev) {
      (async () => callbackFunc(await MockFUOTADeploymentStoreApi.list()))();
      return;
    }

    this.swagger.then(client => {
      client.apis.FUOTADeploymentService.List(filters)
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  listDeploymentDevices(filters, callbackFunc) {
    // Run the following in development environment and early exit from function
    if (isDev) {
      (async () => callbackFunc(await MockFUOTADeploymentStoreApi.listDeploymentDevices()))();
      return;
    }

    this.swagger.then(client => {
      client.apis.FUOTADeploymentService.ListDeploymentDevices(filters)
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  getDeploymentDevice(fuotaDeploymentID, devEUI, callbackFunc) {
    // Run the following in development environment and early exit from function
    if (isDev) {
      (async () => callbackFunc(await MockFUOTADeploymentStoreApi.getDeploymentDevice()))();
      return;
    }

    this.swagger.then(client => {
      client.apis.FUOTADeploymentService.GetDeploymentDevice({
        fuotaDeploymentId: fuotaDeploymentID,
        devEui: devEUI,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "fuota deployment has been " + action,
      },
    });
  }

  emitReload() {
    this.emit("reload");
  }
}

const fuotaDeploymentStore = new FUOTADeploymentStore();
export default fuotaDeploymentStore;

