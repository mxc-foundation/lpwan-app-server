import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../../i18n';
import TitleBar from "../../../components/TitleBar";
import SessionStore from "../../../stores/SessionStore";
import SupernodeStore from "../../../stores/SupernodeStore";
import ModifyEthAccountForm from "../../ethAccount/ModifyEthAccountForm";
import NewEthAccountForm from "../../ethAccount/NewEthAccountForm";
import { ETHER } from "../../../util/CoinType";
import { SUPER_ADMIN } from "../../../util/M2mUtil";
import Admin from "../../../components/Admin";

import breadcrumbStyles from "../../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class SuperNodeEth extends Component {
  constructor() {
    super();
    this.state = {
      activeAccount: '0'
    };
    this.loadData = this.loadData.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  loadData() {
    SupernodeStore.getSuperNodeActiveMoneyAccount(ETHER, SUPER_ADMIN, resp => {
      this.setState({
        activeAccount: resp.supernodeActiveAccount,
      });
    });

  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }

    this.loadData();
  }


  verifyUser(resp) {
    const login = {};
    login.username = resp.username;
    login.password = resp.password;

    return new Promise((resolve, reject) => {
      SessionStore.login(login, (resp) => {
        if (resp) {
          resolve(resp);
        } else {
          alert(`${i18n.t(`${packageNS}:menu.withdraw.incorrect_username_or_password`)}`);
          return false;
        }
      })
    });
  }

  modifyAccount(req, orgId) {
    req.moneyAbbr = ETHER;
    return new Promise((resolve, reject) => {
      SupernodeStore.addSuperNodeMoneyAccount(req, orgId, resp => {
        resolve(resp);
      })
    });
  }

  onSubmit = async (resp) => {
    try {
      if (resp.username !== SessionStore.getUsername()) {
        alert(`${i18n.t(`${packageNS}:menu.withdraw.incorrect_username_or_password`)}`);
        return false;
      }
      const isOK = await this.verifyUser(resp);

      if (isOK) {
        const res = await this.modifyAccount(resp, SUPER_ADMIN);
        if (res.status) {
          window.location.reload();
        }
      }
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  render() {
    const { classes } = this.props;

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <Admin>
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
            </Admin>
            <BreadcrumbItem className={classes.breadcrumbItem}>{i18n.t(`${packageNS}:tr000451`)}</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.eth_account.eth_account`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          {this.state.activeAccount &&
            <ModifyEthAccountForm
              submitLabel={i18n.t(`${packageNS}:menu.common.confirm`)}
              onSubmit={this.onSubmit}
              activeAccount={this.state.activeAccount}
            />
          }
          {!this.state.activeAccount &&
            <NewEthAccountForm
              submitLabel={i18n.t(`${packageNS}:menu.common.confirm`)}
              onSubmit={this.onSubmit}
            />
          }
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(SuperNodeEth));
