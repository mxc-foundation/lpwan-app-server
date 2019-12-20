import React from "react";

import TextField from '@material-ui/core/TextField';
import i18n, { packageNS } from '../../i18n';
import Button from "@material-ui/core/Button";
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";

class ModifyEthAccountForm extends FormComponent {

  state = {
    newaccount: '',
    username: '',
    password: ''
  }
  
  onChange = (event) => {
    const { id, value } = event.target;
    
    this.setState({
      [id]: value
    });
  }

  clear = () => {
    this.setState({
      username: '',
      password: '',
      newaccount: ''
    })
  }

  submit = () => {
    this.props.onSubmit({
      action: 'modifyAccount',  
      currentAccount: this.state.newaccount,
      createAccount: this.state.newaccount,
      username: this.state.username,
      password: this.state.password
    })

    this.setState({
      username: '',
      password: '',
      newaccount: ''
    })
  }

  render() {
    if (this.props.activeAccount === undefined) {
      return(<div></div>);
    }

    const extraButtons = <>
      <Button  variant="outlined" color="inherit" onClick={this.clear} type="button" disabled={false}>{i18n.t(`${packageNS}:menu.staking.reset`)}</Button>
    </>;

    return(
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
        onSubmit={this.submit}
      >
        <TextField
          id="activeAccount"
          label={i18n.t(`${packageNS}:menu.eth_account.current_account`)}
          margin="normal"
          value={this.props.activeAccount || ""}
/*          variant="filled"*/
          InputLabelProps={{
            shrink: true
          }}
          InputProps={{
            readOnly: true,
          }}
          fullWidth
        />

        <TextField
          id="newaccount"//it is defined current account in swagger
          label={i18n.t(`${packageNS}:menu.eth_account.new_account`)}
          margin="normal"
          value={this.state.newaccount}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          placeholder="0x0000000000000000000000000000000000000000" 
          onChange={this.onChange}
          inputProps={{
            pattern: "^0x[a-fA-F0-9]{40}$",
          }}
            
          autoComplete='off'
          required
          fullWidth
        />

        <TextField
          id="username"//it is defined current account in swagger
          label={i18n.t(`${packageNS}:menu.withdraw.username`)}
          margin="normal"
          value={this.state.username}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          placeholder={i18n.t(`${packageNS}:menu.withdraw.type_here`)}
          onChange={this.onChange}
          autoComplete='off'
          required
          fullWidth
        />

        <TextField
          id="password"//it is defined current account in swagger
          label={i18n.t(`${packageNS}:menu.eth_account.password`)}
          margin="normal"
          value={this.state.password}
          variant="filled"
          InputLabelProps={{
            shrink: true
          }}
          placeholder={i18n.t(`${packageNS}:menu.eth_account.type_here`)}
          onChange={this.onChange}
          type="password"
          autoComplete="off"
          required
          fullWidth
        />
        
        {/* <TitleBarButton
            key={1}
            label="Go to Etherscan.io"
            icon={<LinkI />}
            color="secondary"
            onClick={this.deleteOrganization}
        /> */}
          
      </Form>
    );
  }
}

export default ModifyEthAccountForm;
