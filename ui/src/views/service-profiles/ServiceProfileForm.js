import React, { Component } from "react";

import { withStyles } from '@material-ui/core/styles';
import FormControlOrig from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";


import { Button, Form, FormGroup, Label, Input, FormText, Row, Col } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";
import FormSubmit from "../../components/Form";
import FormControl from "../../components/FormControl";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import NetworkServerStore from "../../stores/NetworkServerStore";

import theme from "../../theme";


const styles = {
    a: {
        color: theme.palette.primary.main,
    },
    formLabel: {
        fontSize: 12,
    },
};


class ServiceProfileForm extends FormComponent {
  constructor() {
    super();
    this.getNetworkServerOption = this.getNetworkServerOption.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
  }

    getNetworkServerOption(id, callbackFunc) {
        NetworkServerStore.get(id, resp => {
            callbackFunc({ label: resp.name, value: resp.id });
        });
    }

    getNetworkServerOptions() {
        NetworkServerStore.list(0, 999, 0, resp => {
            const options = resp.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
            let object = this.state.object;
            object.options = options;

            this.setState({
                object
            })
        });
    }

    render() {
        if (this.state.object === undefined) {
            return (<div></div>);
        }

    return(
    <React.Fragment>
      <Form>
          <FormGroup row>
              <Label for="name" sm={2}>{i18n.t(`${packageNS}:tr000149`)}</Label>
              <Col sm={10}>
                  <Input type="text" name="name" id="name" value={this.state.object.name || ""} onChange={this.onChange} />
                  <FormText color="muted">{i18n.t(`${packageNS}:tr000150`)}</FormText>
              </Col>
          </FormGroup>

          {!this.props.update && <FormGroup row>
              <Label for="networkServerID" sm={2}>{i18n.t(`${packageNS}:tr000047`)}</Label>
              <Col sm={10}>
                  <Input type="select" name="networkServerID" id="networkServerID" value={this.state.object.networkServerID || ""} onChange={this.onChange}>
                      <option value={''}>{i18n.t(`${packageNS}:tr000171`)}</option>
                      {this.state.object.options && this.state.object.options.map(project => {
                          return (
                              <option value={project.value}>{project.label}</option>
                          )
                      })}
                  </Input>
              </Col>
          </FormGroup>}


         {/* <FormGroup row>
              <Label for="addGWMetaData" sm={2}>{i18n.t(`${packageNS}:tr000149`)}</Label>
              <Col sm={10}>
                  <Input type="text" name="name" id="name" value={this.state.object.name || ""} onChange={this.onChange} />
                  <FormText color="muted">{i18n.t(`${packageNS}:tr000150`)}</FormText>
              </Col>
          </FormGroup>

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
        />*/}
      </Form>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ServiceProfileForm);
