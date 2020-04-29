import classNames from "classnames";
import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Card, CardBody, Col, Nav, NavItem, NavLink, Row, TabContent, TabPane } from 'reactstrap';
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import TopupCrypto from "./TopupCrypto";
import TopupHistory from "./TopupHistory";






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
            { label: i18n.t(`${packageNS}:tr000568`), active: false },
            { label: i18n.t(`${packageNS}:menu.topup.topup`), active: true }]}></OrgBreadCumb>
      </TitleBar>
      <Row>
        <Col>
          <Card>
            <CardBody className="pb-0">
              <Nav tabs>
                <NavItem>
                  <NavLink
                    className={classNames('nav-link', { active: this.state.activeTab === '0' })} href='#'
                    onClick={(e) => this.onTabToggle("0")}
                  >{i18n.t(`${packageNS}:menu.topup.crypto`)}</NavLink>
                </NavItem>
                <NavItem>
                  <NavLink
                    className={classNames('nav-link', { active: this.state.activeTab === '1' })} href='#'
                    onClick={(e) => this.onTabToggle("1")} disabled
                  >{i18n.t(`${packageNS}:menu.topup.otc`)}</NavLink>
                </NavItem>
              </Nav>

              <TabContent activeTab={this.state.activeTab}>
                <TabPane tabId="0">
                  <TopupCrypto />
                </TabPane>
                <TabPane tabId="1">
                </TabPane>
              </TabContent>

            </CardBody>
          </Card>
        </Col>
      </Row>

      <TopupHistory organizationID={currentOrgID} />
    </React.Fragment>
    );
  }
}

export default withRouter(Topup);