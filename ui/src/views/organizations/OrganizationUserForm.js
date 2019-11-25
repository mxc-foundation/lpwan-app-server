import React from "react";

import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Checkbox from '@material-ui/core/Checkbox';

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";


class OrganizationUserForm extends FormComponent {
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
            label={i18n.t(`${packageNS}:tr000056`)}
            margin="normal"
            value={this.state.object.username || ""}
            required
            fullWidth
            InputProps={{
              readOnly: true,
            }}
          />
          <Typography variant="body1">
            {i18n.t(`${packageNS}:tr000138`)}
          </Typography>
          <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000139`)}
              control={
                <Checkbox
                  id="isAdmin"
                  checked={!!this.state.object.isAdmin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormHelperText>{i18n.t(`${packageNS}:tr000140`)}</FormHelperText>
          </FormControl>
          {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000141`)}
              control={
                <Checkbox
                  id="isDeviceAdmin"
                  checked={!!this.state.object.isDeviceAdmin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormHelperText>{i18n.t(`${packageNS}:tr000142`)}</FormHelperText>
          </FormControl>}
          {!!!this.state.object.isAdmin && <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000143`)}
              control={
                <Checkbox
                  id="isGatewayAdmin"
                  checked={!!this.state.object.isGatewayAdmin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormHelperText>{i18n.t(`${packageNS}:tr000144`)}</FormHelperText>
          </FormControl>}
      </Form>
    );
  }
}

export default OrganizationUserForm;
