import { Component } from "react";

import SessionStore from "../stores/SessionStore";


class NonAdmin extends Component {
  constructor() {
    super();
    this.state = {
      hasOrg: false,
    };

    this.setRole = this.setRole.bind(this);
  }

  componentDidMount() {
    SessionStore.on("change", this.setRole);
    this.setRole();
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setRole);
  }

  componentDidUpdate(prevProps) {
    if (prevProps === this.props) {
      return;
    }

    this.setRole();
  }

  setRole() {
    let orgList = SessionStore.getOrganizations()

    if ((orgList.length === 0) || (orgList.length === 1 && orgList[0].value === "0") || SessionStore.isAdmin() ) {
      this.setState({
        hasOrg: false ,
      });
    }else {
      this.setState({
        hasOrg: true ,
      });
    }
  }

  render() {
    if (this.state.hasOrg) {
      return(this.props.children);
    }

    return(null);
  }
}

export default NonAdmin;
