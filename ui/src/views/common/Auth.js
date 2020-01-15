import React from 'react';
import jwt from "jsonwebtoken";
import SessionStore from "../../stores/SessionStore";

export function CheckJWT() {
    const token = SessionStore.getToken();
    if (token) {
      jwt.verify(token, 'zPWZAeRHJ8aCzrJDr4VuY/isoVnLZ0nNQIgHBbE7nMA=', function(err, decoded) {
        if (err) {
          console.log('err.message', err.message);
          
          localStorage.clear();
        }
      });
    }
    return;
}