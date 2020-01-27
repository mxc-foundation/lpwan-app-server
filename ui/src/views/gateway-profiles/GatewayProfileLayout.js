import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Modal from '../../components/Modal';
import { Button, Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import UpdateGatewayProfile from "./UpdateGatewayProfile";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class GatewayProfileLayout extends Component {
  constructor() {
    super();

    this.state = {
      nsDialog: false
    };

    this.deleteGatewayProfile = this.deleteGatewayProfile.bind(this);
  }

  componentDidMount() {
    GatewayProfileStore.get(this.props.match.params.gatewayProfileID, resp => {
      this.setState({
        gatewayProfile: resp,
      });
    });
  }

  deleteGatewayProfile = () => {
    GatewayProfileStore.delete(this.props.match.params.gatewayProfileID, () => {
      this.props.history.push("/gateway-profiles");
    });
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  }

  render() {
    const { classes } = this.props;

    if (this.state.gatewayProfile === undefined) {
      return (<div></div>);
    }
    const icon = <i className="mdi mdi-delete-empty"></i>;

    return (
      <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:tr000426`)}
          callback={this.deleteGatewayProfile} />}
        <TitleBar
          buttons={[
            <Button color="danger"
              key={1}
              onClick={this.openModal}
              className=""><i className="mdi mdi-delete-empty"></i>{' '}{i18n.t(`${packageNS}:tr000401`)}
            </Button>
          ]}
        >
          <Breadcrumb className={classes.breadcrumb} style={{ fontSize: "1.25rem", margin: "0rem" }}>
            <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
            <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
              `/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)
            }</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
            <BreadcrumbItem active>{`${this.state.gatewayProfile.gatewayProfile.apiGatewayProfile.name}`}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <UpdateGatewayProfile gatewayProfile={this.state.gatewayProfile.gatewayProfile} />
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(GatewayProfileLayout));
