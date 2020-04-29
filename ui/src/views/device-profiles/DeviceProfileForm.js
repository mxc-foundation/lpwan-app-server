import Checkbox from '@material-ui/core/Checkbox';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormHelperText from "@material-ui/core/FormHelperText";
import { withStyles } from "@material-ui/core/styles";
import classnames from 'classnames';
import "codemirror/mode/javascript/javascript";
import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { Button, FormGroup, Label, Nav, NavItem, NavLink, TabContent, TabPane } from 'reactstrap';
import * as Yup from 'yup';
import AutocompleteSelect from "../../components/AutocompleteSelect";
import { ReactstrapInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";





const clone = require('rfdc')();

const styles = {
  formLabel: {
    fontSize: 12,
  },
  codeMirror: {
    zIndex: 1,
  },
};


class DeviceProfileForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      object: props.object || {},
      activeTab: "1",
      loading: true,
    };
  }

  getNetworkServerOptions = async (search, callbackFunc) => {
    const res = await NetworkServerStore.list(this.props.match.params.organizationID, 10, 0);
    const options = res.result.map((ns, i) => { return { label: ns.name, value: ns.id } });
    this.setState({
      loading: false
    });
    callbackFunc(options);
  }

  getMACVersionOptions = (search, callbackFunc) => {
    const macVersionOptions = [
      {value: "1.0.0", label: "1.0.0"},
      {value: "1.0.1", label: "1.0.1"},
      {value: "1.0.2", label: "1.0.2"},
      {value: "1.0.3", label: "1.0.3"},
      {value: "1.1.0", label: "1.1.0"},
    ];

    callbackFunc(macVersionOptions);
  }

  getRegParamsOptions = (search, callbackFunc) => {
    const regParamsOptions = [
      {value: "A", label: "A"},
      {value: "B", label: "B"},
    ];

    callbackFunc(regParamsOptions);
  }

  getPingSlotPeriodOptions = (search, callbackFunc) => {
    const pingSlotPeriodOptions = [
      {value: 32 * 1, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '' })},
      {value: 32 * 2, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '2' })},
      {value: 32 * 4, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '4' })},
      {value: 32 * 8, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '8' })},
      {value: 32 * 16, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '16' })},
      {value: 32 * 32, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '32' })},
      {value: 32 * 64, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '64' })},
      {value: 32 * 128, label: i18n.t(`${packageNS}:tr000200`,  { frequency: '128' })},
    ];

    callbackFunc(pingSlotPeriodOptions);
  }

  getPayloadCodecOptions = (search, callbackFunc) => {
    const payloadCodecOptions = [
      {value: "", label: i18n.t(`${packageNS}:tr000211`)},
      {value: "CAYENNE_LPP", label: i18n.t(`${packageNS}:tr000212`)},
      {value: "CUSTOM_JS", label: i18n.t(`${packageNS}:tr000213`)},
    ];

    callbackFunc(payloadCodecOptions);
  }

  onCodeChange = (field, editor, data, newCode) => {
    let object = this.state.object;
    object[field] = newCode;
    this.setState({
      object: object,
    });
  }

  setActiveTab = (tab) => {
    this.setState({
      activeTab: tab
    })
  }

  toggle = (tab) => {
    const { activeTab } = this.state;
  
    if (activeTab !== tab) {
      this.setActiveTab(tab);
    }
  }

  setValidationErrors = (errors) => {
    this.setState({
      validationErrors: errors
    })
  }

  formikFormSchema = () => {
    let fieldsSchema = {
      object: Yup.object().shape({
        // https://regexr.com/4rg3a
        // name: Yup.string().trim().matches(/^[0-9A-Za-z-]*$/g, i18n.t(`${packageNS}:tr000429`))
        //   .required(i18n.t(`${packageNS}:tr000431`)),
        name: Yup.string().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        maxEIRP: Yup.number()
          .required(i18n.t(`${packageNS}:tr000431`))
      })
    }

    if (this.props.update) {
      fieldsSchema.object.fields.id = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("id");
    }

    if (!this.props.update) {
      fieldsSchema.object.fields.networkServerID = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("networkServerID");
    }

    if (!this.state.object.supportsJoin) {
      fieldsSchema.object.fields.rxDelay1 = Yup.number().required(i18n.t(`${packageNS}:tr000431`)); 
      fieldsSchema.object._nodes.push("rxDelay1");
      fieldsSchema.object.fields.rxDROffset1 = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("rxDROffset1");
      fieldsSchema.object.fields.rxDataRate2 = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("rxDataRate2");
      fieldsSchema.object.fields.rxFreq2 = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("rxFreq2");
      fieldsSchema.object.fields.factoryPresetFreqs = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("factoryPresetFreqs");
      fieldsSchema.object.fields.classBTimeout = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("classBTimeout");
    }

    if (this.state.object.supportsClassB) {
      fieldsSchema.object.fields.pingSlotPeriod = Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("pingSlotPeriod");
      fieldsSchema.object.fields.pingSlotDR = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("pingSlotDR");
      fieldsSchema.object.fields.pingSlotFreq = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("pingSlotFreq");
    }

    if (this.state.object.supportsClassC) {
      fieldsSchema.object.fields.classCTimeout = Yup.number().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("classCTimeout");
    }

    return Yup.object().shape(fieldsSchema);
  }


  render() {
    const { activeTab, loading: loadingState, object } = this.state;
    const { loading: loadingProps } = this.props;
    const isUpdatePage = this.props.update;
    let isLoading = (loadingState || loadingProps);

    if (object === undefined) {
      return null;
    }

    // FIXME - shouldn't this be `isLoading = isUpdatePage ? isLoading : false;`
    isLoading = isUpdatePage ? false : isLoading;

    const codeMirrorOptions = {
      lineNumbers: true,
      mode: "javascript",
      theme: "default",
    };
    
    let payloadEncoderScript = object.payloadEncoderScript;
    let payloadDecoderScript = object.payloadDecoderScript;

    if (payloadEncoderScript === "" || payloadEncoderScript === undefined) {
      payloadEncoderScript = `// Encode encodes the given object into an array of bytes.
//  - fPort contains the LoRaWAN fPort number
//  - obj is an object, e.g. {"temperature": 22.5}
// The function must return an array of bytes, e.g. [225, 230, 255, 0]
function Encode(fPort, obj) {
  return [];
}`;
    }

    if (payloadDecoderScript === "" || payloadDecoderScript === undefined) {
      payloadDecoderScript = `// Decode decodes an array of bytes into an object.
//  - fPort contains the LoRaWAN fPort number
//  - bytes is an array of bytes, e.g. [225, 230, 255, 0]
// The function must return an object, e.g. {"temperature": 22.5}
function Decode(fPort, bytes) {
  return {};
}`;
    }

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                id: object.id || undefined,
                name: object.name || "",
                networkServerID: object.networkServerID || "",
                macVersion: object.macVersion || "",
                regParamsRevision: object.regParamsRevision || "",
                maxEIRP: object.maxEIRP || 0,
                geolocBufferTTL: object.geolocBufferTTL || 0,
                geolocMinBufferSize: object.geolocMinBufferSize || 0,
                supportsJoin: !!object.supportsJoin || false,
                rxDelay1: object.rxDelay1 || 0,
                rxDROffset1: object.rxDROffset1 || 0,
                rxDataRate2: object.rxDataRate2 || 0,
                rxFreq2: object.rxFreq2 || 0,
                factoryPresetFreqs: object.factoryPresetFreqs
                  ? object.factoryPresetFreqs.length
                    ? object.factoryPresetFreqs.join(",")
                    : ","
                  // China, Europe, US (i.e. 433MHz)
                  : "433000000,868000000,915000000",
                supportsClassB: !!object.supportsClassB || false,
                classBTimeout: object.classBTimeout || 0,
                pingSlotPeriod: object.pingSlotPeriod || 0,
                pingSlotDR: object.pingSlotDR || 0,
                pingSlotFreq: object.pingSlotFreq || 0,
                supportsClassC: !!object.supportsClassC || false,
                classCTimeout: object.classCTimeout || 0,
                payloadCodec: object.payloadCodec || ""
              }
            }
          }
          validateOnBlur
          validateOnChange
          // FIXME - temporarily disabled validate on mount, because when
          // the user switches between tabs after correcting an invalid input field
          // (i.e. reducing the detected invalid inputs from say 3 to 2), it restores
          // the validation errors back to 3 again, even though that input field
          // has a correct input field value in the UI. fix this after MVP when
          // we have more time to investigate.

          // validateOnMount
          validationSchema={this.formikFormSchema}
          // Formik Nested Schema Example https://codesandbox.io/s/y7q2v45xqx
          onSubmit={
            (castValues, { setSubmitting }) => {
              const values = this.formikFormSchema().cast(castValues);
              console.log('Submitted values: ', values);

              let newValues = clone(values);
              console.log('Deep copied submitted values: ', newValues !== values);

              let newFactoryPresetFreqsArr;
              if (values.object.factoryPresetFreqs) {
                newFactoryPresetFreqsArr = values.object.factoryPresetFreqs.split(",").filter((v) => v !== '');
                newValues.object.factoryPresetFreqs = newFactoryPresetFreqsArr;
              }

              console.log('Prepared values', newValues);

              this.props.onSubmit(newValues.object);
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
                setValues,
                touched,
                validateForm,
                values
              } = props;
              return (  
                <Form style={{ padding: "0px", backgroundColor: "#ebeff2" }} onSubmit={handleSubmit} noValidate>
                  <Nav tabs>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '1' })}
                        onClick={() => { this.toggle('1'); }}
                      >
                        {i18n.t(`${packageNS}:tr000167`)}
                      </NavLink>
                    </NavItem>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '2' })}
                        onClick={() => { this.toggle('2'); }}
                      >
                        <i className="mdi mdi-radio-tower"></i>
                        &nbsp;{i18n.t(`${packageNS}:tr000184`)}
                      </NavLink>
                    </NavItem>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '3' })}
                        onClick={() => { this.toggle('3'); }}
                      >
                        <i className="mdi mdi-tag-outline"></i>
                        &nbsp;{i18n.t(`${packageNS}:tr000194`)}
                      </NavLink>
                    </NavItem>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '4' })}
                        onClick={() => { this.toggle('4'); }}
                      >
                        <i className="mdi mdi-tag-outline"></i>
                        &nbsp;{i18n.t(`${packageNS}:tr000203`)}
                      </NavLink>
                    </NavItem>
                    <NavItem>
                      <NavLink
                        className={classnames({ active: activeTab === '5' })}
                        onClick={() => { this.toggle('5'); }}
                      >
                        <i className="mdi mdi-code-braces"></i>
                        &nbsp;{i18n.t(`${packageNS}:tr000208`)}
                      </NavLink>
                    </NavItem>
                  </Nav>

                  <TabContent
                      activeTab={activeTab}
                      style={{
                        backgroundColor: "#fff",
                        borderRadius: "2px",
                        borderStyle: "solid",
                        borderWidth: "0 1px 1px 1px",
                        borderColor: "#ddd"
                      }}>
                      <TabPane tabId="1">
                        {isLoading && <Loader light />}

                        {this.props.update &&
                          <>
                            <label htmlFor="object.id" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                              {i18n.t(`${packageNS}:tr000077`)}
                            </label>
                            &nbsp;&nbsp;{values.object.id}

                            <input
                              type="hidden"
                              id="id"
                              disabled
                              name="object.id"
                              value={values.object.id}
                            />
                            {
                              errors.object && errors.object.id
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.id}
                                  </div>
                                ) : null
                            }
                            <br /><br />
                          </>
                        }

                        <Field
                          id="name"
                          name="object.name"
                          type="text"
                          value={values.object.name}
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000168`)}
                          helpText={i18n.t(`${packageNS}:tr000062`)}
                          component={ReactstrapInput}
                          className={
                            errors.object && errors.object.name
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        {
                          errors.object && errors.object.name
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.object.name}
                              </div>
                            ) : null
                        }

                        {!this.props.update &&
                          <> 
                            <label htmlFor="object.networkServerID" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                              {i18n.t(`${packageNS}:tr000047`)}
                            </label>
                            <AutocompleteSelect
                              id="networkServerID"
                              name="object.networkServerID"
                              label={i18n.t(`${packageNS}:tr000115`)}
                              onChange={handleChange}
                              getOptions={this.getNetworkServerOptions}
                              value={values.object.networkServerID}
                              className={
                                errors.object && errors.object.networkServerID
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            <FormHelperText>
                              {i18n.t(`${packageNS}:tr000171`)}
                            </FormHelperText>
                            {
                              errors.object && errors.object.networkServerID
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    <br />
                                    {errors.object.networkServerID}
                                  </div>
                                ) : null
                            }
                            <br />
                          </>
                        }

                        <label htmlFor="object.macVersion" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000172`)}
                        </label>
                        <AutocompleteSelect
                          id="macVersion"
                          name="object.macVersion"
                          label={i18n.t(`${packageNS}:tr000173`)}
                          onChange={handleChange}
                          getOptions={this.getMACVersionOptions}
                          value={values.object.macVersion}
                        />
                        <FormHelperText>
                          {i18n.t(`${packageNS}:tr000174`)}
                        </FormHelperText>

                        <label htmlFor="object.regParamsRevision" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000175`)}
                        </label>
                        <AutocompleteSelect
                          id="regParamsRevision"
                          name="object.regParamsRevision"
                          label={i18n.t(`${packageNS}:tr000176`)}
                          onChange={handleChange}
                          getOptions={this.getRegParamsOptions}
                          value={values.object.regParamsRevision}
                        />
                        <FormHelperText>
                          {i18n.t(`${packageNS}:tr000177`)}
                        </FormHelperText>
                        <br />

                        <Field
                          id="maxEIRP"
                          name="object.maxEIRP"
                          type="number"
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000178`)}
                          helpText={i18n.t(`${packageNS}:tr000179`)}
                          component={ReactstrapInput}
                          className={
                            errors.object && errors.object.maxEIRP
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        {
                          errors.object && errors.object.maxEIRP
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.object.maxEIRP}
                              </div>
                            ) : null
                        }

                        <Field
                          id="geolocBufferTTL"
                          name="object.geolocBufferTTL"
                          type="number"
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000180`)}
                          helpText={i18n.t(`${packageNS}:tr000181`)}
                          component={ReactstrapInput}
                        />

                        <Field
                          id="geolocMinBufferSize"
                          name="object.geolocMinBufferSize"
                          type="number"
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000182`)}
                          helpText={i18n.t(`${packageNS}:tr000183`)}
                          component={ReactstrapInput}
                        />
                      </TabPane>
                      <TabPane tabId="2">
                        <div style={{ marginTop: "10px" }}>
                          <FormGroup>
                            <FormControlLabel
                              label={i18n.t(`${packageNS}:tr000185`)}
                              control={
                                <Checkbox
                                  id="supportsJoin"
                                  name="object.supportsJoin"
                                  onChange={handleChange}
                                  onBlur={handleBlur}
                                  color="primary"
                                  value={!!values.object.supportsJoin}
                                  checked={!!values.object.supportsJoin}
                                  // Note: This approach did not work.
                                  // Instead we just set the default value of
                                  // `factoryPresetFreqs` to ',' if existing
                                  // Device Profile does not have a value, or
                                  // LPWAN default frequencies.
                                  //
                                  // setFieldValue={() =>
                                  //   setFieldValue(
                                  //     "object.factoryPresetFreqs",
                                  //     // If "Device Supports OTAA" is checked, then
                                  //     // the `factoryPresetFreqs` field is no longer shown,
                                  //     // but we'll temporarily give it a value so the validation
                                  //     // error disappears. If it's unchecked again, we'll reset
                                  //     // the value to `""`
                                  //     !!values.object.supportsJoin ? "0" : "",
                                  //     true
                                  //   )
                                  // }
                                />
                              }
                            />
                          </FormGroup>
                        </div>

                        {!values.object.supportsJoin &&
                          <>
                            <Field
                              id="rxDelay1"
                              name="object.rxDelay1"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000186`)}
                              helpText={i18n.t(`${packageNS}:tr000187`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.rxDelay1
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.rxDelay1
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.rxDelay1}
                                  </div>
                                ) : null
                            }

                            <Field
                              id="rxDROffset1"
                              name="object.rxDROffset1"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000188`)}
                              helpText={i18n.t(`${packageNS}:tr000189`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.rxDROffset1
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.rxDROffset1
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.rxDROffset1}
                                  </div>
                                ) : null
                            }

                            <Field
                              id="rxDataRate2"
                              name="object.rxDataRate2"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000190`)}
                              helpText={i18n.t(`${packageNS}:tr000189`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.rxDataRate2
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.rxDataRate2
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.rxDataRate2}
                                  </div>
                                ) : null
                            }

                            <Field
                              id="rxFreq2"
                              name="object.rxFreq2"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000191`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.rxFreq2
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.rxFreq2
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.rxFreq2}
                                  </div>
                                ) : null
                            }

                            <Field
                              id="factoryPresetFreqs"
                              name="object.factoryPresetFreqs"
                              type="string"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000192`)}
                              helpText={i18n.t(`${packageNS}:tr000193`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.factoryPresetFreqs
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.factoryPresetFreqs
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.factoryPresetFreqs}
                                  </div>
                                ) : null
                            }
                          </>
                        }
                      </TabPane>
                      <TabPane tabId="3">
                        <div style={{ marginTop: "10px" }}>
                          <FormGroup>
                            <FormControlLabel
                              label={i18n.t(`${packageNS}:tr000195`)}
                              control={
                                <Checkbox
                                  id="supportsClassB"
                                  name="object.supportsClassB"
                                  onChange={handleChange}
                                  color="primary"
                                  value={!!values.object.supportsClassB}
                                  checked={!!values.object.supportsClassB}
                                />
                              }
                            />
                          </FormGroup>
                        </div>
                        {!!values.object.supportsClassB &&
                          <>
                            <Field
                              id="classBTimeout"
                              name="object.classBTimeout"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000196`)}
                              helpText={i18n.t(`${packageNS}:tr000197`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.classBTimeout
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.classBTimeout
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.classBTimeout}
                                  </div>
                                ) : null
                            }

                            <label htmlFor="object.pingSlotPeriod" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                              {i18n.t(`${packageNS}:tr000198`)}
                            </label>
                            <AutocompleteSelect
                              id="pingSlotPeriod"
                              name="object.pingSlotPeriod"
                              label={i18n.t(`${packageNS}:tr000199`)}
                              onChange={handleChange}
                              getOptions={this.getPingSlotPeriodOptions}
                              value={values.object.pingSlotPeriod}
                              className={
                                errors.object && errors.object.pingSlotPeriod
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            <FormHelperText>
                              {i18n.t(`${packageNS}:tr000198`)}
                            </FormHelperText>
                            {
                              errors.object && errors.object.pingSlotPeriod
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.pingSlotPeriod}
                                  </div>
                                ) : null
                            }
                            <br />

                            <Field
                              id="pingSlotDR"
                              name="object.pingSlotDR"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000201`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.pingSlotDR
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.pingSlotDR
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.pingSlotDR}
                                  </div>
                                ) : null
                            }

                            <Field
                              id="pingSlotFreq"
                              name="object.pingSlotFreq"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000202`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.pingSlotFreq
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.pingSlotFreq
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.pingSlotFreq}
                                  </div>
                                ) : null
                            }
                          </>
                        }
                      </TabPane>
                      <TabPane tabId="4">
                        <div style={{ marginTop: "10px" }}>
                          <FormGroup>
                            <FormControlLabel
                              label={i18n.t(`${packageNS}:tr000204`)}
                              control={
                                <Checkbox
                                  id="supportsClassC"
                                  name="object.supportsClassC"
                                  onChange={handleChange}
                                  color="primary"
                                  value={!!values.object.supportsClassC}
                                  checked={!!values.object.supportsClassC}
                                />
                              }
                            />
                          </FormGroup>
                          <FormHelperText>
                            {i18n.t(`${packageNS}:tr000205`)}
                          </FormHelperText>
                        </div>
                        <br />

                        {!!values.object.supportsClassC &&
                          <>
                            <Field
                              id="classCTimeout"
                              name="object.classCTimeout"
                              type="number"
                              onChange={handleChange}
                              onBlur={handleBlur}
                              label={i18n.t(`${packageNS}:tr000206`)}
                              helpText={i18n.t(`${packageNS}:tr000207`)}
                              component={ReactstrapInput}
                              className={
                                errors.object && errors.object.classCTimeout
                                  ? 'is-invalid form-control'
                                  : ''
                              }
                            />
                            {
                              errors.object && errors.object.classCTimeout
                                ? (
                                  <div
                                    className="invalid-feedback"
                                    style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                  >
                                    {errors.object.classCTimeout}
                                  </div>
                                ) : null
                            }
                            <br />
                          </>
                        }
                      </TabPane>
                      <TabPane tabId="5">
                        <label htmlFor="object.payloadCodec" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000209`)}
                        </label>
                        <AutocompleteSelect
                          id="payloadCodec"
                          name="object.payloadCodec"
                          label={i18n.t(`${packageNS}:tr000210`)}
                          onChange={handleChange}
                          getOptions={this.getPayloadCodecOptions}
                          value={values.object.payloadCodec}
                          className={
                            errors.object && errors.object.payloadCodec
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        <FormHelperText>
                          {i18n.t(`${packageNS}:tr000214`)}
                        </FormHelperText>
                        {
                          errors.object && errors.object.payloadCodec
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.object.payloadCodec}
                              </div>
                            ) : null
                        }
                        <br />

                        {this.state.object.payloadCodec === "CUSTOM_JS" &&
                          <>
                            <Label for="payloadDecoderScript">
                              {i18n.t(`${packageNS}:tr000551`)}
                            </Label>
                            <CodeMirror
                              value={object.payloadDecoderScript}
                              options={codeMirrorOptions}
                              onBeforeChange={this.onCodeChange.bind(this, 'payloadDecoderScript')}
                              className={this.props.classes.codeMirror}
                            />
                            <FormHelperText>
                              {i18n.t(`${packageNS}:tr000215`)}
                            </FormHelperText>

                            <Label for="payloadEncoderScript">
                            {i18n.t(`${packageNS}:tr000552`)}
                            </Label>
                            <CodeMirror
                              value={object.payloadEncoderScript}
                              options={codeMirrorOptions}
                              onBeforeChange={this.onCodeChange.bind(this, 'payloadEncoderScript')}
                              className={this.props.classes.codeMirror}
                            />
                            <FormHelperText>
                              {i18n.t(`${packageNS}:tr000216`)}
                            </FormHelperText>
                          </>
                        }
                      </TabPane>
                    </TabContent>

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
                      {errors.object && Object.keys(errors.object).length
                        ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                          Detected {Object.keys(errors.object).length} errors. Please fix the validation errors shown in each tab before resubmitting.
                        </div>)
                        : null
                      }

                      {/* Show error count when user submits the form */}
                      {this.state.validationErrors && this.state.validationErrors.length
                        ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                          Detected {Object.keys(this.state.validationErrors.object).length} errors. Please fix the validation errors shown in each tab before resubmitting.
                        </div>)
                        : null
                      }
                    </div>
                    <Button
                      type="submit"
                      color="primary"
                      disabled={(errors.object && Object.keys(errors.object).length > 0) || isLoading || isSubmitting}
                      onClick={
                        () => { 
                          validateForm().then((formValidationErrors) => {
                            console.log('Validated form with errors: ', formValidationErrors)
                            this.setValidationErrors(formValidationErrors);
                          })
                        }
                      }
                    >
                      {this.props.submitLabel || (this.props.deviceProfile ? "Update" : "Create")}
                    </Button>
                  {/* </Card> */}
                </Form>
              );
            }
          }
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(DeviceProfileForm);
