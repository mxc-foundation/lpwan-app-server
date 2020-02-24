import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { Row, Col, Card, CardBody } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import OrgBreadCumb from '../../components/OrgBreadcrumb';

import SessionStorage from "../../stores/SessionStore";
import TopupForm from "./TopupForm";
import InfoCard from "./InfoCard";


class Topup extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
  }

  onSubmit = () => {
    if (SessionStorage.getUser().isAdmin) {
      this.props.history.push(`/control-panel/modify-account`);
    } else {
      this.props.history.push(`/modify-account/${this.props.match.params.organizationID}`);
    }
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const path = `/modify-account/${this.props.match.params.organizationID}`;

    return (<React.Fragment>
      <TitleBar>
        <OrgBreadCumb orgListCallback={() => { this.props.switchToSidebarId('DEFAULT'); }} 
          orgNameCallback={() => { this.props.switchToSidebarId('DEFAULT'); }}
          organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000568`), active: false },
            { label: i18n.t(`${packageNS}:menu.topup.topup`), active: true }]}></OrgBreadCumb>
      </TitleBar>

      <Row>
        <Col>
          <Card>
            <CardBody>
              <TopupForm
                reps={this.state.accounts} {...this.props}
                orgId={this.props.match.params.organizationID}
                path={path}
              />

            </CardBody>
          </Card>
        </Col>
        <Col>
          <InfoCard path={path} />
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(Topup);