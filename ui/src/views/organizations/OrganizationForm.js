import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Button, Col, Row } from 'reactstrap';
import * as Yup from 'yup';
import { ReactstrapCheckbox, ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';




const orgSchema = Yup.object().shape({
  name: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
  displayName: Yup.string().required(i18n.t(`${packageNS}:tr000431`)), 
});


class OrganizationForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || {name: "", displayName: ""},
    };
  }


  render() {

    return (<React.Fragment>
      <Row>
        <Col>
          <Formik
              enableReinitialize
              initialValues={this.state.object}
              validationSchema={orgSchema}
              onSubmit={this.props.onSubmit}>
            {({
                handleSubmit,
                handleChange,
                handleBlur,
              }) => (
                <Form onSubmit={handleSubmit} noValidate>
                  <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000030`)+'*'}
                      name="name"
                      id="name"
                      helpText={i18n.t(`${packageNS}:tr000062`)}
                      component={ReactstrapInput}
                      onBlur={handleBlur}
                      onChange={handleChange}
                      required
                  />

                  <Field
                      type="text"
                      label={i18n.t(`${packageNS}:tr000126`)+'*'}
                      name="displayName"
                      id="displayName"
                      helpText={i18n.t(`${packageNS}:tr000031`)}
                      component={ReactstrapInput}
                      onBlur={handleBlur}
                      onChange={handleChange}
                      required
                  />

                  <h4>{i18n.t(`${packageNS}:tr000127`)}</h4>

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
