import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class TopupStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/topup.swagger.json", sessionStore.getClientOpts());
  }

  getTopUpDestination(orgId, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.TopUpService.GetTopUpDestination({
        orgId
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.body);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
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

    getIncome(orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.TopUpService.GetIncome({
        orgId
       
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.body);
      })
      .catch(errorHandler);
    });
  }

  async topupWidget() {
    let res = {
      total: 13000,
      data: [
          { "month": "Jun", "amount": 1000 },
          { "month": "Jul", "amount": 2200 },
          { "month": "Aug", "amount": 2420 },
          { "month": "Sep", "amount": 3400 },
          { "month": "Oct", "amount": 1550 },
          { "month": "Nov", "amount": 1720 },
          { "month": "Dec", "amount": 485 },
      ]
    }

    return res
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: `${i18n.t(`${packageNS}:menu.store.transaction_has_been`)} ` + action,
      },
    });
  }
}

const topupStore = new TopupStore();
export default topupStore;
