import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Card, CardBody, Col, Row } from "reactstrap";
import Loader from "../../components/Loader";
import CommonModal from "../../components/Modal";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from "../../i18n";
import GatewayStore from "../../stores/GatewayStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import GatewayForm from "./GatewayForm";





class CreateGateway extends Component {
  constructor() {
    super();

    this.state = {
      spDialog: false,
      loading: true
    };
  }

  componentDidMount = async () => {
    const resp = await ServiceProfileStore.list(this.props.match.params.organizationID, 10, 0);
    const state = {
      loading: false
    };
    if (resp.totalCount === "0") {
      state.spDialog = true;
    }

    this.setState(state);
  }

  closeDialog = () => {
    this.setState({
      spDialog: false
    });
  };

  onSubmit = (gateway, config, classBConfig) => {
    GatewayStore.create(gateway, resp => {
      this.props.history.push(
        `/organizations/${this.props.match.params.organizationID}/gateways`
      );
    });
  };

  redirectToCreateServiceProfile = () => {
    this.props.history.push(
      `/organizations/${this.props.match.params.organizationID}/service-profiles/create`
    );
  };

  render() {
    const currentOrgID =
      this.props.organizationID || this.props.match.params.organizationID;

    return (
      <React.Fragment>
        <TitleBar>
          <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000063`), active: false, to: `/organizations/${currentOrgID}/gateways` },
            { label: i18n.t(`${packageNS}:tr000277`), active: true }]}></OrgBreadCumb>
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

        <CommonModal
          showToggleButton={false}
          callback={this.redirectToCreateServiceProfile}
          show={this.state.spDialog}
          context={
            <React.Fragment>
              <p>
                {i18n.t(`${packageNS}:tr000165`)}
                {i18n.t(`${packageNS}:tr000326`)}
              </p>
              <p>{i18n.t(`${packageNS}:tr000327`)}</p>
            </React.Fragment>
          }
          title={i18n.t(`${packageNS}:tr000164`)}
          showConfirmButton={true}
          left={i18n.t(`${packageNS}:tr000166`)}
          right={i18n.t(`${packageNS}:tr000277`)}
        ></CommonModal>
      </React.Fragment>
    );
  }
}

export default withRouter(CreateGateway);
