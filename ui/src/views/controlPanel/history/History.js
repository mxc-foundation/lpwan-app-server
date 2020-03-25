import { withStyles } from "@material-ui/core/styles";
import classNames from "classnames";
import React, { Component } from "react";
import { Link, Route, Switch, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Card, CardBody, Col, Nav, NavItem, Row } from 'reactstrap';
import TitleBar from "../../../components/TitleBar";
import i18n, { packageNS } from '../../../i18n';
import { SUPER_ADMIN } from "../../../util/M2mUtil";
import breadcrumbStyles from "../../common/BreadcrumbStyles";
import SuperNodeEthAccount from "./EthAccount";




const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class SupernodeHistory extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeTab: '0',
      loading: false,
      admin: false,
      income:0
    };

    this.onChangeTab = this.onChangeTab.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
  }

  componentDidMount() {
    this.setState({loading:true});
    this.locationToTab();
    this.setState({loading:false});
    //this.getIncome();
  }

  componentDidUpdate(oldProps) {
    if (this.props == oldProps) {
      return;
    }

    this.locationToTab();
  }

  getIncome(){
    /* topupStore.getIncome(0, resp => {
      this.setState({income:resp.amount});
    }); */
  }

  onChangeTab(e, v) {
    this.setState({
      tab: v,
    });
  }

  locationToTab() {
    let tab = 0;
    if (window.location.href.endsWith("/eth-account")) {
      tab = 1;
    } else if (window.location.href.endsWith("/network-activity")) {
      tab = 2;
    }  
    
    this.setState({
      activeTab:tab + '',
    });
  }

  render() {
    const { classes } = this.props;

    return(
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
                onClick={() => {
                  // Change the sidebar content
                  this.props.switchToSidebarId('DEFAULT');
                }}
              >
                Control Panel
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem className={classes.breadcrumbItem}>Wallet</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.history.history`)}</BreadcrumbItem>
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
                      to={`/control-panel/history/`}
                    >{i18n.t(`${packageNS}:menu.history.eth_account`)}</Link>
                  </NavItem>
                </Nav>

                <Row className="pt-3">
                  <Col>
                    <Switch>
                      <Route exact path={`${this.props.match.path}/`} render={props => <SuperNodeEthAccount organizationID={SUPER_ADMIN} {...props} />} />
                    </Switch>
                  </Col>
                </Row>
              </CardBody>
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(SupernodeHistory));