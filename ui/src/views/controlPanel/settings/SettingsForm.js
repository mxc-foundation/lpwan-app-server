import React from "react";

import TextField from '@material-ui/core/TextField';
import i18n, { packageNS } from '../../../i18n';
import Button from "@material-ui/core/Button";
import FormComponent from "../../../classes/FormComponent";
import Form from "../../../components/Form";

import { withRouter, Link } from "react-router-dom";
import WithdrawStore from '../../../stores/WithdrawStore';
import { ETHER } from '../../../util/CoinType';
import NumberFormat from 'react-number-format';
import SettingsStore from '../../../stores/SettingsStore';

const NumberFormatMXC = (props) => {
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
}

const NumberFormatPerc = (props) => {
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
      suffix=" %"
    />
  );
}

class SettingsForm extends FormComponent {

  constructor(props) {
    super(props);

    this.state = {};
  }

  componentDidMount() {
    this.loadSettings();
  }

  loadSettings = async () => {
    try {
      const organizationID = 0;
      //this.setState({loading: true})

      WithdrawStore.getWithdrawFee(ETHER, organizationID, (resp) => {
        this.setState({ withdrawFee: resp.withdrawFee });
      });

      SettingsStore.getSystemSettings((resp) => {
        this.setState({
          downlinkPrice: resp.downlinkFee,
          percentageShare: resp.transactionPercentageShare,
          lbWarning: resp.lowBalanceWarning
        });
      });
    } catch (e) { 
      console.log("Error", e)
    }
  };

  saveSettings = async () => {
    try {
      let bodyWF = {
        moneyAbbr: 'Ether',
        orgId: '0',
        withdrawFee: this.state.withdrawFee
      };

      let bodySettings = {
        downlinkFee: this.state.downlinkPrice,
        lowBalanceWarning: this.state.lbWarning,
        transactionPercentageShare: this.state.percentageShare
      };

      WithdrawStore.setWithdrawFee(ETHER, 0, bodyWF, (resp) => { });

      SettingsStore.setSystemSettings(bodySettings, (resp) => { });
    } catch (e) { 
      console.log("Error", e)
    }
  };

  reset = () => {
    this.loadSettings();
  }

  handleChange = (name, event) => {
    this.setState({
      [name]: event.target.value
    });
  };

  render() {
    const extraButtons = <>
      <Button variant="outlined" color="inherit" onClick={this.reset} type="button" disabled={false}>{i18n.t(`${packageNS}:menu.staking.reset`)}</Button>
    </>;

    return (
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
        onSubmit={this.submit}
      >
        <TextField
          id="withdrawFee"
          label={i18n.t(`${packageNS}:menu.settings.withdraw_fee`)}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          InputProps={{
            inputComponent: NumberFormatMXC
          }}
          margin="normal"
          value={this.state.withdrawFee}
          onChange={(e) => this.handleChange('withdrawFee', e)}
          fullWidth
        />
        <TextField
          id="downlinkPrice"
          label={i18n.t(`${packageNS}:menu.settings.downlink_price`)}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          InputProps={{
            inputComponent: NumberFormatMXC
          }}
          margin="normal"
          value={this.state.downlinkPrice}
          onChange={(e) => this.handleChange('downlinkPrice', e)}
          fullWidth
        />
        <TextField
          id="percentageShare"
          label={i18n.t(`${packageNS}:menu.settings.percentage_share`)}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          InputProps={{
            inputComponent: NumberFormatPerc
          }}
          margin="normal"
          value={this.state.percentageShare}
          onChange={(e) => this.handleChange('percentageShare', e)}
          fullWidth
        />
        <TextField
          id="lbWarning"
          label={i18n.t(`${packageNS}:menu.settings.low_balance`)}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          InputProps={{
            inputComponent: NumberFormatMXC
          }}
          margin="normal"
          value={this.state.lbWarning}
          onChange={(e) => this.handleChange('lbWarning', e)}
          fullWidth
        />


      </Form>
    );
  }
}

export default SettingsForm;
