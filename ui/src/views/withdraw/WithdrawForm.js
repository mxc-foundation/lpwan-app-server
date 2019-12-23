import React from "react";
import { withRouter, Link  } from "react-router-dom";

import TextField from '@material-ui/core/TextField';
import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import TitleBarTitle from "../../components/TitleBarTitle";
//import Button from "@material-ui/core/Button";
import Spinner from "../../components/ScaleLoader"

class WithdrawForm extends FormComponent {
  
  state = {
    amount: ''
  }
 
  onChange = (event) => {
    const { id, value } = event.target;
    
    this.setState({
      [id]: value
    });
  }

  clear() {
    this.setState({
      amount: ''
    })
  }

  render() {
    if (this.props.txinfo === undefined) {
      return(<Spinner on={this.state.loading}/>);
    }
    
    const w_limit = this.props.txinfo.balance - this.props.txinfo.withdrawFee;
    const { txinfo } = this.props;
    
    return(
      <Form
        submitLabel={this.props.submitLabel}
        //extraButtons={extraButtons}
        onSubmit={(e) => this.props.onSubmit(e, {
          amount: parseFloat(this.state.amount),
        })}
      >
        <TextField
          id="amount"
          label={i18n.t(`${packageNS}:menu.withdraw.amount`)}
          margin="normal"
          value={this.state.amount}
          placeholder={i18n.t(`${packageNS}:menu.withdraw.type_here`)}
          onChange={this.onChange}
          autoComplete='off'
          
          required
          fullWidth
          type="number"
          inputProps={{
            min: 0,
            max: w_limit
          }}
        />
        
        <TextField
          id="txFee"
          label={i18n.t(`${packageNS}:menu.withdraw.transaction_fee`)}
          margin="normal"
          
          value={this.props.txinfo.withdrawFee || "0"}
          InputProps={{
            readOnly: true,
          }}
          fullWidth
        />
        
        <TextField
          id="destination"
          label={i18n.t(`${packageNS}:menu.withdraw.to_eth_account`)}
          helperText=""
          margin="normal"
          value={this.props.txinfo.account || ""}
          onChange={this.onChange}
          
          InputProps={{
            readOnly: true,
          }}
          
          fullWidth
        />
        <TitleBarTitle component={Link} to={`/modify-account/${this.props.orgId}`} title={i18n.t(`${packageNS}:menu.withdraw.change_eth_account`)} />
      </Form>
    );
  }
}

export default withRouter(WithdrawForm);
