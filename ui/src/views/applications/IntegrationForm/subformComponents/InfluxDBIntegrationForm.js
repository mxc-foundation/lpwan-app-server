import React from "react";

import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";

import FormComponent from "../../../../classes/FormComponent";
import AutocompleteSelect from "../../../../components/AutocompleteSelect";
import i18n, { packageNS } from '../../../../i18n';

class InfluxDBIntegrationForm extends FormComponent {
  onChange(e) {
    super.onChange(e);
    this.props.onChange(this.state.object);
  }

  getPrecisionOptions(search, callbackFunc) {
    const precisionOptions = [
      {value: "NS", label: "Nanosecond"},
      {value: "U", label: "Microsecond"},
      {value: "MS", label: "Millisecond"},
      {value: "S", label: "Second"},
      {value: "M", label: "Minute"},
      {value: "H", label: "Hour"},
    ];

    callbackFunc(precisionOptions);
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <FormControl fullWidth margin="normal">
        <br />
        <h4>InfluxDB Integration Configuration</h4>
        <TextField
          id="endpoint"
          label="API endpoint (write)"
          placeholder="http://localhost:8086/write"
          value={this.state.object.endpoint || ""}
          onChange={this.onChange}
          margin="normal"
          required
          fullWidth
        />
        <TextField
          id="username"
          label={i18n.t(`${packageNS}:tr000056`)}
          value={this.state.object.username || ""}
          onChange={this.onChange}
          margin="normal"
          fullWidth
        />
        <TextField
          id="password"
          label={i18n.t(`${packageNS}:tr000004`)}
          value={this.state.object.password || ""}
          type="password"
          onChange={this.onChange}
          margin="normal"
          fullWidth
        />
        <TextField
          id="db"
          label="Database name"
          value={this.state.object.db || ""}
          onChange={this.onChange}
          margin="normal"
          fullWidth
          required
        />
        <TextField
          id="retentionPolicyName"
          label="Retention policy name"
          helperText="Sets the target retention policy for the write. InfluxDB writes to the DEFAULT retention policy if you do not specify a retention policy."
          value={this.state.object.retentionPolicyName || ""}
          onChange={this.onChange}
          margin="normal"
          fullWidth
        />
        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>Timestamp precision</FormLabel>
          <AutocompleteSelect
            id="precision"
            label={i18n.t(`${packageNS}:tr000414`)}
            value={this.state.object.precision || ""}
            onChange={this.onChange}
            getOptions={this.getPrecisionOptions}
          />
          <FormHelperText>
            It is recommented to use the least precise precision possible as this can result in significant improvements in compression.
          </FormHelperText>
        </FormControl>
      </FormControl>
    );
  }
}

export default InfluxDBIntegrationForm;
