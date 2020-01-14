import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import Admin from "../../components/Admin";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import SessionStore from "../../stores/SessionStore";
import UpdateServiceProfile from "./UpdateServiceProfile";
import Modal from "../../components/Modal";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class ServiceProfileLayout extends Component {
  constructor() {
    super();
    this.state = {
      admin: false,
      nsDialog: false,
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
    ServiceProfileStore.delete(this.props.match.params.serviceProfileID, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    });
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  };

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      this.state.serviceProfile ? <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:lpwan.service_profiles.delete_service_profile`)}
          callback={this.deleteServiceProfile} />}
        <TitleBar
          buttons={
            <Admin>
              <TitleBarButton
                key={1}
                color="danger"
                label={i18n.t(`${packageNS}:tr000061`)}
                icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                onClick={this.openModal}
              />
            </Admin>
          }
        >
          <Breadcrumb className={classes.breadcrumb}>
            <Admin>
              <BreadcrumbItem>
                <Link
                  className={classes.breadcrumbItemLink}
                  to={`/organizations`}
                >
                  Organizations
              </Link>
              </BreadcrumbItem>
              <BreadcrumbItem>
                <Link
                  className={classes.breadcrumbItemLink}
                  to={`/organizations/${currentOrgID}`}
                >
                  {currentOrgID}
                </Link>
              </BreadcrumbItem>
            </Admin>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}/service-profiles`}
              >
                {i18n.t(`${packageNS}:tr000069`)}
              </Link>
            </BreadcrumbItem>
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

export default withStyles(styles)(withRouter(ServiceProfileLayout));
