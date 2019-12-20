import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class StakeStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/staking.swagger.json", sessionStore.getClientOpts());
  }

   async stake(req) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.StakingService.Stake({
            "orgId": req.orgId,
            body: {
              orgId: req.orgId,
              amount: req.amount
            },
        });
    
        resp = await checkStatus(resp);
        return resp;
      } catch (error) {
        errorHandler(error);
    }
  }

  /* stake(req, callbackFunc) {
    this.swagger.then(client => {
        client.apis.StakingService.Stake({
          "orgId": req.orgId,
          body: {
            orgId: req.orgId,
            amount: req.amount
          },
        })
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.body);
        })
        .catch(errorHandler);
      });
  } */

  async unstake(orgId) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.StakingService.Unstake({
          "orgId": orgId,
          body: {
            orgId
          },
        });
    
        resp = await checkStatus(resp);
        return resp;
      } catch (error) {
        errorHandler(error);
    }
  }

  async getActiveStakes(orgId) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.StakingService.GetActiveStakes({
            orgId
        });
    
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  }

  async getStakingPercentage(orgId) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.StakingService.GetStakingPercentage({
          orgId
      });
    
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  }

  /* async getStakingHistory(orgId, offset, limit) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.StakingService.GetStakingHistory({
            orgId,
            offset,
            limit
        });
    console.log(resp);
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  } */

  getStakingHistory(orgId, offset, limit, callbackFunc) {
    this.swagger.then(client => {
      client.apis.StakingService.GetStakingHistory({
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
        message: "Stake has been " + action,
      },
    });
  }
}

const stakeStore = new StakeStore();
export default stakeStore;
