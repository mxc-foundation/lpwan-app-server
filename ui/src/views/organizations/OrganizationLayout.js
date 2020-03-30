import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, CardBody, Col, Row, Alert } from 'reactstrap';
import Admin from "../../components/Admin";
import Modal from "../../components/Modal";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from '../../i18n';
import OrganizationStore from "../../stores/OrganizationStore";
import UpdateOrganization from "./UpdateOrganization";




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

  loadData = async () => {
    let organization = await OrganizationStore.get(this.props.match.params.organizationID);
    
    this.setState({
      organization
    });
  }

  deleteOrganization = async () => {
    await OrganizationStore.delete(this.props.match.params.organizationID);
    this.props.history.push("/organizations");
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  };

  render() {
    let currentOrgName = this.state.organization ? this.state.organization.organization.name : "";
    let currentOrgFullName = null;
    if (currentOrgName.length > 5) {
      currentOrgFullName = currentOrgName;
      currentOrgName = currentOrgName.slice(0, 5) + "...";
    }

    return (
      this.state.organization ? <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:lpwan.organizations.delete_organization`)}
          closeModal={() => this.setState({ nsDialog: false })}
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
          <Breadcrumb>
            <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
            <BreadcrumbItem><Link to={`/organizations`}>{i18n.t(`${packageNS}:tr000049`)}</Link></BreadcrumbItem>
            <BreadcrumbItem title={currentOrgFullName}>{currentOrgName}</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
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
