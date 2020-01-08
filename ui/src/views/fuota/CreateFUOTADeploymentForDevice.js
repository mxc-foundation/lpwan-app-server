import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Card, Container, Row, Col } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import OrganizationStore from "../../stores/OrganizationStore";
import ApplicationStore from "../../stores/ApplicationStore";
import DeviceStore from "../../stores/DeviceStore";
import FUOTADeploymentStore from "../../stores/FUOTADeploymentStore";
import FUOTADeploymentForm from "./FUOTADeploymentForm";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};


class CreateFUOTADeploymentForDevice extends Component {
  constructor() {
    super();
    this.state = {};
  }

  componentDidMount() {
    const { match } = this.props;

    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0";

    if (isApplication) {
      ApplicationStore.get(match.params.applicationID, resp => {
        this.setState({
          application: resp,
        });
      });
    }

    DeviceStore.get(match.params.devEUI, resp => {
      this.setState({
        device: resp,
      });
    });
  
    OrganizationStore.get(currentOrgID, resp => {
      this.setState({
        organization: resp.organization
      });
    });
  }

  onSubmit = (fuotaDeployment) => {
    const { match } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0";

    FUOTADeploymentStore.createForDevice(match.params.devEUI, fuotaDeployment, resp => {
      isApplication
      ? this.props.history.push(`/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/fuota-deployments`)
      : this.props.history.push(`/organizations/${currentOrgID}/devices/${match.params.devEUI}/fuota-deployments`);
    });
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0";
    const { application, device, organization } = this.state;
    const { classes, match } = this.props;
    const currentOrgName = organization && (organization.name || organization.displayName);

    // if (this.state.application === undefined || this.state.device === undefined) {
    if (this.state.device === undefined) {
      return null;
    }

    return(
      <Container fluid>
        <Row>
          <Col xs={12}>
            <TitleBar noButtons>
              {
                isApplication && application ? (
                  <Breadcrumb className={classes.breadcrumb}>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/applications`
                    }>{i18n.t(`${packageNS}:tr000076`)}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/applications/${currentApplicationID}`
                    }>{application.application.name}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/applications/${currentApplicationID}`
                    }>{i18n.t(`${packageNS}:tr000278`)}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}`
                    }>{device.device.name}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/fuota-deployments`
                    }>Firmware</Link></BreadcrumbItem>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000381`)}</BreadcrumbItem>
                  </Breadcrumb>
                ) : (
                  <Breadcrumb className={classes.breadcrumb}>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations`
                    }>Organizations</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}`
                    }>{currentOrgName || currentOrgID}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/devices`
                    }>{i18n.t(`${packageNS}:tr000278`)}</Link></BreadcrumbItem>
                    <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                      `/organizations/${currentOrgID}/devices/${match.params.devEUI}`
                    }>{device.device.name}</Link></BreadcrumbItem>
                    <BreadcrumbItem>Firmware</BreadcrumbItem>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000381`)}</BreadcrumbItem>
                  </Breadcrumb>
                )
              }      
            </TitleBar>

            <Card body>
              <FUOTADeploymentForm
                submitLabel={i18n.t(`${packageNS}:tr000277`)}
                onSubmit={this.onSubmit}
              />
              <br />
            </Card>
          </Col>       
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(withRouter(CreateFUOTADeploymentForDevice));

