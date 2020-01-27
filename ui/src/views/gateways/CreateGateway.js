import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import CommonModal from '../../components/Modal';

import GatewayForm from "./GatewayForm";
import GatewayStore from "../../stores/GatewayStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class CreateGateway extends Component {
  constructor() {
    super();

    this.state = {
      spDialog: false,
      loading: true
    };
  }

  componentDidMount() {
    ServiceProfileStore.list(this.props.match.params.organizationID, 0, 0, resp => {
      const state = {
        loading: false
      }
      if (resp.totalCount === "0") {
        state.spDialog = true;
      }

      this.setState(state);
    });
  }

  closeDialog = () => {
    this.setState({
      spDialog: false,
    });
  }

  onSubmit = (gateway) => {
    GatewayStore.create(gateway, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    });
  }

  redirectToCreateServiceProfile = () => {
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles/create`);
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
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
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}/gateways`}
              >
                {i18n.t(`${packageNS}:tr000063`)}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <div className="position-relative">
                  {this.state.loading && <Loader />}

                  <GatewayForm
                    match={this.props.match}
                    submitLabel={i18n.t(`${packageNS}:tr000277`)}
                    onSubmit={this.onSubmit}
                  />
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>

        <CommonModal showToggleButton={false} callback={this.redirectToCreateServiceProfile}
          show={this.state.spDialog}
          context={
            <React.Fragment>
              <p>
                {i18n.t(`${packageNS}:tr000165`)}
                {i18n.t(`${packageNS}:tr000326`)}
              </p>
              <p>
                {i18n.t(`${packageNS}:tr000327`)}
              </p>
            </React.Fragment>
          } title={i18n.t(`${packageNS}:tr000164`)}
          showConfirmButton={true} left={i18n.t(`${packageNS}:tr000166`)} right={i18n.t(`${packageNS}:tr000277`)}>    
        </CommonModal>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(CreateGateway));
