import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Button, Col, Row } from 'reactstrap';
import Modal from '../../components/Modal';
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import OrganizationStore from "../../stores/OrganizationStore";
import SessionStore from "../../stores/SessionStore";
import UpdateOrganizationUser from "./UpdateOrganizationUser";




class OrganizationUserLayout extends Component {
  constructor() {
    super();
    this.state = {
      admin: false,
      nsDialog: false
    };
    this.deleteOrganizationUser = this.deleteOrganizationUser.bind(this);
    this.setIsAdmin = this.setIsAdmin.bind(this);
  }

  componentDidMount() {
    OrganizationStore.getUser(this.props.match.params.organizationID, this.props.match.params.userID, resp => {
      this.setState({
        organizationUser: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
  }

  componendWillUnmount() {
    SessionStore.removeListener("on", this.setIsAdmin);
  }

  setIsAdmin() {
    this.setState({
      admin: SessionStore.isAdmin(),
    });
  }

  deleteOrganizationUser() {
    OrganizationStore.deleteUser(this.props.match.params.organizationID, this.props.match.params.userID, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users`);
    });
  }

  gotoUser = () => {
    this.props.history.push(`/users/${this.props.match.params.organizationID}`);
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    if (this.state.organizationUser === undefined) {
      return (<div></div>);
    }

    const titleButtons = [];

    if (this.props.match.params.userID !== SessionStore.getUser().id && currentOrgID === SessionStore.getOrganizationID()) {
      titleButtons.push(<Button color="danger"
        key={1}
        onClick={this.openModal}
        className=""><i className="mdi mdi-delete"></i>{' '}{i18n.t(`${packageNS}:common.delete`)}
      </Button>);
    }

    // <Button color="secondary"
    //   key={1}
    //   onClick={this.gotoUser}
    //   className="btn-rp"><i className="mdi mdi-account-arrow-right-outline"></i>{' '}{i18n.t(`${packageNS}:lpwan.org_users.goto_user`)}
    // </Button> ,


    return (
      <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:lpwan.org_users.delete_user`)}
          closeModal={() => this.setState({ nsDialog: false })}
          callback={this.deleteOrganizationUser} />}

        <TitleBar buttons={titleButtons}>
          <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000068`), active: false, to: `/organizations/${currentOrgID}/users` },
            { label: this.state.organizationUser.organizationUser.username, active: true }]}></OrgBreadCumb>
        </TitleBar>
        <Row>
          <Col>
            <UpdateOrganizationUser organizationUser={this.state.organizationUser.organizationUser} />
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(OrganizationUserLayout);
