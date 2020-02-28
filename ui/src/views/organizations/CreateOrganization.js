import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Form, Row, Col, Card, CardBody } from 'reactstrap';

import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import OrganizationForm from "./OrganizationForm";
import OrganizationStore from "../../stores/OrganizationStore";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

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
    const { classes } = this.props;

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
            <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={`/organizations`}>{i18n.t(`${packageNS}:tr000049`)}</Link></BreadcrumbItem>
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
                    object={{}}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(CreateOrganization));
