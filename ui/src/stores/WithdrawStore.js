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
        resp.obj.totalCount = 2;
       resp.obj.withdrawRequest.push({withdrawId:1, userName: "pepe", availableToken:"ETH", amount: 112, updateAt: "2020-03-03T06:33:06.812528Z"});
       resp.obj.withdrawRequest.push({withdrawId:2, userName: "depa", availableToken:"MXC", amount: 222, updateAt: "2020-03-03T06:33:06.812528Z"});

        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  WithdrawReq(apiWithdrawReqRequest, callbackFunc) {
    this.swagger.then(client => {
      client.apis.WithdrawService.WithdrawReq({
        "orgId": apiWithdrawReqRequest.orgId,
        "moneyAbbr": apiWithdrawReqRequest.moneyAbbr,
        body: {
          amount: apiWithdrawReqRequest.amount,
          moneyAbbr: apiWithdrawReqRequest.moneyAbbr,
          orgId: apiWithdrawReqRequest.orgId
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
