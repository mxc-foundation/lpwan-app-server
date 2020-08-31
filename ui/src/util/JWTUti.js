import jwt from "jsonwebtoken";
import sessionStore from "../stores/SessionStore";


function checkJWT(invalidateSession) {
    const token = sessionStore.getToken();
    if (token) {
        jwt.verify(token, process.env.REACT_APP_JWT_SECRET, function (err, decoded) {
            if (err) {
                invalidateSession(err);
            }
        });
    }
}

export function initJWTTimer(invalidateSession) {
    setInterval(checkJWT.bind(null, invalidateSession), 5000);
}