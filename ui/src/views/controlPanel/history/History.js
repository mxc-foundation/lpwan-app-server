import React, { Component } from "react";
import { Route, Switch, Link, withRouter } from "react-router-dom";
import classNames from "classnames";
import { Breadcrumb, BreadcrumbItem, Nav, NavItem, Row, Col, Card, CardBody } from 'reactstrap';



import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
import TitleBarTitle from "../../../components/TitleBarTitle";
import Spinner from "../../../components/ScaleLoader"
import SuperNodeEthAccount from "./EthAccount";

import topupStore from "../../../stores/TopupStore";
import styles from "./HistoryStyle";
import { SUPER_ADMIN } from "../../../util/M2mUtil";



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
    this.getIncome();
  }

  componentDidUpdate(oldProps) {
    if (this.props == oldProps) {
      return;
    }

    this.locationToTab();
  }

  getIncome(){
    topupStore.getIncome(0, resp => {
      this.setState({income:resp.amount});
    });
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
      
    return(
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
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

export default withRouter(SupernodeHistory);