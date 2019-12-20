import React from "react";

import TextField from '@material-ui/core/TextField';
import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import Divider from '@material-ui/core/Divider';
import Button from "@material-ui/core/Button";
import Typography from '@material-ui/core/Typography';
import InputAdornment from '@material-ui/core/InputAdornment';
import StakeStore from "../../stores/StakeStore";
//import Spinner from "../../components/ScaleLoader"
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";
import {  NumberFormatPerc } from '../../util/M2mUtil';
import NumberFormat from 'react-number-format';
import styles from "./StakeStyle"

const NumberFormatMXC=(props)=> {
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

class StakeForm extends FormComponent {

  state = {
    amount: 0,
    revenue_rate: 0
  }

  componentDidMount() {
    if(this.props.amount > 0){
      this.state.amount = this.props.amount;
    }
    this.loadData();
  }

  loadData = () => {
    const resp = StakeStore.getStakingPercentage(this.props.match.params.organizationID);
    resp.then((res) => {
      let revenue_rate = 0;
      revenue_rate = res.stakingPercentage;
      if (revenue_rate) {
        this.setState({
          revenue_rate
        })
      }
    })
  }
  
  onChange = (event) => {
    const { id, value } = event.target;

    this.setState({
      [id]: value
    });
  }

  handleChange = (name, event) => {
    this.props.onChange(event, name)
		/* this.setState({
			[name]: event.target.value
		}); */
  };
  
  reset = () => {
    this.props.reset();
  }
  
  render() {
    /* if (this.props.txinfo === undefined) {
      return(<Spinner on={this.state.loading}/>);
    } */
    const extraButtons = <>
      <Button variant="outlined" color="inherit" onClick={this.reset} type="button" disabled={false}>{i18n.t(`${packageNS}:menu.staking.reset`)}</Button>
    </>;

    return (
      <Form
        submitLabel={this.props.isUnstake ? i18n.t(`${packageNS}:menu.messages.confirm_unstake`) : i18n.t(`${packageNS}:menu.messages.confirm_stake`)}
        extraButtons={extraButtons}
        onSubmit={(e) => this.props.confirm(e, {
          action: this.props.isUnstake
        })}
      >
        <Typography  /* className={this.props.classes.title} */ gutterBottom>
          {this.props.label}
        </Typography>
        <Divider light={true} />
        <TextField
          id="amount"
          label={i18n.t(`${packageNS}:menu.common.amount`)}
          margin="normal"
          required={!this.props.isUnstake}
          variant="filled"
          
          //value={this.props.amount}
          //onChange={this.props.onChange}
          autoComplete='off'
          //value={this.props.amount !== 0 ? this.props.amount :this.state.amount}
          value={this.props.amount}
          onChange={(e) => this.handleChange('amount', e)}

          InputProps={{
            inputComponent: NumberFormatMXC,
            min: 0,
            readOnly: this.props.isUnstake
          }}
          fullWidth
        />
        <TextField
          id="revRate"
          label={i18n.t(`${packageNS}:menu.messages.revenue_rate`)}
          margin="normal"

          value={this.state.revenue_rate}
          InputProps={{
            inputComponent: NumberFormatPerc,
            readOnly: true,
            append: 'Monthly'
          }}
          fullWidth
        />
      </Form>
    );
  }
}

export default withStyles(styles)(withRouter(StakeForm));
