import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class ServerInfoStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/serverInfo.swagger.json", sessionStore.getClientOpts());
  }

  getAppserverVersion(callbackFunc) {
    this.swagger.then(client => {
      client.apis.ServerInfoService.GetAppserverVersion()
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.data);
      })
      .catch(errorHandler);
    });
  }

  async getServerRegion() {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServerInfoService.GetServerRegion({});
        
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "server has been " + action,
      },
    });
  }
}

const profileStore = new ServerInfoStore();
export default profileStore;
