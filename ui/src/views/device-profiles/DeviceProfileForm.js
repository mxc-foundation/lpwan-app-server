import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from "@material-ui/core/FormControl";
import FormHelperText from "@material-ui/core/FormHelperText";
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

import {Controlled as CodeMirror} from "react-codemirror2";
import "codemirror/mode/javascript/javascript";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import NetworkServerStore from "../../stores/NetworkServerStore";
import { FormLabel } from "../../../node_modules/@material-ui/core";


const styles = {
  formLabel: {
    fontSize: 12,
  },
  codeMirror: {
    zIndex: 1,
  },
};


class DeviceProfileForm extends FormComponent {
  constructor() {
    super();
    this.state = {
      tab: 0,
    };

    this.onTabChange = this.onTabChange.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
    this.getMACVersionOptions = this.getMACVersionOptions.bind(this);
    this.getRegParamsOptions = this.getRegParamsOptions.bind(this);
    this.getPingSlotPeriodOptions = this.getPingSlotPeriodOptions.bind(this);
    this.getPayloadCodecOptions = this.getPayloadCodecOptions.bind(this);
    this.onCodeChange = this.onCodeChange.bind(this);
  }

  getNetworkServerOptions(search, callbackFunc) {
    NetworkServerStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((ns, i) => {return {label: ns.name, value: ns.id}});
      callbackFunc(options);
    });
  }

  getMACVersionOptions(search, callbackFunc) {
    const macVersionOptions = [
      {value: "1.0.0", label: "1.0.0"},
      {value: "1.0.1", label: "1.0.1"},
      {value: "1.0.2", label: "1.0.2"},
      {value: "1.0.3", label: "1.0.3"},
      {value: "1.1.0", label: "1.1.0"},
    ];

    callbackFunc(macVersionOptions);
  }

  getRegParamsOptions(search, callbackFunc) {
    const regParamsOptions = [
      {value: "A", label: "A"},
      {value: "B", label: "B"},
    ];

    callbackFunc(regParamsOptions);
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

  getPayloadCodecOptions(search, callbackFunc) {
    const payloadCodecOptions = [
      {value: "", label: i18n.t(`${packageNS}:tr000211`)},
      {value: "CAYENNE_LPP", label: i18n.t(`${packageNS}:tr000212`)},
      {value: "CUSTOM_JS", label: i18n.t(`${packageNS}:tr000213`)},
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

  onTabChange(e, v) {
    this.setState({
      tab: v,
    });
  }

  onChange(e) {
    super.onChange(e);
    if (e.target.id === "factoryPresetFreqsStr") {
      let object = this.state.object;
      let freqsStr = e.target.value.split(",");
      object["factoryPresetFreqs"] = freqsStr.map((f, i) => parseInt(f, 10));
      this.setState({
        object: object,
      });
    }
  }

  render() {
    if (this.state.object === undefined) {
      return null;
    }

    let factoryPresetFreqsStr = "";
    if (this.state.object.factoryPresetFreqsStr !== undefined) {
      factoryPresetFreqsStr = this.state.object.factoryPresetFreqsStr;
    } else if (this.state.object.factoryPresetFreqs !== undefined) {
      factoryPresetFreqsStr = this.state.object.factoryPresetFreqs.join(", ");
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
        disabled={this.props.disabled}
      >
        <Tabs value={this.state.tab} onChange={this.onTabChange} indicatorColor="primary">
          <Tab label={i18n.t(`${packageNS}:tr000167`)} />
          <Tab label={i18n.t(`${packageNS}:tr000184`)} />
          <Tab label={i18n.t(`${packageNS}:tr000194`)} />
          <Tab label={i18n.t(`${packageNS}:tr000203`)} />
          <Tab label={i18n.t(`${packageNS}:tr000208`)} />
        </Tabs>

        {this.state.tab === 0 && <div>
          <TextField
            id="name"
            label={i18n.t(`${packageNS}:tr000168`)}
            margin="normal"
            value={this.state.object.name || ""}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000169`)}
            required
            fullWidth
          />
          {!this.props.update && <FormControl fullWidth margin="normal">
            <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000047`)}</FormLabel>
            <AutocompleteSelect
              id="networkServerID"
              label={i18n.t(`${packageNS}:tr000115`)}
              value={this.state.object.networkServerID || ""}
              onChange={this.onChange}
              getOptions={this.getNetworkServerOptions}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000171`)}
            </FormHelperText>
          </FormControl>}
          <FormControl fullWidth margin="normal">
            <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000172`)}</FormLabel>
            <AutocompleteSelect
              id="macVersion"
              label={i18n.t(`${packageNS}:tr000173`)}
              value={this.state.object.macVersion || ""}
              onChange={this.onChange}
              getOptions={this.getMACVersionOptions}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000174`)}
            </FormHelperText>
          </FormControl>
          <FormControl fullWidth margin="normal">
            <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000175`)}</FormLabel>
            <AutocompleteSelect
              id="regParamsRevision"
              label={i18n.t(`${packageNS}:tr000176`)}
              value={this.state.object.regParamsRevision || ""}
              onChange={this.onChange}
              getOptions={this.getRegParamsOptions}
            />
            <FormHelperText>
              {i18n.t(`${packageNS}:tr000177`)}
            </FormHelperText>
          </FormControl>
          <TextField
            id="maxEIRP"
            label={i18n.t(`${packageNS}:tr000178`)}
            type="number"
            margin="normal"
            value={this.state.object.maxEIRP || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000179`)}
            required
            fullWidth
          />
          <TextField
            id="geolocBufferTTL"
            label={i18n.t(`${packageNS}:tr000180`)}
            type="number"
            margin="normal"
            value={this.state.object.geolocBufferTTL || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000181`)}
            fullWidth
          />
          <TextField
            id="geolocMinBufferSize"
            label={i18n.t(`${packageNS}:tr000182`)}
            type="number"
            margin="normal"
            value={this.state.object.geolocMinBufferSize || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000183`)}
            fullWidth
          />
        </div>}

        {this.state.tab === 1 && <div>
          <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000185`)}
              control={
                <Checkbox
                  id="supportsJoin"
                  checked={!!this.state.object.supportsJoin}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
          </FormControl>
          {!this.state.object.supportsJoin && <TextField
            id="rxDelay1"
            label={i18n.t(`${packageNS}:tr000186`)}
            type="number"
            margin="normal"
            value={this.state.object.rxDelay1 || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000187`)}
            required
            fullWidth
          />}
          {!this.state.object.supportsJoin && <TextField
            id="rxDROffset1"
            label={i18n.t(`${packageNS}:tr000188`)}
            type="number"
            margin="normal"
            value={this.state.object.rxDROffset1 || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000189`)}
            required
            fullWidth
          />}
          {!this.state.object.supportsJoin && <TextField
            id="rxDataRate2"
            label={i18n.t(`${packageNS}:tr000190`)}
            type="number"
            margin="normal"
            value={this.state.object.rxDataRate2 || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000189`)}
            required
            fullWidth
          />}
          {!this.state.object.supportsJoin && <TextField
            id="rxFreq2"
            label={i18n.t(`${packageNS}:tr000191`)}
            type="number"
            margin="normal"
            value={this.state.object.rxFreq2 || 0}
            onChange={this.onChange}
            required
            fullWidth
          />}
          {!this.state.object.supportsJoin && <TextField
            id="factoryPresetFreqsStr"
            label={i18n.t(`${packageNS}:tr000192`)}
            margin="normal"
            value={factoryPresetFreqsStr}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000193`)}
            required
            fullWidth
          />}
        </div>}

        {this.state.tab === 2 && <div>
          <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000195`)}
              control={
                <Checkbox
                  id="supportsClassB"
                  checked={!!this.state.object.supportsClassB}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
          </FormControl>

          {this.state.object.supportsClassB && <TextField
            id="classBTimeout"
            label={i18n.t(`${packageNS}:tr000196`)}
            type="number"
            margin="normal"
            value={this.state.object.classBTimeout || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000197`)}
            required
            fullWidth
          />}
          {this.state.object.supportsClassB && <FormControl
              fullWidth
              margin="normal"
            >
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
          {this.state.object.supportsClassB && <TextField
            id="pingSlotDR"
            label={i18n.t(`${packageNS}:tr000201`)}
            type="number"
            margin="normal"
            value={this.state.object.pingSlotDR || 0}
            onChange={this.onChange}
            required
            fullWidth
          />}
          {this.state.object.supportsClassB && <TextField
            id="pingSlotFreq"
            label={i18n.t(`${packageNS}:tr000202`)}
            type="number"
            margin="normal"
            value={this.state.object.pingSlotFreq || 0}
            onChange={this.onChange}
            required
            fullWidth
          />}
        </div>}

        {this.state.tab === 3 && <div>
          <FormControl fullWidth margin="normal">
            <FormControlLabel
              label={i18n.t(`${packageNS}:tr000204`)}
              control={
                <Checkbox
                  id="supportsClassC"
                  checked={!!this.state.object.supportsClassC}
                  onChange={this.onChange}
                  color="primary"
                />
              }
            />
            <FormHelperText>{i18n.t(`${packageNS}:tr000205`)}</FormHelperText>
          </FormControl>

          <TextField
            id="classCTimeout"
            label={i18n.t(`${packageNS}:tr000206`)}
            type="number"
            margin="normal"
            value={this.state.object.classCTimeout || 0}
            onChange={this.onChange}
            helperText={i18n.t(`${packageNS}:tr000207`)}
            required
            fullWidth
          />
        </div>}

        {this.state.tab === 4 && <div>
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
      </Form>
    );
  }
}

export default withStyles(styles)(DeviceProfileForm);
