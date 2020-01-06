import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../i18n';
import { ReactstrapInput } from '../../components/FormInputs';
import AESKeyField from "../../components/FormikAESKeyField";
import Loader from "../../components/Loader";
import DeviceStore from "../../stores/DeviceStore";


class DeviceKeys extends Component {
  constructor() {
    super();
    this.state = {
      loading: true,
      update: false,
    };
  }

  componentDidMount() {
    const { match } = this.props;

    DeviceStore.getKeys(match.params.devEUI, resp => {
      if (resp === null) {
        this.setState({
          object: {
            deviceKeys: {},
          },
          loading: false
        });
      } else {
        this.setState({
          update: true,
          object: resp,
          loading: false
        });
      }
    });
  }

  onSubmit = (deviceKeys) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    const isApplication = currentApplicationID && currentApplicationID !== "0"; 

    const { history, match } = this.props;

    if (this.state.update) {
      DeviceStore.updateKeys(deviceKeys, resp => {
        isApplication
        ? history.push(`/organizations/${currentOrgID}/applications/${currentApplicationID}`)
        : history.push(`/organizations/${currentOrgID}`);
      });
    } else {
      let keys = deviceKeys;
      keys.devEUI = match.params.devEUI;
      DeviceStore.createKeys(keys, resp => {
        isApplication
        ? history.push(`/organizations/${currentOrgID}/applications/${currentApplicationID}`)
        : history.push(`/organizations/${currentOrgID}`);
      });
    }
  }

  formikFormSchema = () => {
    let fieldsSchema = {
      object: Yup.object().shape({
        deviceKeys: Yup.object().shape({
          // https://regexr.com/4rg3a
          nwkKey: Yup.string()
            .required(i18n.t(`${packageNS}:tr000431`)),
          devEUI: Yup.string()
            .required("DevEUI is required"),
        })
      })
    }

    if (this.props.deviceProfile.macVersion.startsWith("1.1")) {
      fieldsSchema['object.deviceKeys.genAppKey'] = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
    }

    return Yup.object().shape(fieldsSchema);
  }

  render() {
    const { loading, object } = this.state;
    const { deviceProfile } = this.props;

    if (object === undefined) {
      return <React.Fragment>{loading && <Loader light />}</React.Fragment>
    }

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                deviceKeys: {
                  devEUI: object.deviceKeys.devEUI || undefined,
                  nwkKey: object.deviceKeys.nwkKey || undefined,
                  genAppKey: object.deviceKeys.genAppKey || undefined,
                  appKey: object.deviceKeys.appKey || undefined,
                }
              }
            }
          }
          validateOnBlur
          validateOnChange
          validationSchema={this.formikFormSchema}
          onSubmit={
            (values, { setSubmitting }) => {
              console.log('Submitted values: ', values);

              this.onSubmit(values.object.deviceKeys);
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
                  {object && (deviceProfile.macVersion.startsWith("1.0") || deviceProfile.macVersion.startsWith("1.1")) && (
                    <>
                      <span style={{ display: 'block', fontSize: "16px", fontWeight: "700" }}>
                        { deviceProfile.macVersion.startsWith("1.0") ? "LPWAN 1.0 Device Keys" : "" }
                        { deviceProfile.macVersion.startsWith("1.1") ? "LPWAN 1.1 Device Keys" : "" }
                      </span>
                      <label htmlFor="object.deviceKeys.nwkKey" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                        {i18n.t(`${packageNS}:tr000388`)}
                      </label>
                      <AESKeyField
                        id="nwkKey"
                        name="object.deviceKeys.nwkKey"
                        helperText={i18n.t(`${packageNS}:tr000397`)}
                        onChange={handleChange}
                        value={object.deviceKeys.nwkKey || ""}
                        required
                        random
                        // FIXME - we need input field validation styles to work
                        // className={
                        //   errors.object && errors.object.deviceKeys.nwkKey
                        //     ? 'is-invalid form-control'
                        //     : ''
                        // }
                      />

                      {
                        errors.object && errors.object.deviceKeys.nwkKey
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.object.deviceKeys.nwkKey}
                            </div>
                          ) : null
                      }
                    </>
                  )}

                  {object && deviceProfile.macVersion.startsWith("1.0") && (
                    <>
                      <label htmlFor="object.deviceKeys.genAppKey" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                        {i18n.t(`${packageNS}:tr000389`)}
                      </label>
                      <AESKeyField
                        id="genAppKey"
                        name="object.deviceKeys.genAppKey"
                        helperText={i18n.t(`${packageNS}:tr000398`)}
                        onChange={handleChange}
                        value={object.deviceKeys.genAppKey || ""}
                        random
                      />
                    </>
                  )}
                  {object && deviceProfile.macVersion.startsWith("1.1") && (
                    <>
                      <label htmlFor="object.deviceKeys.appKey" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                        {i18n.t(`${packageNS}:tr000387`)}
                      </label>
                      <AESKeyField
                        id="appKey"
                        name="object.deviceKeys.appKey"
                        helperText={i18n.t(`${packageNS}:tr000386`)}
                        onChange={handleChange}
                        value={object.deviceKeys.appKey || ""}
                        required
                        random
                      />

                      {/* FIXME - this validation should be built-in to each input field */}
                      {
                        errors.object && errors.object.deviceKeys.appKey
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.object.deviceKeys.appKey}
                            </div>
                          ) : null
                      }
                    </>
                  )}

                  <>
                    <label htmlFor="object.deviceKeys.devEUI" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                      {i18n.t(`${packageNS}:tr000371`)}
                    </label>
                    &nbsp;&nbsp;{object.deviceKeys.devEUI}

                    <input
                      type="hidden"
                      id="devEUI"
                      disabled
                      name="object.deviceKeys.devEUI"
                      value={object.deviceKeys.devEUI || ""}
                    />
                    {
                      errors.object && errors.object.deviceKeys.devEUI
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.deviceKeys.devEUI}
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
                    {this.props.submitLabel || i18n.t(`${packageNS}:tr000292`)}
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

export default withRouter(DeviceKeys);
