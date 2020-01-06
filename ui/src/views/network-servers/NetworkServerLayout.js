import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Card } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import Delete from "mdi-material-ui/Delete";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import NetworkServerStore from "../../stores/NetworkServerStore";
import UpdateNetworkServer from "./UpdateNetworkServer";

const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  }
});

class NetworkServerLayout extends Component {
  constructor() {
    super();

    this.state = {};

    this.deleteNetworkServer = this.deleteNetworkServer.bind(this);
  }

  componentDidMount() {
    NetworkServerStore.get(this.props.match.params.networkServerID, (resp) => {
      this.setState({
        networkServer: resp,
      });
    });
  }

  deleteNetworkServer() {
    if (window.confirm("Are you sure you want to delete this network-server?")) {
      NetworkServerStore.delete(this.props.match.params.networkServerID, () => {
        this.props.history.push("/network-servers");
      });
    }
  }

  render() {
    const { classes } = this.props;

    if (this.state.networkServer === undefined) {
      return(<div></div>);
    }

    return(
      <React.Fragment>
        <TitleBar
          buttons={[
            <TitleBarButton
              color="danger"
              key={1}
              icon={<Delete />}
              label={i18n.t(`${packageNS}:tr000061`)}
              onClick={this.deleteNetworkServer}
            />,
          ]}
        >
          <Breadcrumb className={classes.breadcrumb} style={{ fontSize: "1.25rem", margin: "0rem" }}>
            <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
              `/network-servers`}>{i18n.t(`${packageNS}:tr000040`)
            }</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
            <BreadcrumbItem active>{`${this.state.networkServer.networkServer.id}`}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <UpdateNetworkServer networkServer={this.state.networkServer.networkServer} />
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(NetworkServerLayout));
