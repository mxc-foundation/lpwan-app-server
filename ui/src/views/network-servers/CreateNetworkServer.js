import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Card,
  CardSubtitle,
  CardTitle,
  Col,
  Container,
  Row
} from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import NetworkServerForm from "./NetworkServerForm";
import NetworkServerStore from "../../stores/NetworkServerStore";


class CreateNetworkServer extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(networkServer) {
    NetworkServerStore.create(networkServer, resp => {
      this.props.history.push("/network-servers");
    });
  }

  render() {
    return(
      <Container>
        <Row>
          <Col md="12" sm="12">
            <Card className="card-box" style={{ minWidth: "25rem" }}>
              <Row>
                <Col md="12" xs="12">
                  <Breadcrumb>
                    <BreadcrumbItem><a href="#">Home</a></BreadcrumbItem>
                    <BreadcrumbItem><a href="/network-servers">{i18n.t(`${packageNS}:tr000040`)}</a></BreadcrumbItem>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
                  </Breadcrumb>
                </Col>
              </Row>
              <Row>
                <Col md="12" xs="12">
                  <CardTitle className="mt-0 header-title">{i18n.t(`${packageNS}:tr000040`)}</CardTitle>
                  <CardSubtitle className="text-muted font-14 mb-3">
                    Create a network server.
                  </CardSubtitle>
                </Col>
              </Row>
              <Row className="md-center">
                <Col md="12" sm="12">
                  <NetworkServerForm
                    submitLabel={i18n.t(`${packageNS}:tr000277`)}
                    onSubmit={this.onSubmit}
                  />
                </Col>
              </Row>
            </Card>
          </Col>
        </Row>
      </Container>
    );
  }
}

export default withRouter(CreateNetworkServer);
