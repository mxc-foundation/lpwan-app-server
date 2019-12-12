import React from "react";

import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from "@material-ui/core/FormGroup";
import FormHelperText from '@material-ui/core/FormHelperText';
import Checkbox from '@material-ui/core/Checkbox';
import { Button, Label, Input, FormText } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import FormControl from "../../components/FormControl";
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";



class OrganizationForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }
    return (
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <div class="form-group row">
          <label class="col-sm-2  col-form-label" for="example-helping">{i18n.t(`${packageNS}:tr000030`)}</label>
          <div class="col-sm-10">
            <input type="text" id="name" class="form-control" placeholder="Helping text" value={this.state.object.name || ""} onChange={this.onChange} />
            <FormText color="muted">{i18n.t(`${packageNS}:tr000062`)}</FormText>
          </div>
        </div>
        {/* pattern: "[\\w-]+", */}
        <div class="form-group row">
          <label class="col-sm-2  col-form-label" for="example-helping">{i18n.t(`${packageNS}:tr000126`)}</label>
          <div class="col-sm-10">
            <input type="text" id="displayName" class="form-control" placeholder="Helping text" value={this.state.object.displayName || ""} onChange={this.onChange} />
          </div>
        </div>

        <FormControl
          label={i18n.t(`${packageNS}:tr000063`)}
        >
          <div class="mt-3">
            <div class="custom-control custom-checkbox">
              <input type="checkbox" class="custom-control-input" id="canHaveGateways" checked={!!this.state.object.canHaveGateways} onChange={this.onChange} value="true" />
              <label class="custom-control-label" for="canHaveGateways">{i18n.t(`${packageNS}:tr000064`)}</label>
            </div>
          </div>
          <FormHelperText>{i18n.t(`${packageNS}:tr000065`)}</FormHelperText>
        </FormControl>
      </Form>
    );
  }
}

export default OrganizationForm;
