import React from "react";

import TextField from '@material-ui/core/TextField';
import { Redirect } from 'react-router-dom'
import clsx from 'clsx'
import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Grid from "@material-ui/core/Grid";
import TitleBar from "../../components/TitleBar";


import TitleBarTitle from "../../components/TitleBarTitle";
import ExtLink from "../../components/ExtLink";
import Typography from '@material-ui/core/Typography';
import StakeForm from "./StakeForm";
import StakeStore from "../../stores/StakeStore";
import Button from "@material-ui/core/Button";
import Modal from "../common/Modal";
import ModalTimer from "../common/ModalTimer";
//import Button from "@material-ui/core/Button";
import Spinner from "../../components/ScaleLoader";
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { MXC } from "../../util/Messages"
import { EXT_URL_STAKE } from "../../util/Data"
import { withStyles } from "@material-ui/core/styles";
import { withRouter, Link } from "react-router-dom";
import styles from "./StakeStyle"
import { Divider } from "@material-ui/core";
import InfoCard from "../topup/InfoCard";

class SetStake extends FormComponent {

  state = {
    amount: 0,
    revRate: 0,
    event: null,
    isUnstake: false,
    info: '',
    modal: null,
    modalTimer: null,
    infoStatus: 0,
    notice: {
      succeed: i18n.t(`${packageNS}:menu.messages.congratulations_stake_set`),
      unstakeSucceed: i18n.t(`${packageNS}:menu.messages.unstake_successful`),
      warning: i18n.t(`${packageNS}:menu.messages.close_to_acquiring`)
    },
  }

  componentDidMount() {
    this.loadData();
  }

  componentWillReceiveProps() {
    this.loadStakeTextTranslation();
  }

  loadData = () => {
    const resp = StakeStore.getActiveStakes(this.props.match.params.organizationID);
    resp.then((res) => {
      let amount = 0;
      let isUnstake = false;

      if (res.actStake !== null) {
        amount = res.actStake.Amount;
        isUnstake = true;
      }

      this.setState({
        amount,
        isUnstake,
        info: i18n.t(`${packageNS}:menu.messages.staking_enhances`)
      })
    })

  }

  loadStakeTextTranslation = () => {
    this.setState({
      info: i18n.t(`${packageNS}:menu.messages.staking_enhances`)
    })
  }

  onChange = (event, name) => {
    this.setState({
      [name]: event.target.value
    });
  }

  reset = () => {
    this.setState({
      amount: 0
    })
  }

  onSubmit = (amt) => {
    const amount = parseFloat(amt);
    const orgId = this.props.match.params.organizationID;
    const req = {
      orgId,
      amount
    }

    this.setState({ modal: true });

    if (this.state.isUnstake) {
      this.unstake(this.state.event, orgId);
    } else {
      this.stake(this.state.event, req);
    }

    this.setState({ modal: false });
  }

  confirm = (e) => {
    if (this.state.amount === 0) {
      return false;
    }

    if (this.state.isUnstake) {
      const orgId = this.props.match.params.organizationID;
      this.unstake(e, orgId);
    } else {
      this.setState({
        modal: true,
        event: e,
        amount: this.state.amount
      });
    }
  }

  openModalTimer = (data) => {
    this.setState({
      modalTimer: true,
      modal: null
    });
  }

  stake = (e, req) => {
    e.preventDefault();
    const resp = StakeStore.stake(req);
    resp.then((res) => {
      if (res.body.status === i18n.t(`${packageNS}:menu.staking.stake_success`)) {
        this.setState({
          isUnstake: true,
          info: i18n.t(`${packageNS}:menu.messages.congratulations_stake_set`),
          infoStatus: 1,
        });
        setInterval(() => this.displayInfo(), 5000);
      } else {
        this.setState({
          info: res.body.status,
          infoStatus: 2,
        });
        setInterval(() => this.displayInfo(), 5000);
      }
    })
  }

  unstake = (e, orgId) => {
    e.preventDefault();
    const resp = StakeStore.unstake(orgId);
    resp.then((res) => {
      this.setState({
        isUnstake: false,
        amount: 0,
        info: i18n.t(`${packageNS}:menu.messages.unstake_successful`),
        infoStatus: 1,
      });
      setInterval(() => this.displayInfo(), 5000);
    })
  }

