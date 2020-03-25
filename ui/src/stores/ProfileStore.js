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

  getUserOrganizationList(orgId, callbackFunc) {
    this.profileSwagger.then(client => {
      client.apis.InternalService.GetUserOrganizationList({
        orgId
      })
      .then(checkStatus)
      .then(updateOrganizations)
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
        message: `${i18n.t(`${packageNS}:menu.store.profile_has_been`)} ` + action,
      },
    });
  }
}

const profileStore = new ProfileStore();
export default profileStore;
