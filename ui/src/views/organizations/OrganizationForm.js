import React from "react";

import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from "@material-ui/core/FormGroup";
import FormHelperText from '@material-ui/core/FormHelperText';
import Checkbox from '@material-ui/core/Checkbox';

import i18n, { packageNS } from '../../i18n';
import FormControl from "../../components/FormControl";
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";



class OrganizationForm extends FormComponent {
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
          id="name"
          label={i18n.t(`${packageNS}:tr000030`)}
          helperText={i18n.t(`${packageNS}:tr000062`)}
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
          id="displayName"
          label={i18n.t(`${packageNS}:tr000126`)}
          margin="normal"
          value={this.state.object.displayName || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <FormControl
          label={i18n.t(`${packageNS}:tr000063`)}
        >
          <FormGroup>
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000064`)}
              control={
                <Checkbox
                  id="canHaveGateways"
                  checked={!!this.state.object.canHaveGateways}
                  onChange={this.onChange}
                  value="true"
                  color="primary"
                />
              }
            />
          </FormGroup>
          <FormHelperText>{i18n.t(`${packageNS}:tr000065`)}</FormHelperText>
        </FormControl>
      </Form>
    );
  }
}

export default OrganizationForm;
