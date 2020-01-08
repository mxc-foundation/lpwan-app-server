import React, { Component } from "react";

import { Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import Modal from "../../components/Modal";
import ModalWithProgress from "../../components/ModalWithProgress";

import StakeStore from "../../stores/StakeStore";
//import Spinner from "../../components/ScaleLoader"
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

import NumberFormat from 'react-number-format';
import styles from "./StakeStyle"

/* const NumberFormatMXC = (props) => {
  const { inputRef, onChange, ...other } = props;

  return (
    <NumberFormat
      {...other}
      getInputRef={inputRef}
      onValueChange={(values) => {
        onChange({
          target: {
            value: values.value
          }
        });
      }}
      suffix=" MXC"
    />
  );
} */

class StakeForm extends Component {

  constructor(props) {
    super(props);

    this.state = {
      object: {
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
    };
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ amount: nextProps.amount });
  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props.revRate === oldProps.revRate) {
      return;
    }

    this.loadData();
  }

  loadData = async () => {
    let res = await StakeStore.getActiveStakes(this.props.match.params.organizationID);
    let amount = 0;
    let isUnstake = false;

    if (res && res.actStake !== null) {
      amount = res.actStake.Amount;
      if (res.actStake.StakeStatus === 'ACTIVE') {
        isUnstake = true;
      }
    }

    res = await StakeStore.getStakingPercentage(this.props.match.params.organizationID);
    let revRate = 0;
    if (res) {
      revRate = res.stakingPercentage;
    }

    const object = this.state.object;
    object.amount = amount;
    object.revRate = revRate;
    object.isUnstake = isUnstake;

    this.setState({
      object
    });
    this.props.setTitle(this.state.object.isUnstake);
  }

  onChange = (event) => {
    const { id, value } = event.target;

    this.setState({
      object: {
        [id]: value
      }
    });
  }

  reset = () => {
    const object = this.state.object;
    object.amount = 0;

    this.setState({
      object
    });
  }

  onSubmit = (amount) => {
    const orgId = this.props.match.params.organizationID;
    const req = {
      orgId,
      amount: parseFloat(amount)
    }

    if (this.state.object.isUnstake) {
      this.unstake(orgId);
    } else {
      this.stake(req);
    }
    const object = this.state.object;
    object.modal = false
    this.setState({ object });
  }

  confirm = (data) => {
    if (data.amount === 0) {
      return false;
    }

    if (this.state.object.isUnstake) {
      const orgId = this.props.match.params.organizationID;
      this.unstake(orgId);
    } else {
      const { object } = this.state;
      this.setState({
        object: {
          ...object,
          amount: data.amount,
          modal: true
        }
      });
    }
  }

  openModalTimer = () => {
    const object = this.state.object;
    object.modalTimer = true;
    object.modal = null

    this.setState({
      object
    });
  }

  stake = (req) => {
    const resp = StakeStore.stake(req);
    resp.then((res) => {
      const object = this.state.object;

      if (!res) {
        object.info = "Service unavailable. Try again later.";
        object.infoStatus = 3;
        object.infoModal = true;
        this.setState({
          object
        });
        return;
      }

      if (res.body.status === 'Stake successful.') {
        object.isUnstake = true;
        object.info = i18n.t(`${packageNS}:menu.messages.congratulations_stake_set`);
        object.infoStatus = 1;
        object.infoModal = true;
        this.setState({
          object
        });
        this.props.setTitle(this.state.object.isUnstake);

        //setInterval(() => this.displayInfo(), 8000);
      } else {
        object.info = res.body.status;
        object.infoStatus = 2;
        object.infoModal = true;
        this.setState({
          object
        });
      }
    })
  }

  unstake = (orgId) => {
    const resp = StakeStore.unstake(orgId);
    resp.then((res) => {
      const object = this.state.object;

      if (!res) {
        object.info = "Service unavailable. Try again later.";
        object.infoStatus = 3;
        object.infoModal = true;
        this.setState({
          object
        });
        return;
      }

      if (res.body.status === 'Unstake successful.') {
        object.isUnstake = false;
        object.amount = 0;

        this.setState({
          object,
        });
        this.props.setTitle(this.state.object.isUnstake);
      } else {
        object.info = res.body.status;
        object.infoStatus = 2;
        object.infoModal = true;
        this.setState({
          object
        });
      }
    })
  }

  displayInfo = () => {
    const object = this.state.object;
    object.info = i18n.t(`${packageNS}:menu.messages.staking_enhances`);
    object.infoStatus = 0;

    this.setState({
      object
    });
  }

  showModal = (modal) => {
    const object = this.state.object;
    object.modal = modal;
    this.setState({ object });
  }

  handleCloseModal = () => {
    const object = this.state.object;
    object.modal = null;
    this.setState({
      object
    })
  }

  closeInfoModal = () => {
    const object = this.state.object;
    object.infoModal = null;
    this.setState({
      object
    })
  }

  handleOnclick = () => {
    this.props.history.push(`/history/${this.props.match.params.organizationID}/stake`);
  }

  handleProgress = (oldCompleted) => {
    const object = this.state.object;
    object.modalTimer = null
    this.setState({
      object
    })

    if (oldCompleted === 100) {
      this.onSubmit(this.state.object.amount);
    }
  }

  render() {
    let fieldsSchema = {
      amount: Yup.number().moreThan(0).required(),
      revRate: Yup.number(),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <React.Fragment>
        {this.state.object.infoModal && <Modal
          title={i18n.t(`${packageNS}:menu.topup.notice`)}
          left={i18n.t(`${packageNS}:menu.staking.cancel`)}
          right={i18n.t(`${packageNS}:menu.staking.confirm`)}
          context={this.state.object.info}
          callback={this.closeInfoModal}
        />}

        {this.state.object.modal && <Modal
          title={i18n.t(`${packageNS}:menu.messages.confirmation`)}
          left={i18n.t(`${packageNS}:menu.staking.cancel`)}
          right={i18n.t(`${packageNS}:menu.staking.confirm`)}
          onProgress={this.handleProgress}
          onCancelProgress={this.handleCancel}
          onClose={this.handleCloseModal}
          context={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)}
          callback={this.openModalTimer} />}

        {this.state.object.modalTimer && <ModalWithProgress
          title={i18n.t(`${packageNS}:menu.messages.confirmation`)}
          left={i18n.t(`${packageNS}:menu.staking.cancel`)}
          right={i18n.t(`${packageNS}:menu.staking.confirm`)}
          handleProgress={this.handleProgress}
          onClose={this.handleCloseModal}
          context={i18n.t(`${packageNS}:menu.messages.stake_confirmation_text`)}
          callback={this.onSubmit} />}

        <Formik
          enableReinitialize
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={this.confirm}>
          {({
            handleSubmit,
            handleChange,
            setFieldValue,
            values,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="number"
                  label={i18n.t(`${packageNS}:menu.common.amount`)}
                  name="amount"
                  id="amount"
                  value={this.state.object.amount || ""}
                  autoComplete='off'
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  readOnly={this.state.object.isUnstake}
                  min={0}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="number"
                  label={i18n.t(`${packageNS}:menu.messages.revenue_rate`)}
                  name="revRate"
                  id="revRate"
                  value={this.state.object.revRate || ""}
                  component={ReactstrapInput}
                  onChange={handleChange}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Button className="btn-block" onClick={this.reset}>{i18n.t(`${packageNS}:common.reset`)}</Button>
                <Button type="submit" className="btn-block" color="primary">{this.state.object.isUnstake ? i18n.t(`${packageNS}:menu.messages.confirm_unstake`) : i18n.t(`${packageNS}:menu.messages.confirm_stake`)}</Button>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(StakeForm));
