import React, { Component } from "react";

import { Row, Col, Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field, FieldArray } from 'formik';
import * as Yup from 'yup';

import { Map, Marker } from 'react-leaflet';
import FormHelperText from "@material-ui/core/FormHelperText";

import { ReactstrapInput, ReactstrapCheckbox, AsyncAutoComplete } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';

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

const clone = require('rfdc')();

class GatewayForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mapZoom: 15,
      object: this.props.object || {},
      loading: true,
    };

    this.markerRef = React.createRef(null);
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
                    location: {altitude:object.location.altitude || 0},
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

                  this.props.onSubmit(newValues);
                  setSubmitting(false);
                }
              }
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
                        id="description"
                        name="description"
                        type="textarea"
                        value={values.description}
                        onChange={handleChange}
                        onBlur={handleBlur}
                        label={i18n.t(`${packageNS}:tr000219`)}
                        component={ReactstrapInput}
                        className={
                          errors && errors.description
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.description
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.description}
                            </div>
                          ) : null
                      }

                      {!this.props.update &&
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
                      }

                      <Field
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
                      }
                      <br />

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
                          errors && errors.location && errors.location.altitude
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.location && errors.location.altitude
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {object.location.altitude}
                            </div>
                          ) : null
                      }

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
                              onClick={() => {arrayHelpers.push({ fpgaID: '', fineTimestampKey: '' });}}
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
                      <Button
                        type="submit"
                        color="primary"
                        className="btn-block"
                        disabled={(errors && Object.keys(errors).length > 0) || isLoading || isSubmitting}
                        onClick={
                          () => { 
                            validateForm().then(() => {
                            })
                          }
                        }
                      >
                        {this.props.submitLabel || (this.props.update ? "Update" : "Create")}
                      </Button>
                    </Form>
                  );
                }
              }
            </Formik>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayForm;
