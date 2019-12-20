import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import MoneyStore from "../../stores/MoneyStore";
import SessionStore from "../../stores/SessionStore";
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import ModifyEthAccountForm from "./ModifyEthAccountForm";
import NewEthAccountForm from "./NewEthAccountForm";
import styles from "./EthAccountStyle";
import { ETHER } from "../../util/CoinType";


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
          if(resp){
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
        if(resp.username !== SessionStore.getUsername() ){
          alert(`${i18n.t(`${packageNS}:menu.withdraw.incorrect_username_or_password`)}`);
          return false;
        }
        const isOK = await this.verifyUser(resp);
        
        if(isOK) {
          const res = await this.modifyAccount(resp, orgId);
          if(res.status){
            window.location.reload();
          }
        } 
      } catch (error) {
        console.error(error);
        this.setState({ error });
      }
    } 

  render() {
    return(
      <Grid container spacing={24}>
        <Grid item xs={12} md={12} lg={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
              <TitleBar className={this.props.classes.padding}>
                <TitleBarTitle title={i18n.t(`${packageNS}:menu.eth_account.eth_account`)} />
              </TitleBar>
{/*               <Divider light={true}/>
              <div className={this.props.classes.between}>
              <TitleBar>
                <TitleBarTitle component={Link} to="#" title="M2M Wallet" className={this.props.classes.link}/>
                <TitleBarTitle component={Link} to="#" title="/" className={this.props.classes.link}/>
                <TitleBarTitle component={Link} to="#" title={i18n.t(`${packageNS}:menu.withdraw.eth_account`)} className={this.props.classes.link}/>
              </TitleBar>
              </div>*/}
          </div>
        </Grid>
        <Grid item xs={12} md={12} lg={6} className={this.props.classes.column}>
          {/* <Card className={this.props.classes.card}>
            <CardContent> */}
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
            {/* </CardContent>
          </Card> */}
          
        </Grid>
        <Grid item xs={6}>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(ModifyEthAccount));