import { EventEmitter } from "events";

import Swagger from "swagger-client";

import i18n, { packageNS } from '../i18n';
import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";

class MoneyStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/ext_account.swagger.json", sessionStore.getClientOpts());
  }

  getActiveMoneyAccount(moneyAbbr, orgId, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.MoneyService.GetActiveMoneyAccount({
        moneyAbbr,
        orgId,
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

  modifyMoneyAccount(req, callbackFunc) {
    this.swagger.then(client => {
      client.apis.MoneyService.ModifyMoneyAccount({
        "orgId": req.orgId,
        "moneyAbbr": req.moneyAbbr,
        body: {
          currentAccount: req.newAccount,
          orgId: req.orgId,
          moneyAbbr: req.moneyAbbr
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        callbackFunc(resp.body);
      })
      .catch(errorHandler);
    });
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: `${i18n.t(`${packageNS}:menu.store.account_has_been`)} ` + action,
      },
    });
  }
}

const moneyStore = new MoneyStore();
export default moneyStore;
