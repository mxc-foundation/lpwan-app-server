import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormLabel from "@material-ui/core/FormLabel";
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from "@material-ui/core/FormGroup";
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from "@material-ui/core/FormControl";
import FormHelperText from "@material-ui/core/FormHelperText";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import NetworkServerStore from "../../stores/NetworkServerStore";


const styles = {
  fontSize: 12,
};


class ServiceProfileForm extends FormComponent {
  constructor() {
    super();
    this.getNetworkServerOption = this.getNetworkServerOption.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
  }

  getNetworkServerOption(id, callbackFunc) {
    NetworkServerStore.get(id, resp => {
      callbackFunc({label: resp.networkServer.name, value: resp.networkServer.id});
    });
  }

  getNetworkServerOptions(search, callbackFunc) {
    NetworkServerStore.list(0, 999, 0, resp => {
      const options = resp.result.map((ns, i) => {return {label: ns.name, value: ns.id}});
      callbackFunc(options);
    });
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
        disabled={this.props.disabled}
      >
        <TextField
          id="name"
          label={i18n.t(`${packageNS}:tr000149`)}
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000150`)}
          required
          fullWidth
        />
        {!this.props.update && <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.FormLabel} required>{i18n.t(`${packageNS}:tr000047`)}</FormLabel>
          <AutocompleteSelect
            id="networkServerID"
            label={i18n.t(`${packageNS}:tr000047`)}
            value={this.state.object.networkServerID || null}
            onChange={this.onChange}
            getOption={this.getNetworkServerOption}
            getOptions={this.getNetworkServerOptions}
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000171`)}
          </FormHelperText>
        </FormControl>}
        <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000151`)}
            control={
              <Checkbox
                id="addGWMetaData"
                checked={!!this.state.object.addGWMetaData}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000152`)}
          </FormHelperText>
        </FormControl>
        <FormControl fullWidth margin="normal">
          <FormControlLabel
            label={i18n.t(`${packageNS}:tr000153`)}
            control={
              <Checkbox
                id="nwkGeoLoc"
                checked={!!this.state.object.nwkGeoLoc}
                onChange={this.onChange}
                color="primary"
              />
            }
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000154`)}
          </FormHelperText>
        </FormControl>
        <TextField
          id="devStatusReqFreq"
          label={i18n.t(`${packageNS}:tr000155`)}
          margin="normal"
          type="number"
          value={this.state.object.devStatusReqFreq || 0}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000156`)}
          fullWidth
        />
        {this.state.object.devStatusReqFreq > 0 && <FormControl fullWidth margin="normal">
          <FormGroup>
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000157`)}
              control={
                <Checkbox
                  id="reportDevStatusBattery"
                  checked={!!this.state.object.reportDevStatusBattery}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000158`)}
              control={
                <Checkbox
                  id="reportDevStatusMargin"
                  checked={!!this.state.object.reportDevStatusMargin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
          </FormGroup>
        </FormControl>}
        <TextField
          id="drMin"
          label={i18n.t(`${packageNS}:tr000159`)}
          margin="normal"
          type="number"
          value={this.state.object.drMin || 0}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000160`)}
          fullWidth
          required
        />
        <TextField
          id="drMax"
          label="Maximum allowed data-rate"
          margin="normal"
          type="number"
          value={this.state.object.drMax || 0}
          onChange={this.onChange}
          helperText="Maximum allowed data rate. Used for ADR."
          fullWidth
          required
        />
      </Form>
    );
  }
}

export default withStyles(styles)(ServiceProfileForm);
