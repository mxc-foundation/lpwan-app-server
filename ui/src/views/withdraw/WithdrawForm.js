import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Button, FormGroup, Label, Input } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "./WithdrawStyle"
import i18n, { packageNS } from "../../i18n";
import NumberFormat from 'react-number-format';
import { ETHER } from "../../util/CoinType"

import breadcrumbStyles from "../common/BreadcrumbStyles";
import Modal from './Modal';
import ModalCom from './ModalComplete';
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
      comDialog: false,
      activeAccount: "",
      balance: 0,
      amount: 0,
      fee: 0
    }
  }

  loadData = async () => {
    const orgId = this.props.match.params.organizationID;

    var result = await getWalletBalance(orgId);
    const fee = result.withdrawFee;
    
    WithdrawStore.getWithdrawFee(0, resp => {
      const object = this.state;
      object.fee = resp.withdrawFee;

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

  openModal = (modal) => {
    this.setState({
      [modal]: true
    });
  };

  handleChange = (event) => {
    this.setState({ [event.target.id]: event.target.value });
  }

  submitWithdrawReq = () => {
    const orgId = this.props.match.params.organizationID;

    let req = {};
    req.orgId = orgId;
    req.moneyAbbr = ETHER;
    req.amount = this.state.amount;
    req.ethAddress = this.state.activeAccount;
    req.availableBalance = this.state.balance;

    WithdrawStore.withdrawReq(req, (res) => {
      const object = this.state;
      object.loading = false;

      this.props.history.push(`/withdraw/${this.props.match.params.organizationID}`);
      this.openModal('comDialog');
    });
  }

  render() {
    const { classes } = this.props;

    return (
      <React.Fragment>
        {this.state.nsDialog && <Modal
          title={i18n.t(`${packageNS}:menu.withdraw.confirm_modal_title`)}
          context={(this.state.status) ? i18n.t(`${packageNS}:menu.withdraw.confirm_text`) : i18n.t(`${packageNS}:menu.withdraw.deny_text`)}
          status={this.state.status}
          amount={this.state.amount}
          handleChange={this.handleChange}
          closeModal={() => this.setState({ nsDialog: false })}
          callback={() => { this.submitWithdrawReq() }} />}
        {this.state.comDialog && <ModalCom
          closeModal={() => this.setState({ comDialog: false })}
        />}
        <Row>
          <Col xs="4">
            <FormGroup>
              <Label for="amount">{i18n.t(`${packageNS}:menu.withdraw.amount`)}</Label>
              <NumberFormat id="amount" className={classes.s_input} value={this.state.amount} onChange={this.handleChange} thousandSeparator={true} suffix={' MXC'} />
            </FormGroup>
            <FormGroup>
              <Label for="activeAccount">{i18n.t(`${packageNS}:menu.withdraw.destination`)}</Label>
              <Input type="text" name="activeAccount" id="activeAccount" value={this.state.activeAccount} placeholder={'0x00000000000000000000000000000000'} readOnly />
            </FormGroup>
            <FormGroup>
              <Label for="fee">{i18n.t(`${packageNS}:menu.withdraw.currentFee`)}:</Label>{' '}
              <NumberFormat id="fee" displayType={'text'} value={this.state.fee} onChange={this.handleChange} thousandSeparator={true} readOnly={true}  suffix={' MXC'} />
            </FormGroup>
          </Col>
          <Col xs="2"></Col>
          <Col xs="6">
            <FormGroup>
              <Label for="exampleEmail">{i18n.t(`${packageNS}:menu.withdraw.available`)} MXC</Label>
              <NumberFormat id="amount" className={classes.t_input}  value={this.state.balance} onChange={this.handleChange} thousandSeparator={true} readOnly={true} suffix={' MXC'} />
            </FormGroup>
          </Col>
        </Row>
        <Row>
          <Col>
            <FormGroup>
              <Button onClick={() => { this.openModal('nsDialog') }} color="primary">{i18n.t(`${packageNS}:menu.withdraw.submit_withdrawal_request`)}</Button>
            </FormGroup>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));
