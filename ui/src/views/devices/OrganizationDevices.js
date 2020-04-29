import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Col, Container, Nav, NavItem, NavLink, Row, TabContent, TabPane } from 'reactstrap';
import DeviceAdmin from "../../components/DeviceAdmin";
import Loader from "../../components/Loader";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from "../../i18n";
import OrganizationStore from "../../stores/OrganizationStore";
import theme from "../../theme";
import breadcrumbStyles from "../common/BreadcrumbStyles";





const localStyles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
  tabContent: {
    backgroundColor: "#FFFFFF",
    borderRadius: "5px"
  },
  tabPane: {
    backgroundColor: "#EBEFF2",
    borderRadius: "5px",
    padding: "20px"
  }
};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class OrganizationDevices extends Component {
  constructor() {
    super();

    this._isMounted = false;

    this.state = {
      activeMainTabIndex: 0,
      loading: true
    };
  }

  componentDidMount() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const { mainTabIndex } = this.props;
    this._isMounted = true;

    if (mainTabIndex && this._isMounted) {
      this.setState({
        activeMainTabIndex: mainTabIndex
      });
    }

    this.loadOrganization(currentOrgID);
  }

  loadOrganization = async (id) => {
    let resp = await OrganizationStore.get(id);  
    this.setState({
      organization: resp.organization,
      loading: false
    });
  }

  componentWillUnmount() {
    this._isMounted = false;
  }

  toggleMainTabIndex = mainTabIndex => {
    const { activeMainTabIndex } = this.state;
    if (activeMainTabIndex !== mainTabIndex && this._isMounted) {
      this.setState({
        activeMainTabIndex: mainTabIndex
      });
    }
  }

  render() {
    const { activeMainTabIndex, loading } = this.state;
    const { children } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <Container fluid>
        <Row>
          <Col xs={12}>
            <TitleBar
              buttons={<DeviceAdmin organizationID={currentOrgID}>
                {/* TODO - this will take the user to a form where there'll be a selection box where
                          they choose the application id to associate with the new device
                          (prior to being able to access it from the url parameters)
                */}
                <TitleBarButton
                  key={1}
                  label={i18n.t(`${packageNS}:tr000503`)}
                  icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
                  to={`/organizations/${currentOrgID}/devices/create`}
                />
              </DeviceAdmin>}
            >

              <OrgBreadCumb organizationID={currentOrgID} items={[
                { label: i18n.t(`${packageNS}:tr000278`), active: false }]}></OrgBreadCumb>
            </TitleBar>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            {/* <Card>
              <CardBody> */}
                <Nav tabs>
                  <NavItem>
                    <NavLink
                      active={activeMainTabIndex === 0}
                      onClick={() => { this.toggleMainTabIndex(0); }}
                      tag={Link}
                      to={`/organizations/${currentOrgID}/devices`}
                    >
                      <i className="mdi mdi-information-outline"></i>
                      <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000554`)}</span>
                    </NavLink>
                  </NavItem>
                  <NavItem>
                    <NavLink
                      active={activeMainTabIndex === 1}
                      onClick={() => { this.toggleMainTabIndex(1); }}
                      tag={Link}
                      to={`/organizations/${currentOrgID}/applications`}
                    >
                      <i className="mdi mdi-apps"></i>
                      <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000076`)}</span>
                    </NavLink>
                  </NavItem>
                  <NavItem>
                    <NavLink
                      active={activeMainTabIndex === 2}
                      onClick={() => { this.toggleMainTabIndex(2); }}
                      tag={Link}
                      to={`/organizations/${currentOrgID}/device-profiles`}
                    >
                      <i className="mdi mdi-folder-multiple"></i>
                      <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000501`)}</span>
                    </NavLink>
                  </NavItem>
                </Nav>
                <TabContent className={this.props.classes.tabContent} activeTab={activeMainTabIndex}>
                  {loading && <Loader light />}
                  <TabPane tabId={0} className={this.props.classes.tabPane}>
                    {children}
                  </TabPane>
                  <TabPane tabId={1} className={this.props.classes.tabPane}>
                    {children}
                  </TabPane>
                  <TabPane tabId={2} className={this.props.classes.tabPane}>
                    {children}
                  </TabPane>
                </TabContent>
              {/* </CardBody>
            </Card> */}
          </Col>
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(OrganizationDevices);
