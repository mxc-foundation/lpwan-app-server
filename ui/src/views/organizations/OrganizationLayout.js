import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import OrganizationStore from "../../stores/OrganizationStore";
import UpdateOrganization from "./UpdateOrganization";
import Admin from "../../components/Admin";
import UpdateServiceProfile from "../service-profiles/UpdateServiceProfile";
import Modal from "../../components/Modal";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class OrganizationLayout extends Component {
  constructor() {
    super();
    this.state = {
      nsDialog: false,
    };
    this.loadData = this.loadData.bind(this);
    this.deleteOrganization = this.deleteOrganization.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(prevProps) {
    if (prevProps === this.props) {
      return;
    }

    this.loadData();
  }

  loadData() {
    OrganizationStore.get(this.props.match.params.organizationID, resp => {
      this.setState({
        organization: resp,
      });
    });
  }

  deleteOrganization() {
    OrganizationStore.delete(this.props.match.params.organizationID, () => {
      this.props.history.push("/organizations");
    });
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  };

  render() {
    const { classes } = this.props;

    return (
      this.state.organization ? <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:lpwan.organizations.delete_organization`)}
          callback={this.deleteOrganization} />}
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
              <BreadcrumbItem className={classes.breadcrumbItem}>Control Panel</BreadcrumbItem>
              <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                `/organizations`}>{i18n.t(`${packageNS}:tr000049`)
                }</Link></BreadcrumbItem>
            </Admin>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
            <BreadcrumbItem active>{this.state.organization.organization.name}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <UpdateOrganization organization={this.state.organization} />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment> : <div></div>
    );
  }
}


export default withStyles(styles)(withRouter(OrganizationLayout));
