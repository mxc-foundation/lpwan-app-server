import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";
import { withStyles } from '@material-ui/core/styles';
import SessionStore from "../../stores/SessionStore";

class Logout extends Component {
  constructor() {
    super();
  }

  componentDidMount() {
    SessionStore.logout(() => {
        this.props.history.push("/login");
    });
  }

  render() {
    return(
      <>
        
      </>
    );
  }
}

export default withRouter(Logout);
