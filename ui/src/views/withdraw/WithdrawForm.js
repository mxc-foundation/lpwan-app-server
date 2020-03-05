import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Button,FormGroup, Label, Input } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "./WithdrawStyle"
import i18n, { packageNS } from "../../i18n";
import NumberFormat from 'react-number-format';
import {ETHER} from "../../util/CoinType"

import breadcrumbStyles from "../common/BreadcrumbStyles";
import Modal from './Modal';
import WithdrawStore from "../../stores/WithdrawStore";
import MoneyStore from "../../stores/MoneyStore";
import WalletStore from "../../stores/WalletStore";

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

function getWalletBalance(organizationId) {
  return new Promise((resolve, reject) => {
    WalletStore.getWalletBalance(organizationId, resp => {
      return resolve(resp);
    });
  });
}

class Withdraw extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nsDialog: false,
      activeAccount: "",
      balance: 0,
      amount: 0
    }
  }

  loadData = async () => {
    const orgId = this.props.match.params.organizationID;

    var result = await getWalletBalance(orgId);
    const balance = result.balance;

    MoneyStore.getActiveMoneyAccount(0, orgId, resp => {
      const activeAccount = resp.activeAccount;

      const object = this.state;
      object.balance = balance;
      object.activeAccount = activeAccount;

      this.setState({
        object
      });
    });
    
  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState !== this.state && prevState.data !== this.state.data) {

    }
  }

  openModal = (row, status) => {
    this.setState({
      nsDialog: true,
      row,
      status
    });
  };

  handleChange = (event) => {
    this.setState({ [event.target.id]: event.target.value });
  }

  submitWithdrawReq = () => {
    const orgId = this.props.match.params.organizationID;

    let req = {};
    req.orgId= orgId;
    req.moneyAbbr= ETHER;
    req.amount= this.state.amount;
    req.ethAddress= this.state.activeAccount;
    req.availableBalance= this.state.balance;

    WithdrawStore.withdrawReq(req, (res) => {
      const object = this.state;
      object.loading = false;
      this.props.history.push(`/withdraw/${this.props.match.params.organizationID}`);
    });
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (

      <React.Fragment>
        {this.state.nsDialog && <Modal
          title={i18n.t(`${packageNS}:menu.withdraw.confirm_modal_title`)}
          context={(this.state.status) ? i18n.t(`${packageNS}:menu.withdraw.confirm_text`) : i18n.t(`${packageNS}:menu.withdraw.deny_text`)}
          status={this.state.status}
          row={this.state.row}
          handleChange={this.handleChange}
          closeModal={() => this.setState({ nsDialog: false })}
          callback={() => { this.confirm(this.state.row, this.state.status) }} />}
        <Row>
          <Col xs="4">
            <FormGroup>
              <Label for="activeAccount">Destination</Label>
              <Input type="text" name="activeAccount" id="activeAccount" value={this.state.activeAccount} readOnly/>
            </FormGroup>
            <FormGroup>
              <Label for="amount">Amount</Label>
              <NumberFormat id="amount" className={classes.s_input} value={this.state.amount} onChange={this.handleChange} thousandSeparator={true} suffix={' MXC'} />
            </FormGroup>
            <FormGroup>
              <Label for="exampleEmail">Available MXC</Label>
              <NumberFormat id="amount" className={classes.s_input} value={this.state.balance} onChange={this.handleChange} thousandSeparator={true} readOnly={true} suffix={' MXC'} />
            </FormGroup>
          </Col>
          <Col xs="4"></Col>
          <Col xs="4"></Col>
        </Row>
        <Row>
          <Col>
            <FormGroup>
              <Button onClick={this.submitWithdrawReq} color="primary">Submit Withdrawal Request</Button>
            </FormGroup>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));
