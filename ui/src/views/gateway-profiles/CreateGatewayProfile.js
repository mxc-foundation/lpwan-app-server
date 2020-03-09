import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import Modal from '../../components/Modal';
import { Breadcrumb, BreadcrumbItem, Form, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import GatewayProfileForm from "./GatewayProfileForm";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import NetworkServerStore from "../../stores/NetworkServerStore";


class CreateGatewayProfile extends Component {
  constructor() {
    super();
    this.state = {
      nsDialog: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  componentDidMount() {
    NetworkServerStore.list(0, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          nsDialog: true,
        });
      }
    });
  }

  closeDialog = () => {
    this.setState({
      nsDialog: false,
    });
  }

  onSubmit(gatewayProfile) {
    GatewayProfileStore.create(gatewayProfile, resp => {
      this.props.history.push("/gateway-profiles");
    });
  }

  render() {

    return (<>
      <TitleBar>
        <Breadcrumb>
          <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
          <BreadcrumbItem><Link to={`/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)}</Link></BreadcrumbItem>
          <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
        </Breadcrumb>
      </TitleBar>

      <Form>
        {this.state.nsDialog && <Modal
          title={""}
          left={"DISMISS"}
          right={"ADD"}
          closeModal={this.closeDialog}
          context={i18n.t(`${packageNS}:tr000377`)}
          callback={this.deleteGatewayProfile} />}
        <Row>
          <Col>
            <Card>
              <CardBody>
                <GatewayProfileForm
                  submitLabel={i18n.t(`${packageNS}:tr000277`)}
                  onSubmit={this.onSubmit}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </Form>
    </>
    );
  }
}

export default withRouter(CreateGatewayProfile);
