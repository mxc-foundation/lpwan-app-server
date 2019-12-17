import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Modal from '../../components/Modal';
import { Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import UpdateGatewayProfile from "./UpdateGatewayProfile";


class GatewayProfileLayout extends Component {
  constructor() {
    super();

    this.state = {};

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
    if (window.confirm("Are you sure you want to delete this gateway-profile?")) {
      GatewayProfileStore.delete(this.props.match.params.gatewayProfileID, () => {
        this.props.history.push("/gateway-profiles");
      });
    }
  }

  render() {
    if (this.state.gatewayProfile === undefined) {
      return (<div></div>);
    }

    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Modal buttonLabel={i18n.t(`${packageNS}:tr000401`)} title={""} context={i18n.t(`${packageNS}:tr000426`)} callback={this.deleteGatewayProfile} />
          ]}
        >

          <TitleBarTitle title={i18n.t(`${packageNS}:tr000063`)} />
          <Breadcrumb>
            <BreadcrumbItem><Link to={`/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{this.state.gatewayProfile.gatewayProfile.name}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <UpdateGatewayProfile gatewayProfile={this.state.gatewayProfile.gatewayProfile} />
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(GatewayProfileLayout);
