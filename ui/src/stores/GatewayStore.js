import {EventEmitter} from "events";
import RobustWebSocket from "robust-websocket";
import Swagger from "swagger-client";
import dispatcher from "../dispatcher";
import {checkStatus, errorHandler, errorHandlerIgnoreNotFound} from "./helpers";
import sessionStore from "./SessionStore";


class GatewayStore extends EventEmitter {
    constructor() {
        super();
        this.wsStatus = null;
        this.swagger = new Swagger("/swagger/gateway.swagger.json", sessionStore.getClientOpts());
    }

    getWSStatus() {
        return this.wsStatus;
    }

    getGatewayList(orgId, offset, limit, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.GetGatewayList({
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

    getGatewayProfile(gwId, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.GetGatewayProfile({
                gwId,
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    getGatewayHistory(orgId, gwId, offset, limit, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.GetGatewayHistory({
                orgId,
                gwId,
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

    setGatewayMode(orgId, gwId, gwMode, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.SetGatewayMode({
                "orgId": orgId,
                "gwId": gwId,
                body: {
                    orgId,
                    gwId,
                    gwMode
                },
            })
                .then(checkStatus)
                .then(resp => {
                    this.emit("update");
                    this.notify("updated");
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    create(gateway, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.Create({
                body: {
                    gateway: gateway,
                },
            })
                .then(checkStatus)
                .then(resp => {
                    this.notify("created");
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    register(gateway, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.Register({
                body: {
                    organizationId: gateway.organizationId,
                    sn: gateway.sn.serial
                },
            })
                .then(checkStatus)
                .then(resp => {
                    this.notify("registered");
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    get(id, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.Get({
                id: id,
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    async getConfig(gatewayId) {
        try {
            const client = await this.swagger;
            let resp = await client.apis.GatewayService.GetGwConfig({
                gatewayId
            });

            resp = await checkStatus(resp);
            return resp.body.conf;
        } catch (error) {
            errorHandler(error);
        }
    }

    async update(gateway) {
        try {
            const client = await this.swagger;
            let resp = await client.apis.GatewayService.Update({
                "gateway.id": gateway.id,
                body: {
                    gateway,
                },
            });

            resp = await checkStatus(resp);
            this.notify("updated");

            return resp.obj;
        } catch (error) {
            errorHandler(error);
        }
    }

    async updateConfig(gateway, config) {
        try {
            const client = await this.swagger;
            let resp = await client.apis.GatewayService.UpdateGwConfig({
                "gatewayId": gateway.id,
                body: {
                    gatewayId: gateway.id,
                    conf: JSON.stringify(config)
                },
            });

            resp = await checkStatus(resp);
            this.notify("updated");

            return resp.obj;
        } catch (error) {
            errorHandler(error);
        }
    }

    delete(id, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.Delete({
                id: id,
            })
                .then(checkStatus)
                .then(resp => {
                    this.notify("deleted");
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    list(search, organizationID, limit, offset, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.List({
                limit: limit,
                offset: offset,
                organizationID: organizationID,
                search: search,
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    listLocations(callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.ListLocations()
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    getStats(gatewayID, start, end, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.GetStats({
                gatewayID: gatewayID,
                interval: "DAY",
                startTimestamp: start,
                endTimestamp: end,
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }

    getLastPing(gatewayID, callbackFunc) {
        this.swagger.then(client => {
            client.apis.GatewayService.GetLastPing({
                gatewayID: gatewayID,
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandlerIgnoreNotFound);
        });
    }

    getFrameLogsConnection(gatewayID, onOpen, onClose, onData) {
        const loc = window.location;
        const wsURL = (() => {
            if (loc.host === "localhost:3000" || loc.host === "localhost:3001") {
                return `wss://localhost:8080/api/gateways/${gatewayID}/frames`;
            }

            const wsProtocol = loc.protocol === "https:" ? "wss:" : "ws:";
            return `${wsProtocol}//${loc.host}/api/gateways/${gatewayID}/frames`;
        })();

        const conn = new RobustWebSocket(wsURL, ["Bearer", sessionStore.getToken()], {});

        conn.addEventListener("open", () => {
            this.wsStatus = "CONNECTED";
            this.emit("ws.status.change");
            onOpen();
        });

        conn.addEventListener("message", (e) => {
            const msg = JSON.parse(e.data);
            if (msg.error !== undefined) {
                dispatcher.dispatch({
                    type: "CREATE_NOTIFICATION",
                    notification: {
                        type: "error",
                        message: msg.error.message,
                    },
                });
            } else if (msg.result !== undefined) {
                onData(msg.result);
            }
        });

        conn.addEventListener("close", () => {
            this.wsStatus = null;
            this.emit("ws.status.change");
            onClose();
        });

        conn.addEventListener("error", () => {
            this.wsStatus = "ERROR";
            this.emit("ws.status.change");
        });

        return conn;
    }

    notify(action) {
        dispatcher.dispatch({
            type: "CREATE_NOTIFICATION",
            notification: {
                type: "success",
                message: "gateway has been " + action,
            },
        });
    }

    async getRootConfig(gatewayId, sn) {
        try {
            const client = await this.swagger.then((client) => client);
            let resp = await client.apis.GatewayService.GetGwPwd({
                gatewayId,
                sn
            });

            resp = await checkStatus(resp);
            return resp.obj;
        } catch (error) {
            errorHandler(error);
        }
    }

    async setAutoUpdateFirmware(gatewayId, autoUpdate) {
        try {
            const client = await this.swagger;
            let resp = await client.apis.GatewayService.SetAutoUpdateFirmware({
                gatewayId,
                body: {
                    gatewayId,
                    autoUpdate
                }
            });

            resp = await checkStatus(resp);
            this.emit("update");
            this.notify("updated");

            return resp.obj;
        } catch (error) {
            errorHandler(error);
        }
    }

    async getGatewayConfig(gatewayId) {
        /* try {
            const client = await this.swagger;
            let resp = await client.apis.GatewayService.GetGatewayConfig({
              gatewayId
            });

            resp = await checkStatus(resp);
            this.emit("update");
            this.notify("updated");

            return resp.obj;
          } catch (error) {
            errorHandler(error);
        } */
    }

}


const gatewayStore = new GatewayStore();
export default gatewayStore;
