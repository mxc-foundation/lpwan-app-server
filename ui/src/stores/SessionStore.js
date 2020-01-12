import { EventEmitter } from "events";

import Swagger from "swagger-client";
import { checkStatus, errorHandler, errorHandlerLogin } from "./helpers";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import MockSessionStoreApi from '../api/mockSessionStoreApi';
import isDev from '../util/isDev';

class SessionStore extends EventEmitter {
  constructor() {
    super();
    this.client = null;
    this.user = null;
    this.organizations = [];
    this.settings = {};
    this.branding = {};

    this.swagger = Swagger("/swagger/internal.swagger.json", this.getClientOpts())
    
    this.swagger.then(client => {
      this.client = client;
      const token = this.getToken();
      if (token) {// !== null && !history.location.pathname.includes('/registration-confirm/')) {
        this.fetchProfile(() => {});
      }
    });
  }

  getClientOpts() {
    return {
      requestInterceptor: (req) => {
        if (this.getToken() !== null) {
          req.headers["Grpc-Metadata-Authorization"] = "Bearer " + this.getToken();
        }
        return req;
      },
    }
  }

  setToken(token) {
    localStorage.setItem("jwt", token);
  }

  getToken() {
    return localStorage.getItem("jwt");
  }

  setSupportedLanguages(languages) {
    localStorage.setItem("languages-supported", JSON.stringify(languages));
  }

  getSupportedLanguages() {
    return JSON.parse(localStorage.getItem("languages-supported"));
  }

  setLanguage(language) {
    localStorage.setItem("language", JSON.stringify(language));
  }

  getLanguage() {
    return JSON.parse(localStorage.getItem("language"));
  }

  getOrganizationID() {
    const orgID = localStorage.getItem("organizationID");
    if (!orgID) {
      return null;
    }

    return orgID;
  }

  setOrganizationID(id) {
    localStorage.setItem("organizationID", id);
    this.emit("organization.change");
  }

  setUser(user) {
    localStorage.setItem("user", JSON.stringify(user));
  }

  setOrganizations(organizations) {
    localStorage.setItem("organizations", JSON.stringify(organizations));
  }

  getUser() {
    // Run the following in development environment and early exit from function
    // Uncomment to show mock profile pic 
    // if (isDev) {
    //   return MockSessionStoreApi.getUser();
    // }

    let user = this.user;
    if (!user) {
      user = localStorage.getItem("user");
      if (user) user = JSON.parse(user);
    }
    return user;
  }

  getOrganizations() {
    let organizations = this.organizations;
    if (!organizations || (organizations && organizations.length === 0)) {
      organizations = localStorage.getItem("organizations");
      if (organizations) organizations = JSON.parse(organizations);
    }
    return organizations || [];
  }

  getSettings() {
    return this.settings;
  }

  isAdmin() {
    if (this.user === undefined || this.user === null) {
      return false;
    }
    return this.user.isAdmin;
  }

  isOrganizationAdmin(organizationID) {
    for (let i = 0; i < this.organizations.length; i++) {
      if (this.organizations[i].organizationID === organizationID) {
        return this.organizations[i].isAdmin;
      }
    }
  }

  isOrganizationDeviceAdmin(organizationID) {
    for (let i = 0; i < this.organizations.length; i++) {
      if (this.organizations[i].organizationID === organizationID) {
        return this.organizations[i].isAdmin || this.organizations[i].isDeviceAdmin;
      }
    }
  }

  isOrganizationGatewayAdmin(organizationID) {
    for (let i = 0; i < this.organizations.length; i++) {
      if (this.organizations[i].organizationID === organizationID) {
        return this.organizations[i].isAdmin || this.organizations[i].isGatewayAdmin;
      }
    }
  }

  login(login, callBackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.Login({body: login})
        .then(checkStatus)
        .then(resp => {
          this.setToken(resp.obj.jwt);
          this.fetchProfile(callBackFunc);
        })
        .catch(errorHandlerLogin);
    });
  }

  logout(callBackFunc) {
    localStorage.clear();
    this.user = null;
    this.organizations = [];
    this.settings = {};
    this.emit("change");
    callBackFunc && callBackFunc();
  }

  async getProfile(){
    try {
      const client = await this.swagger.then((client) => client);
      let resp = await client.apis.InternalService.Profile();

      resp = await checkStatus(resp);
      return resp;
    } catch (error) {
      errorHandler(error);
    }
  }

  fetchProfile(callBackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.Profile({})
        .then(checkStatus)
        .then(resp => {
          this.user = resp.obj.user;
          this.setUser(this.user);

          if(resp.obj.organizations !== undefined) {
            this.organizations = resp.obj.organizations;
            this.setOrganizations(this.organizations);
          }

          if(resp.obj.settings !== undefined) {
            this.settings = resp.obj.settings;
          }

          this.emit("change");
          callBackFunc(resp);
        })
        .catch(errorHandler);
    });
  }

  globalSearch(search, limit, offset, callbackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.GlobalSearch({
        search: search,
        limit: limit,
        offset: offset,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
      });
  }

  getBranding(callbackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.Branding({})
        .then(checkStatus)
        .then(resp => {
          callbackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }
  
  register(data, callbackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.RegisterUser({
        body: {
          email: data.username,
          language: data.language
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notifyActivation();
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  confirmRegistration(securityToken, callbackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.ConfirmRegistration({
        body: {
          token: securityToken,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.setToken(resp.obj.jwt);
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  finishRegistration(data, callbackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.FinishRegistration({
        body: {
          userId: data.userId,
          password: data.password,
          organizationName: data.organizationName,
          organizationDisplayName: data.organizationDisplayName,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.fetchProfile(callbackFunc);
      })
      .catch(errorHandler);
    });
  }

  getVerifyingGoogleRecaptcha(req, callBackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.GetVerifyingGoogleRecaptcha({body: req})
        .then(checkStatus)
        .then(resp => {
          callBackFunc(resp.obj);
        })
        .catch(errorHandler);
    });
  }

  notifyActivation() {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: i18n.t(`${packageNS}:tr000018`),
      },
    });
  }

}

const sessionStore = new SessionStore();
export default sessionStore;
