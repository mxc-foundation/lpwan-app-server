import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Col, Container, Nav, NavItem, NavLink, Row, TabContent, TabPane } from 'reactstrap';
import Admin from "../../components/Admin";
import Modal from "../../components/Modal";
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
    backgroundColor: "#FFFFFF",
    borderRadius: "5px"
  }
};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class ApplicationDevices extends Component {
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

  openModal = () => {
    this.setState({ openModal: true });
  }

  closeModal = () => {
    this.setState({ openModal: false });
  }

  render() {
    const { activeMainTabAppIndex } = this.state;
    const { admin, application, children, classes, deleteApplication, mainTabAppIndex } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (
      <Container fluid>
        <Row>
          <Col xs={12}>
            <TitleBar
              buttons={
                <Admin organizationID={currentOrgID}>
                  {this.state.openModal && <Modal
                    title={""}
                    closeModal={this.closeModal}
                    left={i18n.t(`${packageNS}:menu.common.cancel`)}
                    right={i18n.t(`${packageNS}:menu.common.confirm`)}
                    context={i18n.t(`${packageNS}:menu.application.del_application`)}
                    callback={() => {
                      this.setState({ openModal: false });
                      deleteApplication(application.application.id);
                    }
                    }
                  />}

                  <TitleBarButton
                    label={i18n.t(`${packageNS}:tr000061`)}
                    color="danger"
                    icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                    //onClick={()=>(deleteApplication(application.application.id))}
                    onClick={this.openModal}
                  />
                </Admin>
              }
            >
              <Breadcrumb className={classes.breadcrumb}>
                <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                  `/organizations/${currentOrgID}/applications`
                }>{i18n.t(`${packageNS}:tr000076`)}</Link></BreadcrumbItem>
                <BreadcrumbItem active>{application.application.name}</BreadcrumbItem>
              </Breadcrumb>
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
                  active={activeMainTabAppIndex === 0}
                  onClick={() => { this.toggleMainTabAppIndex(0); }}
                  tag={Link}
                  to={`/organizations/${currentOrgID}/applications/${this.props.match.params.applicationID}`}
                >
                  <i className="mdi mdi-memory"></i>
                  <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000278`)}</span>
                </NavLink>
              </NavItem>
              {admin &&
                <NavItem>
                  <NavLink
                    active={activeMainTabAppIndex === 1}
                    onClick={() => { this.toggleMainTabAppIndex(1); }}
                    tag={Link}
                    to={`/organizations/${currentOrgID}/applications/${this.props.match.params.applicationID}/edit`}
                  >
                    <i className="mdi mdi-pencil"></i>
                    <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000395`)}</span>
                  </NavLink>
                </NavItem>
              }
              {admin &&
                <NavItem>
                  <NavLink
                    active={activeMainTabAppIndex === 2}
                    onClick={() => { this.toggleMainTabAppIndex(2); }}
                    tag={Link}
                    to={`/organizations/${currentOrgID}/applications/${this.props.match.params.applicationID}/integrations`}
                  >
                    <i className="mdi mdi-cloud"></i>
                    <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000384`)}</span>
                  </NavLink>
                </NavItem>
              }
              {admin &&
                <NavItem>
                  <NavLink
                    active={activeMainTabAppIndex === 3}
                    onClick={() => { this.toggleMainTabAppIndex(3); }}
                    tag={Link}
                    to={`/organizations/${currentOrgID}/applications/${this.props.match.params.applicationID}/fuota-deployments`}
                  >
                    <i className="mdi mdi-cloud-upload"></i>
                    <span>&nbsp;&nbsp;{i18n.t(`${packageNS}:tr000555`)}</span>
                  </NavLink>
                </NavItem>
              }
            </Nav>
            <TabContent className={this.props.classes.tabContent} activeTab={activeMainTabAppIndex}>
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
            </TabContent>
            {/* </CardBody>
            </Card> */}
          </Col>
        </Row>
      </Container>
    );
  }
}

export default withStyles(styles)(ApplicationDevices);
