import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";
import updateOrganizations from "./SetUserProfile";




class ProfileStore extends EventEmitter {
  constructor() {
    super();
    this.profileSwagger = new Swagger("/swagger/profile.swagger.json", sessionStore.getClientOpts());
  }

  async getUserOrganizationList(orgId) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.OrganizationService.GetUserOrganizationList({
          orgId
        });
        
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: `${i18n.t(`${packageNS}:menu.store.profile_has_been`)} ` + action,
      },
    });
  }
}

const profileStore = new ProfileStore();
export default profileStore;
