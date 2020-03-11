import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";
import MockGatewayProfileStoreApi from '../api/mockGatewayProfileStoreApi';
import isDev from '../util/isDev';

class GatewayProfileStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/gatewayProfile.swagger.json", sessionStore.getClientOpts());
  }

  create(gatewayProfile, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.GatewayProfileService.Create({
        body: {
          gatewayProfile: gatewayProfile,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("created");
        callbackFunc(resp.obj);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  get(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.GatewayProfileService.Get({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  update(gatewayProfile, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.GatewayProfileService.Update({
        "gatewayProfile.id": gatewayProfile.id,
        body: {
          gatewayProfile: gatewayProfile,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        callbackFunc(resp.obj);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  delete(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.GatewayProfileService.Delete({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("deleted");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  list(networkServerID, limit, offset, callbackFunc, errorCallbackFunc) {
    this.swagger.then((client) => {
      client.apis.GatewayProfileService.List({
        networkServerID: networkServerID,
        limit: limit,
        offset: offset,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "gateway-profile has been " + action,
      },
    });
  }
}

const gatewayProfileStore = new GatewayProfileStore();
export default gatewayProfileStore;
