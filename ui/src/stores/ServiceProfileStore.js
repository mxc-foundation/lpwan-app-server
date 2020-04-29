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

  async get(id) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServiceProfileService.Get({
          id,
        });
  
        resp = await checkStatus(resp);
        
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  async update(serviceProfile) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServiceProfileService.Update({
          "serviceProfile.id": serviceProfile.id,
          body: {
            serviceProfile: serviceProfile,
          },
        });
  
        resp = await checkStatus(resp);
        this.notify("updated");
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  async delete(id) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServiceProfileService.Delete({
          id
        });

        resp = await checkStatus(resp);
        this.notify(i18n.t(`${packageNS}:tr000326`));

        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  async list(organizationID, limit, offset) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.ServiceProfileService.List({
          organizationID,
          limit,
          offset,
        });

        resp = await checkStatus(resp);

        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
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
