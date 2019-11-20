import React from "react";

import { withStyles } from "@material-ui/core/styles";
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";
import Button from "@material-ui/core/Button";

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import DurationField from "../../components/DurationField";
import AutocompleteSelect from "../../components/AutocompleteSelect";

const styles = {
  formLabel: {
    fontSize: 12,
  },
};

class FUOTADeploymentForm extends FormComponent {
  constructor() {
    super();

    this.state.file = null;

    this.onFileChange = this.onFileChange.bind(this);
  }

  getGroupTypeOptions(search, callbackFunc) {
    const options = [
      {value: "CLASS_C", label: i18n.t(`${packageNS}:tr000203`)},
    ];

    callbackFunc(options);
  }

  getMulticastTimeoutOptions(search, callbackFunc) {
    let options = [];

    for (let i = 0; i < (1 << 4); i++) {
      options.push({
        label: `${1 << i} ${i18n.t(`${packageNS}:tr000357`)}`,
        value: i,
      });
    }

    callbackFunc(options);
  }

  onFileChange(e) {
    let object = this.state.object;

    if (e.target.files.length !== 1) {
      object.payload = "";

      this.setState({
        file: null,
        object: object,
      });
    } else {
      this.setState({
        file: e.target.files[0],
      });

      const reader = new FileReader();
      reader.onload = () => {
        const encoded = reader.result.replace(/^data:(.*;base64,)?/, '');
        object.payload = encoded;

        this.setState({
          object: object,
        });
      };
      reader.readAsDataURL(e.target.files[0]);
    }
  }

  render() {
    if (this.state.object === undefined) {
      return null;
    }

    let fileLabel = "";
    if (this.state.file !== null) {
      fileLabel = `${this.state.file.name} (${this.state.file.size} bytes)`
    } else {
      fileLabel = i18n.t(`${packageNS}:tr000370`)
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="name"
          label={i18n.t(`${packageNS}:tr000369`)}
          helperText={i18n.t(`${packageNS}:tr000368`)}
          margin="normal"
          value={this.state.object.name || ""}
          onChange={this.onChange}
          fullWidth
          required
        />

        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000367`)}</FormLabel>
          <Button component="label">
            {fileLabel}
            <input type="file" style={{display: "none"}} onChange={this.onFileChange} />
          </Button>
          <FormHelperText>
            {i18n.t(`${packageNS}:tr000366`)}
          </FormHelperText>
        </FormControl>

        <TextField
          id="redundancy"
          label={i18n.t(`${packageNS}:tr000344`)}
          helperText={i18n.t(`${packageNS}:tr000364`)}
          margin="normal"
          type="number"
          value={this.state.object.redundancy || 0}
          onChange={this.onChange}
          required
          fullWidth
        />

        <DurationField
          id="unicastTimeout"
          label={i18n.t(`${packageNS}:tr000362`)}
          helperText={i18n.t(`${packageNS}:tr000363`)}
          value={this.state.object.unicastTimeout}
          onChange={this.onChange}
        />

        <TextField
          id="dr"
          label="Data-rate"
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

        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>{i18n.t(`${packageNS}:tr000349`)}</FormLabel>
          <AutocompleteSelect
            id="multicastTimeout"
            label={i18n.t(`${packageNS}:tr000361`)}
            value={this.state.object.multicastTimeout || ""}
            onChange={this.onChange}
            getOptions={this.getMulticastTimeoutOptions}
          />
        </FormControl>

      </Form>
    );
  }
}

export default withStyles(styles)(FUOTADeploymentForm);

