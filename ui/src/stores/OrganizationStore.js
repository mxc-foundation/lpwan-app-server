import { EventEmitter } from "events";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import { checkStatus, errorHandler } from "./helpers";
import sessionStore from "./SessionStore";




class OrganizationStore extends EventEmitter {
  constructor() {
    super();
    this.swagger = new Swagger("/swagger/organization.swagger.json", sessionStore.getClientOpts());
  }

  
  async create(organization) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.OrganizationService.Create({
          body: {
            organization: organization,
          },
        });
  
        resp = await checkStatus(resp);
        this.emit("create", organization);
        this.notify("created");
        
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  async get(id) {
    try {
        const client = await this.swagger.then((client) => client);
        let resp = await client.apis.OrganizationService.Get({
          id
        });
    
        resp = await checkStatus(resp);
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  update(organization, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.Update({
        "organization.id": organization.id,
        body: {
          organization: organization,
        },
      })
      .then(checkStatus)
      .then(resp => {
        this.emit("change", organization);
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
      client.apis.OrganizationService.Delete({
        id: id,
      })
      .then(checkStatus)
      .then(resp => {
        this.emit("delete", id);
        this.notify("deleted");
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  async list(search, limit, offset) {
    try {
        const client = await this.swagger;
        let resp = await client.apis.OrganizationService.List({
          search,
          limit,
          offset,
        });
        
        resp = await checkStatus(resp);
        return resp.obj;
      } catch (error) {
        errorHandler(error);
    }
  }

  addUser(organizationID, user, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.AddUser({
        "organizationUser.organizationId": organizationID,
        body: {
          organizationUser: user,
        },
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch((error) => {
        errorHandler(error);
        if (errorCallbackFunc) errorCallbackFunc(error);
      });
    });
  }

  getUser(organizationID, userID, callbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.GetUser({
        organizationID: organizationID,
        userID: userID,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  deleteUser(organizationID, userID, callbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.DeleteUser({
        organizationID: organizationID,
        userID: userID,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  updateUser(organizationUser, callbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.UpdateUser({
        "organizationUser.organizationId": organizationUser.organizationID,
        "organizationUser.userId": organizationUser.userID,
        body: {
          organizationUser: organizationUser,
        },
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
      })
      .catch(errorHandler);
    });
  }

  listUsers(organizationID, limit, offset, callbackFunc, errorCallbackFunc) {
    this.swagger.then(client => {
      client.apis.OrganizationService.ListUsers({
        organizationID: organizationID,
        limit: limit,
        offset: offset,
      })
      .then(checkStatus)
      .then(resp => {
        callbackFunc(resp.obj);
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
        message: "organization has been " + action,
      },
    });
  }
}

const organizationStore = new OrganizationStore();
export default organizationStore;
