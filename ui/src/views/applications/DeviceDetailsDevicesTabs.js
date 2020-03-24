import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Col, Container, Nav, NavItem, NavLink, Row, TabContent, TabPane } from 'reactstrap';
import DeviceAdmin from "../../components/DeviceAdmin";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from "../../i18n";
import theme from "../../theme";
import breadcrumbStyles from "../common/BreadcrumbStyles";




const localStyles = {
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
  tabContent: {
    backgroundColor: "#fff",
    borderRadius: "0 5px 5px 5px",
    borderStyle: "solid",
    borderWidth: "1px",
    // No border at top under active tab since it is white
    borderColor: "transparent #dee2e6 #dee2e6 #dee2e6"
    // padding: "0px"
  }
};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class DeviceDetailsDevicesTabs extends Component {
  constructor() {
    super();

    this.state = {
      activeMainTabDeviceIndex: 0
    };
  }

  componentDidMount() {
    const { mainTabDeviceIndex } = this.props;

    if (mainTabDeviceIndex) {
      this.setState({
        activeMainTabDeviceIndex: mainTabDeviceIndex
      });
    }
  }

  toggleMainTabDeviceIndex = mainTabDeviceIndex => {
    const { activeMainTabDeviceIndex } = this.state;
    if (activeMainTabDeviceIndex !== mainTabDeviceIndex) {
      this.setState({
        activeMainTabDeviceIndex: mainTabDeviceIndex
      });
    }
  }

  render() {
    const { activeMainTabDeviceIndex } = this.state;
    const { admin, application, children, classes, device, deviceProfile, deleteDevice, mainTabDeviceIndex, match, organization } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0" && application;
    const currentOrgName = organization && (organization.name || organization.displayName);

    return(
      <Container fluid>
        <Row>
          <Col xs={12}>
            <TitleBar
              buttons={
                <DeviceAdmin organizationID={match.params.organizationID}>
                  <TitleBarButton
                    label={i18n.t(`${packageNS}:tr000061`)}
                    icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                    onClick={deleteDevice}
                  />
                </DeviceAdmin>
              }
            >
               {
                  isApplication ? (
                    <Breadcrumb className={classes.breadcrumb}>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/applications`
                      }>{i18n.t(`${packageNS}:tr000076`)}</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/applications/${currentApplicationID}`
                        }>{application.application.name}</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/applications/${currentApplicationID}`
                      }>{i18n.t(`${packageNS}:tr000278`)}</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}`
                      }>{device.device.name}</Link></BreadcrumbItem>
                      <BreadcrumbItem active>Show</BreadcrumbItem>
                    </Breadcrumb>
                  ) : (
                    <Breadcrumb className={classes.breadcrumb}>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations`
                      }>Organizations</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}`
                      }>{currentOrgName || currentOrgID}</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/devices`
                      }>{i18n.t(`${packageNS}:tr000278`)}</Link></BreadcrumbItem>
                      <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                        `/organizations/${currentOrgID}/devices/${match.params.devEUI}`
                      }>{match.params.devEUI}</Link></BreadcrumbItem>
                      <BreadcrumbItem active>Show</BreadcrumbItem>
                    </Breadcrumb>
                  )
                }
            </TitleBar>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <Nav tabs>
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 0}
                  onClick={() => { this.toggleMainTabDeviceIndex(0); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}`
                  }
                >
                  <i className="mdi mdi-information-outline"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000280`)}</span>
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 1}
                  onClick={() => { this.toggleMainTabDeviceIndex(1); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/edit`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/edit`
                  }
                >
                  <i className="mdi mdi-pencil"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000298`)}</span>
                </NavLink>
              </NavItem>
              {/* Only available to admins and where the devices has an associated device profile that supports OTAA */}
              {admin && deviceProfile && deviceProfile.deviceProfile.supportsJoin &&
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 2}
                  onClick={() => { this.toggleMainTabDeviceIndex(2); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/keys`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/keys`
                  }
                >
                  <i className="mdi mdi-key"></i>
                  <span>&nbsp;&nbsp;Keys (OTAA)</span>
                </NavLink>
              </NavItem>
              }
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 3}
                  onClick={() => { this.toggleMainTabDeviceIndex(3); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/activation`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/activation`
                  }
                >
                  <i className="mdi mdi-cloud-check"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000311`)}</span>
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 4}
                  onClick={() => { this.toggleMainTabDeviceIndex(4); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/data`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/data`
                  }
                >
                  <i className="mdi mdi-poll"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000317`)}</span>
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 5}
                  onClick={() => { this.toggleMainTabDeviceIndex(5); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/frames`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/frames`
                  }
                >
                  <i className="mdi mdi-video"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000318`)}</span>
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  active={activeMainTabDeviceIndex === 6}
                  onClick={() => { this.toggleMainTabDeviceIndex(6); }}
                  tag={Link}
                  to={
                    currentApplicationID
                    ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/${match.params.devEUI}/fuota-deployments`
                    : `/organizations/${currentOrgID}/devices/${match.params.devEUI}/fuota-deployments`
                  }
                >
                  <i className="mdi mdi-cloud-upload"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000319`)}</span>
                </NavLink>
              </NavItem>
            </Nav>
            <TabContent className={this.props.classes.tabContent} activeTab={activeMainTabDeviceIndex}>
              <TabPane tabId={0}>
                {children}
              </TabPane>
              <TabPane tabId={1}>
                {children}
              </TabPane>
              <TabPane tabId={2}>
                {children}
              </TabPane>
              <TabPane tabId={3}>
                {children}
              </TabPane>
              <TabPane tabId={4}>
                {children}
              </TabPane>
              <TabPane tabId={5}>
                {children}
              </TabPane>
              <TabPane tabId={6}>
                {children}
              </TabPane>
            </TabContent>
          </Col>
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(DeviceDetailsDevicesTabs);
