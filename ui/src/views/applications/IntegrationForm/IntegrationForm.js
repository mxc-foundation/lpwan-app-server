import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import { withStyles } from "@material-ui/core/styles";
import React from "react";
import FormComponent from "../../../classes/FormComponent";
import AutocompleteSelect from "../../../components/AutocompleteSelect";
import Form from "../../../components/Form";
import i18n, { packageNS } from '../../../i18n';
import theme from "../../../theme";
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
    const { classes } = this.props;
    const { object } = this.state;

    if (object === undefined) {
      return(<div></div>);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        {!this.props.update && <FormControl fullWidth margin="normal">
          <FormLabel className={classes.formLabel} required>Integration Kind</FormLabel>
          <AutocompleteSelect
            id="kind"
            style={{ width: '100px' }}
            label={i18n.t(`${packageNS}:tr000413`)}
            value={object.kind || ""}
            onChange={this.onChange}
            getOptions={this.getKindOptions}
          />
        </FormControl>}
        {object.kind === "http" && <HTTPIntegrationForm classes={classes} object={object} onChange={this.onFormChange} />}
        {object.kind === "influxdb" && <InfluxDBIntegrationForm classes={classes} object={object} onChange={this.onFormChange} />}
        {object.kind === "thingsboard" && <ThingsBoardIntegrationForm classes={classes} object={object} onChange={this.onFormChange} />}
      </Form>
    );
  }
}

export default withStyles(styles)(IntegrationForm);
