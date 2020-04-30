import classNames from 'classnames';
import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Button, Card, CardBody, Col, Nav, NavItem, Row } from 'reactstrap';
import * as Yup from 'yup';
import { AsyncAutoComplete, ReactstrapCheckbox, ReactstrapInput, ReactstrapPasswordInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import OrgBreadCumb from '../../components/OrgBreadcrumb';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import OrganizationStore from "../../stores/OrganizationStore";
import SessionStore from "../../stores/SessionStore";
import UserStore from "../../stores/UserStore";



class AssignUserForm extends Component {
  constructor() {
    super();
    // we need combo box
    // this.getUserOption = this.getUserOption.bind(this);
    this.getUserOptions = this.getUserOptions.bind(this);
  }

  getUserOptions(search, callbackFunc) {
    UserStore.list(search, 99999999, 0, resp => {
      const options = resp.result.map((u, i) => {return {label: u.username, value: u.id}});
      callbackFunc(options);
    });
  }

  render() {

    const fieldsSchema = Yup.object().shape({
      userID: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`)),
        isAdmin: Yup.bool(),
        isDeviceAdmin: Yup.bool(),
        isGatewayAdmin: Yup.bool(),
    });
  

    return(
      <React.Fragment>
        <Row>
          <Col>
            <Formik
              enableReinitialize
              initialValues={{}}
              validateOnBlur
              validateOnChange
              validationSchema={fieldsSchema}
              onSubmit={
                (values, { setSubmitting }) => {
                  let newValues = {...values};
                  if(newValues.isAdmin) {
                    newValues.isGatewayAdmin = false;
                    newValues.isDeviceAdmin = false;
                  }
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
                      <Field
                            id="userID"
                            name="userID"
                            type="text"
                            value={values.userID}
                            onBlur={handleBlur}
                            label={i18n.t(`${packageNS}:tr000056`) + ' *'}
                            helpText={i18n.t(`${packageNS}:tr000138`)}
                            getOptions={this.getUserOptions}
                            setFieldValue={(field, val) => setFieldValue('userID', val)}
                            inputProps={{
                              clearable: true,
                              cache: false,
                              classNamePrefix: 'react-select-validation'
                            }}
                            component={AsyncAutoComplete}
                            className={
                              errors && errors.userID && touched && touched && touched['react-select-userID-input']
                                ? 'is-invalid'
                                : ''
                            }
                          />
                          {
                            errors && errors.userID && touched && touched && touched['react-select-userID-input']
                              ? (
                                <div
                                  className="invalid-feedback"
                                  style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                                >
                                  {errors.userID}
                                </div>
                              ) : null
                          }
                        
                      <Field
                        id="isAdmin"
                        name="isAdmin"
                        type="checkbox"
                        value={values.isAdmin}
                        label={i18n.t(`${packageNS}:tr000139`)}
                        helpText={i18n.t(`${packageNS}:tr000140`)}
                        component={ReactstrapCheckbox}
                        onChange={handleChange}
                        className={
                          errors && errors.isAdmin
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.isAdmin
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.isAdmin}
                            </div>
                          ) : null
                      }

                      {!!!values.isAdmin ? <React.Fragment>
                        <Field
                          id="isDeviceAdmin"
                          name="isDeviceAdmin"
                          type="checkbox"
                          value={values.isDeviceAdmin}
                          label={i18n.t(`${packageNS}:tr000141`)}
                          helpText={i18n.t(`${packageNS}:tr000142`)}
                          onChange={handleChange}
                          component={ReactstrapCheckbox}
                          className={
                            errors && errors.isDeviceAdmin
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        {
                          errors && errors.isDeviceAdmin
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.isDeviceAdmin}
                              </div>
                            ) : null
                        }

                        <Field
                          id="isGatewayAdmin"
                          name="isGatewayAdmin"
                          type="checkbox"
                          value={values.isGatewayAdmin}
                          onChange={handleChange}
                          label={i18n.t(`${packageNS}:tr000143`)}
                          helpText={i18n.t(`${packageNS}:tr000144`)}
                          component={ReactstrapCheckbox}
                          className={
                            errors && errors.isGatewayAdmin
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />
                        {
                          errors && errors.isGatewayAdmin
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.isGatewayAdmin}
                              </div>
                            ) : null
                        }
                      </React.Fragment>: null}
                        
                      <Button type="submit" color="primary" className="btn-block"
                        disabled={(errors && errors.userID) || !values.userID}>
                        {i18n.t(`${packageNS}:tr000041`)}</Button>
                    </Form>
                  );
                }
              }
            </Formik>
          </Col>
        </Row>
      </React.Fragment>
    );
  };
}


class CreateUserForm extends Component {
  render() {
    const fieldsSchema = Yup.object().shape({
      username: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`)),
      email: Yup.string()
          .email(i18n.t(`${packageNS}:tr000431`))
          .required(i18n.t(`${packageNS}:tr000431`)),
      note: Yup.string(),
      password: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
      isAdmin: Yup.bool(),
      isDeviceAdmin: Yup.bool(),
      isGatewayAdmin: Yup.bool(),
    });

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
            enableReinitialize
            initialValues={{}}
            validateOnBlur
            validateOnChange
            validationSchema={fieldsSchema}
            onSubmit={
              (values, { setSubmitting }) => {
                let newValues = { ...values };
                if (newValues.isAdmin) {
                  newValues.isGatewayAdmin = false;
                  newValues.isDeviceAdmin = false;
                }
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

                    <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000056`) + ' *'}
                      name="username"
                      id="username"
                      component={ReactstrapInput}
                      onBlur={handleBlur}
                    />

                    <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000147`) + ' *'}
                      name="email"
                      id="email"
                      component={ReactstrapInput}
                      onBlur={handleBlur}
                    />

                    <Field
                      type="textarea"
                      label={i18n.t(`${packageNS}:tr000129`)}
                      name="note"
                      id="note"
                      helpText={i18n.t(`${packageNS}:tr000130`)}
                      component={ReactstrapInput}
                      onBlur={handleBlur}
                    />

                    <Field
                      label={i18n.t(`${packageNS}:tr000004`) + ' *'}
                      name="password"
                      id="password"
                      helpText={i18n.t(`${packageNS}:tr000138`)}
                      component={ReactstrapPasswordInput}
                      onBlur={handleBlur}
                    />

                    <Field
                      id="isAdmin"
                      name="isAdmin"
                      type="checkbox"
                      value={values.isAdmin}
                      label={i18n.t(`${packageNS}:tr000139`)}
                      helpText={i18n.t(`${packageNS}:tr000140`)}
                      component={ReactstrapCheckbox}
                      onChange={handleChange}
                      className={
                        errors && errors.isAdmin
                          ? 'is-invalid form-control'
                          : ''
                      }
                    />
                    {
                      errors && errors.isAdmin
                        ? (
                          <div
                            className="invalid-feedback"
                            style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                          >
                            {errors.isAdmin}
                          </div>
                        ) : null
                    }

                    {!!!values.isAdmin ? <React.Fragment>
                      <Field
                        id="isDeviceAdmin"
                        name="isDeviceAdmin"
                        type="checkbox"
                        value={values.isDeviceAdmin}
                        label={i18n.t(`${packageNS}:tr000141`)}
                        helpText={i18n.t(`${packageNS}:tr000142`)}
                        onChange={handleChange}
                        component={ReactstrapCheckbox}
                        className={
                          errors && errors.isDeviceAdmin
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.isDeviceAdmin
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.isDeviceAdmin}
                            </div>
                          ) : null
                      }

                      <Field
                        id="isGatewayAdmin"
                        name="isGatewayAdmin"
                        type="checkbox"
                        value={values.isGatewayAdmin}
                        onChange={handleChange}
                        label={i18n.t(`${packageNS}:tr000143`)}
                        helpText={i18n.t(`${packageNS}:tr000144`)}
                        component={ReactstrapCheckbox}
                        className={
                          errors && errors.isGatewayAdmin
                            ? 'is-invalid form-control'
                            : ''
                        }
                      />
                      {
                        errors && errors.isGatewayAdmin
                          ? (
                            <div
                              className="invalid-feedback"
                              style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                            >
                              {errors.isGatewayAdmin}
                            </div>
                          ) : null
                      }
                    </React.Fragment> : null}
                    
                    <Button type="submit" color="primary" className="btn-block" 
                      disabled={(errors && Object.keys(errors).length)|| (values && (!values.username || !values.email || !values.password))}>
                        {i18n.t(`${packageNS}:tr000277`)}</Button>
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


