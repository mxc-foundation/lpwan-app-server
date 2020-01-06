import React from "react";

import { withStyles } from "@material-ui/core/styles";
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";

import FormComponent from "../../../classes/FormComponent";
import Form from "../../../components/Form";
import AutocompleteSelect from "../../../components/AutocompleteSelect";
import theme from "../../../theme";
import i18n, { packageNS } from '../../../i18n';

import HTTPIntegrationForm from "./subformComponents/HTTPIntegrationForm";
import InfluxDBIntegrationForm from "./subformComponents/InfluxDBIntegrationForm";
import ThingsBoardIntegrationForm from "./subformComponents/ThingsBoardIntegrationForm";

const styles = {
  delete: {
    marginTop: 3 * theme.spacing(1),
  },
  formLabel: {
    fontSize: 12,
  },
};

class IntegrationForm extends FormComponent {
  constructor() {
    super();
    this.getKindOptions = this.getKindOptions.bind(this);
    this.onFormChange = this.onFormChange.bind(this);
  }

  onFormChange(object) {
    this.setState({
      object: object,
    });
  }

  getKindOptions(search, callbackFunc) {
    const kindOptions = [
      {value: "http", label: "HTTP integration"},
      {value: "influxdb", label: "InfluxDB integration"},
      {value: "thingsboard", label: "ThingsBoard.io"},
    ];

    callbackFunc(kindOptions);
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        {!this.props.update && <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.formLabel} required>Integration kind</FormLabel>
          <AutocompleteSelect
            id="kind"
            style={{ width: '100px' }}
            label={i18n.t(`${packageNS}:tr000413`)}
            value={this.state.object.kind || ""}
            onChange={this.onChange}
            getOptions={this.getKindOptions}
          />
        </FormControl>}
        {this.state.object.kind === "http" && <HTTPIntegrationForm object={this.state.object} onChange={this.onFormChange} />}
        {this.state.object.kind === "influxdb" && <InfluxDBIntegrationForm classes={this.props.classes} object={this.state.object} onChange={this.onFormChange} />}
        {this.state.object.kind === "thingsboard" && <ThingsBoardIntegrationForm object={this.state.object} onChange={this.onFormChange} />}
      </Form>
    );
  }
}

export default withStyles(styles)(IntegrationForm);