  displayInfo = () => {
    this.setState({
      info: i18n.t(`${packageNS}:menu.messages.staking_enhances`),
      infoStatus: 0
    });
  }
  showModal = (modal) => {
    this.setState({ modal });
  }

  handleCloseModal = () => {
    this.setState({
      modal: null
    })
  }

  handleOnclick = () => {
    this.props.history.push(`/history/${this.props.match.params.organizationID}/stake`);
  }

  handleProgress = (oldCompleted) => {
    this.setState({
      modalTimer: null
    })
    if (oldCompleted === 100) {
      this.onSubmit(this.state.amount);
    }
  }

  render() {
    /* if (this.props.txinfo === undefined) {
      return(<Spinner on={this.state.loading}/>);
    } */
    const info = this.state.info;
    const infoBoxCss = [this.props.classes.infoBox,
    this.props.classes.infoBoxSucceed,
    this.props.classes.infoBoxError];
    const title = this.state.isUnstake ? i18n.t(`${packageNS}:menu.messages.unstake`) : i18n.t(`${packageNS}:menu.messages.set_stake`);
    let path = null;
    if(this.props.match.params.organizationID === process.env.REACT_APP_SUPER_ADMIN_LPWAN){
      path = '/control-panel/modify-account/';
    }else{
      path = `/modify-account/${this.props.match.params.organizationID}`;
    }

    return (
      <Grid container spacing={24} className={this.props.classes.backgroundColor}>
        <Grid item xs={12} md={12} lg={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
            <TitleBar className={this.props.classes.padding}>
              <TitleBarTitle title={title} />
            </TitleBar>
            {/* <div className={this.props.classes.between}>
              <TitleBar>
                <TitleBarTitle component={Link} to="#" title="M2M Wallet" className={this.props.classes.link} />
                <TitleBarTitle component={Link} to="#" title="/" className={this.props.classes.link} />
                <TitleBarTitle component={Link} to="#" title={title} className={this.props.classes.link} />
              </TitleBar>
            </div> */}
          </div>
        </Grid>

        {this.state.modal &&
          <Modal title={i18n.t(`${packageNS}:menu.messages.confirmation`)} description={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)} onProgress={this.handleProgress} onCancelProgress={this.handleCancel} onClose={this.handleCloseModal} open={!!this.state.modal} data={this.state.modal} onSubmit={this.openModalTimer} />}

        {this.state.modalTimer &&
          <ModalTimer title={i18n.t(`${packageNS}:menu.messages.stake_proc_tit`)} description={i18n.t(`${packageNS}:menu.messages.stake_proc_desc`)} onProgress={this.handleProgress} onCancelProgress={this.handleProgress} onProcClose={this.handleCloseProcModal} open={!!this.state.modalTimer} data={this.state.modalTimer} onProgress={this.handleProgress} onSubmit={this.onSubmit} />}

        <Grid item xs={12} md={12} lg={6} className={this.props.classes.pRight}>
          {/* <Card className={this.props.classes.card}>
            <CardContent> */}
              <StakeForm isUnstake={this.state.isUnstake} label={this.state.isUnstake ? i18n.t(`${packageNS}:menu.messages.withdraw_stake`) : i18n.t(`${packageNS}:menu.messages.set_stake`)} onChange={this.onChange} amount={this.state.amount} revRate={this.state.revRate} reset={this.reset} confirm={this.confirm} />
            {/* </CardContent>
          </Card> */}

        </Grid>
        <Grid item xs={12} md={12} lg={6} className={this.props.classes.pLeft}>
          <div className={infoBoxCss[this.state.infoStatus]}>
            <Typography gutterBottom>
              {this.state.info}
            </Typography>
            <div className={this.props.classes.between}>
              <ExtLink to={EXT_URL_STAKE} context={i18n.t(`${packageNS}:menu.common.learn_more`)} />
            </div>
          </div>
          <InfoCard path={path}/>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(SetStake));
