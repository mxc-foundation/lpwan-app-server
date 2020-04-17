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

  async getWalletMiningIncome(userId, orgId) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.StakingService.GetWalletMiningIncome({
            userId,
            orgId
        });
    
        resp = await checkStatus(resp);
        console.log('getMining', resp);
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
        message: "balance has been " + action,
      },
    });
  }
}

const walletStore = new WalletStore();
export default walletStore;
