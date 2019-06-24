import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class WithdrawStore extends EventEmitter {
  constructor() {
    super();
    this.withdrawSwagger = new Swagger("/swagger/withdraw.swagger.json", sessionStore.getClientOpts());
  }
  
  getWithdrawFee(money_abbr, callbackFunc) {
    this.withdrawSwagger.then(client => {
      client.apis.WithdrawService.GetWithdrawFee({money_abbr})
      .then(checkStatus)
      .then(resp => {
        console.log("withdrawFee:", resp);
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
