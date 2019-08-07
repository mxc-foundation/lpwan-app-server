import { EventEmitter } from "events";

import Swagger from "swagger-client";
import { checkStatus, errorHandler, errorHandlerLogin } from "./helpers";
import dispatcher from "../dispatcher";

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

  getUser() {
    return this.user;
  }

  getOrganizations() {
    return this.organizations;
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
    callBackFunc();
  }

  fetchProfile(callBackFunc) {
    this.swagger.then(client => {
      client.apis.InternalService.Profile({})
        .then(checkStatus)
        .then(resp => {
          this.user = resp.obj.user;

          if(resp.obj.organizations !== undefined) {
            this.organizations = resp.obj.organizations;
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

  notifyActivation() {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "Confirmation email has been sent.",
      },
    });
  }

}

const sessionStore = new SessionStore();
export default sessionStore;
