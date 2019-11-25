import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import AESKeyField from "../../components/AESKeyField";
import DevAddrField from "../../components/DevAddrField";
import Form from "../../components/Form";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import theme from "../../theme";


const styles = {
  formLabel: {
    fontSize: 12,
  },
  link: {
    color: theme.palette.primary.main,
  },
};


class MulticastGroupForm extends FormComponent {
  constructor() {
    super();
    this.getServiceProfileOption = this.getServiceProfileOption.bind(this);
    this.getServiceProfileOptions = this.getServiceProfileOptions.bind(this);
  }

  getServiceProfileOption(id, callbackFunc) {
    ServiceProfileStore.get(id, resp => {
      callbackFunc({label: resp.serviceProfile.name, value: resp.serviceProfile.id});
    });
  }

  getServiceProfileOptions(search, callbackFunc) {
    ServiceProfileStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((sp, i) => {return {label: sp.name, value: sp.id}});
      callbackFunc(options);
    });
  }

  getRandomKey(len) {
    let key = "";
    const possible = 'abcdef0123456789';

    for(let i = 0; i < len; i++){
      key += possible.charAt(Math.floor(Math.random() * possible.length));
    }

    return key;
  }

  getRandomMcAddr = (cb) => {
    cb(this.getRandomKey(8));
  }

  getRandomSessionKey = (cb) => {
    cb(this.getRandomKey(32));
  }


  getGroupTypeOptions(search, callbackFunc) {
    const options = [
      {value: "CLASS_B", label: i18n.t(`${packageNS}:tr000194`)},
      {value: "CLASS_C", label: i18n.t(`${packageNS}:tr000203`)},
    ];

    callbackFunc(options);
  }

  getPingSlotPeriodOptions(search, callbackFunc) {
    const pingSlotPeriodOptions = [
      {value: 32 * 1, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '' })},
      {value: 32 * 2, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '2' })},
      {value: 32 * 4, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '4' })},
      {value: 32 * 8, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '8' })},
      {value: 32 * 16, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '16' })},
      {value: 32 * 32, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '32' })},
      {value: 32 * 64, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '64' })},
      {value: 32 * 128, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '128' })},
    ];

    callbackFunc(pingSlotPeriodOptions);
  }

  render() {
    if (this.state.object === undefined) {
      return null;
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="name"
          label={i18n.t(`${packageNS}:tr000261`)}
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000262`)}
          fullWidth
          required
        />
        {!this.props.update && <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000078`)}</FormLabel> 
          <AutocompleteSelect
            id="serviceProfileID"
            label={i18n.t(`${packageNS}:tr000256`)}
            value={this.state.object.serviceProfileID || ""}
            onChange={this.onChange}
            getOption={this.getServiceProfileOption}
            getOptions={this.getServiceProfileOptions}
            margin="none"
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000264`)}
          </FormHelperText>
        </FormControl>}
        <DevAddrField
          id="mcAddr"
          label={i18n.t(`${packageNS}:tr000265`)}
          margin="normal"
          value={this.state.object.mcAddr || ""}
          onChange={this.onChange}
          disabled={this.props.disabled}
          randomFunc={this.getRandomMcAddr}
          fullWidth
          required
          random
        />
        <AESKeyField
          id="mcNwkSKey"
          label={i18n.t(`${packageNS}:tr000266`)}
          margin="normal"
          value={this.state.object.mcNwkSKey || ""}
          onChange={this.onChange}
          disabled={this.props.disabled}
          fullWidth
          required
          random
        />
        <AESKeyField
          id="mcAppSKey"
          label={i18n.t(`${packageNS}:tr000267`)}
          margin="normal"
          value={this.state.object.mcAppSKey || ""}
          onChange={this.onChange}
          disabled={this.props.disabled}
          fullWidth
          required
          random
        />
        <TextField
          id="fCnt"
          label={i18n.t(`${packageNS}:tr000268`)}
          margin="normal"
          type="number"
          value={this.state.object.fCnt || 0}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="dr"
          label={i18n.t(`${packageNS}:tr000269`)}
          helperText={i18n.t(`${packageNS}:tr000270`)}
          margin="normal"
          type="number"
          value={this.state.object.dr || 0}
          onChange={this.onChange}
          required
          fullWidth
        />
        <TextField
          id="frequency"
          label={i18n.t(`${packageNS}:tr000271`)}
          helperText={i18n.t(`${packageNS}:tr000272`)}
          margin="normal"
          type="number"
          value={this.state.object.frequency || 0}
          onChange={this.onChange}
          required
          fullWidth
        />
        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000273`)}</FormLabel>
          <AutocompleteSelect
            id="groupType"
            label={i18n.t(`${packageNS}:tr000274`)}
            value={this.state.object.groupType || ""}
            onChange={this.onChange}
            getOptions={this.getGroupTypeOptions}
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000275`)}
          </FormHelperText>
        </FormControl>
        {this.state.object.groupType === "CLASS_B" && <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000198`)}</FormLabel>
          <AutocompleteSelect
            id="pingSlotPeriod"
            label={i18n.t(`${packageNS}:tr000199`)}
            value={this.state.object.pingSlotPeriod || ""}
            onChange={this.onChange}
            getOptions={this.getPingSlotPeriodOptions}
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000198`)}</FormHelperText>
        </FormControl>}
      </Form>
    );
  }
}

export default withStyles(styles)(MulticastGroupForm);
