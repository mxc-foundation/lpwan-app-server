import {EventEmitter} from "events";
import Swagger from "swagger-client";
import {checkStatus, errorHandler} from "./helpers";
import sessionStore from "./SessionStore";


class WithdrawStore extends EventEmitter {
    constructor() {
        super();
        this.swagger = new Swagger("/swagger/withdraw.swagger.json", sessionStore.getClientOpts());
    }

    getWithdrawFee(moneyAbbr, callbackFunc) {
        this.swagger.then(client => {
            client.apis.WithdrawService.GetWithdrawFee({
                moneyAbbr
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.body);
                })
                .catch(errorHandler);
        });
    }

    setWithdrawFee(moneyAbbr, orgId, body, callbackFunc) {
        this.swagger.then(client => {
            client.apis.WithdrawService.ModifyWithdrawFee({
                moneyAbbr,
                orgId,
                body
            })
                .then(checkStatus)
                .then(resp => {
                    callbackFunc(resp.obj);
                })
                .catch(errorHandler);
        });
    }
}

const withdrawStore = new WithdrawStore();
export default withdrawStore;
