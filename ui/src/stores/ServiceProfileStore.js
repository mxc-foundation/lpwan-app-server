import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class ServiceProfileStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/serviceProfile.swagger.json", sessionStore.getClientOpts());
  }

  async create(serviceProfile) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServiceProfileService.Create({
          body: {
            serviceProfile,
          },
        });
  
        resp = await checkStatus(resp);
        this.notify("created");
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  get(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.ServiceProfileService.Get({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  update(serviceProfile, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.ServiceProfileService.Update({
        "serviceProfile.id": serviceProfile.id,
        body: {
          serviceProfile: serviceProfile,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        callbackFunc(resp.obj);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  delete(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.ServiceProfileService.Delete({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        this.notify(i18n.t(`${packageNS}:tr000326`));
        callbackFunc(resp.ojb);
      })
      .catch(errorHandler);
    });
  }

  list(organizationID, limit, offset, callbackFunc, errorCallbackFunc) {
    return this.swagger.then(client => {
      client.apis.ServiceProfileService.List({
        organizationID: organizationID,
        limit: limit,
        offset: offset,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc && callbackFunc(resp.obj);
      })
      .catch(error => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "service-profile has been " + action,
      },
    });
  }
}

const serviceProfileStore = new ServiceProfileStore();
export default serviceProfileStore;
