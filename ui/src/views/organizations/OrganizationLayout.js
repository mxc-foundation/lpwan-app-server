import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';


import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";

import OrganizationStore from "../../stores/OrganizationStore";
import UpdateOrganization from "./UpdateOrganization";
import Admin from "../../components/Admin";
import UpdateServiceProfile from "../service-profiles/UpdateServiceProfile";
import Modal from "../../components/Modal";


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

          <TitleBarTitle title={i18n.t(`${packageNS}:tr000049`)} />
          <Breadcrumb>
            <BreadcrumbItem><Link to={`/organizations`}>{i18n.t(`${packageNS}:tr000049`)}</Link></BreadcrumbItem>
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


export default withRouter(OrganizationLayout);
