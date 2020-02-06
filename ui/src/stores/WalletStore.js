import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import { checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";
import MockWalletStoreApi from '../api/mockWalletStoreApi';
import isDev from '../util/isDev';


class WalletStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/wallet.swagger.json", sessionStore.getClientOpts());
  }

  getDlPrice(orgId, callbackFunc) {
    // Run the following in development environment and early exit from function
    if (isDev) {
      (async () => callbackFunc(await MockWalletStoreApi.getDlPrice(orgId)))();
      return;
    }

    this.swagger.then(client => {
      client.apis.ProxyRequest.GetDlPrice({
        orgId,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
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
