import { EventEmitter } from "events";

import Swagger from "swagger-client";

import i18n, { packageNS } from '../i18n';
import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class NetworkServerStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/networkServer.swagger.json", sessionStore.getClientOpts());
  }

  create(networkServer, callbackFunc) {
    this.swagger.then(client => {
      client.apis.NetworkServerService.Create({
        body: {
          networkServer: networkServer,
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
    this.swagger.then((client) => {
      client.apis.NetworkServerService.Get({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  update(networkServer, callbackFunc) {
    this.swagger.then(client => {
      client.apis.NetworkServerService.Update({
        "networkServer.id": networkServer.id,
        body: {
          networkServer: networkServer,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
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
        message: `${i18n.t(`${packageNS}:tr000356`)} ` + action,
      },
    });
  }

  delete(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.NetworkServerService.Delete({
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
  
  list(organizationID, limit, offset, callbackFunc) {
    this.swagger.then((client) => {
      client.apis.NetworkServerService.List({
        organizationID: organizationID,
        limit: limit,
        offset: offset,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }
}

const networkServerStore = new NetworkServerStore();
export default networkServerStore;
window.test = networkServerStore;
