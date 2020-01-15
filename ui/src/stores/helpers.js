import dispatcher from "../dispatcher";
import history from '../history';
import SessionStore from '../stores/SessionStore';

export function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response
  } else {
    throw response.json();
  }
};

export function errorHandler(error) {
  if(error.response === undefined) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "error",
        message: error.message,
      },
    });
  } else {
    console.error('Stores errorHandler error', error.response);
    if (error.response.obj && error.response.obj.code === 16) {
      // TODO: handle this error properly. do NOT route or logout here (since it can cause logout bugs)!
      setTimeout(() => {
        SessionStore.logout();
      }, 1000);
    } else {
      dispatcher.dispatch({
        type: "CREATE_NOTIFICATION",
        notification: {
          type: "error",
          message: error.response.obj ? (
            error.response.obj.error + " (code: " + error.response.obj.code + ")"
            // Internal Server Error 500 returns object with the following
          ) : error.message + " (code: " + error.status + ")",
        },
      });
    }
  }
};

export function errorHandlerLogin(error) {
  if(error.response === undefined) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "error",
        message: error.message,
      },
    });
  } else {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "error",
        message: error.response.obj.error + " (code: " + error.response.obj.code + ")",
      },
    });
  }
};

export function errorHandlerIgnoreNotFound(error) {
  if (error.response === undefined) {
    dispatcher.dispatch({
      type: "CREATE_NOTIFICATION",
      notification: {
        type: "error",
        message: error.message,
      },
    });
  } else {
    if (error.response.obj.code === 16 && history.location.pathname !== "/login") {
      history.push("/login");
    } else if (error.response.obj.code !== 5) {
      dispatcher.dispatch({
        type: "CREATE_NOTIFICATION",
        notification: {
          type: "error",
          message: error.response.obj.error + " (code: " + error.response.obj.code + ")",
        },
      });
    }
  }
};

export function errorHandlerIgnoreNotFoundWithCallback(callbackFunc) {
  return function(error) {
    if (error.response.obj.code === 5) {
      callbackFunc(null);
    } else {
      errorHandlerIgnoreNotFound(error);
    }
  }
}
