import { EventEmitter } from "events";

import Swagger from "swagger-client";

import i18n, { packageNS } from '../i18n';
import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class SupernodeStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/super_node.swagger.json", sessionStore.getClientOpts());
  }

  getSuperNodeActiveMoneyAccount(moneyAbbr, orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.SuperNodeService.GetSuperNodeActiveMoneyAccount({
        orgId,
        moneyAbbr
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.body);
      })
      .catch(errorHandler);
    });
  }

  addSuperNodeMoneyAccount(req, orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.SuperNodeService.AddSuperNodeMoneyAccount({
        "orgId": orgId,
        "moneyAbbr": req.moneyAbbr,
        body: {
            moneyAbbr: req.moneyAbbr,
            accountAddr: req.createAccount
        },
      })
      .then(checkStatus)
      .then(resp => {
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
        message: `${i18n.t(`${packageNS}:menu.store.user_has_been`)} ` + action,
      },
    });
  }
}

const supernodeStore = new SupernodeStore();
export default supernodeStore;