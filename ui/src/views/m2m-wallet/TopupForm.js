import React from "react";

import TextField from '@material-ui/core/TextField';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import TitleBarButton from "../../components/TitleBarButton";
import LinkI from "mdi-material-ui/Link";


class TopupForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="amount"
          label="Amount"
          helperText="Send MXC amount from."
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          inputProps={{
            pattern: "[\\w-]+",
          }}
          required
          fullWidth
        />
        <TextField
          id="to"
          label="To"
          helperText="Ethereum address to."
          margin="normal"
          value={this.state.object.displayName || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        
        <TitleBarButton
            key={1}
            label="Go to AXS"
            icon={<LinkI />}
            color="secondary"
            /* onClick={this.deleteOrganization} */
        />
          
      </Form>
    );
  }
}

export default TopupForm;
