import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem, Container, Row, Col, Card, CardBody,
  TabContent, TabPane, Nav, NavItem, NavLink } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import theme from "../../theme";
import i18n, { packageNS } from "../../i18n";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import DeviceAdmin from "../../components/DeviceAdmin";
import Admin from "../../components/Admin";

const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  },
  tabs: {
    borderBottom: "1px solid " + theme.palette.divider,
    height: "48px",
    overflow: "visible",
  },
  tabContent: {
    backgroundColor: "#FFFFFF",
    borderRadius: "5px"
  }
});

class ApplicationFUOTADeploymentTabs extends Component {
  constructor() {
    super();

    this.state = {
      activeMainTabAppIndex: 0
    };
  }

  componentDidMount() {
    const { mainTabAppIndex } = this.props;

    if (mainTabAppIndex) {
      this.setState({
        activeMainTabAppIndex: mainTabAppIndex
      });
    }
  }

  toggleMainTabAppIndex = mainTabAppIndex => {
    const { activeMainTabAppIndex } = this.state;
    if (activeMainTabAppIndex !== mainTabAppIndex) {
      this.setState({
        activeMainTabAppIndex: mainTabAppIndex
      });
    }
  }

  render() {
    const { activeMainTabAppIndex } = this.state;
    const { admin, application, children, classes, deleteApplication, fuotaDeployment, mainTabAppIndex } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    return(
      <Container fluid>
        <Row>
          <Col xs={12}>
            <TitleBar
              buttons={
                <Admin organizationID={currentOrgID}>
                  <TitleBarButton
                    label={i18n.t(`${packageNS}:tr000061`)}
                    icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                    onClick={deleteApplication}
                  />
                </Admin>
              }
            >
              <Breadcrumb className={classes.breadcrumb}>
                <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                  `/organizations/${currentOrgID}/applications`
                }>{i18n.t(`${packageNS}:tr000076`)}</Link></BreadcrumbItem>
                <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                  `/organizations/${currentOrgID}/applications/${currentApplicationID}`
                }>{application.application.name}</Link></BreadcrumbItem>
                <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                  `/organizations/${currentOrgID}/applications/${currentApplicationID}/fuota-deployments`
                }>FUOTA (Firmware update jobs)</Link></BreadcrumbItem>
                <BreadcrumbItem active>{fuotaDeployment.fuotaDeployment.name}</BreadcrumbItem>
            </Breadcrumb>
            </TitleBar>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <Nav tabs>
              <NavItem>
                <NavLink
                  active={activeMainTabAppIndex === 0}
                  onClick={() => { this.toggleMainTabAppIndex(0); }}
                  tag={Link}
                  to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/fuota-deployments/${this.props.match.params.fuotaDeploymentID}`}
                >
                  <i className="mdi mdi-information-outline"></i>
                  <span>&nbsp;&nbsp;Information</span>
                </NavLink>
              </NavItem>
              {/* FIXME - temporarily allow any user to access for debugging purposes */}
              {/* {admin && */}
                <NavItem>
                  <NavLink
                    active={activeMainTabAppIndex === 1}
                    onClick={() => { this.toggleMainTabAppIndex(1); }}
                    tag={Link}
                    to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/fuota-deployments/${this.props.match.params.fuotaDeploymentID}/devices`}
                  >
                    <i className="mdi mdi-memory"></i>
                    <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000278`)}</span>
                  </NavLink>
                </NavItem>
              {/* } */}
            </Nav>
            <TabContent className={this.props.classes.tabContent} activeTab={activeMainTabAppIndex}>
              <TabPane tabId={0}>
                {children}
              </TabPane>
              <TabPane tabId={1}>
                {children}
              </TabPane>
            </TabContent>
          </Col>
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(ApplicationFUOTADeploymentTabs);
