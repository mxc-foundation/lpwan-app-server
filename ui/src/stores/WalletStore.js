import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class WalletStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/wallet.swagger.json", sessionStore.getClientOpts());
  }

  getDlPrice(orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WalletService.GetDlPrice({
        orgId,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  getWalletBalance(orgId, userId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WalletService.GetWalletBalance({
        userId,
        orgId
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