class CreateOrganizationUser extends Component {
  constructor() {
    super();

    this.state = {
      tab: 0,
      assignUser: false,
      loading: false
    };

    this.onChangeTab = this.onChangeTab.bind(this);
    this.onAssignUser = this.onAssignUser.bind(this);
    this.onCreateUser = this.onCreateUser.bind(this);
    this.setAssignUser = this.setAssignUser.bind(this);
  }

  componentDidMount() {
    this.setAssignUser();

    SessionStore.on("change", this.setAssignUser);
  }

  comomentWillUnmount() {
    SessionStore.removeListener("change", this.setAssignUser);
  }

  setAssignUser() {
    const settings = SessionStore.getSettings();
    this.setState({
      assignUser: !settings.disableAssignExistingUsers || SessionStore.isAdmin(),
    });
  }

  onChangeTab(v) {
    this.setState({
      tab: v,
    });
  }

  onAssignUser = async (user) => {
    this.setState({loading: true});
    const res = await OrganizationStore.addUser(this.props.match.params.organizationID, user);
    this.setState({loading: false});
    this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users`);
  };

  onCreateUser(user) {
    const orgs = [
      {isAdmin: user.isAdmin, isDeviceAdmin: user.isDeviceAdmin, isGatewayAdmin: user.isGatewayAdmin, organizationID: this.props.match.params.organizationID},
    ];

    let u = user;
    u.isActive = true;

    delete u.isAdmin;
    delete u.isDeviceAdmin;
    delete u.isGatewayAdmin;

    this.setState({loading: true});
    // on success or error handling loading
    UserStore.create({user: u, password: user.password, organizations: orgs}, resp => {
      this.setState({loading: false});
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/users`);
    }, error => {
      this.setState({loading: false});
    });
  };

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <React.Fragment>
        <TitleBar>
          <OrgBreadCumb organizationID={currentOrgID} items={[
            { label: i18n.t(`${packageNS}:tr000068`), active: false, to: `/organizations/${currentOrgID}/users` },
            { label: i18n.t(`${packageNS}:tr000277`), active: true }]}></OrgBreadCumb>
        </TitleBar>

        <Row>
          <Col>
            <Card>
              <CardBody>
                <div className="position-relative">
                  {this.state.loading ? <Loader /> : null}
                  <Nav tabs>
                    {this.state.assignUser ? <NavItem>
                      <Link
                        className={classNames('nav-link', { active: this.state.tab === 0 })}
                        onClick={() => this.onChangeTab(0)}
                        to='#'>{i18n.t(`${packageNS}:tr000136`)}</Link>
                    </NavItem> : null}

                    <NavItem>
                      <Link
                        className={classNames('nav-link', { active: this.state.tab === 1 })}
                        onClick={() => this.onChangeTab(1)}
                        to='#'
                      >{i18n.t(`${packageNS}:tr000146`)}</Link>
                    </NavItem>
                  </Nav>

                  <Row className="pt-2">
                    <Col>
                      {(this.state.tab === 0 && this.state.assignUser) && <AssignUserForm onSubmit={this.onAssignUser} />}
                      {(this.state.tab === 1 || !this.state.assignUser) && <CreateUserForm onSubmit={this.onCreateUser} />}
                    </Col>
                  </Row>
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>

      </React.Fragment>
    );
  }
}

export default withRouter(CreateOrganizationUser);
