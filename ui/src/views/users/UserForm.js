import React from "react";

import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from "@material-ui/core/FormGroup";
import Checkbox from '@material-ui/core/Checkbox';

import FormComponent from "../../classes/FormComponent";
import FormControl from "../../components/FormControl";
import Form from "../../components/Form";
import i18n, { packageNS } from '../../i18n';


class UserForm extends FormComponent {
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
          id="username"
          label={i18n.t(`${packageNS}:tr000056`)}
          margin="normal"
          value={this.state.object.username || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="email"
          label={i18n.t(`${packageNS}:tr000147`)}
          margin="normal"
          value={this.state.object.email || ""}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="note"
          label={i18n.t(`${packageNS}:tr000129`)}
          helperText={i18n.t(`${packageNS}:tr000130`)}
          margin="normal"
          value={this.state.object.note || ""}
          onChange={this.onChange}
          rows={4}
          fullWidth
          multiline
        />
        {this.state.object.id === undefined && <TextField
          id="password"
          label={i18n.t(`${packageNS}:tr000004`)}
          type="password"
          margin="normal"
          value={this.state.object.password || ""}
          onChange={this.onChange}
          required
          fullWidth
        />}
        <FormControl label={i18n.t(`${packageNS}:tr000131`)}>
          <FormGroup>
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000132`)}
              control={
                <Checkbox
                  id="isActive"
                  checked={!!this.state.object.isActive}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000133`)}
              control={
                <Checkbox
                  id="isAdmin"
                  checked={!!this.state.object.isAdmin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
          </FormGroup>
        </FormControl>
      </Form>
    );
  }
}

export default UserForm;
