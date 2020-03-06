import { EventEmitter } from "events";

import Swagger from "swagger-client";

import i18n, { packageNS } from '../i18n';
import sessionStore from "./SessionStore";
import { checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class HistoryStore extends EventEmitter {
  constructor() {
    super();
    this.topupSwagger = new Swagger("/swagger/topup.swagger.json", sessionStore.getClientOpts());
    this.withdrawSwagger = new Swagger("/swagger/withdraw.swagger.json", sessionStore.getClientOpts());
    this.walletSwagger = new Swagger("/swagger/wallet.swagger.json", sessionStore.getClientOpts());
    this.moneySwagger = new Swagger("/swagger/ext_account.swagger.json", sessionStore.getClientOpts());
  }

  getTopUpHistory(orgId, offset, limit, callbackFunc) {
    this.swagger.then(client => {
      client.apis.TopUpService.GetTopUpHistory({
        orgId,
        offset,
        limit
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.body);
        })
        .catch(errorHandler);
    });
  }

  getWithdrawHistory(moneyAbbr, orgId, limit, offset, callbackFunc) {
    this.withdrawSwagger.then((client) => {
      client.apis.WithdrawService.GetWithdrawHistory({
        moneyAbbr,
        orgId,
        limit,
        offset,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  getVmxcTxHistory(orgId, limit, offset, callbackFunc) {
    this.walletSwagger.then((client) => {
      client.apis.WalletService.GetVmxcTxHistory({
        orgId,
        limit,
        offset,
      })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  getChangeMoneyAccountHistory(moneyAbbr, orgId, limit, offset, callbackFunc, errorCallbackFunc) {
    this.moneySwagger.then((client) => {
      client.apis.MoneyService.GetChangeMoneyAccountHistory({
        moneyAbbr,
        orgId,
        limit,
        offset,
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

  getWalletUsageHist(orgId, offset, limit, callbackFunc) {
    this.walletSwagger.then(client => {
      client.apis.WalletService.GetWalletUsageHist({
        orgId,
        offset,
        limit
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

const historyStore = new HistoryStore();
export default historyStore;
