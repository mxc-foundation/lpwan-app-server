import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class NetworkServerStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/networkServer.swagger.json", sessionStore.getClientOpts());
  }

  create(networkServer, callbackFunc, errorCallbackFunc) {
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
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  async get(id) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.NetworkServerService.Get({
          id
        });
    
        resp = await checkStatus(resp);
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  update(networkServer, callbackFunc, errorCallbackFunc) {
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
        message: `${i18n.t(`${packageNS}:tr000356`)} ` + action,
      },
    });
  }

  async delete(id) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.NetworkServerService.Delete({
          id
        });

        resp = await checkStatus(resp);
        this.notify("deleted");
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }
  
  list(organizationID, limit, offset, callbackFunc, errorCallbackFunc) {
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
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }
}

const networkServerStore = new NetworkServerStore();
export default networkServerStore;
window.test = networkServerStore;
