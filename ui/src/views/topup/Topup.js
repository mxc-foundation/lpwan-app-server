import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Card, CardBody, Col, Nav, NavItem, NavLink, Row, TabContent, TabPane } from 'reactstrap';
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import DataDash from "../../components/DataDash"






class Topup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeTab: "0",
    };

    this.onTabToggle = this.onTabToggle.bind(this);
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
  }

  onTabToggle(tab) {
    this.setState({ activeTab: tab });
  }

  onSubmit = () => {
    /* if (SessionStorage.getUser().isAdmin) {
      this.props.history.push(`/control-panel/modify-account`);
    } else {
      this.props.history.push(`/modify-account/${this.props.match.params.organizationID}`);
    } */
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (<React.Fragment>
      <TitleBar>
        <OrgBreadCumb orgListCallback={() => { this.props.switchToSidebarId('DEFAULT'); }}
          orgNameCallback={() => { this.props.switchToSidebarId('DEFAULT'); }}
          organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:m2m_redirect.link`), active: true },]}></OrgBreadCumb>
      </TitleBar>
      <Row>
          <DataDash />
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(Topup);