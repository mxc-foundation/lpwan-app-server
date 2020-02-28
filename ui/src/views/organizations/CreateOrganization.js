import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import OrganizationForm from "./OrganizationForm";
import OrganizationStore from "../../stores/OrganizationStore";


class CreateOrganization extends Component {
  constructor() {
    super();
    this.state = {
      loading: false
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(organization) {
    this.setState({ loading: true });
    OrganizationStore.create(organization, resp => {
      this.setState({ loading: false });
      this.props.history.push("/organizations");
    }, error => { this.setState({ loading: false }) });
  }

  render() {

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
            <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
            <BreadcrumbItem><Link to={`/organizations`}>{i18n.t(`${packageNS}:tr000049`)}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
              {this.state.loading && <Loader />}
                <OrganizationForm
                    match={this.props.match}
                    submitLabel={i18n.t(`${packageNS}:tr000277`)}
                    onSubmit={this.onSubmit}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(CreateOrganization);
