import React from "react";

import { Breadcrumb, BreadcrumbItem, Container, Row, Col, Card, CardBody, CardHeader, CardFooter, CardText } from 'reactstrap';
import { withRouter, Link } from "react-router-dom";


import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import TitleBar from "../../components/TitleBar";


import ExtLink from "../../components/ExtLink";
import Typography from '@material-ui/core/Typography';
import StakeForm from "./StakeForm";
import StakeStore from "../../stores/StakeStore";
import Modal from "../../components/Modal";
import ModalWithProgress from "../../components/ModalWithProgress";
//import Modal from "../common/Modal";
import ModalTimer from "../common/ModalTimer";
//import Button from "@material-ui/core/Button";
import Spinner from "../../components/ScaleLoader";

import { EXT_URL_STAKE } from "../../util/Data"

import InfoCard from "../topup/InfoCard";

class SetStake extends FormComponent {

  state = {
    amount: 0,
    revRate: 0,
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

  onSubmit = () => {
    const amount = parseFloat(this.state.amount);
    const orgId = this.props.match.params.organizationID;
    const req = {
      orgId,
      amount
    }

    //this.setState({ modal: true });

    if (this.state.isUnstake) {
      this.unstake(orgId);
    } else {
      this.stake(req);
    }

    this.setState({ modal: false });
  }

  confirm = (amount) => {
    this.setState({ modal: false });
    if (amount === 0) {
      return false;
    }

    if (this.state.isUnstake) {
      const orgId = this.props.match.params.organizationID;
      this.unstake(orgId);
    } else {
      this.setState({
        modal: true,
        //event: e,
        amount: amount.amount
      });
    }
  }

  openModalTimer = (data) => {
    this.setState({
      modalTimer: true,
      modal: null
    });
  }

  stake = (req) => {
    const resp = StakeStore.stake(req);
    resp.then((res) => {
      if (res.body.status === i18n.t(`${packageNS}:menu.staking.stake_success`)) {
        this.setState({
          isUnstake: true,
          info: i18n.t(`${packageNS}:menu.messages.congratulations_stake_set`),
          infoStatus: 1,
        });
        setInterval(() => this.displayInfo(), 8000);
      } else {
        this.setState({
          info: res.body.status,
          infoStatus: 2,
        });
        setInterval(() => this.displayInfo(), 8000);
      }
    })
  }

  unstake = (orgId) => {
    const resp = StakeStore.unstake(orgId);
    resp.then((res) => {
      if (res.body.status === i18n.t(`${packageNS}:menu.staking.unstake_success`)) {
        this.setState({
          isUnstake: false,
          amount: 0,
          info: i18n.t(`${packageNS}:menu.messages.unstake_successful`),
          infoStatus: 1,
        });
        setInterval(() => this.displayInfo(), 8000);
      } else {
        this.setState({
          info: res.body.status,
          infoStatus: 2,
        });
        setInterval(() => this.displayInfo(), 8000);
      }
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
    const title = this.state.isUnstake ? i18n.t(`${packageNS}:menu.messages.unstake`) : i18n.t(`${packageNS}:menu.messages.set_stake`);
    let path = null;
    if (this.props.match.params.organizationID === process.env.REACT_APP_SUPER_ADMIN_LPWAN) {
      path = '/control-panel/modify-account/';
    } else {
      path = `/modify-account/${this.props.match.params.organizationID}`;
    }

    return (
      <>
        {this.state.modal && <Modal
          title={i18n.t(`${packageNS}:menu.messages.confirmation`)}
          left={i18n.t(`${packageNS}:menu.staking.cancel`)}
          right={i18n.t(`${packageNS}:menu.staking.confirm`)}
          onProgress={this.handleProgress}
          onCancelProgress={this.handleCancel}
          onClose={this.handleCloseModal}
          context={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)}
          callback={this.openModalTimer} />}

        {this.state.modalTimer && <ModalWithProgress
          title={i18n.t(`${packageNS}:menu.messages.confirmation`)}
          left={i18n.t(`${packageNS}:menu.staking.cancel`)}
          right={i18n.t(`${packageNS}:menu.staking.confirm`)}
          handleProgress={this.handleProgress}
          onClose={this.handleCloseModal}
          context={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)}
          callback={this.onSubmit} />}

        {/* {this.state.modal &&
          <Modal title={i18n.t(`${packageNS}:menu.messages.confirmation`)} description={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)} onProgress={this.handleProgress} onCancelProgress={this.handleCancel} onClose={this.handleCloseModal} open={!!this.state.modal} data={this.state.modal} onSubmit={this.openModalTimer} />} */}

        {/* {this.state.modalTimer &&
          <ModalTimer title={i18n.t(`${packageNS}:menu.messages.stake_proc_tit`)} description={i18n.t(`${packageNS}:menu.messages.stake_proc_desc`)} onProgress={this.handleProgress} onCancelProgress={this.handleProgress} onProcClose={this.handleCloseProcModal} open={!!this.state.modalTimer} data={this.state.modalTimer} onProgress={this.handleProgress} onSubmit={this.onSubmit} />} */}
        <TitleBar>
          <Breadcrumb>
            <BreadcrumbItem active>{title}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Row xs="1">
          <Col>
            <Card>
              <CardBody>
                {this.state.info}
              </CardBody>
              <CardFooter>
                <ExtLink to={EXT_URL_STAKE} context={i18n.t(`${packageNS}:menu.common.learn_more`)} />
              </CardFooter>
            </Card>
          </Col>
        </Row>
        <Row xs="1" lg="2">
          <Col>
            <Card>
              <CardBody>
                <StakeForm isUnstake={this.state.isUnstake} label={this.state.isUnstake ? i18n.t(`${packageNS}:menu.messages.withdraw_stake`) : i18n.t(`${packageNS}:menu.messages.set_stake`)} onChange={this.onChange} amount={this.state.amount} revRate={this.state.revRate} reset={this.reset} confirm={this.confirm} />
              </CardBody>
            </Card>
          </Col>
          <Col><InfoCard path={path} /></Col>
        </Row>
      </>
    );
  }
}

export default withRouter(SetStake);
