import React, { Component } from "react";
import { withRouter } from "react-router-dom";
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
