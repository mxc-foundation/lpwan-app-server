import React, { Component } from "react";

import {
  Row,
  Col,
  Button,
  FormGroup,
  Label,
  FormText,
  Card,
  CardBody,
  CardImg,
  Input
} from "reactstrap";
import { Formik, Form, Field, FieldArray } from "formik";
import * as Yup from "yup";

import { Map, Marker } from "react-leaflet";
import FormHelperText from "@material-ui/core/FormHelperText";

import {
  ReactstrapInput,
  ReactstrapCheckbox,
  AsyncAutoComplete
} from "../../components/FormInputs";
import i18n, { packageNS } from "../../i18n";

import NetworkServerStore from "../../stores/NetworkServerStore";
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import LocationStore from "../../stores/LocationStore";
import MapTileLayer from "../../components/MapTileLayer";
import EUI64Field from "../../components/FormikEUI64Field";
import AESKeyField from "../../components/FormikAESKeyField";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";

import {
  getLBTChannels,
  getChannelsWithFrequency,
  getAntennaGain,
  getLBTConfigStatus
} from "./utils";
import GatewayFormLBT from "./GatewayFormLBT";
import GatewayFormMacChannels from "./GatewayFormMacChannels";
import GatewayFormClassB from "./GatewayFormClassB";

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

  componentDidMount() {
    // Create Gateway
    if (!this.props.update) {
      this.setCurrentPosition();
      return;
      // Update Gateway
    } else {
      this.setKVArrayBoards();
    }

    this.loadGatewayConfig();
    this.loadClassBConfig();
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

  getNetworkServerOptions = (search, callbackFunc) => {
    this.setState({ loading: true });
    NetworkServerStore.list(this.props.match.params.organizationID, 999, 0, resp => {
      const options = resp.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
      this.setState({ loading: false });
      callbackFunc(options);
    });
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
        id: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
        networkServerID: Yup.string()
      }
      // fieldsSchema.object.fields.id = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
      // fieldsSchema.object._nodes.push("id");

      // fieldsSchema.object.fields.networkServerID = Yup.string();
      // fieldsSchema.object._nodes.push("networkServerID");
    }

    return Yup.object().shape(fieldsSchema);
  };

  loadGatewayConfig() {
    // TODO - call actual api here for now working with dummy conf object
    const conf = {
      SX1301_conf: {
        lorawan_public: true,
        clksrc: 1,
        lbt_cfg: {
          enable: true,
          rssi_target: -81,
          chan_cfg: [
            { freq_hz: 868100000, scan_time_us: 5000 },
            { freq_hz: 868300000, scan_time_us: 5000 },
            { freq_hz: 868500000, scan_time_us: 5000 },
            { freq_hz: 868800000, scan_time_us: 5000 },
            { freq_hz: 864700000, scan_time_us: 5000 },
            { freq_hz: 864900000, scan_time_us: 5000 },
            { freq_hz: 865100000, scan_time_us: 5000 },
            { freq_hz: 869525000, scan_time_us: 5000 }
          ],
          sx127x_rssi_offset: -7
        },
        antenna_gain: 2.5,
        radio_0: {
          enable: true,
          type: "SX1257",
          freq: 864900000,
          rssi_offset: -166,
          tx_enable: true,
          tx_notch_freq: 129000,
          tx_freq_min: 863000000,
          tx_freq_max: 870000000
        },
        radio_1: {
          enable: true,
          type: "SX1257",
          freq: 868500000,
          rssi_offset: -166,
          tx_enable: false
        },
        chan_multiSF_0: { enable: true, radio: 1, if: -400000 },
        chan_multiSF_1: { enable: true, radio: 1, if: -200000 },
        chan_multiSF_2: { enable: true, radio: 1, if: 0 },
        chan_multiSF_3: { enable: true, radio: 1, if: 300000 },
        chan_multiSF_4: { enable: true, radio: 0, if: -200000 },
        chan_multiSF_5: { enable: true, radio: 0, if: 0 },
        chan_multiSF_6: { enable: true, radio: 0, if: 200000 },
        chan_multiSF_7: { enable: true, radio: 0, if: 400000 },
        chan_Lora_std: {
          enable: true,
          radio: 1,
          if: -200000,
          bandwidth: 250000,
          spread_factor: 7
        },
        chan_FSK: {
          enable: true,
          radio: 1,
          if: 300000,
          bandwidth: 125000,
          datarate: 50000
        },
        tx_lut_0: { pa_gain: 0, mix_gain: 8, rf_power: -6, dig_gain: 2 },
        tx_lut_1: { pa_gain: 0, mix_gain: 11, rf_power: -3, dig_gain: 3 },
        tx_lut_2: { pa_gain: 0, mix_gain: 11, rf_power: 0, dig_gain: 1 },
        tx_lut_3: { pa_gain: 0, mix_gain: 14, rf_power: 3, dig_gain: 0 },
        tx_lut_4: { pa_gain: 1, mix_gain: 11, rf_power: 6, dig_gain: 3 },
        tx_lut_5: { pa_gain: 1, mix_gain: 11, rf_power: 10, dig_gain: 0 },
        tx_lut_6: { pa_gain: 1, mix_gain: 13, rf_power: 11, dig_gain: 2 },
        tx_lut_7: { pa_gain: 1, mix_gain: 13, rf_power: 12, dig_gain: 1 },
        tx_lut_8: { pa_gain: 1, mix_gain: 14, rf_power: 13, dig_gain: 1 },
        tx_lut_9: { pa_gain: 1, mix_gain: 14, rf_power: 14, dig_gain: 0 },
        tx_lut_10: { pa_gain: 2, mix_gain: 9, rf_power: 16, dig_gain: 0 },
        tx_lut_11: { pa_gain: 2, mix_gain: 12, rf_power: 20, dig_gain: 1 },
        tx_lut_12: { pa_gain: 2, mix_gain: 13, rf_power: 23, dig_gain: 0 },
        tx_lut_13: { pa_gain: 1, mix_gain: 10, rf_power: 25, dig_gain: 1 },
        tx_lut_14: { pa_gain: 3, mix_gain: 12, rf_power: 26, dig_gain: 2 },
        tx_lut_15: { pa_gain: 3, mix_gain: 14, rf_power: 27, dig_gain: 0 }
      },
      gateway_conf: {
        server_address: "192.168.0.7",
        serv_port_up: 1700,
        serv_port_down: 1700,
        keepalive_interval: 10,
        stat_interval: 30,
        push_timeout_ms: 100,
        forward_crc_valid: true,
        forward_crc_error: false,
        forward_crc_disabled: false,
        gps_tty_path: "/dev/ttyS1",
        ref_latitude: 0,
        ref_longitude: 0,
        ref_altitude: 0
      }
    };

    this.setState({
      gatewayConfig: conf,
      gatewayConfigAntenna: getAntennaGain(conf)
    });
  }

  loadClassBConfig() {
    const conf = [
      {
        beacon_period: "0 868.2 9",
        beacon_freq: "125 0 14",
        beacon_datarate: "",
        beacon_bandwidth: "",
        beacon_power: "",
        beacon_info: ""
      }
    ];
    this.setState({ classBConfig: conf });
  }
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

  render() {
    const { object, loading } = this.state;
    let isLoading = loading;

    if (object === undefined) {
      return (<div></div>);
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
                  gatewayProfileID: object.gatewayProfileID || '',
                  networkServerID: object.networkServerID || '',
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
                  console.log('Submitted values: ', values);

                  // Deep copy is required otherwise we can change the original values of
                  // 'boards' (and we will not be able to render the different format in the UI)
                  // Reference: https://medium.com/javascript-in-plain-english/how-to-deep-copy-objects-and-arrays-in-javascript-7c911359b089
                  let newValues = clone(values);
                  console.log('Deep copied submitted values: ', newValues !== values);

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

                  console.log('Prepared values: ', newValues);

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
                      <Col lg={5}>
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
                            getOption={this.getGatewayProfileOption}
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
                      </Col>
                      <Col lg={4}>
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
                              // value={values.object.id}
                              // onBlur={handleBlur}
                              required
                              random
                              // className={
                              //   errors.object && errors.object.id
                              //     ? 'is-invalid form-control'
                              //     : ''
                              // }
                            />
                            {/* {
                            errors.object && errors.object.id
                              ? (
                                <div
                                  className="invalid-feedback"
                                  style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                >
                                  {errors.object.id}
                                </div>
                              ) : null
                          } */}
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
                      <Col lg={3} className="text-right">
                        <Button
                          type="submit"
                          color="primary"
                          className="mr-2"
                          disabled={
                            (errors && Object.keys(errors).length > 0) ||
                            isLoading ||
                            isSubmitting
                          }
                          onClick={() => {
                            validateForm().then(() => {});
                          }}
                        >
                          {this.props.submitLabel ||
                            (this.props.update
                              ? i18n.t(`${packageNS}:tr000614`)
                              : i18n.t(`${packageNS}:tr000277`))}
                        </Button>
                        <Button
                          type="submit"
                          color="secondary"
                          className="d-inline"
                          // onClick={
                          //   () => {
                          //       resetLaraConfig().then(() => {
                          //     })
                          //   }
                          // }
                        >
                          {"Reset Lara Config"}
                        </Button>
                      </Col>
                    </Row>

                    <GatewayFormLBT
                      records={getLBTChannels(this.state.gatewayConfig)}
                      onDataChanged={this.onLBTDataChanged}
                      status={getLBTConfigStatus(this.state.gatewayConfig)}
                      onStatusChanged={this.onLBTStatusChanged}
                    />
                    <GatewayFormMacChannels
                      records={getChannelsWithFrequency(
                        this.state.gatewayConfig
                      )}
                      onDataChanged={this.onLoraMacChannelsChanged}
                    />
                    <Row>
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
                    </Row>

                    <GatewayFormClassB
                      records={this.state.classBConfig}
                      onDataChanged={this.onClassBDataChanged}
                    />

                    <Row>
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
                    </Row>
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
                              onClick={() => { arrayHelpers.push({ fpgaID: '', fineTimestampKey: '' }); }}
                            >
                              {i18n.t(`${packageNS}:tr000234`)}
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
