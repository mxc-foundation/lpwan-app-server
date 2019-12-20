import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import Loader from "../../components/Loader";
import CommonModal from '../../components/Modal';

import GatewayForm from "./GatewayForm";
import GatewayStore from "../../stores/GatewayStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";


class CreateGateway extends Component {
  constructor() {
    super();

    this.state = {
      spDialog: false,
      loading: false
    };

    this.onSubmit = this.onSubmit.bind(this);
    this.redirectToCreateServiceProfile = this.redirectToCreateServiceProfile.bind(this);
  }

  componentDidMount() {
    ServiceProfileStore.list(this.props.match.params.organizationID, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          spDialog: true,
        });
      }
    });
  }

  closeDialog = () => {
    this.setState({
      spDialog: false,
    });
  }

  onSubmit(gateway) {
    let gw = gateway;
    gw.organizationID = this.props.match.params.organizationID;

    GatewayStore.create(gateway, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    });
  }

  redirectToCreateServiceProfile = () => {
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles/create`);
  }

  render() {
    return (<React.Fragment>

      <TitleBar>
        <TitleBarTitle title={i18n.t(`${packageNS}:tr000277`) + " " + i18n.t(`${packageNS}:tr000072`)} />
        <Breadcrumb>
          <BreadcrumbItem><Link to={`/organizations/${this.props.match.params.organizationID}/gateways`}>{i18n.t(`${packageNS}:tr000063`)}</Link></BreadcrumbItem>
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
                  object={{ name: '', description: '', id: '', location: { altitude: 0 } }}
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
        showConfirmButton={true} left={i18n.t(`${packageNS}:tr000166`)} right={i18n.t(`${packageNS}:tr000277`)}></CommonModal>
    </React.Fragment>
    );
  }
}

export default withRouter(CreateGateway);
