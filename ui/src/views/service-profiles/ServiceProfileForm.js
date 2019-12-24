import React, { Component } from "react";

import { Row, Col, Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field, FieldArray } from 'formik';
import * as Yup from 'yup';

import { ReactstrapInput, ReactstrapCheckbox, AsyncAutoComplete } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';

import NetworkServerStore from "../../stores/NetworkServerStore";


class ServiceProfileForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || {},
    };

    this.getNetworkServerOption = this.getNetworkServerOption.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
  }

  

  getNetworkServerOption(id, callbackFunc) {
    NetworkServerStore.get(id, resp => {
      callbackFunc({ label: resp.networkServer.name, value: resp.networkServer.id });
    });
  }

  getNetworkServerOptions(search, callbackFunc) {
    NetworkServerStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
      callbackFunc(options);
    });
  }


  onNetworkSelect = (v) => {
    if (!this.state.object.networkServerID || (this.state.object.networkServerID && this.state.object.networkServerID !== v.id)) {
      let object = this.state.object;
      object.gatewayProfileID = null;
      object.networkServerID = v.value;
      this.setState({
        object: object,
      });
    }
  }

  render() {
    if (this.state.object === undefined) {
        return (<div></div>);
    }

    return(<React.Fragment>
      <Row>
        <Col>
          <Formik
            enableReinitialize
            initialValues={this.state.object}
            onSubmit={this.props.onSubmit}>
            {({
              handleSubmit,
              setFieldValue,
              handleChange,
              handleBlur,
            }) => (
                <Form onSubmit={handleSubmit} noValidate>
                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000149`)+"*"}
                    name="name"
                    id="name"
                    helpText={i18n.t(`${packageNS}:tr000150`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                    required
                  />

                    {!this.props.update && <Field
                        type="text"
                        label={i18n.t(`${packageNS}:tr000047`)+"*"}
                        name="networkServerID"
                        id="networkServerID"
                    // value={this.state.object.networkServerID || ""}
                    // getOption={this.getNetworkServerOption}
                        getOptions={this.getNetworkServerOptions}
                        setFieldValue={setFieldValue}
                        helpText={i18n.t(`${packageNS}:tr000223`)}
                        onBlur={handleBlur}
                        inputProps={{
                            clearable: true,
                            cache: false,
                        }}
                        onChange={this.onNetworkSelect}
                        component={AsyncAutoComplete}
                        required
                    />}

                    <Field
                        type="checkbox"
                        label={i18n.t(`${packageNS}:tr000151`)}
                        name="addGWMetaData"
                        id="addGWMetaData"
                        component={ReactstrapCheckbox}
                        onChange={handleChange}
                        helpText={i18n.t(`${packageNS}:tr000152`)}
                        onBlur={handleBlur}
                    />

                    <Field
                        type="checkbox"
                        label={i18n.t(`${packageNS}:tr000153`)}
                        name="nwkGeoLoc"
                        id="nwkGeoLoc"
                        component={ReactstrapCheckbox}
                        onChange={handleChange}
                        onBlur={handleBlur}
                        helpText={i18n.t(`${packageNS}:tr000154`)}
                    />

                    <Field
                        type="number"
                        label={i18n.t(`${packageNS}:tr000155`)}
                        name="devStatusReqFreq"
                        id="devStatusReqFreq"
                        value={this.state.object.devStatusReqFreq || 0}
                        helpText={i18n.t(`${packageNS}:tr000156`)}
                        component={ReactstrapInput}
                        onBlur={handleBlur}
                    />

                    {this.state.object.devStatusReqFreq > 0 && <FormGroup>
                        <Field
                            type="checkbox"
                            label={i18n.t(`${packageNS}:tr000157`)}
                            name="reportDevStatusBattery"
                            id="reportDevStatusBattery"
                            component={ReactstrapCheckbox}
                            onChange={handleChange}
                        />

                        <Field
                            type="checkbox"
                            label={i18n.t(`${packageNS}:tr000158`)}
                            name="reportDevStatusMargin"
                            id="reportDevStatusMargin"
                            component={ReactstrapCheckbox}
                            onChange={handleChange}
                            />

                    </FormGroup>}

                    <Field
                        type="number"
                        label={i18n.t(`${packageNS}:tr000159`)+"*"}
                        name="drMin"
                        id="drMin"
                        value={this.state.object.drMin || 0}
                        helpText={i18n.t(`${packageNS}:tr000160`)}
                        component={ReactstrapInput}
                        required
                    />

                    <Field
                        type="number"
                        label={i18n.t(`${packageNS}:tr000161`)+"*"}
                        name="drMax"
                        id="drMax"
                        value={this.state.object.drMax || 0}
                        helpText={i18n.t(`${packageNS}:tr000162`)}
                        component={ReactstrapInput}
                        required
                    />

                  <Button type="submit" color="primary">{this.props.submitLabel}</Button>
                </Form>
            )}
          </Formik>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default ServiceProfileForm;
