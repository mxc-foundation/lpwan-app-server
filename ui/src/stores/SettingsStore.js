import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class SettingsStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/settings.swagger.json", sessionStore.getClientOpts());
  }

  getSystemSettings(callbackFunc) {
    this.swagger.then(client => {
      client.apis.SettingsService.GetSettings()
      .then(checkStatus)
      //.then(updateOrganizations)
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
      //.then(updateOrganizations)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }


}

const settingsStore = new SettingsStore();
export default settingsStore;
