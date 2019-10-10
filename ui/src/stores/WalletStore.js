import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class WalletStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/wallet.swagger.json", sessionStore.getClientOpts());
  }

  getWalletBalance(orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WalletService.GetWalletBalance({
        orgId,
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
        message: "balance has been " + action,
      },
    });
  }
}

const walletStore = new WalletStore();
export default walletStore;
