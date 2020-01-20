import React, { Component } from "react";

import { withStyles } from "@material-ui/core/styles";
import { NavLink, Card, Button, Row, Col } from 'reactstrap';
import { ErrorMessage, Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
import FormHelperText from "@material-ui/core/FormHelperText";

import i18n, { packageNS } from '../../i18n';
import { ReactstrapInput } from '../../components/FormInputs';
import AutocompleteSelect from "../../components/AutocompleteSelect";
import Loader from "../../components/Loader";

const styles = {
  formLabel: {
    fontSize: 12,
  },
};

class FUOTADeploymentForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      file: null,
      object: {}
    }
  }

  getGroupTypeOptions = (search, callbackFunc) => {
    const options = [
      {value: "CLASS_C", label: i18n.t(`${packageNS}:tr000203`)},
    ];

    callbackFunc(options);
  }

  getMulticastTimeoutOptions = (search, callbackFunc) => {
    let options = [];

    for (let i = 0; i < (1 << 4); i++) {
      options.push({
        label: `${1 << i} ${i18n.t(`${packageNS}:tr000357`)}`,
        value: i,
      });
    }

    callbackFunc(options);
  }

  onFileChange = (e) => {
    let object = this.state.object;

    if (e.target.files.length !== 1) {
      object.payload = "";

      this.setState({
        file: null,
        object: object,
      });
    } else {
      this.setState({
        file: e.target.files[0],
      });

      const reader = new FileReader();
      reader.onload = () => {
        const encoded = reader.result.replace(/^data:(.*;base64,)?/, '');
        object.payload = encoded;

        this.setState({
          object: object,
        });
      };
      reader.readAsDataURL(e.target.files[0]);
    }
  }

  formikFormSchema = () => {
    let fieldsSchema = {
      object: Yup.object().shape({
        // https://regexr.com/4rg3a
        name: Yup.string().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        redundancy: Yup.number().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        unicastTimeout: Yup.string()
          .trim().matches(/^[0-9]*$/, "Requires a number")
          .max(19, 'Requires number less than 19 digits')
          .required(i18n.t(`${packageNS}:tr000431`)),
        dr: Yup.number().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        frequency: Yup.number().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        groupType: Yup.string().trim()
          .required(i18n.t(`${packageNS}:tr000431`)),
        multicastTimeout: Yup.number().trim()
          .required(i18n.t(`${packageNS}:tr000431`))
      })
    }

    return Yup.object().shape(fieldsSchema);
  }

  render() {
    const { loading, object } = this.state;

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                name: object && object.name || "",
                redundancy: object && object.redundancy || 0,
                unicastTimeout: object && object.unicastTimeout || "0",
                dr: object && object.dr || 0,
                frequency: object && object.frequency || 0,
                groupType: object && object.groupType || "",
                multicastTimeout: object && object.multicastTimeout || 0
              }
            }
          }
          validateOnBlur
          validateOnChange
          validationSchema={this.formikFormSchema}
          // Formik Nested Schema Example https://codesandbox.io/s/y7q2v45xqx
          onSubmit={
            (castValues, { setSubmitting }) => {
              const values = this.formikFormSchema().cast(castValues);
              console.log('Submitted values: ', values);
              this.props.onSubmit(values.object);
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
                  {loading && <Loader light />}

                  <Field
                    id="name"
                    name="object.name"
                    type="text"
                    value={values.object.name}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    label={i18n.t(`${packageNS}:tr000369`)}
                    helpText={i18n.t(`${packageNS}:tr000368`)}
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

                  <label htmlFor="object.file" style={{ display: 'block', fontWeight: "700" }}>
                    {i18n.t(`${packageNS}:tr000367`)}
                  </label>
                  <input id="file" name="object.file" type="file" className="fuota-input-file" onChange={this.onFileChange} />
                  {this.state.file !== null ? (
                      <label htmlFor="object.file">
                        {this.state.file.name}&nbsp;{this.state.file.size} bytes
                      </label>
                    ) : null
                  }
                  <FormHelperText>
                    {i18n.t(`${packageNS}:tr000366`)}
                  </FormHelperText>
                  <br />

                  <Field
                    id="redundancy"
                    name="object.redundancy"
                    type="number"
                    value={values.object.redundancy}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    label={i18n.t(`${packageNS}:tr000344`)}
                    helpText={i18n.t(`${packageNS}:tr000364`)}
                    component={ReactstrapInput}
                    className={
                      errors.object && errors.object.redundancy
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  {
                    errors.object && errors.object.redundancy
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          {errors.object.redundancy}
                        </div>
                      ) : null
                  }

                  <Field
                    id="unicastTimeout"
                    name="object.unicastTimeout"
                    type="string"
                    value={values.object.unicastTimeout}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    label={i18n.t(`${packageNS}:tr000362`)}
                    helpText={i18n.t(`${packageNS}:tr000363`)}
                    component={ReactstrapInput}
                    className={
                      errors.object && errors.object.unicastTimeout
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  {
                    errors.object && errors.object.unicastTimeout
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          {errors.object.unicastTimeout}
                        </div>
                      ) : null
                  }

                  <Field
                    id="dr"
                    name="object.dr"
                    type="number"
                    value={values.object.dr}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    label="Data Rate"
                    helpText={i18n.t(`${packageNS}:tr000270`)}
                    component={ReactstrapInput}
                    className={
                      errors.object && errors.object.dr
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  {
                    errors.object && errors.object.dr
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          {errors.object.dr}
                        </div>
                      ) : null
                  }

                  <Field
                    id="frequency"
                    name="object.frequency"
                    type="number"
                    value={values.object.dr}
                    onChange={handleChange}
                    onBlur={handleBlur}
                    label={i18n.t(`${packageNS}:tr000271`)}
                    helpText={i18n.t(`${packageNS}:tr000272`)}
                    component={ReactstrapInput}
                    className={
                      errors.object && errors.object.frequency
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  {
                    errors.object && errors.object.frequency
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          {errors.object.frequency}
                        </div>
                      ) : null
                  }

                  <label htmlFor="object.groupType" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                    {i18n.t(`${packageNS}:tr000273`)}
                  </label>
                  <AutocompleteSelect
                    id="groupType"
                    name="object.groupType"
                    label={i18n.t(`${packageNS}:tr000274`)}
                    value={values.object.groupType}
                    onChange={handleChange}
                    getOptions={this.getGroupTypeOptions}
                    className={
                      errors.object && errors.object.redundancy
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  <FormHelperText>
                    {i18n.t(`${packageNS}:tr000275`)}
                  </FormHelperText>
                  {
                    errors.object && errors.object.groupType
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          <br />
                          {errors.object.groupType}
                        </div>
                      ) : null
                  }
                  <br />

                  <label htmlFor="object.multicastTimeout" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                    {i18n.t(`${packageNS}:tr000349`)}
                  </label>
                  <AutocompleteSelect
                    id="multicastTimeout"
                    name="object.multicastTimeout"
                    label={i18n.t(`${packageNS}:tr000361`)}
                    value={values.object.multicastTimeout}
                    onChange={handleChange}
                    getOptions={this.getMulticastTimeoutOptions}
                    className={
                      errors.object && errors.object.redundancy
                        ? 'is-invalid form-control'
                        : ''
                    }
                  />
                  {
                    errors.object && errors.object.multicastTimeout
                      ? (
                        <div
                          className="invalid-feedback"
                          style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                        >
                          <br />
                          {errors.object.multicastTimeout}
                        </div>
                      ) : null
                  }
                  <br />

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
                </Form>
              );
            }
          }
        </Formik>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(FUOTADeploymentForm);

