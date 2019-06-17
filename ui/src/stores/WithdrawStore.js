import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class WithdrawStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/withdraw.swagger.json", sessionStore.getClientOpts());
  }
  
  get(id, callbackFunc) {
    //console.log("Organication id")
    //console.log(id)
    this.swagger.then(client => {
      client.apis.OrganizationService.Get({
        id: id,
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
        message: "organization has been " + action,
      },
    });
  }
}

const withdrawStore = new WithdrawStore();
export default withdrawStore;
