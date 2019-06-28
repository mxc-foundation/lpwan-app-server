import React from "react";

import TextField from '@material-ui/core/TextField';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import Button from "@material-ui/core/Button";

import { withRouter } from "react-router-dom";

class WithdrawForm extends FormComponent {
  
  render() {
    if (this.props.organization === undefined) {
      return(<div>loading...</div>);
    }

    const extraButtons = [
      <Button color="primary" type="button" disabled={false}>Cancel</Button>
    ]

    return(
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
        onSubmit={this.onSubmit}
      >
        
        <TextField
          id="amount"
          //bgcolor="primary.main"
          label="Amount"
          //helperText="The name may only contain words, numbers and dashes."
          margin="normal"
          value={this.props.organization.balance || ""}
          onChange={this.onChange}
          className={this.props.classes.root}
          
          required
          fullWidth
        />
        
        <TextField
          id="txFee"
          label="Transaction fee"
          margin="normal"
          value={this.props.organization.displayName || ""}
          onChange={this.onChange}
          className={this.props.classes.root}

          required
          fullWidth
        />
        
        <TextField
          id="destination"
          label="Destination"
          helperText="ETH Account."
          margin="normal"
          value={this.props.organization.name || ""}
          onChange={this.onChange}
          className={this.props.classes.root}
          InputProps={{
            pattern: "[\\w-]+",
          }}
          
          required
          fullWidth
        />
      </Form>
    );
  }
}

export default withRouter(WithdrawForm);
