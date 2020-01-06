import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import classNames from "classnames";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';


import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import GatewayAdmin from "../../components/GatewayAdmin";
import GatewayStore from "../../stores/GatewayStore";
import SessionStore from "../../stores/SessionStore";
import GatewayDetails from "./GatewayDetails";
import UpdateGateway from "./UpdateGateway";
import GatewayDiscovery from "./GatewayDiscovery";
import GatewayFrames from "./GatewayFrames";


class GatewayLayout extends Component {
  constructor() {
    super();
    this.state = {
      activeTab: '0',
      admin: false,
    };
    this.deleteGateway = this.deleteGateway.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
    this.setIsAdmin = this.setIsAdmin.bind(this);
  }

  componentDidMount() {
    GatewayStore.get(this.props.match.params.gatewayID, resp => {
      this.setState({
        gateway: resp,
      });
    });

    SessionStore.on("change", this.setIsAdmin);
    this.setIsAdmin();
    this.locationToTab();
  }

  componentDidUpdate(prevProps) {
    if (this.props === prevProps) {
      return;
    }

    this.locationToTab();
  }

  componentWillUnmount() {
    SessionStore.removeListener("change", this.setIsAdmin);
  }

  setIsAdmin() {
    this.setState({
      admin: SessionStore.isAdmin() || SessionStore.isOrganizationGatewayAdmin(this.props.match.params.organizationID),
    });
  }

  deleteGateway() {
    if (window.confirm("Are you sure you want to delete this gateway?")) {
      GatewayStore.delete(this.props.match.params.gatewayID, () => {
        this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
      });
    }
  }

  locationToTab() {
    let tab = 0;

    if (window.location.href.endsWith("/edit")) {
      tab = 1;
    } else if (window.location.href.endsWith("/discovery")) {
      tab = 2;
    } else if (window.location.href.endsWith("/frames")) {
      tab = 3;
    }

    if (tab > 0 && !this.state.admin) {
      tab = tab - 1;
    }

    this.setState({
      activeTab: tab + '',
    });
  }

  render() {
    return (
      this.state.gateway ? <React.Fragment>
        <TitleBar
          buttons={<GatewayAdmin organizationID={this.props.match.params.organizationID}>
            <TitleBarButton
              key={1}
              color="danger"
              label={i18n.t(`${packageNS}:tr000061`)}
              icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
              onClick={this.deleteGateway}
            />
          </GatewayAdmin>}
        >

          <TitleBarTitle title={i18n.t(`${packageNS}:tr000063`)} />
          <Breadcrumb>
            <BreadcrumbItem><Link to={`/organizations/${this.props.match.params.organizationID}/gateways`}>{i18n.t(`${packageNS}:tr000063`)}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{this.state.gateway.gateway.name}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <Nav tabs>
                  <NavItem>
                    <Link
                      className={classNames('nav-link', { active: this.state.activeTab === '0' })}
                      to={`/organizations/${this.props.match.params.organizationID}/gateways/${this.props.match.params.gatewayID}`}
                    >{i18n.t(`${packageNS}:tr000423`)}</Link>
                  </NavItem>
                  {this.state.admin && <NavItem>
                    <Link
                      className={classNames('nav-link', { active: this.state.activeTab === '1' })}
                      to={`/organizations/${this.props.match.params.organizationID}/gateways/${this.props.match.params.gatewayID}/edit`}
                    >{i18n.t(`${packageNS}:tr000298`)}</Link>
                  </NavItem>}
                  <NavItem>
                    <Link
                      className={classNames('nav-link', { active: this.state.activeTab === '2' })}
                      disabled={!this.state.gateway.gateway.discoveryEnabled} 
                      to={`/organizations/${this.props.match.params.organizationID}/gateways/${this.props.match.params.gatewayID}/discovery`}
                    >{i18n.t(`${packageNS}:tr000095`)}</Link>
                  </NavItem>
                  <NavItem>
                    <Link
                      className={classNames('nav-link', { active: this.state.activeTab === '3' })}
                      to={`/organizations/${this.props.match.params.organizationID}/gateways/${this.props.match.params.gatewayID}/frames`}
                    >{i18n.t(`${packageNS}:tr000247`)}</Link>
                  </NavItem>
                </Nav>

                <Row className="pt-3">
                  <Col>
                    <Switch>
                      <Route exact path={`${this.props.match.path}`} render={props => <GatewayDetails gateway={this.state.gateway.gateway} lastSeenAt={this.state.gateway.lastSeenAt} {...props} />} />
                      <Route exact path={`${this.props.match.path}/edit`} render={props => <UpdateGateway gateway={this.state.gateway.gateway} {...props} />} />
                      <Route exact path={`${this.props.match.path}/discovery`} render={props => <GatewayDiscovery gateway={this.state.gateway.gateway} {...props} />} />
                      <Route exact path={`${this.props.match.path}/frames`} render={props => <GatewayFrames gateway={this.state.gateway.gateway} {...props} />} />
                    </Switch>
                  </Col>
                </Row>
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment> : <div></div>
    );
  }
}

export default withRouter(GatewayLayout);
