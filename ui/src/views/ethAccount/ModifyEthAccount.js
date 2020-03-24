import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Row } from 'reactstrap';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import MoneyStore from "../../stores/MoneyStore";
import SessionStore from "../../stores/SessionStore";
import { ETHER } from "../../util/CoinType";
import ModifyEthAccountForm from "./ModifyEthAccountForm";
import NewEthAccountForm from "./NewEthAccountForm";



class ModifyEthAccount extends Component {
  constructor() {
    super();
    this.state = {
      activeAccount: '0'
    };
    this.loadData = this.loadData.bind(this);
  }

  componentDidMount() {
    /*window.analytics.page();*/
    this.loadData();
  }

  loadData() {
    const orgId = this.props.match.params.organizationID;
    
    MoneyStore.getActiveMoneyAccount(ETHER, orgId, resp => {
      this.setState({
        activeAccount: resp.activeAccount,
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
    const loginBody = {};
    loginBody.username = resp.username;
    loginBody.password = resp.password;

    return new Promise((resolve, reject) => {
      SessionStore.login(loginBody, (resp) => {
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
    req.orgId = orgId;
    return new Promise((resolve, reject) => {
      MoneyStore.modifyMoneyAccount(req, resp => {
        resolve(resp);
      })
    });
  }

  onSubmit = async (resp) => {
    const orgId = this.props.match.params.organizationID;

    try {
      const userProfile = await SessionStore.getProfile();
      
      let username = '';
      if(userProfile.body.user.username){
        username = userProfile.body.user.username;
      } 

      if (resp.username !== username) {
        alert(`${i18n.t(`${packageNS}:menu.withdraw.incorrect_username_or_password`)}`);
        return false;
      }

      const isOK = await this.verifyUser(resp);

      if (isOK) {
        const res = await this.modifyAccount(resp, orgId);
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
    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.eth_account.eth_account`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          {this.state.activeAccount &&
            <ModifyEthAccountForm
              submitLabel={i18n.t(`${packageNS}:menu.eth_account.confirm`)}
              onSubmit={this.onSubmit}
              activeAccount={this.state.activeAccount}
            />
          }
          {!this.state.activeAccount &&
            <NewEthAccountForm
              submitLabel={i18n.t(`${packageNS}:menu.eth_account.confirm`)}
              onSubmit={this.onSubmit}
            />
          }
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(ModifyEthAccount);