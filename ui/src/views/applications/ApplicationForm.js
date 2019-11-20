import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";
import Typography from '@material-ui/core/Typography';

import {Controlled as CodeMirror} from "react-codemirror2";
import "codemirror/mode/javascript/javascript";

import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import i18n, { packageNS } from '../../i18n';


const styles = {
  codeMirror: {
    zIndex: 1,
  },
  formLabel: {
    fontSize: 12,
  },
};


class ApplicationForm extends FormComponent {
  constructor() {
    super();
    this.getServiceProfileOption = this.getServiceProfileOption.bind(this);
    this.getServiceProfileOptions = this.getServiceProfileOptions.bind(this);
    this.getPayloadCodecOptions = this.getPayloadCodecOptions.bind(this);
    this.onCodeChange = this.onCodeChange.bind(this);
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

  getPayloadCodecOptions(search, callbackFunc) {
    const payloadCodecOptions = [
      {value: "", label: i18n.t(`${packageNS}:tr000211`)},
      {value: "CAYENNE_LPP", label: i18n.t(`${packageNS}:tr000212`)},
      {value: "CUSTOM_JS", label: i18n.t(`${packageNS}:tr000212`)},
    ];

    callbackFunc(payloadCodecOptions);
  }

  onCodeChange(field, editor, data, newCode) {
    let object = this.state.object;
    object[field] = newCode;
    this.setState({
      object: object,
    });
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    const codeMirrorOptions = {
      lineNumbers: true,
      mode: "javascript",
      theme: "default",
    };
    
    let payloadEncoderScript = this.state.object.payloadEncoderScript;
    let payloadDecoderScript = this.state.object.payloadDecoderScript;

    if (payloadEncoderScript === "" || payloadEncoderScript === undefined) {
      payloadEncoderScript = `// Encode encodes the given object into an array of bytes.
//  - fPort contains the LoRaWAN fPort number
//  - obj is an object, e.g. {"temperature": 22.5}
// The function must return an array of bytes, e.g. [225, 230, 255, 0]
function Encode(fPort, obj) {
  return [];
}`;
    }

    if (payloadDecoderScript === "" || payloadDecoderScript === undefined) {
      payloadDecoderScript = `// Decode decodes an array of bytes into an object.
//  - fPort contains the LoRaWAN fPort number
//  - bytes is an array of bytes, e.g. [225, 230, 255, 0]
// The function must return an object, e.g. {"temperature": 22.5}
function Decode(fPort, bytes) {
  return {};
}`;
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="name"
          label={i18n.t(`${packageNS}:tr000254`)}
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          helperText={i18n.t(`${packageNS}:tr000062`)}
          fullWidth
          required
        />
        <TextField
          id="description"
          label={i18n.t(`${packageNS}:tr000255`)}
          margin="normal"
          value={this.state.object.description || ""}
          onChange={this.onChange}
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
          />
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000257`)}
          </FormHelperText>
        </FormControl>}
        {this.state.object.payloadCodec !== "" && <div>
          <FormControl fullWidth margin="normal">
            <FormLabel className={this.props.classes.formLabel}>{i18n.t(`${packageNS}:tr000209`)}</FormLabel>
            <AutocompleteSelect
              id="payloadCodec"
              label={i18n.t(`${packageNS}:tr000210`)}
              value={this.state.object.payloadCodec || ""}
              onChange={this.onChange}
              getOptions={this.getPayloadCodecOptions}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000214`)}
            </FormHelperText>
          </FormControl>
          {this.state.object.payloadCodec === "CUSTOM_JS" && <FormControl fullWidth margin="normal">
            <CodeMirror
              value={payloadDecoderScript}
              options={codeMirrorOptions}
              onBeforeChange={this.onCodeChange.bind(this, 'payloadDecoderScript')}
              className={this.props.classes.codeMirror}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000215`)}
            </FormHelperText>
          </FormControl>}
          {this.state.object.payloadCodec === "CUSTOM_JS" && <FormControl fullWidth margin="normal">
            <CodeMirror
              value={payloadEncoderScript}
              options={codeMirrorOptions}
              onBeforeChange={this.onCodeChange.bind(this, 'payloadEncoderScript')}
              className={this.props.classes.codeMirror}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000216`)}
            </FormHelperText>
          </FormControl>}
        </div>}
        {this.state.object.payloadCodec === "" && <FormControl fullWidth margin="normal">
          <Typography variant="body1">
            {i18n.t(`${packageNS}:tr000259`)}
          </Typography>
        </FormControl>}
      </Form>
    );
  }
}

export default withStyles(styles)(ApplicationForm);
