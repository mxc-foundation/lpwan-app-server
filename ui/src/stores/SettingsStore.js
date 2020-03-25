import { EventEmitter } from "events";
import Swagger from "swagger-client";
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class SettingsStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/settings.swagger.json", sessionStore.getClientOpts());
  }

  getSystemSettings(callbackFunc) {
    this.swagger.then(client => {
      client.apis.SettingsService.GetSettings()
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  setSystemSettings(body, callbackFunc) {
    this.swagger.then(client => {
      client.apis.SettingsService.ModifySettings({
        body
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }


}

const settingsStore = new SettingsStore();
export default settingsStore;
