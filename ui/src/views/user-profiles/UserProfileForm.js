import React from "react";
import {
  Button, Col, Container, FormFeedback, FormGroup, FormText, Input, Label, 
  TabContent, TabPane, Nav, NavItem, NavLink, Row,
} from 'reactstrap';
import { Form, Field } from "react-final-form";
import classnames from 'classnames';

import profilePic from '../../assets/images/users/profile-icon.png';

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";

const submitButton = (submitting, submitLabel) => {
  return (
    <Button
      aria-label={submitLabel}
      block
      color="primary"
      disabled={submitting}
      size="md"
      type="submit"
    >
      {submitLabel}
    </Button>
  );
};

class UserProfileForm extends FormComponent {
  constructor() {
    super();

    this.state = {
    };
  }

  render() {
    const { onSubmit, submitLabel } = this.props;

    if (this.state.object === undefined) {
      return(null);
    }

    return(
      <Form
        onSubmit={onSubmit}
        initialValues={{
          id: this.state.object.id,
          picture: this.state.object.picture,
          username: this.state.object.username,
          email: this.state.object.email,
          oldPassword: this.state.object.oldPassword,
          newPassword: this.state.object.newPassword,
          confirmPassword: this.state.object.confirmPassword
        }}
        validate={values => {
          console.log('validateForm values/activeTab: ', values);
          if (!values) {
            return {};
          }

          const errors = {};
        
          if (!values.username) {
            errors.username = "Required";
          }
        
          return errors;
        }}
        render={({ handleSubmit, form, submitting, pristine, values }) => (
          <form onSubmit={handleSubmit}>
            <Container>
              <Row>
                <Col sm="12">
                  <h4>{i18n.t(`${packageNS}:tr000431`)}</h4>
                  <br />
                  <FormGroup row>
                    <Field name="id">
                      {({ input, meta }) => (
                        <div>
                          <Input
                            {...input}
                            id="id"
                            name="id"
                            type="hidden"
                          />
                        </div>
                      )}
                    </Field>
                  </FormGroup>
                  <FormGroup row>
                    <Label for="name" sm={3}>
                      {i18n.t(`${packageNS}:tr000437`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="picture">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              autoFocus
                              id="picture"
                              name="picture"
                              type="text"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>
                    </Col>
                  </FormGroup>
                  <FormGroup row>
                    <Label for="name" sm={3}>
                      {i18n.t(`${packageNS}:tr000056`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="username">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              autoFocus
                              id="username"
                              name="username"
                              type="text"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>
                    </Col>
                  </FormGroup>
                  <FormGroup row>
                    <Label for="email" sm={3}>
                      {i18n.t(`${packageNS}:tr000147`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="email">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="email"
                              name="email"
                              type="email"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>  
                    </Col>
                  </FormGroup>
                  <br />
                  <h4>{i18n.t(`${packageNS}:tr000436`)}</h4>
                  <br />
                  <FormGroup row>
                    <Label for="oldPassword" sm={3}>
                      {i18n.t(`${packageNS}:tr000432`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="oldPassword">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="oldPassword"
                              name="oldPassword"
                              type="password"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>  
                    </Col>
                  </FormGroup>
                  <FormGroup row>
                    <Label for="newPassword" sm={3}>
                      {i18n.t(`${packageNS}:tr000433`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="newPassword">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="newPassword"
                              name="newPassword"
                              type="password"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>
                    </Col>
                  </FormGroup>
                  <FormGroup row>
                    <Label for="confirmPassword" sm={3}>
                      {i18n.t(`${packageNS}:tr000434`)}
                    </Label>
                    <Col sm={9}>
                      <Field name="confirmPassword">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="confirmPassword"
                              name="confirmPassword"
                              type="password"
                              invalid={meta.error && meta.touched}
                            />
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </div>
                        )}
                      </Field>
                    </Col>
                  </FormGroup>
                  <FormGroup row>
                    <Col sm="12">
                      <br />
                      {submitButton(submitting, submitLabel)}
                    </Col>
                  </FormGroup>
                </Col>
              </Row>
            </Container>
              
          </form>
        )}
      />
    );
  }
}

export default UserProfileForm;
