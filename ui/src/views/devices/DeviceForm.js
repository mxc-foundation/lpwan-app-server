import React, { Component } from "react";

import { withStyles } from "@material-ui/core/styles";
import { TabContent, TabPane, Nav, NavItem, NavLink, Card, Button, Row, Col } from 'reactstrap';
import { ErrorMessage, Formik, Form, Field, FieldArray, getIn } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import FormControl from "@material-ui/core/FormControl";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormHelperText from "@material-ui/core/FormHelperText";
import Checkbox from "@material-ui/core/Checkbox";
import FormGroup from "@material-ui/core/FormGroup";
import IconButton from '@material-ui/core/IconButton';
import Typography from "@material-ui/core/Typography";

import Delete from "mdi-material-ui/Delete";

import i18n, { packageNS } from '../../i18n';
import { ReactstrapInput, ReactstrapCheckbox } from '../../components/FormInputs';
import EUI64Field from "../../components/FormikEUI64Field";
import AutocompleteSelect from "../../components/AutocompleteSelect";
import Loader from "../../components/Loader";
import ApplicationStore from "../../stores/ApplicationStore";
import DeviceProfileStore from "../../stores/DeviceProfileStore";

import theme from "../../theme";

const clone = require('rfdc')();

const styles = {
  formLabel: {
    fontSize: 12,
  },
  delete: {
    marginTop: 3 * theme.spacing(1),
  },
};


class DeviceForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || {},
      activeTab: "1",
      loading: true,
    };
  }

  componentDidMount() {
    // New Device
    if (!this.props.object) {
      return;
    }
    this.setKVArrayVariables();
    this.setKVArrayTags();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.object !== this.props.object) {
      this.setKVArrayVariables();
      this.setKVArrayTags();
    }
  }

  // Storage has the 'variables' and 'tags' stored as follows:
  // variables: { my_var_key1: "my var value1", my_var_key2: "my var value2" }
  //
  // But we're leveraging FormikArray, so locally we're converting it into format:
  // variables: [ { key: "my_var_key1", value: "my var value1" }, { key: "my_var_key2", value: "my var value2" } ]
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

  convertArrayToObj = (arr, key) => {
    const formatKey = (k) => k.trim().split(' ').join('_');

    let asObject = {};
    for (const el of arr.object[key]) {
      if (el.key !== "") {
        asObject[formatKey(el.key)] = el.value;
      }
    };

    return asObject;
  }

  setKVArrayVariables = () => {
    if (this.props.object && this.props.object.variables.length === 0) {
      return;
    }

    const propAsArray = this.convertObjToArray(this.props.object.variables);

    this.setState(prevState => {
      if (prevState.object && prevState.object.variables.length === 0) {
        return;
      }

      // Obtain the existing variables that are already in the local state
      let existingStateVariables = prevState.object.variables;
      let existingStateVariablesAsArray = this.convertObjToArray(existingStateVariables);

      // Retrieve the variables array passed as props from the parent component
      let propVariables = propAsArray; //this.props.object.variables;

      // Iterate through the key value pairs
      let updatedVariables = propVariables.map(
        el => {
          let resObj = existingStateVariablesAsArray.find(x => x.key === el.key);
          const resIndex = existingStateVariablesAsArray.indexOf(resObj);
  
          // Assuming that all keys are unique. If the current key passed from props
          // is not already in state, then we want to add that new element key value pair to state,
          // otherwise update the value of that key if the key exists in state already.
          if (resIndex === -1) {
            return el;
          // Otherwise retain existing state key value pair
          } else {
            resObj.value = el.value;
            return resObj;
          }
        }
      )

      return {
        object: {
          ...prevState.object,
          variables: updatedVariables
        }
      }
    })
  }

  setKVArrayTags = () => {
    if (this.props.object !== undefined && this.props.object.tags.length === 0) {
      return;
    }

    const propAsArray = this.convertObjToArray(this.props.object.tags);

    this.setState(prevState => {
      if (prevState.object !== undefined && prevState.object.tags.length === 0) {
        return;
      }

      // Obtain the existing tags that are already in the local state
      let existingStateTags = prevState.object.tags;
      let existingStateTagsAsArray = this.convertObjToArray(existingStateTags);

      // Retrieve the tags array passed as props from the parent component
      let propTags = propAsArray; // this.props.object.tags;

      // Iterate through the key value pairs
      let updatedTags = propTags.map(
        el => {
          let resObj = existingStateTagsAsArray.find(x => x.key === el.key);
          const resIndex = existingStateTagsAsArray.indexOf(resObj);
  
          // Assuming that all keys are unique. If the current key passed from props
          // is not already in state, then we want to add that new element key value pair to state,
          // otherwise update the value of that key if the key exists in state already.
          if (resIndex === -1) {
            return el;
          // Otherwise retain existing state key value pair
          } else {
            resObj.value = el.value;
            return resObj;
          }
        }
      )

      return {
        object: {
          ...prevState.object,
          tags: updatedTags
        }
      }
    })
  }

  getApplicationOption = (id, callbackFunc) => {
    ApplicationStore.get(id, resp => {
      this.setState({
        loading: false
      })
      callbackFunc({label: resp.application.name, value: resp.application.id});
    });
  }

  getApplicationOptions = (search, callbackFunc) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    ApplicationStore.list("", currentOrgID, 999, 0, resp => {
      const options = resp.result.map((app, i) => {return {label: app.name, value: app.id}});
      this.setState({
        loading: false
      })
      callbackFunc(options);
    });
  }

  getDeviceProfileOption = (id, callbackFunc) => {
    DeviceProfileStore.get(id, resp => {
      this.setState({
        loading: false
      })
      callbackFunc({label: resp.deviceProfile.name, value: resp.deviceProfile.id});
    });
  }

  getDeviceProfileOptions = (search, callbackFunc) => {
    DeviceProfileStore.list(0, this.props.match.params.applicationID, 999, 0, resp => {
      const options = resp.result.map((dp, i) => {return {label: dp.name, value: dp.id}});
      this.setState({
        loading: false
      })
      callbackFunc(options);
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

  formikFormSchema = () => {
    let fieldsSchema = {
      object: Yup.object().shape({
        // https://regexr.com/4rg3a
        name: Yup.string().trim().matches(/^[0-9A-Za-z-]*$/g, i18n.t(`${packageNS}:tr000429`))
          .required(i18n.t(`${packageNS}:tr000431`)),
        description: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`)),
        deviceProfileID: Yup.string()
          .required(i18n.t(`${packageNS}:tr000431`))
      })
    }

    return Yup.object().shape(fieldsSchema);
  }

  render() {
    const { activeTab, loading: loadingState, object } = this.state;
    const { classes, loading: loadingProps } = this.props;
    const isLoading = (loadingState || loadingProps);

    if (object === undefined) {
      return null;
    }

    return(
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={
            {
              object: {
                name: object.name || undefined,
                description: object.description || undefined,
                devEUI: object.devEUI || undefined,
                applicationID: object.applicationID || undefined,
                deviceProfileID: object.deviceProfileID || undefined,
                skipFCntCheck: !!object.skipFCntCheck || false,
                variables: (
                  (object.variables !== undefined && object.variables.length > 0 && object.variables) || []
                ),
                tags: (
                  (object.tags !== undefined && object.tags.length > 0 && object.tags) || []
                )
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

              // Deep copy is required otherwise we can change the original values of
              // 'variables' and 'tags' (and we will not be able to render the different format in the UI)
              // Reference: https://medium.com/javascript-in-plain-english/how-to-deep-copy-objects-and-arrays-in-javascript-7c911359b089
              let newValues = clone(values);
              console.log('Deep copied submitted values: ', newValues !== values);
              let variablesAsObject;
              let tagsAsObject;
              if (Array.isArray(values.object.variables)) {
                variablesAsObject = this.convertArrayToObj(values, "variables");
                newValues.object.variables = variablesAsObject;
              }

              if (Array.isArray(values.object.tags)) {
                tagsAsObject = this.convertArrayToObj(values, "tags");
                newValues.object.tags = tagsAsObject;
              }

              console.log('Prepared values', newValues);

              // return;
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
                isSubmitting,
                isValidating,
                setFieldValue,
                touched,
                validateForm,
                values
              } = props;
              return (  
                <Form style={{ padding: "0px", backgroundColor: "#ebeff2" }} onSubmit={handleSubmit} noValidate>
                  {/* <Card body style={{ backgroundColor: "#ebeff2" }}> */}
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
                          <i className="mdi mdi-code-braces"></i>
                          &nbsp;{i18n.t(`${packageNS}:tr000305`)}
                        </NavLink>
                      </NavItem>
                      <NavItem>
                        <NavLink
                          className={classnames({ active: activeTab === '3' })}
                          onClick={() => { this.toggle('3'); }}
                        >
                          <i className="mdi mdi-tag-multiple"></i>
                          &nbsp;{i18n.t(`${packageNS}:tr000308`)}
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

                        <Field
                          id="name"
                          name="object.name"
                          type="text"
                          value={values.object.name}
                          onChange={handleChange}
                          onBlur={handleBlur}
                          label={i18n.t(`${packageNS}:tr000300`)}
                          helpText={i18n.t(`${packageNS}:tr000062`)}
                          component={ReactstrapInput}
                          // FIXME - to show form validation errors this approach isn't usually necessary
                          // but they aren't appearing automatically so i've had to do it manually
                          className={
                            errors.object && errors.object.name
                            // && touched.object && touched.object.name
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />

                        {/* FIXME - to show form validation errors this approach isn't usually necessary
                            but they aren't appearing automatically so i've had to do it manually
                        */}
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
                        <br />

                        <Field
                          id="description"
                          name="object.description"
                          type="text"
                          value={values.object.description}
                          label={i18n.t(`${packageNS}:tr000301`)}
                          component={ReactstrapInput}
                          className={
                            errors.object && errors.object.description
                            // && touched.object && touched.object.description
                              ? 'is-invalid form-control'
                              : ''
                          }
                        />

                        {
                          errors.object && errors.object.description
                            ? (
                              <div
                                className="invalid-feedback"
                                style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}
                              >
                                {errors.object.description}
                              </div>
                            ) : null
                        }
                        <br />

                        {!this.props.update && 
                          <EUI64Field
                            id="devEUI"
                            name="object.devEUI"
                            label={i18n.t(`${packageNS}:tr000371`)}
                            random
                          />
                        }

                        <label htmlFor="object.applicationID" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000407`)}
                        </label>
                        <AutocompleteSelect
                          id="applicationID"
                          name="object.applicationID"
                          label={i18n.t(`${packageNS}:tr000407`)}
                          // FIXME - show loading until an option is available
                          onChange={handleChange}
                          getOption={this.getApplicationOption}
                          getOptions={this.getApplicationOptions}
                        />

                        <label htmlFor="object.deviceProfileID" style={{ display: 'block', fontWeight: "700", marginTop: 16 }}>
                          {i18n.t(`${packageNS}:tr000281`)}
                        </label>
                        <AutocompleteSelect
                          id="deviceProfileID"
                          name="object.deviceProfileID"
                          label={i18n.t(`${packageNS}:tr000281`)}
                          // FIXME - show loading until an option is available
                          onChange={handleChange}
                          getOption={this.getDeviceProfileOption}
                          getOptions={this.getDeviceProfileOptions}
                        />

                        <div style={{ marginTop: "10px" }}>
                          <FormGroup>
                            <FormControlLabel
                              label={i18n.t(`${packageNS}:tr000303`)}
                              control={
                                <Checkbox
                                  id="skipFCntCheck"
                                  name="object.skipFCntCheck"
                                  onChange={handleChange}
                                  color="primary"
                                />
                              }
                            />
                          </FormGroup>
                          <FormHelperText>
                            {i18n.t(`${packageNS}:tr000304`)}
                          </FormHelperText>
                        </div>

                        {/* FIXME - unable to click this checkbox for some reason when try to implement it */}
                        {/* <Field
                          type="checkbox"
                          label={i18n.t(`${packageNS}:tr000303`)}
                          id="object.skipFCntCheck"
                          name="object.skipFCntCheck"
                          onChange={handleChange}
                          component={ReactstrapCheckbox}
                        /> */}

                      </TabPane>
                      <TabPane tabId="2">
                        <Typography variant="body1">
                          {i18n.t(`${packageNS}:tr000306`)}
                        </Typography>
                        <br />

                        {/* TODO - we could refactor the 'variables' and 'tags' FieldArrays into a subcomponent
                          since the only thing that changes is the key, but it may make using Formik more complex
                        */}
                        <FieldArray
                          id="variables"
                          name="object.variables"
                          value={values.object.variables}
                          render={arrayHelpers => (
                            <div>
                              {/* { JSON.stringify(values.object) } */}
                              {
                                values.object && values.object.variables !== undefined && 
                                values.object.variables.length > 0 &&
                                values.object.variables.map((variable, index) => (
                                  variable && Object.keys(variable).length == 2 ? (
                                    <div key={index}>
                                      {/* Debug Row */}
                                      {/* <Row>
                                        <Col xs={4} md={4}>
                                          { JSON.stringify(variable) }
                                        </Col>
                                      </Row> */}
                                      <Row>
                                        <Col xs={4} md={4}>
                                          <Field
                                            type="text"
                                            id={`variables[${index}].key`}
                                            name={`object.variables[${index}].key`}
                                            label={i18n.t(`${packageNS}:tr000042`)}
                                            value={variable.key}
                                            onChange={handleChange}
                                            component={ReactstrapInput}
                                          />
                                        </Col>
                                        <Col xs={5} md={7}>
                                          <Field
                                            type="text"
                                            id={`variables[${index}].value`}
                                            name={`object.variables[${index}].value`}
                                            label="Value"
                                            value={variable.value}
                                            onChange={handleChange}
                                            component={ReactstrapInput}
                                          />
                                        </Col>
                                        <Col xs={3} md={1} className={classes.delete}>
                                          <IconButton aria-label="delete" onClick={() => arrayHelpers.remove(index)}>
                                            <Delete />
                                          </IconButton>
                                        </Col>
                                      </Row>
                                    </div>
                                  ) : <div key={index}></div>
                                ))}
                              <Button
                                type="button"
                                variant="outlined"
                                onClick={() => arrayHelpers.push({ key: '', value: '' })}
                              >
                                {i18n.t(`${packageNS}:tr000307`)}
                              </Button>
                            </div>
                          )}
                        />
                      </TabPane>
                      <TabPane tabId="3">
                        <Typography variant="body1">
                          {i18n.t(`${packageNS}:tr000309`)}
                        </Typography>
                        <br />

                        <FieldArray
                          id="tags"
                          name="object.tags"
                          value={values.object.tags}
                          render={arrayHelpers => (
                            <div>
                              {
                                values.object && values.object.tags !== undefined && 
                                values.object.tags.length > 0 &&
                                values.object.tags.map((tag, index) => (
                                  tag && Object.keys(tag).length == 2 ? (
                                    <div key={index}>
                                      {/* Debug Row */}
                                      {/* <Row>
                                        <Col xs={4} md={4}>
                                          { JSON.stringify(tag) }
                                        </Col>
                                      </Row> */}
                                      <Row>
                                        <Col xs={4} md={4}>
                                          <Field
                                            type="text"
                                            id={`tags[${index}].key`}
                                            name={`object.tags[${index}].key`}
                                            label={i18n.t(`${packageNS}:tr000042`)}
                                            value={tag.key}
                                            onChange={handleChange}
                                            component={ReactstrapInput}
                                          />
                                        </Col>
                                        <Col xs={5} md={7}>
                                          <Field
                                            type="text"
                                            id={`tags[${index}].value`}
                                            name={`object.tags[${index}].value`}
                                            label="Value"
                                            value={tag.value}
                                            onChange={handleChange}
                                            component={ReactstrapInput}
                                          />
                                        </Col>
                                        <Col xs={3} md={1} className={classes.delete}>
                                          <IconButton aria-label="delete" onClick={() => arrayHelpers.remove(index)}>
                                            <Delete />
                                          </IconButton>
                                        </Col>
                                      </Row>
                                    </div>
                                  ) : <div key={index}></div>
                                ))}
                              <Button
                                type="button"
                                variant="outlined"
                                onClick={() => arrayHelpers.push({ key: '', value: '' })}
                              >
                                {i18n.t(`${packageNS}:tr000307`)}
                              </Button>
                            </div>
                          )}
                        />
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
                      {errors.object
                        ? <div style={{ display: "block", color: "#ff5b5b", fontSize: "0.75rem", marginTop: "-0.75rem" }}>
                            Form Validation Errors. Please enter valid inputs and try again...
                          </div>
                        : ''
                      }
                    </div>
                    {/* <Button
                      type="button"
                      color="secondary"
                      onClick={handleReset}
                      disabled={!dirty || isSubmitting}
                    >
                      Reset
                    </Button> */}
                    <Button
                      type="submit"
                      color="primary"
                      disabled={(errors.object && Object.keys(errors.object).length > 0) || isLoading || isSubmitting}
                      onClick={
                        () => validateForm().then((formValidationErrors) =>
                          console.log('Validated form with errors: ', formValidationErrors))
                      }
                    >
                      {this.props.submitLabel || i18n.t(`${packageNS}:tr000292`)}
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

export default withStyles(styles)(DeviceForm);
