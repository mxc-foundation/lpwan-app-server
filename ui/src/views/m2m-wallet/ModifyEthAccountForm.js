import React from "react";

import TextField from '@material-ui/core/TextField';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";

class ModifyEthAccountForm extends FormComponent {
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
          id="newaccount"
          label="Your current ETH account is"
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          inputProps={{
            pattern: "[\\w-]+",
          }}
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
