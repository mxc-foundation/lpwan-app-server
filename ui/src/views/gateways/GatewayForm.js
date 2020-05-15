import { Field, FieldArray, Form, Formik } from "formik";
import React, { Component } from "react";
import { Map, Marker } from "react-leaflet";
import { Button, Card, CardBody, Col, CustomInput, FormGroup, FormText, Input, Label, Row, Alert } from "reactstrap";
import * as Yup from "yup";
import AESKeyField from "../../components/FormikAESKeyField";
import EUI64Field from "../../components/FormikEUI64Field";
import { AsyncAutoComplete, ReactstrapCheckbox, ReactstrapRootPasswordInput, ReactstrapPasswordInput, ReactstrapInput } from "../../components/FormInputs";
import Loader from "../../components/Loader";
import MapTileLayer from "../../components/MapTileLayer";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from "../../i18n";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import LocationStore from "../../stores/LocationStore";
import NetworkServerStore from "../../stores/NetworkServerStore";
import GatewayFormClassB from "./GatewayFormClassB";
import GatewayFormLBT from "./GatewayFormLBT";
import GatewayFormMacChannels from "./GatewayFormMacChannels";
import GatewayStore from "../../stores/GatewayStore";
import { getAntennaGain, getChannelsWithFrequency, getLBTChannels, getLBTConfigStatus } from "./utils";






const clone = require("rfdc")();

class GatewayForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mapZoom: 15,
      object: this.props.object || { location: { altitude: 0 } },
      loading: true,
      gatewayConfig: {},
      classBConfig: {},
      gatewayConfigAntenna: "",
      statistics: "Frames received and frames sent",
      specturalImage: "/img/world-map.png"
    };

    this.markerRef = React.createRef(null);
    this.onLBTDataChanged = this.onLBTDataChanged.bind(this);
    this.onLBTStatusChanged = this.onLBTStatusChanged.bind(this);
    this.onLoraMacChannelsChanged = this.onLoraMacChannelsChanged.bind(this);
    this.onClassBDataChanged = this.onClassBDataChanged.bind(this);
    this.onAntennaValueChange = this.onAntennaValueChange.bind(this);
  }

  componentDidMount = async () => {
    // Create Gateway
    if (!this.props.update) {
      this.setCurrentPosition();
      return;
      // Update Gateway
    } else {
      this.setKVArrayBoards();
    }
    
    const gatewayId = this.props.object.id;
    let name = '';
    if (this.props.object) {
      name = this.props.object.name;
    }
    const sn = name.split("_")[1];
    
    let conf = await GatewayStore.getConfig(gatewayId);
    const rootPassword = await GatewayStore.getRootConfig(gatewayId, sn);
    const object = this.state.object;
    if(rootPassword !== undefined){
      object.password = rootPassword.password;
    }

    var json_conf = JSON.parse(conf.trim());

    let classBConfig = { beacon_period: json_conf.gateway_conf.beacon_period };
    classBConfig.beacon_freq_hz = json_conf.gateway_conf.beacon_freq_hz;
    classBConfig.beacon_datarate = json_conf.gateway_conf.beacon_datarate;
    classBConfig.beacon_bw_hz = json_conf.gateway_conf.beacon_bw_hz;
    classBConfig.beacon_power = json_conf.gateway_conf.beacon_power;
    classBConfig.beacon_infodesc = json_conf.gateway_conf.beacon_infodesc;

    this.setState({
      object,
      gatewayConfig: json_conf,
      gatewayConfigAntenna: getAntennaGain(json_conf),
      classBConfig
    });
  }

  componentDidUpdate(prevProps) {
    if (prevProps.object !== this.props.object) {
      this.setKVArrayBoards();
    }
  }

  // Storage has the 'boards' stored as follows:
  // variables: { my_var_key1: "my var value1", my_var_key2: "my var value2" }
  //
  // But we're leveraging FormikArray, so locally we're converting it into format:
  // variables: [ { fpgaID: "my_key1", fineTimestampKey: "my value1" }, { fpgaID: "my_key2", fineTimestampKey: "my value2" } ]
  convertObjToArray = (obj) => {
    let arr = [];

    for (let [key, value] of Object.entries(obj)) {
      let el = {};
      el.key = key;
      el.value = value;
      arr.push(el);
    }

    return arr;
  }

  // // key: fpgaID, value: fineTimestampKey
  // convertArrayToObj = (arr, key) => {
  //   const formatKey = (k) => k.trim().split(' ').join('_');

  //   let asObject = {};
  //   for (const el of arr.object[key]) {
  //     if (el.fpgaID !== "") {
  //       asObject[formatKey(el.fpgaID)] = el.fineTimestampKey;
  //     }
  //   };

  //   return asObject;
  // }

  setKVArrayBoards = () => {
    if (this.props.object && !Array.isArray(this.props.object.boards)) {
      return;
    }

    if (this.props.object && this.props.object.boards.length === 0) {
      return;
    }

    const propAsArray = this.convertObjToArray(this.props.object.boards);

    this.setState(prevState => {
      if (prevState.object && prevState.object.boards.length === 0) {
        return;
      }

      // Obtain the existing boards that are already in the local state
      let existingStateBoards = prevState.object.boards;
      let existingStateBoardsAsArray = this.convertObjToArray(existingStateBoards);

      // Retrieve the boards array passed as props from the parent component
      let propBoards = propAsArray; //this.props.object.boards;

      // Iterate through the key value pairs
      let updatedBoards = propBoards.map(
        el => {
          let resObj = existingStateBoardsAsArray.find(x => x.fpgaID === el.fpgaID);
          const resIndex = existingStateBoardsAsArray.indexOf(resObj);

          // Assuming that all keys (fpgaID) are unique. If the current key (fpgaID) passed from props
          // is not already in state, then we want to add that new element key value pair to state,
          // otherwise update the value (fineTimestampKey) of that key (fpgaID) if the key exists in state already.
          if (resIndex === -1) {
            return el;
            // Otherwise retain existing state key value pair
          } else {
            resObj.fineTimestampKey = el.fineTimestampKey;
            return resObj;
          }
        }
      )

      return {
        object: {
          ...prevState.object,
          boards: updatedBoards
        }
      }
    })
  }

  setCurrentPosition = (e) => {
    if (e !== undefined) {
      e.preventDefault();
    }
    this.setState({ loading: true });

    LocationStore.getLocation(position => {
      let object = this.state.object;
      object.location = {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
      }
      this.setState({
        object: object,
        loading: false
      });
    });
  }

  updatePosition = () => {
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

  updateZoom = (e) => {
    this.setState({
      mapZoom: e.target.getZoom(),
    });
  }

  getNetworkServerOption = (id, callbackFunc) => {
    NetworkServerStore.get(id, resp => {
      callbackFunc({ label: resp.networkServer.name, value: resp.networkServer.id });
    });
  }

  getNetworkServerOptions = async (search, callbackFunc) => {
    this.setState({ loading: true });

    const res = await NetworkServerStore.list(this.props.match.params.organizationID, 10, 0);
    const options = res.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
    this.setState({ loading: false });
    callbackFunc(options);
  }

  getGatewayProfileOption = (id, callbackFunc) => {
    GatewayProfileStore.get(id, resp => {
      callbackFunc({ label: resp.gatewayProfile.name, value: resp.gatewayProfile.id });
    });
  }

  getGatewayProfileOptions = (search, callbackFunc) => {
    this.setState({ loading: true });
    // Only fetch Gateway Profiles associated with the Network Server that the
    // user must have chosen first.
    if (this.state.object === undefined || this.state.object.networkServerID === undefined) {
      callbackFunc([]);
      return;
    }

    GatewayProfileStore.list(this.state.object.networkServerID, 999, 0, resp => {
      const options = resp.result.map((gp, i) => { return { label: gp.name, value: gp.id } });
      this.setState({ loading: false });
      callbackFunc(options);
    });
  }

  getLocationSourceOptions = (search, callbackFunc) => {
    const options = [
      { value: "UNKNOWN", label: "Unknown" },
      { value: "GPS", label: "GPS" },
      { value: "CONFIG", label: "Manually configured" },
      { value: "GEO_RESOLVER", label: "Geo resolver" },
    ];

    callbackFunc(options);
  }

  onNetworkSelect = (v) => {
    if (!this.state.object.networkServerID || (this.state.object.networkServerID && this.state.object.networkServerID !== v.id)) {
      let object = this.state.object;
      object.gatewayProfileID = null;
      object.networkServerID = v.value;
      this.setState({
        object,
      });
    }
  }

  onGatewayProfileSelect = (v) => {
    if (!this.state.object.gatewayProfileID || (this.state.object.gatewayProfileID && this.state.object.gatewayProfileID !== v.id)) {
      let object = this.state.object;
      object.gatewayProfileID = v.value;
      this.setState({
        object
      });
    }
  }

  onLocationSourceSelect = (v) => {
    const { object } = this.state;
    if (object.location && !object.location.source || (object.location && object.location.source && object.location.source !== v.id)) {
      let object = this.state.object;
      object.location.source = v.value;
      this.setState({
        object
      });
    }
  }

  setValidationErrors = (errors) => {
    this.setState({
      validationErrors: errors
    })
  }

  formikFormSchema = () => {
    let fieldsSchema = {
      // object: Yup.object().shape({
      id: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
      name: Yup.string() //.trim().matches(/[\\w-]+/, i18n.t(`${packageNS}:tr000429`))
        .required(i18n.t(`${packageNS}:tr000431`)),
      description: Yup.string()
        .required(i18n.t(`${packageNS}:tr000431`)),
      // FIXME - for some reason, only 'name', and 'description' are
      // showing as 'required' fields in the UI, but the others aren't
      gatewayProfileID: Yup.string(),
      discoveryEnabled: Yup.bool(),
      location: Yup.object().shape({
        altitude: Yup.number()
          .required(i18n.t(`${packageNS}:tr000431`))
        /* accuracy: Yup.number()
          .required(i18n.t(`${packageNS}:tr000431`)), 
        source: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`))*/
      })
      // })
    }

    if (this.props.update) {
      fieldsSchema = {
        ...fieldsSchema,
        //id: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
        networkServerID: Yup.string(),
        password: Yup.string(),
        server_address: Yup.string(),
        keepalive_interval: Yup.number(),
        stat_interval: Yup.number(),
        push_timeout_ms: Yup.number(),
        serv_port_up: Yup.number(),
        gps_tty_path: Yup.string(),
        serv_port_down: Yup.number(),
        forward_crc_disabled: Yup.bool(),
        forward_crc_error: Yup.bool(),
        forward_crc_valid: Yup.bool(),
      }
      // fieldsSchema.object.fields.id = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
      // fieldsSchema.object._nodes.push("id");

      // fieldsSchema.object.fields.networkServerID = Yup.string();
      // fieldsSchema.object._nodes.push("networkServerID");
    }

    return Yup.object().shape(fieldsSchema);
  };

  /**
   * On lbt data changed
   * @param {*} changedData
   */
  onLBTDataChanged(changedData) {
    let conf = { ...this.state.gatewayConfig };
    let formattedData = [];
    for (const record of changedData) {
      let formattedRecord = { ...record };
      delete formattedRecord["channel"];
      formattedData.push(formattedRecord);
    }
    conf[Object.keys(conf)[0]]["lbt_cfg"]["chan_cfg"] = formattedData;
    this.setState({ gatewayConfig: conf });
  }

  onLBTStatusChanged(status) {
    let conf = { ...this.state.gatewayConfig };
    conf[Object.keys(conf)[0]]["lbt_cfg"]["enable"] = status;

    this.setState({ gatewayConfig: conf });
  }

  onClassBDataChanged(changedData) {
    this.setState({ classBConfig: changedData });
  }

  /**
   * On lora mac channel changed
   * @param {*} changedData
   */
  onLoraMacChannelsChanged(changedData) {

    console.log('changedData', changedData);
    let conf = { ...this.state.gatewayConfig };
    for (const record of changedData) {
      conf[Object.keys(conf)[0]][record.channel]["enable"] = record.enable;
    }
    this.setState({ gatewayConfig: conf });
  }

  /**
   * On antenna value changed
   * @param {*} e
   */
  onAntennaValueChange(e) {
    const antennaVal = e.target.value;
    let conf = { ...this.state.gatewayConfig };
    conf[Object.keys(conf)[0]]["antenna_gain"] = antennaVal;
    this.setState({ gatewayConfig: conf, gatewayConfigAntenna: antennaVal });
  }

  /**
   * On switch toggle
   * @param {*} idx
   * @param {*} e
   */
  onToggle(idx, e) {
    let records = this.state.gatewayConfig;
    
    records.gateway_conf[idx] = e.target.checked;
    this.setState({ records });

    /* if (this.props.onDataChanged) {
      this.props.onDataChanged(records);
    } else {
      this.setState({ records });
    }  */
  }

  render() {
    const { object, loading, gatewayConfig } = this.state;
    let isLoading = loading;

    if (object === undefined) {
      return (<div></div>);
    }

    let gateway_conf = {
      server_address: '',
      keepalive_interval: '',
      stat_interval: '',
      push_timeout_ms: '',
      serv_port_up: '',
      gps_tty_path: '',
      serv_port_down: '',
      forward_crc_disabled: false,
      forward_crc_error: false,
      forward_crc_valid: false,
    };

    if (gatewayConfig !== undefined) {
      if (gatewayConfig.gateway_conf !== undefined) {
        gateway_conf = gatewayConfig.gateway_conf;
      }
    }

    const style = {
      height: 400,
      zIndex: 1,
    };

    let position = [];
    if (object.location && object.location.latitude !== undefined && object.location.longitude !== undefined) {
      position = [object.location.latitude, object.location.longitude];
    } else {
      position = [0, 0];
    }
    // console.log(object.discoveryEnabled);
    const discoveryEnabled = object.discoveryEnabled;

    return (
      <React.Fragment>
        <Row>
          <Col>
            <Formik
              enableReinitialize
              initialValues={
                {
                  id: object.id || undefined,
                  name: object.name || '',
                  description: object.description || '',
                  discoveryEnabled: object.discoveryEnabled || false,
                  location: { altitude: object.location.altitude || 0 },
                  autoUpdate: object.autoUpdate || false,
                  gatewayProfileID: object.gatewayProfileID || '',
                  networkServerID: object.networkServerID || '',
                  server_address: gateway_conf.server_address,
                  keepalive_interval: gateway_conf.keepalive_interval,
                  stat_interval: gateway_conf.stat_interval,
                  push_timeout_ms: gateway_conf.push_timeout_ms,
                  serv_port_up: gateway_conf.serv_port_up,
                  gps_tty_path: gateway_conf.gps_tty_path,
                  serv_port_down: gateway_conf.serv_port_down,
                  forward_crc_disabled: gateway_conf.forward_crc_disabled,
                  forward_crc_error: gateway_conf.forward_crc_error,
                  forward_crc_valid: gateway_conf.forward_crc_valid,
                  password: object.password || '',
                  boards: (
                    (object.boards !== undefined && object.boards.length > 0 && object.boards) || []
                  ),
                }
              }
              validateOnBlur
              validateOnChange
              validationSchema={this.formikFormSchema}
              onSubmit={
                (values, { setSubmitting }) => {
                  const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
                  //console.log('Submitted values: ', values);

                  // Deep copy is required otherwise we can change the original values of
                  // 'boards' (and we will not be able to render the different format in the UI)
                  // Reference: https://medium.com/javascript-in-plain-english/how-to-deep-copy-objects-and-arrays-in-javascript-7c911359b089
                  let newValues = clone(values);
                  //console.log('Deep copied submitted values: ', newValues !== values);

                  // let boardsAsObject;
                  // if (Array.isArray(values.object.boards)) {
                  //   boardsAsObject = this.convertArrayToObj(values, "boards");
                  //   newValues.object.boards = boardsAsObject;
                  // e.g.
                  // newValues.object.boards = [{ fpgaID: "9999999999999999", fineTimestampKey: "99999999999999999999999999999999"}];
                  // }

                  newValues.organizationID = currentOrgID;
                  // delete newValues.object.location.source;
                  // delete newValues.object.location.accuracy;

                  //console.log('Prepared values: ', newValues);

                  this.props.onSubmit(
                    newValues,
                    this.state.gatewayConfig,
                    this.state.classBConfig
                  );
                  setSubmitting(false);
                }}
            >
              {
                props => {
                  const {
                    dirty,
                    errors,
                    handleBlur,
                    handleChange,
                    handleReset,
                    handleSubmit,
                    initialErrors,
                    isSubmitting,
                    isValidating,
                    setFieldValue,
                    touched,
                    validateForm,
                    values
                  } = props;
                  // errors && console.error('validation errors', errors);
                  return (
                    <Form onSubmit={handleSubmit} noValidate>
                      {isLoading && <Loader light />}

                      {!this.props.update &&
                        <>
                          <Field
                            id="networkServerID"
                            name="networkServerID"
                            type="text"
                            value={values.networkServerID}
                            onChange={this.onNetworkSelect}
                            onBlur={handleBlur}
                            label={i18n.t(`${packageNS}:tr000047`)}
                            helpText={i18n.t(`${packageNS}:tr000223`)}
                            // value={values.networkServerID}
                            // getOption={this.getNetworkServerOption}
                            getOptions={this.getNetworkServerOptions}
                            // Hack: we want to trigger Gateway Profild ID list to populate
                            // whenever the Network Server ID changes
                            setFieldValue={() => setFieldValue("gatewayProfileID", "0", false)}
                            inputProps={{
                              clearable: true,
                              cache: false,
                            }}
                            component={AsyncAutoComplete}
                            className={
                              errors && errors.networkServerID
                                ? 'is-invalid form-control'
                                : ''
                            }
                          />
                          {
                            errors && errors.networkServerID
                              ? (
                                <div
                                  className="invalid-feedback"
                                  style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                >
                                  {errors.networkServerID}
                                </div>
                              ) : null
                          }
                        </>
                      }
                      <Row>
                        <Col sm={12} lg={6}>
                          {values.networkServerID &&
                            <>
                              <Field
                                id="gatewayProfileID"
                                name="gatewayProfileID"
                                type="text"
                                value={values.gatewayProfileID}
                                // onChange={handleChange}
                                onChange={this.onGatewayProfileSelect}
                                onBlur={handleBlur}
                                label={i18n.t(`${packageNS}:tr000224`)}
                                helpText={i18n.t(`${packageNS}:tr000227`)}
                                // value={values.gatewayProfileID}
                                //getOption={this.getGatewayProfileOption}
                                getOptions={this.getGatewayProfileOptions}
                                setFieldValue={setFieldValue}
                                inputProps={{
                                  clearable: true,
                                  cache: false,
                                }}
                                component={AsyncAutoComplete}
                                className={
                                  errors && errors.gatewayProfileID
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                              {
                                errors && errors.gatewayProfileID
                                  ? (
                                    <div
                                      className="invalid-feedback"
                                      style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                    >
                                      {errors.gatewayProfileID}
                                    </div>
                                  ) : null
                              }
                            </>
                          }

                          <Field
                            id="name"
                            name="name"
                            type="text"
                            value={values.name}
                            onChange={handleChange}
                            onBlur={handleBlur}
                            label={i18n.t(`${packageNS}:tr000218`)}
                            helpText={i18n.t(`${packageNS}:tr000062`)}
                            component={ReactstrapInput}
                            className={
                              errors && errors.name
                                ? 'is-invalid form-control'
                                : ''
                            }
                          />
                          {
                            errors && errors.name
                              ? (
                                <div
                                  className="invalid-feedback"
                                  style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                >
                                  {errors.name}
                                </div>
                              ) : null
                          }

                          <Field
                            id="location-altitude"
                            name="location.altitude"
                            type="number"
                            value={values.location.altitude}
                            onChange={handleChange}
                            onBlur={handleBlur}
                            label={i18n.t(`${packageNS}:tr000230`)}
                            helpText={i18n.t(`${packageNS}:tr000231`)}
                            component={ReactstrapInput}
                            className={
                              errors &&
                                errors.location &&
                                errors.location.altitude
                                ? "is-invalid form-control"
                                : ""
                            }
                          />
                          {errors &&
                            errors.location &&
                            errors.location.altitude ? (
                              <div
                                className="invalid-feedback"
                                style={{
                                  display: "block",
                                  color: "#ff5b5b",
                                  fontSize: "0.75rem",
                                  marginTop: "-0.75rem"
                                }}
                              >
                                {object.location.altitude}
                              </div>
                            ) : null}
                          <Field
                            type="checkbox"
                            label={'Auto update firmware'}
                            name="autoUpdate"
                            id="autoUpdate"
                            value={values.isAdmin}

                            component={ReactstrapCheckbox}
                            onChange={handleChange}

                            onBlur={handleBlur}
                            helpText={'The firmware will be updated automatically.'}
                          />
                        </Col>
                        <Col sm={12} lg={6}>
                          <Field
                            id="description"
                            name="description"
                            type="textarea"
                            rows={10}
                            value={values.description}
                            onChange={handleChange}
                            onBlur={handleBlur}
                            label={i18n.t(`${packageNS}:tr000219`)}
                            component={ReactstrapInput}
                            className={
                              errors && errors.description
                                ? "is-invalid form-control"
                                : ""
                            }
                          />
                          {errors && errors.description ? (
                            <div
                              className="invalid-feedback"
                              style={{
                                display: "block",
                                color: "#ff5b5b",
                                fontSize: "0.75rem",
                                marginTop: "-0.75rem"
                              }}
                            >
                              {errors.description}
                            </div>
                          ) : null}
                          {!this.props.update && (
                            <>
                              <EUI64Field
                                id="id"
                                name="id"
                                label={i18n.t(`${packageNS}:tr000074`)}
                                value={values.id}
                                onBlur={handleBlur}
                                required
                                random
                                className={
                                  errors && errors.id
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                              {
                                errors && errors.id
                                  ? (
                                    <div
                                      className="invalid-feedback"
                                      style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                    >
                                      {errors.id}
                                    </div>
                                  ) : null
                              }
                            </>
                          )}

                          {/* commented as per new ui */}
                          {/* <Field
                        id="discoveryEnabled"
                        name="discoveryEnabled"
                        type="checkbox"
                        // value={!!values.object.discoveryEnabled}
                        onChange={handleChange}
                        // onBlur={handleBlur}
                        label={i18n.t(`${packageNS}:tr000228`)}
                        helpText={i18n.t(`${packageNS}:tr000229`)}
                        component={ReactstrapCheckbox}
                        className={
                          errors && errors.discoveryEnabled
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.discoveryEnabled
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.discoveryEnabled}
                            </div>
                          ) : null
                      } */}
                        </Col>

                      </Row>
                      {this.props.update &&
                        <>
                          <Row>
                            <Col sm={12} lg={6}>
                              <Field
                                id="server_address"
                                name="server_address"
                                type="text"
                                value={values.server_address}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Server Address'}
                                /* helpText={'server_address'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.server_address
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="keepalive_interval"
                                name="keepalive_interval"
                                type="number"
                                value={values.keepalive_interval}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Keepalive Interval'}
                                /* helpText={'keepalive_interval'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.keepalive_interval
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="stat_interval"
                                name="stat_interval"
                                type="number"
                                value={values.stat_interval}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Stat Interval'}
                                /* helpText={'stat_interval'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.stat_interval
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="push_timeout_ms"
                                name="push_timeout_ms"
                                type="number"
                                value={values.push_timeout_ms}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Push Timeout(ms)'}
                                /* helpText={'push_timeout_ms'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.push_timeout_ms
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="serv_port_up"
                                name="serv_port_up"
                                type="number"
                                value={values.serv_port_up}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Serv Port Up'}
                                /* helpText={'serv_port_up'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.serv_port_up
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="serv_port_down"
                                name="serv_port_down"
                                type="number"
                                value={values.serv_port_down}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'Serv Port Down'}
                                /* helpText={'serv_port_down'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.serv_port_down
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                id="gps_tty_path"
                                name="gps_tty_path"
                                type="text"
                                value={values.gps_tty_path}
                                onChange={handleChange}
                                onBlur={handleBlur}
                                label={'GPS TTY Path'}
                                /* helpText={'gps_tty_path'} */
                                component={ReactstrapInput}
                                className={
                                  errors && errors.gps_tty_path
                                    ? 'is-invalid form-control'
                                    : ''
                                }
                              />
                            </Col>
                            <Col sm={12} lg={6}>
                              <Field
                                style={{ color: 'red' }}
                                helpText={this.state.object.helpText}
                                label={(<span style={{ color: 'red' }}>{i18n.t(`${packageNS}:tr000619`)}</span>)}
                                name="password"
                                id="password"
                                component={ReactstrapRootPasswordInput}
                                onBlur={handleBlur}
                              />
                            </Col>
                          </Row>

                          <Row>
                            <Col sm={12} lg={4}>
                              <CustomInput
                                type="switch"
                                id={`forward_crc_valid`}
                                name="forward_crc_valid"
                                label="forward_crc_valid"
                                checked={values.forward_crc_valid}
                                onChange={e => this.onToggle('forward_crc_valid', e)}
                              />
                            </Col>
                            <Col sm={12} lg={4}>
                              <CustomInput
                                type="switch"
                                id={`forward_crc_error`}
                                name="forward_crc_error"
                                label="forward_crc_error"
                                checked={values.forward_crc_error}
                                onChange={e => this.onToggle('forward_crc_error', e)}
                              />
                            </Col>
                            <Col sm={12} lg={4}>
                              <CustomInput
                                type="switch"
                                id={`forward_crc_disabled`}
                                name="forward_crc_disabled"
                                label="forward_crc_disabled"
                                checked={values.forward_crc_disabled}
                                onChange={e => this.onToggle('forward_crc_disabled', e)}
                              />
                            </Col>
                          </Row>
                        </>
                      }
                      <Row>&nbsp;</Row>
                      {/* <GatewayFormLBT
                        records={getLBTChannels(this.state.gatewayConfig)}
                        onDataChanged={this.onLBTDataChanged}
                        status={getLBTConfigStatus(this.state.gatewayConfig)}
                        onStatusChanged={this.onLBTStatusChanged}
                      /> */}
                      <GatewayFormMacChannels
                        records={getChannelsWithFrequency(
                          this.state.gatewayConfig
                        )}
                        onDataChanged={this.onLoraMacChannelsChanged}
                      />
                      {/* <Row>
                        <Col lg={3} sm={6} xs={12}>
                          <FormGroup>
                            <Label>{i18n.t(`${packageNS}:tr000600`)}</Label>
                            <Input
                              id="Antenna Gain"
                              name="antenna_gain"
                              value={this.state.gatewayConfigAntenna}
                              onChange={this.onAntennaValueChange}
                              type="text"
                              className={
                                errors && errors.antenna_gain
                                  ? "is-invalid form-control"
                                  : ""
                              }
                            ></Input>

                            {errors && errors.antenna_gain ? (
                              <div
                                className="invalid-feedback"
                                style={{
                                  display: "block",
                                  color: "#ff5b5b",
                                  fontSize: "0.75rem",
                                  marginTop: "-0.75rem"
                                }}
                              >
                                {errors.antenna_gain}
                              </div>
                            ) : null}
                          </FormGroup>
                        </Col>
                      </Row> */}

                      <GatewayFormClassB
                        records={this.state.classBConfig}
                        onDataChanged={this.onClassBDataChanged}
                      />

                      {/* <Row>
                        <Col lg={12}>
                          <FormGroup>
                            <Input
                              id="statistics"
                              name="statistics"
                              value={this.state.statistics}
                              onChange={handleChange}
                              type="textarea"
                              rows={3}
                              readOnly={true}
                              className={
                                errors && errors.statistics
                                  ? "is-invalid form-control"
                                  : ""
                              }
                            ></Input>

                            {errors && errors.statistics ? (
                              <div
                                className="invalid-feedback"
                                style={{
                                  display: "block",
                                  color: "#ff5b5b",
                                  fontSize: "0.75rem",
                                  marginTop: "-0.75rem"
                                }}
                              >
                                {errors.statistics}
                              </div>
                            ) : null}
                          </FormGroup>
                        </Col>
                        <Col lg={4}>
                          <FormGroup>
                            <Label>{i18n.t(`${packageNS}:tr000602`)}</Label>
                            <Card>
                              <CardImg
                                top
                                width="100%"
                                src={this.state.specturalImage}
                                alt="Card image cap"
                              />
                            </Card>
                          </FormGroup>
                        </Col>
                      </Row>
                      <Row>
                        <Col lg={5}>
                          <FormGroup>
                            <Button type="button" className="mb-2 mr-2">
                              {i18n.t(`${packageNS}:tr000602`)}
                            </Button>
                            <Button type="button" className="mb-2">
                              {i18n.t(`${packageNS}:tr000603`)}
                            </Button>
                          </FormGroup>
                        </Col>
                      </Row> */}
                      {/* <Field
                        id="location-accuracy"
                        name="object.location.accuracy"
                        type="number"
                        value={values.object.location.accuracy}
                        onChange={handleChange}
                        onBlur={handleBlur}
                        label="Gateway Location Accuracy"
                        helpText="Accuracy (meters)"
                        component={ReactstrapInput}
                        className={
                          errors.object && errors.object.location && errors.object.location.accuracy
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.location && errors.object.location.accuracy
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.object.location.accuracy}
                            </div>
                          ) : null
                      } */}

                      {/* <Field
                        id="locationSource"
                        name="object.location.source"
                        type="text"
                        value={values.object.location.source}
                        // onChange={handleChange}
                        onChange={this.onLocationSourceSelect}
                        onBlur={handleBlur}
                        label="Gateway Location Source"
                        helpText="UNKNOWN: Unknown; GPS: GPS; CONFIG: Manually configured; GEO_RESOLVER: Geo resolver;"
                        // value={values.object.location.source}
                        getOption={this.getLocationSourceOption}
                        getOptions={this.getLocationSourceOptions}
                        setFieldValue={setFieldValue}
                        inputProps={{
                          clearable: true,
                          cache: false,
                        }}
                        component={AsyncAutoComplete}
                        className={
                          errors.object && errors.object.location && errors.object.location.source
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.location && errors.object.location.source
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.object.location.source}
                            </div>
                          ) : null
                      } */}
                      <br />

                      <FormGroup>
                        <TitleBar
                          buttons={
                            <TitleBarButton
                              color="primary"
                              label={i18n.t(`${packageNS}:tr000328`)}
                              icon={<i className="mdi mdi-crosshairs-gps mr-1 align-middle"></i>}
                              to="#getlocation"
                              onClick={() => this.setCurrentPosition}
                            />
                          }
                        >
                          <h5>{i18n.t(`${packageNS}:tr000232`)}</h5>
                        </TitleBar>
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

                      <FieldArray
                        id="boards"
                        name="boards"
                        value={values.boards}
                        render={arrayHelpers => (
                          <div>
                            {
                              values && values.boards &&
                              values.boards.length > 0 &&
                              values.boards.map((board, index) => (
                                board && Object.keys(board).length == 2 ? (
                                  <React.Fragment key={index}>
                                    <Row>
                                      <Col>
                                        <Card className="shadow-none border">
                                          <CardBody>
                                            <TitleBar
                                              buttons={
                                                <TitleBarButton
                                                  color="danger"
                                                  label={i18n.t(`${packageNS}:tr000061`)}
                                                  icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
                                                  onClick={() => arrayHelpers.remove(index)}
                                                />
                                              }
                                            >
                                              <h4>{i18n.t(`${packageNS}:tr000400`)} #{index}</h4>
                                            </TitleBar>

                                            <EUI64Field
                                              id={`boards-${index}-fpgaID`}
                                              name={`boards[${index}].fpgaID`}
                                              label={i18n.t(`${packageNS}:tr000236`)}
                                              value={board.fpgaID}
                                              helpText={i18n.t(`${packageNS}:tr000237`)}
                                            />

                                            <AESKeyField
                                              name={`boards[${index}].fineTimestampKey`}
                                              id={`boards-${index}-fineTimestampKey`}
                                              label={i18n.t(`${packageNS}:tr000238`)}
                                              value={board.fineTimestampKey}
                                              helpText={i18n.t(`${packageNS}:tr000239`)}
                                            />
                                          </CardBody>
                                        </Card>
                                      </Col>
                                    </Row>
                                  </React.Fragment>
                                ) : <div key={index}></div>
                              ))
                            }

                            <Button
                              type="button"
                              variant="outlined"
                              className="mb-2"
                              style={{ width: '275px' }}
                              onClick={() => { arrayHelpers.push({ fpgaID: '', fineTimestampKey: '' }); }}
                            >
                              {i18n.t(`${packageNS}:tr000234`)}
                            </Button>
                            {' '}
                            <Button
                              type="submit"
                              color="secondary"
                              className="mb-2"
                              style={{ width: '275px' }}
                            //className="d-inline"
                            // onClick={
                            //   () => {
                            //       resetLaraConfig().then(() => {
                            //     })
                            //   }
                            // }
                            >
                              {i18n.t(`${packageNS}:menu.gateways.reset_lora_config`)}
                            </Button>
                          </div>
                        )}
                      />

                      <div style={{ margin: "20px 0 10px 20px" }}>
                        {isValidating
                          ? <div style={{ display: "block", color: "orange", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                            Validating. Please wait...
                            </div>
                          : ''
                        }
                        {isSubmitting
                          ? <div style={{ display: "block", color: "orange", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                            Submitting. Please wait...
                            </div>
                          : ''
                        }
                        {/* `initialErrors` does not work for some reason */}
                        {/* {initialErrors.length && JSON.stringify(initialErrors)} */}

                        {/* Show error count when page loads, before user submits the form */}
                        {errors && Object.keys(errors).length
                          ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                            Detected {Object.keys(errors).length} errors. Please fix the validation errors shown in each tab before resubmitting.
                          </div>)
                          : null
                        }

                        {/* Show error count when user submits the form */}
                        {/* {this.state.validationErrors && this.state.validationErrors.length
                          ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                            Detected {Object.keys(this.state.validationErrors).length} errors. Please fix the validation errors shown in each tab before resubmitting.
                          </div>)
                          : null
                        } */}
                      </div>
                      <Row>
                        <Col lg={12} className="text-right">
                          <Button
                            type="submit"
                            color="primary"
                            //className="mr-2"
                            style={{ width: '100%' }}
                            disabled={
                              (errors && Object.keys(errors).length > 0) ||
                              isLoading ||
                              isSubmitting
                            }
                            onClick={() => {
                              validateForm().then(() => { });
                            }}
                          >
                            {this.props.submitLabel ||
                              (this.props.update
                                ? i18n.t(`${packageNS}:tr000614`)
                                : i18n.t(`${packageNS}:tr000277`))}
                          </Button>

                        </Col>
                      </Row>
                    </Form>
                  );
                }}
            </Formik>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayForm;
