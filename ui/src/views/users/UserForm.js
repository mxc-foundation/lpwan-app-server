import React, { Component } from "react";
import { Button, FormGroup, Card } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';

import i18n, { packageNS } from '../../i18n';
import Admin from '../../components/Admin';
import { ReactstrapInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import defaultProfilePic from '../../assets/images/users/profile-icon.png';
import UserProfilePicFile from './UserProfilePicFile';

class UserForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      object: props.object || {},
    };
  }

  handleUploadedProfilePic = (output) => {
    const { result, successMessage, errorMessage } = output;

    if (errorMessage) {
      this.setState({
        errorMessageUploadingProfilePic: errorMessage
      });
    }

    this.setState({
      successMessageUploadingProfilePic: successMessage,
      uploadedProfilePic: result
    });
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
        // FIXME - get validation for email format to work
        // email: Yup.string().trim().matches(/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/, i18n.t(`${packageNS}:tr000455`))
        //   .required(i18n.t(`${packageNS}:tr000431`))
        email: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`))
      })
    }

    if (this.props.update) {
      fieldsSchema.object.fields.id = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("id");
    }

    if (!this.props.update) {
      fieldsSchema.object.fields.password = Yup.string().required(i18n.t(`${packageNS}:tr000431`));
      fieldsSchema.object._nodes.push("password");
    }

    return Yup.object().shape(fieldsSchema);
  }

  render() {
    const { uploadedProfilePic, errorMessageUploadingProfilePic, successMessageUploadingProfilePic } = this.state;

    const { object } = this.state;
    const { loading, update } = this.props;

    const isLoading = loading;

    if (object === undefined) {
      return null;
    }

    return (
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                id: object.id || undefined,
                profilePic: object.profilePic || uploadedProfilePic || defaultProfilePic,
                username: object.username || "",
                email: object.email || "",
                note: object.note || "",
                password: object.password || "",
                isAdmin: object.isAdmin || false,
                isActive: object.isActive || false
              }
            }
          }
          validateOnBlur
          validateOnChange
          validationSchema={this.formikFormSchema}
          // Formik Nested Schema Example https://codesandbox.io/s/y7q2v45xqx
          onSubmit={
            (values, { setSubmitting }) => {
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
                initialErrors,
                isSubmitting,
                isValidating,
                setFieldValue,
                touched,
                validateForm,
                values
              } = props;
              return ( 
                <Form style={{ padding: "0px", backgroundColor: "#ebeff2" }} onSubmit={handleSubmit} noValidate>
                  <Card body style={{ backgroundColor: "#fff" }}>
                    {isLoading && <Loader light />}
                    {this.props.update &&
                      <>
                        {/* <label htmlFor="object.id" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000077`)}
                        </label>
                        &nbsp;&nbsp;{values.object.id} */}

                        <input
                          type="hidden"
                          id="id"
                          name="object.id"
                          disabled
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
                      </>
                    }

                    <label htmlFor="object.profilePic" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                      {i18n.t(`${packageNS}:tr000454`)}
                    </label>
                    <UserProfilePicFile
                      profilePicImage={
                        <img
                          src={(object && object.profilePic) || uploadedProfilePic || defaultProfilePic}
                          className="rounded-circle"
                          alt="Profile Picture"
                          style={{ width: "100px", height: "100px" }}
                        />
                      }
                      onChange={this.handleUploadedProfilePic}
                    />
                    {successMessageUploadingProfilePic}
                    {errorMessageUploadingProfilePic}

                    <Field
                      id="profilePic"
                      name="object.profilePic"
                      type="hidden"
                      value={values.object.profilePic}
                      onChange={handleChange}
                      onBlur={handleBlur}
                      component={ReactstrapInput}
                      className={
                        errors.object && errors.object.profilePic
                          ? 'is-invalid form-control'
                          : ''
                      }
                    />
                    {
                      errors.object && errors.object.profilePic
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.profilePic}
                          </div>
                        ) : null
                    }

                    <Field
                      id="username"
                      name="object.username"
                      type="text"
                      value={values.object.username}
                      onChange={handleChange}
                      onBlur={handleBlur}
                      helpText="Username must only contain letters, and digits"
                      label={i18n.t(`${packageNS}:tr000056`)}
                      component={ReactstrapInput}
                      className={
                        errors.object && errors.object.username
                          ? 'is-invalid form-control'
                          : ''
                      }
                    />
                    {
                      errors.object && errors.object.username
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.username}
                          </div>
                        ) : null
                    }

                    <Field
                      id="email"
                      name="object.email"
                      type="email"
                      value={values.object.email}
                      onChange={handleChange}
                      onBlur={handleBlur}
                      label={i18n.t(`${packageNS}:tr000147`)}
                      component={ReactstrapInput}
                      className={
                        errors.object && errors.object.email
                          ? 'is-invalid form-control'
                          : ''
                      }
                    />
                    {
                      errors.object && errors.object.email
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.email}
                          </div>
                        ) : null
                    }

                    <Field
                      id="note"
                      name="object.note"
                      type="textarea"
                      multiline="true"
                      rows="4"
                      value={values.object.note}
                      onChange={handleChange}
                      onBlur={handleBlur}
                      label={i18n.t(`${packageNS}:tr000129`)}
                      helpText={i18n.t(`${packageNS}:tr000130`)}
                      component={ReactstrapInput}
                      className={
                        errors.object && errors.object.note
                          ? 'is-invalid form-control'
                          : ''
                      }
                    />
                    {
                      errors.object && errors.object.note
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.object.note}
                          </div>
                        ) : null
                    }

                    {!this.props.update &&
                      <>
                        <Field
                          id="password"
                          name="object.password"
                          type="password"
                          value={values.object.password}
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000004`)}
                          component={ReactstrapInput}
                          className={
                            errors.object && errors.object.password
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        {
                          errors.object && errors.object.password
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.object.password}
                              </div>
                            ) : null
                        }
                      </>
                    }

                    <Admin>
                      <FormGroup check>
                        <FormControlLabel
                          label={i18n.t(`${packageNS}:tr000133`)}
                          control={
                            <Checkbox
                              id="isAdmin"
                              name="object.isAdmin"
                              onChange={handleChange}
                              color="primary"
                              value={!!values.object.isAdmin}
                              checked={!!values.object.isAdmin}
                            />
                          }
                        />
                      </FormGroup>

                      <FormGroup check>
                        <FormControlLabel
                          label={i18n.t(`${packageNS}:tr000132`)}
                          control={
                            <Checkbox
                              id="isActive"
                              name="object.isActive"
                              onChange={handleChange}
                              color="primary"
                              value={!!values.object.isActive}
                              checked={!!values.object.isActive}
                            />
                          }
                        />
                      </FormGroup>
                    </Admin>

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
                      {/* {`${JSON.stringify(errors.object)}`} */}
                      {/* Show error count when page loads, before user submits the form */}
                      {errors.object && Object.keys(errors.object).length
                        ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                          Detected {Object.keys(errors.object).length} errors. Please fix the validation errors shown before resubmitting.
                        </div>)
                        : null
                      }

                      {/* Show error count when user submits the form */}
                      {this.state.validationErrors && this.state.validationErrors.length
                        ? (<div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                          Detected {Object.keys(this.state.validationErrors.object).length} errors. Please fix the validation errors shown before resubmitting.
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
                  </Card>
                </Form>
              );
            }
          }
        </Formik>
      </React.Fragment>
    );
  }
}

export default UserForm;
