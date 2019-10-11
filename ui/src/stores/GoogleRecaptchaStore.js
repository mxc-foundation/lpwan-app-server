import { EventEmitter } from "events";

import Swagger from "swagger-client";

import sessionStore from "./SessionStore";
import {checkStatus, errorHandler } from "./helpers";
import dispatcher from "../dispatcher";


class GoogleRecaptchaStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/recaptcha.swagger.json", sessionStore.getClientOpts());
  }

  getVerifyingGoogleRecaptcha(req, callBackFunc) {
    console.log('req', req);  
    this.swagger.then(client => {
      client.apis.InternalService.getVerifyingGoogleRecaptcha({body: req})
        .then(checkStatus)
        .then(resp => {
          console.log('google-recaptcha', resp);
        })
        .catch(errorHandler);
    });
  }

  notify(action) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "success",
        message: "balance has been " + action,
      },
    });
  }
}

const googleRecaptchaStore = new GoogleRecaptchaStore();
export default googleRecaptchaStore;
