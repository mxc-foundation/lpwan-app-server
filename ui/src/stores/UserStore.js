import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import i18n, { packageNS } from '../i18n';
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";



class UserStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/user.swagger.json", sessionStore.getClientOpts());
  }

  create(newUserObject, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.UserService.Create({
        body: newUserObject,
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("created");
        callbackFunc(resp.obj);
      })
      .catch((error) => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  get(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.UserService.Get({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  update(user, callbackFunc) {
    this.swagger.then(client => {
      client.apis.UserService.Update({
        "user.id": user.id,
        body: {
          "user": user,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  delete(id, callbackFunc) {
    this.swagger.then(client => {
      client.apis.UserService.Delete({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        this.notify(i18n.t(`${packageNS}:tr000326`));
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  updatePassword(id, password, callbackFunc) {
    this.swagger.then(client => {
      client.apis.UserService.UpdatePassword({
        "userId": id,
        body: {
          password: password,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.notify("updated");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  list(search, limit, offset, callbackFunc) {
    this.swagger.then((client) => {
      client.apis.UserService.List({
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

  getOTPCode = async (userEmail) => {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.UserService.GetOTPCode({
          userEmail
      });
    
        resp = await checkStatus(resp);
        return resp.body;
      } catch (error) {
        errorHandler(error);
    }
  }

  async getUserEmail(userEmail) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.UserService.GetUserEmail({
          userEmail
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
        message: "user has been " + action,
      },
    });
  }
}

const userStore = new UserStore();
export default userStore;
