import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";
import Typograhy from "@material-ui/core/Typography";
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../i18n';
import AESKeyField from "../../components/FormikAESKeyField";
import DevAddrField from "../../components/FormikDevAddrField";
import Loader from "../../components/Loader";
import { ReactstrapInput } from '../../components/FormInputs';
import DeviceStore from "../../stores/DeviceStore";
import theme from "../../theme";


const styles = {
  link: {
    color: theme.palette.primary.main,
  },
};


class DeviceActivation extends Component {
  constructor() {
    super();
    this.state = {
      loading: true,
    };
  }
  
  componentDidMount() {
    DeviceStore.getActivation(this.props.match.params.devEUI, resp => {
      if (resp === null) {
        this.setState({
          loading: false,
          object: {
            deviceActivation: {},
          },
        });
      } else {
        this.setState({
          loading: false,
          object: resp,
        });
      }
    });
  }

  getRandomDevAddr = (cb) => {
    DeviceStore.getRandomDevAddr(this.props.match.params.devEUI, resp => {
      cb(resp.devAddr);
    });
  }

  onSubmit = (deviceActivation) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0"; 

    const { history, match } = this.props;
    let act = deviceActivation;
    act.devEUI = match.params.devEUI;

    if (this.props.deviceProfile.macVersion.startsWith("1.0")) {
      act.fNwkSIntKey = act.nwkSEncKey;
      act.sNwkSIntKey = act.nwkSEncKey;
    }

