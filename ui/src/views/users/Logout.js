import React, { Component } from "react";
import { withRouter, Link, Redirect } from "react-router-dom";
import { withStyles } from '@material-ui/core/styles';
import SessionStore from "../../stores/SessionStore";

class Logout extends Component {
  constructor() {
    super();
  }

  componentDidMount() {
    setTimeout(() => { 
      SessionStore.logout(() => {
      });
    }, 300);
  }

  render() {
    return(
      <>
        {/* <Redirect path="/login" /> */}
      </>
    );
  }
}

export default withRouter(Logout);
