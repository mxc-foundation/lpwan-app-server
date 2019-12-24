import React, { Component } from "react";

import { Row, Col, Button, FormGroup, Label, FormText, Card, CardBody } from 'reactstrap';
import { Formik, Form, Field, FieldArray } from 'formik';
import * as Yup from 'yup';

import { ReactstrapInput, ReactstrapCheckbox, AsyncAutoComplete } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';
import TitleBarTitle from "../../components/TitleBarTitle";



class OrganizationForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || {},
    };
  }


  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
              enableReinitialize
              initialValues={this.state.object}
              onSubmit={this.props.onSubmit}>
            {({
                handleSubmit,
                handleChange,
              }) => (
                <Form onSubmit={handleSubmit} noValidate>
                  <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000030`)+'*'}
                      name="name"
                      id="name"
                      value={this.state.object.name || ""}
                      helpText={i18n.t(`${packageNS}:tr000062`)}
                      component={ReactstrapInput}
                      required
                  />

                  <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000126`)+'*'}
                      name="displayName"
                      id="displayName"
                      value={this.state.object.displayName || ""}
                      helpText={i18n.t(`${packageNS}:tr000031`)}
                      component={ReactstrapInput}
                      required
                  />

                  <TitleBarTitle title={i18n.t(`${packageNS}:tr000127`)} />

                  <Field
                      type="checkbox"
                      label={i18n.t(`${packageNS}:tr000064`)}
                      name="canHaveGateways"
                      id="canHaveGateways"
                      helpText={i18n.t(`${packageNS}:tr000065`)}
                      component={ReactstrapCheckbox}
                      onChange={handleChange}
                  />

                  <Button type="submit" color="primary">{this.props.submitLabel}</Button>
                </Form>
            )}
          </Formik>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default OrganizationForm;
