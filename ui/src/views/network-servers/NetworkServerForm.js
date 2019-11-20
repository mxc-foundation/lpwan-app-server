import React from "react";

import TextField from '@material-ui/core/TextField';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from "@material-ui/core/FormGroup";
import FormHelperText from '@material-ui/core/FormHelperText';
import Checkbox from '@material-ui/core/Checkbox';

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import FormControl from "../../components/FormControl";


class NetworkServerForm extends FormComponent {
  constructor() {
    super();
    this.state = {
      tab: 0,
    };

    this.onChangeTab = this.onChangeTab.bind(this);
  }

  onChangeTab(event, value) {
    this.setState({
      tab: value,
    });
  }

  render() {
    if (this.state.object === undefined) {
      return(null);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
            <Tabs
              value={this.state.tab}
              indicatorColor="primary"
              textColor="primary"
              onChange={this.onChangeTab}
            >
              <Tab label={i18n.t(`${packageNS}:tr000167`)} />
              <Tab label={i18n.t(`${packageNS}:tr000095`)} />
              <Tab label={i18n.t(`${packageNS}:tr000104`)} />
            </Tabs>
          {this.state.tab === 0 && <div>
            <TextField
              id="name"
              label={i18n.t(`${packageNS}:tr000090`)}
              fullWidth={true}
              margin="normal"
              value={this.state.object.name || ""}
              onChange={this.onChange}
              helperText={i18n.t(`${packageNS}:tr000091`)}
              required={true}
            />
            <TextField
              id="server"
              label={i18n.t(`${packageNS}:tr000092`)}
              fullWidth={true}
              margin="normal"
              value={this.state.object.server || ""}
              onChange={this.onChange}
              helperText={i18n.t(`${packageNS}:tr000093`)}
              required={true}
            />
          </div>}
          {this.state.tab === 1 && <div>
            <FormControl
              label={i18n.t(`${packageNS}:tr000095`)}
            >
              <FormGroup>
                <FormControlLabel
                  control={
                    <Checkbox
                      id="gatewayDiscoveryEnabled"
                      checked={!!this.state.object.gatewayDiscoveryEnabled}
                      onChange={this.onChange}
                      value="true"
                      color="primary"
                    />
                  }
                  label={i18n.t(`${packageNS}:tr000096`)}
                />
              </FormGroup>
              <FormHelperText>{i18n.t(`${packageNS}:tr000097`)}</FormHelperText>
            </FormControl>
            {this.state.object.gatewayDiscoveryEnabled && <TextField
              id="gatewayDiscoveryInterval"
              label={i18n.t(`${packageNS}:tr000098`)}
              type="number"
              fullWidth={true}
              margin="normal"
              value={this.state.object.gatewayDiscoveryInterval}
              onChange={this.onChange}
              helperText={i18n.t(`${packageNS}:tr000099`)}
              required={true}
            />}
            {this.state.object.gatewayDiscoveryEnabled && <TextField
              id="gatewayDiscoveryTXFrequency"
              label={i18n.t(`${packageNS}:tr000100`)}
              type="number"
              fullWidth={true}
              margin="normal"
              value={this.state.object.gatewayDiscoveryTXFrequency}
              onChange={this.onChange}
              helperText={i18n.t(`${packageNS}:tr000101`)}
              required={true}
            />}
            {this.state.object.gatewayDiscoveryEnabled && <TextField
              id="gatewayDiscoveryDR"
              label={i18n.t(`${packageNS}:tr000102`)}
              type="number"
              fullWidth={true}
              margin="normal"
              value={this.state.object.gatewayDiscoveryDR}
              onChange={this.onChange}
              helperText={i18n.t(`${packageNS}:tr000103`)}
              required={true}
            />}
          </div>}
          {this.state.tab === 2 && <div>
            <FormControl
              label={i18n.t(`${packageNS}:tr000105`)}
            >
              <FormGroup>
                <TextField
                  id="caCert"
                  label={i18n.t(`${packageNS}:tr000106`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.caCert || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000107`)}
                  multiline
                  rows="4"
                />
                <TextField
                  id="tlsCert"
                  label={i18n.t(`${packageNS}:tr000110`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.tlsCert || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000109`)}
                  multiline
                  rows="4"
                />
                <TextField
                  id="tlsKey"
                  label={i18n.t(`${packageNS}:tr000108`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.tlsKey || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000109`)}
                  multiline
                  rows="4"
                />
              </FormGroup>
            </FormControl>

            <FormControl
              label={i18n.t(`${packageNS}:tr000105`)}
            >
              <FormGroup>
                <TextField
                  id="routingProfileCACert"
                  label={i18n.t(`${packageNS}:tr000106`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.routingProfileCACert || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000107`)}
                  multiline
                  rows="4"
                />
                <TextField
                  id="routingProfileTLSCert"
                  label={i18n.t(`${packageNS}:tr000110`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.routingProfileTLSCert || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000107`)}
                  multiline
                  rows="4"
                />
                <TextField
                  id="routingProfileTLSKey"
                  label={i18n.t(`${packageNS}:tr000108`)}
                  fullWidth={true}
                  margin="normal"
                  value={this.state.object.routingProfileTLSKey || ""}
                  onChange={this.onChange}
                  helperText={i18n.t(`${packageNS}:tr000109`)}
                  multiline
                  rows="4"
                />
              </FormGroup>
            </FormControl>
          </div>}
      </Form>
    );
  }
}

export default NetworkServerForm;
