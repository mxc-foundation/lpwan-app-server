import React, { Component } from "react";

import { Row, Col, Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field, FieldArray } from 'formik';
import * as Yup from 'yup';

import { Map, Marker } from 'react-leaflet';

import { ReactstrapInput, ReactstrapCheckbox, AsyncAutoComplete } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';

import NetworkServerStore from "../../stores/NetworkServerStore";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import LocationStore from "../../stores/LocationStore";
import MapTileLayer from "../../components/MapTileLayer";
import EUI64Field from "../../components/FormikEUI64Field";
import AESKeyField from "../../components/FormikAESKeyField";


class GatewayForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mapZoom: 15,
      object: this.props.object || {},
    };

    this.getNetworkServerOption = this.getNetworkServerOption.bind(this);
    this.getNetworkServerOptions = this.getNetworkServerOptions.bind(this);
    this.getGatewayProfileOption = this.getGatewayProfileOption.bind(this);
    this.getGatewayProfileOptions = this.getGatewayProfileOptions.bind(this);
    this.setCurrentPosition = this.setCurrentPosition.bind(this);
    this.updatePosition = this.updatePosition.bind(this);
    this.updateZoom = this.updateZoom.bind(this);

    this.markerRef = React.createRef(null);
  }

  componentDidMount() {
    if (!this.props.update) {
      this.setCurrentPosition();
    }
  }

  setCurrentPosition(e) {
    if (e !== undefined) {
      e.preventDefault();
    }

    LocationStore.getLocation(position => {
      let object = this.state.object;
      object.location = {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
      }
      this.setState({
        object: object,
      });
    });
  }

  updatePosition() {
    const position = this.markerRef.leafletElement.getLatLng();
    let object = this.state.object;
    object.location = {
      latitude: position.lat,
      longitude: position.lng,
    }
    this.setState({
      object: object,
    });
  }

  updateZoom(e) {
    this.setState({
      mapZoom: e.target.getZoom(),
    });
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

  getGatewayProfileOption(id, callbackFunc) {
    GatewayProfileStore.get(id, resp => {
      callbackFunc({ label: resp.gatewayProfile.name, value: resp.gatewayProfile.id });
    });
  }

  getGatewayProfileOptions(search, callbackFunc) {
    if (this.state.object === undefined || this.state.object.networkServerID === undefined) {
      callbackFunc([]);
      return;
    }

    GatewayProfileStore.list(this.state.object.networkServerID, 999, 0, resp => {
      const options = resp.result.map((gp, i) => { return { label: gp.name, value: gp.id } });
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

    const style = {
      height: 400,
      zIndex: 1,
    };

    let position = [];
    if (this.state.object.location.latitude !== undefined && this.state.object.location.longitude !== undefined) {
      position = [this.state.object.location.latitude, this.state.object.location.longitude];
    } else {
      position = [0, 0];
    }

    let fieldsSchema = {
      name: Yup.string().trim().matches(/^[a-zA-Z0-9\-]+$/, i18n.t(`${packageNS}:tr000429`)).required(i18n.t(`${packageNS}:tr000431`)),
      description: Yup.string()
        .required(i18n.t(`${packageNS}:tr000431`)),
      gatewayProfileID: Yup.string(),
      discoveryEnabled: Yup.bool(),
      location: Yup.object().shape({
        altitude: Yup.number().required(i18n.t(`${packageNS}:tr000431`))
      })
    }

    if (!this.props.update) {
      fieldsSchema['id'] = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema['networkServerID'] = Yup.string();
    }
    const formSchema = Yup.object().shape(fieldsSchema);

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
            enableReinitialize
            initialValues={this.state.object}
            validationSchema={formSchema}
            onSubmit={this.props.onSubmit}>
            {({
              handleSubmit,
              setFieldValue,
              values,
              handleBlur,
            }) => (
                <Form onSubmit={handleSubmit} noValidate>
                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000218`)}
                    name="name"
                    id="name"
                    helpText={i18n.t(`${packageNS}:tr000062`)}
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                  />

                  <Field
                    type="textarea"
                    label={i18n.t(`${packageNS}:tr000219`)}
                    name="description"
                    id="description"
                    component={ReactstrapInput}
                    onBlur={handleBlur}
                  />

                  {!this.props.update && <EUI64Field
                    id="id"
                    label={i18n.t(`${packageNS}:tr000074`)}
                    name="id"
                    value={this.state.object.id || ""}
                    onBlur={handleBlur}
                    required
                    random
                  />}

                  {!this.props.update && <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000047`)}
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
                  />}

                  <Field
                    type="text"
                    label={i18n.t(`${packageNS}:tr000224`)}
                    name="gatewayProfileID"
                    id="gatewayProfileID"
                    triggerReload={this.state.object.networkServerID || values.networkServerID || ""}
                    // value={this.state.object.gatewayProfileID || ""}
                    getOption={this.getGatewayProfileOption}
                    getOptions={this.getGatewayProfileOptions}
                    setFieldValue={setFieldValue}
                    helpText={i18n.t(`${packageNS}:tr000227`)}
                    onBlur={handleBlur}
                    inputProps={{
                      clearable: true,
                      cache: false,
                    }}
                    component={AsyncAutoComplete}
                  />

                  <Field
                    type="checkbox"
                    label={i18n.t(`${packageNS}:tr000228`)}
                    name="discoveryEnabled"
                    id="discoveryEnabled"
                    component={ReactstrapCheckbox}
                    onBlur={handleBlur}
                    helpText={i18n.t(`${packageNS}:tr000229`)}
                  />

                  <Field
                    type="number"
                    label={i18n.t(`${packageNS}:tr000230`)}
                    name="location.altitude"
                    id="location-altitude"
                    component={ReactstrapInput}
                    helpText={i18n.t(`${packageNS}:tr000231`)}
                    onBlur={handleBlur}
                  />

                  <FormGroup>
                    <Label>{i18n.t(`${packageNS}:tr000232`)} (<a onClick={this.setCurrentPosition} href="#getlocation">{i18n.t(`${packageNS}:tr000328`)}</a>)</Label>
                    <Map
                      center={position}
                      zoom={this.state.mapZoom}
                      style={style}
                      animate={true}
                      scrollWheelZoom={false}
                      onZoomend={this.updateZoom}
                    >
                      <MapTileLayer />
                      <Marker position={position} draggable={true} onDragend={this.updatePosition} ref={this.markerRef} />
                    </Map>
                    <FormText color="muted">
                      {i18n.t(`${packageNS}:tr000233`)}
                    </FormText>
                  </FormGroup>


                  { /* boards */}

                  <FieldArray
                    name="boards"
                    render={arrayHelpers => (
                      <div>
                        {values.boards && values.boards.length > 0 && values.boards.map((b, index) => (
                          <React.Fragment key={index}>
                            <Row>
                              <Col>
                                <Card className="shadow-none border">
                                  <CardBody>
                                    <h5>{i18n.t(`${packageNS}:tr000400`)} #{index} (<Button color="link" className="p-0" onClick={() => arrayHelpers.remove(index)}>{i18n.t(`${packageNS}:tr000401`)}</Button>)</h5>

                                    <EUI64Field
                                      label={i18n.t(`${packageNS}:tr000236`)}
                                      name={`boards[${index}].fpgaID`}
                                      id={`boards-${index}-fpgaID`}
                                      value={b.fpgaID || ""}
                                      helpText={i18n.t(`${packageNS}:tr000237`)}
                                    />

                                    <AESKeyField
                                      name={`boards[${index}].fineTimestampKey`}
                                      id={`boards-${index}-fineTimestampKey`}
                                      label={i18n.t(`${packageNS}:tr000238`)}
                                      value={b.fineTimestampKey || ""}
                                      helpText={i18n.t(`${packageNS}:tr000239`)}
                                    />
                                  </CardBody>
                                </Card>
                              </Col>
                            </Row>
                          </React.Fragment>
                        ))}

                        <Button type="button" color="primary" outline className="mb-2" 
                          onClick={() => {arrayHelpers.push({});}}>{i18n.t(`${packageNS}:tr000234`)}</Button>
                      </div>)}
                    ></FieldArray>
              
                  <Button type="submit" color="primary">{this.props.submitLabel || i18n.t(`${packageNS}:tr000066`)}</Button>
                </Form>
              )}
          </Formik>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default GatewayForm;
