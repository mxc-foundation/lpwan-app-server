import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
import FormControl from "@material-ui/core/FormControl";
import FormHelperText from "@material-ui/core/FormHelperText";
import FormLabel from "@material-ui/core/FormLabel";
import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import FormComponent from "../../classes/FormComponent";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import Form from "../../components/Form";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';
import DeviceStore from "../../stores/DeviceStore";
import MulticastGroupStore from "../../stores/MulticastGroupStore";




const styles = {
  card: {
    overflow: "visible",
  },
  formLabel: {
    fontSize: 12,
  },
};


class AddDeviceForm extends FormComponent {
  constructor() {
    super();

    this.getDeviceOption = this.getDeviceOption.bind(this);
    this.getDeviceOptions = this.getDeviceOptions.bind(this);
  }

  getDeviceOption(devEUI, callbackFunc) {
    DeviceStore.get(devEUI, resp => {
      callbackFunc({label: `${resp.device.name} (${resp.device.devEUI})`, value: resp.device.devEUI});
    });
  }

  getDeviceOptions(search, callbackFunc) {
    DeviceStore.list({serviceProfileID: this.props.serviceProfileID, search: search, limit: 10}, resp => {
      const options = resp.result.map((d, i) => {return {label: `${d.name} (${d.devEUI})`, value: d.devEUI}});
      callbackFunc(options);
    });
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
        <FormControl fullWidth margin="normal">
          <FormLabel className={this.props.classes.FormLabel} required>{i18n.t(`${packageNS}:tr000403`)}</FormLabel>
          <AutocompleteSelect
            id="devEUI"
            label={i18n.t(`${packageNS}:tr000402`)}
            value={this.state.object.devEUI || ""}
            onChange={this.onChange}
            getOption={this.getDeviceOption}
            getOptions={this.getDeviceOptions}
            margin="none"
          />
          <FormHelperText>{i18n.t(`${packageNS}:tr000404`)}</FormHelperText>
        </FormControl>
      </Form>
    );
  }
}

AddDeviceForm = withStyles(styles)(AddDeviceForm);


class AddDeviceToMulticastGroup extends Component {
  constructor() {
    super();
    this.state = {};
    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    MulticastGroupStore.get(this.props.match.params.multicastGroupID, resp => {
      this.setState({
        multicastGroup: resp.multicastGroup,
      });
    });
  }

  onSubmit(device) {
    MulticastGroupStore.addDevice(this.props.match.params.multicastGroupID, device.devEUI, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/multicast-groups/${this.props.match.params.multicastGroupID}`);
    });
  }

  render() {
    if (this.state.multicastGroup === undefined) {
      return null;
    }

    return(
      <Grid container spacing={4}>
        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000083`)} to={`/organizations/${this.props.match.params.organizationID}/multicast-groups`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={this.state.multicastGroup.name} to={`/organizations/${this.props.match.params.organizationID}/multicast-groups/${this.state.multicastGroup.id}`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000405`)} />
        </TitleBar>

        <Grid item xs={12}>
          <Card className={this.props.classes.card}>
            <CardContent>
              <AddDeviceForm submitLabel={i18n.t(`${packageNS}:tr000406`)} onSubmit={this.onSubmit} serviceProfileID={this.state.multicastGroup.serviceProfileID} />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(AddDeviceToMulticastGroup))
