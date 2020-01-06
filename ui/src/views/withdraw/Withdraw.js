import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Grid from "@material-ui/core/Grid";
import i18n, { packageNS } from '../../i18n';
import TitleBarTitle from "../../components/TitleBarTitle";
import MoneyStore from "../../stores/MoneyStore";
import WithdrawStore from "../../stores/WithdrawStore";
import SupernodeStore from "../../stores/SupernodeStore";
import WalletStore from "../../stores/WalletStore";
import Modal from "./Modal";
import { withStyles } from "@material-ui/core/styles";
import styles from "./WithdrawStyle"
import { ETHER } from "../../util/CoinType"
import { SUPER_ADMIN } from "../../util/M2mUtil"
import theme from "../../theme";
import TableCell from "@material-ui/core/TableCell";

function formatNumber(number) {
  return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function loadWithdrawFee(ETHER, organizationID) {
  return new Promise((resolve, reject) => {
    WithdrawStore.getWithdrawFee(ETHER, organizationID,
      resp => {
        resp.moneyAbbr = ETHER;
        resolve(resp);
      })
  });
}  

function loadCurrentAccount(ETHER, orgId) {
  return new Promise((resolve, reject) => {
    if (orgId === SUPER_ADMIN) {
      SupernodeStore.getSuperNodeActiveMoneyAccount(ETHER, orgId, resp => {
        resolve(resp.supernodeActiveAccount);
        
      });
    }else{
      MoneyStore.getActiveMoneyAccount(ETHER, orgId, resp => {
        resolve(resp.activeAccount);
        
      });
    }
  });
}

      
function loadWalletBalance(orgId) {
  return new Promise((resolve, reject) => {
    WalletStore.getWalletBalance(orgId,
      resp => {
        /* Object.keys(resp).forEach(attr => {
          const value = resp[attr];
  
          if (typeof value === 'number') {
            resp[attr] = formatNumber(value);
          }
        }); */
        resolve(resp);
      });
  });
}

class Withdraw extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      modal: null,
    };
  }

  loadData = async () => {
    try {
      const orgId = this.props.match.params.organizationID;
      this.setState({loading: true})
      var result = await loadWithdrawFee(ETHER, orgId);
      var wallet = await loadWalletBalance(orgId);
      var account = await loadCurrentAccount(ETHER, orgId);
      
      /* this.setState({
        activeAccount: resp.supernodeActiveAccount,
      }); */

      const txinfo = {};
      txinfo.withdrawFee = result.withdrawFee;
      txinfo.balance = wallet.balance;
      
      txinfo.account = account;

      this.setState({
        txinfo
      });
      this.setState({loading: false})
    } catch (error) {
      this.setState({loading: false})
      console.error(error);
      this.setState({ error });
    }
  }

  componentDidMount() {
    //this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
    this.loadData();
  }
  
  showModal(modal) {
    this.setState({ modal });
  }

  onSubmit = (e, apiWithdrawReqRequest) => {
    e.preventDefault();
    this.showModal(apiWithdrawReqRequest);
  }

  handleCloseModal = () => {
    this.setState({
      modal: null
    })
  }

  onConfirm = (data) => {
    data.moneyAbbr = ETHER;
    data.orgId = this.props.match.params.organizationID;
    if(data.amount === 0){
      alert(i18n.t(`${packageNS}:menu.messages.invalid_amount`));
      return false;
    } 

    if(data.destination){
      alert(i18n.t(`${packageNS}:menu.messages.invalid_account`));
      return false;
    }
    
    this.setState({loading: true});
    WithdrawStore.WithdrawReq(data, resp => {
      this.setState({loading: false});
    });

  }

  render() {
    return (
      <Grid container spacing={24} className={this.props.classes.backgroundColor}>
        {this.state.modal && 
          <Modal title={i18n.t(`${packageNS}:menu.messages.confirmation`)} description={i18n.t(`${packageNS}:menu.messages.confirmation_text`)} onClose={this.handleCloseModal} open={!!this.state.modal} data={this.state.modal} onConfirm={this.onConfirm} />}
        <Grid item xs={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
            <TitleBarTitle title={i18n.t(`${packageNS}:menu.withdraw.withdraw`)} />
          </div>

        </Grid>
        <Grid item xs={6}>
          <TableCell align={this.props.align}>
                    <span style={
                      {
                        textDecoration: "none",
                        color: theme.palette.primary.main,
                        cursor: "pointer",
                        padding: 0,
                        fontWeight: "bold",
                        fontSize: 20,
                        opacity: 0.7,
                        "&:hover": {
                          opacity: 1,
                        }
                      }
                    } className={this.props.classes.link} >
                        {i18n.t(`${packageNS}:menu.messages.coming_soon`)}
                    </span>
          </TableCell>
          {/*<WithdrawForm
            submitLabel={i18n.t(`${packageNS}:menu.withdraw.withdraw`)}
            txinfo={this.state.txinfo} {...this.props}
            onSubmit={this.onSubmit}
          />*/}
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));