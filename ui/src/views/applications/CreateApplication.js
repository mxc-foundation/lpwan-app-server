import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Button, Card, Container, Modal, ModalHeader, ModalBody, ModalFooter, NavLink, Row, Col } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import ApplicationForm from "./ApplicationForm";
import ApplicationStore from "../../stores/ApplicationStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";


const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  },
  card: {
    overflow: "visible",
  },
});


class CreateApplication extends Component {
  constructor() {
    super();
    this.state = {
      spDialog: false,
    };
  }

  componentDidMount() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    ServiceProfileStore.list(currentOrgID, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          spDialog: true,
        });
      }
    });
  }

  toggleDialog = () => {
    this.setState({
      spDialog: !this.state.spDialog,
    });
  }

  onSubmit = (application) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    let app = application;
    app.organizationID = currentOrgID;

    ApplicationStore.create(app, resp => {
      this.props.history.push(`/organizations/${currentOrgID}/applications`);
    });
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const closeBtn = <button className="close" onClick={this.toggleDialog}>&times;</button>;
    const { classes } = this.props;

    return(
      <Container fluid>
        <Row>
          <Col xs={12}>
            <Modal
              isOpen={this.state.spDialog}
              toggle={this.toggleDialog}
              aria-labelledby="help-dialog-title"
              aria-describedby="help-dialog-description"
            >
              <ModalHeader
                toggle={this.toggleDialog}
                close={closeBtn}
                id="help-dialog-title"
              >
                {i18n.t(`${packageNS}:tr000164`)}
              </ModalHeader>
              <ModalBody id="help-dialog-description">
                <p>
                  {i18n.t(`${packageNS}:tr000165`)}
                  {i18n.t(`${packageNS}:tr000326`)}
                </p>
                <p>
                  {i18n.t(`${packageNS}:tr000327`)}
                </p>
              </ModalBody>
              <ModalFooter>
                <Button variant="outlined">
                  <NavLink
                    style={{ color: "#fff", padding: "0" }}
                    tag={Link}
                    to={`/organizations/${currentOrgID}/service-profiles/create`}
                  >
                    {i18n.t(`${packageNS}:tr000277`)}
                  </NavLink>
                </Button>
                <Button color="primary" onClick={this.toggleDialog}>{i18n.t(`${packageNS}:tr000166`)}</Button>{' '}
              </ModalFooter>
            </Modal>

            <TitleBar>
              <Breadcrumb className={classes.breadcrumb}>
                <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                  `/organizations/${currentOrgID}/applications`
                }>{i18n.t(`${packageNS}:tr000076`)}</Link></BreadcrumbItem>
                <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
              </Breadcrumb>
            </TitleBar>

            <Card body>
              <ApplicationForm
                match={this.props.match}
                onSubmit={this.onSubmit}
                submitLabel={i18n.t(`${packageNS}:tr000277`)}
              />
              <br />
            </Card>
          </Col>       
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(withRouter(CreateApplication));
