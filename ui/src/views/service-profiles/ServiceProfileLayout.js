import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';


import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";

import Admin from "../../components/Admin";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import SessionStore from "../../stores/SessionStore";
import UpdateServiceProfile from "./UpdateServiceProfile";


class ServiceProfileLayout extends Component {
  constructor() {
    super();
    this.state = {
      admin: false,
    };
    this.deleteServiceProfile = this.deleteServiceProfile.bind(this);
    this.setIsAdmin = this.setIsAdmin.bind(this);
  }

  componentDidMount() {
    ServiceProfileStore.get(this.props.match.params.serviceProfileID, resp => {
      this.setState({
        serviceProfile: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  setIsAdmin() {
    this.setState({
      admin: SessionStore.isAdmin(),
    });
  }

  deleteServiceProfile() {
    if (window.confirm("Are you sure you want to delete this service-profile?")) {
      ServiceProfileStore.delete(this.props.match.params.serviceProfileID, resp => {
        this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
      });
    }
  }

  render() {
    return (
      this.state.serviceProfile ? <React.Fragment>
        <TitleBar
          buttons={
            <Admin>
              <TitleBarButton
                key={1}
              color="danger"
                label={i18n.t(`${packageNS}:tr000061`)}
              icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                onClick={this.deleteServiceProfile}
              />
            </Admin>
          }
        >

          <TitleBarTitle title={i18n.t(`${packageNS}:tr000069`)} />
          <Breadcrumb>
            <BreadcrumbItem><Link to={`/organizations/${this.props.match.params.organizationID}/service-profiles`}>{i18n.t(`${packageNS}:tr000069`)}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{this.state.serviceProfile.serviceProfile.name}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <UpdateServiceProfile serviceProfile={this.state.serviceProfile.serviceProfile} admin={this.state.admin} />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment> : <div></div>
    );
  }
}

export default withRouter(ServiceProfileLayout);
