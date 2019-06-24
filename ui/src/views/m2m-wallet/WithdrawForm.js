import React from "react";

import TextField from '@material-ui/core/TextField';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import purple from "@material-ui/core/colors/purple";
import green from "@material-ui/core/colors/green";

import { createMuiTheme } from "@material-ui/core/styles";
import { ThemeProvider } from "@material-ui/styles";
import { withRouter } from "react-router-dom";
//import { withStyles } from "@material-ui/core/styles";

const theme = createMuiTheme({
  palette: {
    primary: purple,
    secondary: green
  },
  overrides: {
    MuiInputLabel: { // Name of the component ⚛️ / style sheet
      root: { // Name of the rule
        color: "orange",
        "&$focused": { // increase the specificity for the pseudo class
          color: "purple"
        }
      }
    }
  }
});

class WithdrawForm extends FormComponent {
  
  render() {
    if (this.props.organization === undefined) {
      return(<div>loading...</div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <ThemeProvider theme={theme}>
        <TextField
          id="amount"
          bgcolor="primary.main"
          label="Amount to Withdraw"
          //helperText="The name may only contain words, numbers and dashes."
          margin="normal"
          value={this.props.organization.balance || ""}
          onChange={this.onChange}
          
          required
          fullWidth
        />
        </ThemeProvider>
        <TextField
          id="txFee"
          label="Transaction fee is"
          margin="normal"
          value={this.props.organization.displayName || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="destination"
          label="Withdraw destination"
          //helperText="The name may only contain words, numbers and dashes."
          margin="normal"
          value={this.props.organization.name || ""}
          onChange={this.onChange}
          inputProps={{
            pattern: "[\\w-]+",
          }}
          required
          fullWidth
        />
      </Form>
    );
  }
}

export default (withRouter(WithdrawForm));