    DeviceStore.activate(act, resp => {
      isApplication
      ? history.push(`/organizations/${currentOrgID}/applications/${currentApplicationID}`)
      : history.push(`/organizations/${currentOrgID}`);
    });
  }

  formikFormSchema = () => {
    let fieldsSchema = {
      object: Yup.object().shape({
        deviceActivation: Yup.object().shape({
          // https://regexr.com/4rg3a
          devEUI: Yup.string().trim()
            // .trim().matches(/([A-Fa-f0-9]){16}/, "Must be length 16")
            .required(i18n.t(`${packageNS}:tr000431`)),
          devAddr: Yup.string().trim()
            // FIXME - changes to DevAddr component required to get these to work
            // since the length of the value when debugging `values.object.deviceActivation.devAddr`
            // changes from 8 to 11 if you change the value in the UI from 8 to 7.
            // .trim().matches(/([A-Fa-f0-9]){8}/, "Must be length 8")
            .required(i18n.t(`${packageNS}:tr000431`)),
          appSKey: Yup.string().trim()
            // .trim().matches(/([A-Fa-f0-9]){32}/, "Must be length 32")
            .required(i18n.t(`${packageNS}:tr000431`)),
          nwkSEncKey: Yup.string().trim()
            // .trim().matches(/([A-Fa-f0-9]){32}/, "Must be length 32")
            .required(i18n.t(`${packageNS}:tr000431`)),
          fCntUp: Yup.number().trim()
            .required(i18n.t(`${packageNS}:tr000431`)),
          nFCntDown: Yup.number().trim()
            .required(i18n.t(`${packageNS}:tr000431`)),
        })
      })
    }

    if (this.props.deviceProfile.macVersion.startsWith("1.1")) {
      fieldsSchema.object.fields.deviceActivation.fields.sNwkSIntKey = Yup.string().trim()
        // .trim().matches(/([A-Fa-f0-9]){32}/, "Must be length 32")
        .required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object.fields.deviceActivation._nodes.push("sNwkSIntKey");
      fieldsSchema.object.fields.deviceActivation.fields.fNwkSIntKey = Yup.string().trim()
        // .trim().matches(/([A-Fa-f0-9]){32}/, "Must be length 32")
        .required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object.fields.deviceActivation._nodes.push("fNwkSIntKey");
      fieldsSchema.object.fields.deviceActivation.fields.aFCntDown = Yup.number().trim()
        .required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object.fields.deviceActivation._nodes.push("aFCntDown");
    }

    return Yup.object().shape(fieldsSchema);
  }

  render() {
    const { loading, object } = this.state;
    const { deviceProfile } = this.props;

    if (object === undefined) {
      return <React.Fragment>{loading && <Loader light />}</React.Fragment>
    }

    let submitLabel = deviceProfile && !deviceProfile.supportsJoin
      ? "(Re)activate Device"
      : "Activate Device";

    let showForm = false;
    if (
      deviceProfile && !deviceProfile.supportsJoin ||
      (deviceProfile && deviceProfile.supportsJoin && object.deviceActivation.devAddr !== undefined)
    ) {
      showForm = true;
    }

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                deviceActivation: {
                  devEUI: object.deviceActivation.devEUI || undefined,
                  devAddr: object.deviceActivation.devAddr || undefined,
                  appSKey: object.deviceActivation.appSKey || undefined,
                  nwkSEncKey: object.deviceActivation.nwkSEncKey || undefined,
                  sNwkSIntKey: object.deviceActivation.sNwkSIntKey || undefined,
                  fNwkSIntKey: object.deviceActivation.fNwkSIntKey || undefined,
                  fCntUp: object.deviceActivation.fCntUp || undefined,
                  nFCntDown: object.deviceActivation.nFCntDown || undefined,
                  aFCntDown: object.deviceActivation.aFCntDown || undefined
                }
              }
            }
          }
          validateOnBlur
          validateOnChange
          validationSchema={this.formikFormSchema}
          onSubmit={
            (castValues, { setSubmitting }) => {
              const values = this.formikFormSchema().cast(castValues);
              console.log('Submitted values: ', values);

              this.onSubmit(values.object.deviceActivation);
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
                isSubmitting,
                isValidating,
                setFieldValue,
                touched,
                validateForm,
                values
              } = props;
              return (
                <Form style={{ padding: "0px", backgroundColor: "#fff" }} onSubmit={handleSubmit} noValidate>

                  {deviceProfile && (deviceProfile.macVersion.startsWith("1.0") || deviceProfile.macVersion.startsWith("1.1")) && (
                    <span style={{ display: 'block', fontSize: "16px", fontWeight: "700" }}>
                      { deviceProfile.macVersion.startsWith("1.0") ? `LPWAN 1.0 Device ${!deviceProfile.supportsJoin ? '(Re)Activation' : 'Activation'}` : "" }
                      { deviceProfile.macVersion.startsWith("1.1") ? `LPWAN 1.1 Device ${!deviceProfile.supportsJoin ? '(Re)Activation' : 'Activation'}` : "" }
                    </span>
                  )}
                  <br />

                  {showForm && this.props.deviceProfile.macVersion.startsWith("1.0") &&
                    <>
                      <DevAddrField
                        id="devAddr"
                        name="object.deviceActivation.devAddr"
                        label={i18n.t(`${packageNS}:tr000312`)}
                        value={object.deviceActivation.devAddr || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        randomFunc={this.getRandomDevAddr}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.devAddr
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.devAddr}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="nwkSEncKey"
                        name="object.deviceActivation.nwkSEncKey"
                        label={i18n.t(`${packageNS}:tr000313`)}
                        value={object.deviceActivation.nwkSEncKey || ""}
                        onChange={handleChange}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.nwkSEncKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.nwkSEncKey}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="appSKey"
                        name="object.deviceActivation.appSKey"
                        label={i18n.t(`${packageNS}:tr000314`)}
                        value={object.deviceActivation.appSKey || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.appSKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.appSKey}
                            </div>
                          : null
                      }
                      <Field
                        id="fCntUp"
                        name="object.deviceActivation.fCntUp"
                        label={i18n.t(`${packageNS}:tr000315`)}
                        type="number"
                        value={object.deviceActivation.fCntUp || 0}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        component={ReactstrapInput}
                        required
                        className={
                          errors.object && errors.object.deviceActivation.fCntUp
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.deviceActivation.fCntUp
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.fCntUp}
                            </div>
                          : null
                      }
                      <Field
                        id="nFCntDown"
                        name="object.deviceActivation.nFCntDown"
                        label={i18n.t(`${packageNS}:tr000316`)}
                        type="number"
                        value={object.deviceActivation.nFCntDown || 0}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        component={ReactstrapInput}
                        required
                        className={
                          errors.object && errors.object.deviceActivation.nFCntDown
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.deviceActivation.nFCntDown
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.nFCntDown}
                            </div>
                          : null
                      }
                    </>
                  }

                  {showForm && this.props.deviceProfile.macVersion.startsWith("1.1") &&
                    <>
                      <DevAddrField
                        id="devAddr"
                        name="object.deviceActivation.devAddr"
                        label={i18n.t(`${packageNS}:tr000312`)}
                        value={object.deviceActivation.devAddr || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        randomFunc={this.getRandomDevAddr}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.devAddr
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.devAddr}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="nwkSEncKey"
                        name="object.deviceActivation.nwkSEncKey"
                        label={i18n.t(`${packageNS}:tr000332`)}
                        value={object.deviceActivation.nwkSEncKey || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.nwkSEncKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.nwkSEncKey}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="sNwkSIntKey"
                        name="object.deviceActivation.sNwkSIntKey"
                        label={i18n.t(`${packageNS}:tr000333`)}
                        value={object.deviceActivation.sNwkSIntKey || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.sNwkSIntKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.sNwkSIntKey}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="fNwkSIntKey"
                        name="object.deviceActivation.fNwkSIntKey"
                        label={i18n.t(`${packageNS}:tr000334`)}
                        value={object.deviceActivation.fNwkSIntKey || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.fNwkSIntKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.fNwkSIntKey}
                            </div>
                          : null
                      }
                      <AESKeyField
                        id="appSKey"
                        name="object.deviceActivation.appSKey"
                        label={i18n.t(`${packageNS}:tr000335`)}
                        value={object.deviceActivation.appSKey || ""}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        required
                        random
                      />
                      {
                        errors.object && errors.object.deviceActivation.appSKey
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.appSKey}
                            </div>
                          : null
                      }
                      <Field
                        id="fCntUp"
                        name="object.deviceActivation.fCntUp"
                        label={i18n.t(`${packageNS}:tr000336`)}
                        type="number"
                        value={object.deviceActivation.fCntUp || 0}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        component={ReactstrapInput}
                        required
                        className={
                          errors.object && errors.object.deviceActivation.fCntUp
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.deviceActivation.fCntUp
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.fCntUp}
                            </div>
                          : null
                      }
                      <Field
                        id="nFCntDown"
                        name="object.deviceActivation.nFCntDown"
                        label={i18n.t(`${packageNS}:tr000337`)}
                        type="number"
                        value={object.deviceActivation.nFCntDown || 0}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        component={ReactstrapInput}
                        required
                        className={
                          errors.object && errors.object.deviceActivation.nFCntDown
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.deviceActivation.nFCntDown
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.nFCntDown}
                            </div>
                          : null
                      }
                      <Field
                        id="aFCntDown"
                        name="object.deviceActivation.aFCntDown"
                        label={i18n.t(`${packageNS}:tr000338`)}
                        type="number"
                        value={object.deviceActivation.aFCntDown || 0}
                        onChange={handleChange}
                        disabled={deviceProfile && deviceProfile.supportsJoin}
                        component={ReactstrapInput}
                        required
                        className={
                          errors.object && errors.object.deviceActivation.aFCntDown
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors.object && errors.object.deviceActivation.aFCntDown
                          ? <div className="invalid-feedback" style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                              {errors.object.deviceActivation.aFCntDown}
                            </div>
                          : null
                      }
                    </>
                  }
                  {!showForm &&
                    <Typograhy variant="body1">
                      This device has not (yet) been activated.
                    </Typograhy>
                  }

                  <>
                    <label htmlFor="object.deviceActivation.devEUI" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                      {i18n.t(`${packageNS}:tr000371`)}
                    </label>
                    &nbsp;&nbsp;{object.deviceActivation.devEUI}

                    <input
                      type="hidden"
                      id="devEUI"
                      disabled
                      name="object.deviceActivation.devEUI"
                      value={object.deviceActivation.devEUI || ""}
                    />
                    {
                      errors.object && errors.object.deviceActivation.devEUI
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.deviceActivation.devEUI}
                          </div>
                        ) : null
                    }
                  </>
 
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
                    {errors.object
                      ? <div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                          Form Validation Errors. Please enter valid inputs and try again...
                        </div>
                      : ''
                    }
                  </div>
                  <Button
                    type="submit"
                    color="primary"
                    disabled={(errors.object && Object.keys(errors.object).length > 0) || loading || isSubmitting}
                    onClick={
                      () => validateForm().then((formValidationErrors) =>
                        console.log('Validated form with errors: ', formValidationErrors))
                    }
                  >
                    {submitLabel || this.props.submitLabel || i18n.t(`${packageNS}:tr000292`)}
                  </Button>
                  <br />
                </Form>
              );
            }
          }
        </Formik>
      </React.Fragment>
 
    );
  }
}

export default withRouter(withStyles(styles)(DeviceActivation));
