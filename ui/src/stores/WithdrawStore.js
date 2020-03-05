import { EventEmitter } from "events";

import Swagger from "swagger-client";

import i18n, { packageNS } from '../i18n';
import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class WithdrawStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/withdraw.swagger.json", sessionStore.getClientOpts());
  }

  getWithdrawFee(moneyAbbr, orgId, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.GetWithdrawFee({
        moneyAbbr,
        orgId
      })
      .then(checkStatus)
      //.then(updateOrganizations)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  setWithdrawFee(moneyAbbr, orgId, body, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.ModifyWithdrawFee({
        moneyAbbr,
        orgId,
        body
      })
      .then(checkStatus)
      //.then(updateOrganizations)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getWithdrawRequestList(limit, offset, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.GetWithdrawRequestList({
        offset,
        limit
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  getWithdrawHistory(moneyAbbr, orgId, limit, offset, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.GetWithdrawHistory({
        moneyAbbr,
        orgId,
        offset,
        limit
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  withdrawReq(req, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.WithdrawReq({
        "moneyAbbr": req.moneyAbbr,
        body: {
          orgId: req.orgId,
          moneyAbbr: req.moneyAbbr,
          amount: req.amount,
          ethAddress: req.ethAddress,
          availableBalance: req.availableBalance
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        this.emit("withdraw");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  confirmWithdraw(req, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.ConfirmWithdraw({
        body: {
          orgId: req.orgId,
          confirmStatus:req.confirmStatus,
          denyComment: req.denyComment,
          withdrawId: req.withdrawId
        },
      })
      .then(checkStatus)
      //.then(updateOrganizations)
      .then(resp => {
        this.notify("updated");
        this.emit("withdraw");
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
        message: `${i18n.t(`${packageNS}:menu.store.successful_withdrawal`)}`
      },
    });
  }
}

const withdrawStore = new WithdrawStore();
export default withdrawStore;
