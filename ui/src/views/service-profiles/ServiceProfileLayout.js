import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import Admin from "../../components/Admin";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import SessionStore from "../../stores/SessionStore";
import UpdateServiceProfile from "./UpdateServiceProfile";
import Modal from "../../components/Modal";
import OrgBreadCumb from '../../components/OrgBreadcrumb';


class ServiceProfileLayout extends Component {
  constructor() {
    super();
    this.state = {
      admin: false,
      nsDialog: false,
    };
    this.deleteServiceProfile = this.deleteServiceProfile.bind(this);
    this.setIsAdmin = this.setIsAdmin.bind(this);
  }

  componentDidMount() {
    ServiceProfileStore.get(this.props.match.params.serviceProfileID, resp => {
      this.setState({
        serviceProfile: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  setIsAdmin() {
    this.setState({
      admin: SessionStore.isAdmin(),
    });
  }

  deleteServiceProfile() {
    ServiceProfileStore.delete(this.props.match.params.serviceProfileID, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    });
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  };

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      this.state.serviceProfile ? <React.Fragment>
        
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:lpwan.service_profiles.delete_service_profile`)}
          closeModal={() => this.setState({ nsDialog: false })}
          callback={this.deleteServiceProfile} />}

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
          <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000078`), active: false, to: `/organizations/${currentOrgID}/service-profiles` },
            { label: this.state.serviceProfile.serviceProfile.name, active: false }]}></OrgBreadCumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <UpdateServiceProfile serviceProfile={this.state.serviceProfile.serviceProfile} admin={this.state.admin} />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment> : <div></div>
    );
  }
}

export default withRouter(ServiceProfileLayout);
