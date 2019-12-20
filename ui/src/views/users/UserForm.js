import React from "react";
import { Button, Col, Container, FormFeedback, FormGroup, FormText, Input, Label, Row } from 'reactstrap';
import { Form, Field } from "react-final-form";

import FormComponent from "../../classes/FormComponent";
import i18n, { packageNS } from '../../i18n';
import defaultProfilePic from '../../assets/images/users/profile-icon.png';
import UserProfilePicFile from './UserProfilePicFile';

const submitButton = (submitting, submitLabel) => {
  return (
    <Button
      aria-label={submitLabel}
      block
      color="primary"
      disabled={submitting}
      size="md"
      style={{ marginTop: "0px" }}
      type="submit"
    >
      {submitLabel}
    </Button>
  );
};


class UserForm extends FormComponent {
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

  render() {
    const { uploadedProfilePic, errorMessageUploadingProfilePic, successMessageUploadingProfilePic } = this.state;
    const { onSubmit, submitLabel, object } = this.props;
    // Obtain userID from URL parameters if parent component does not provide it via props
    const user = object || this.props.match.params.userID;

    if (user === undefined) {
      return (<div></div>);
    }

    return (
      <Form
        onSubmit={onSubmit}
        initialValues={{
          id: user.id,
          profilePic: user.profilePic || uploadedProfilePic,
          username: user.username,
          email: user.email,
          note: user.note,
          password: user.password,
          isAdmin: user.isAdmin
        }}
        validate={values => {
          console.log('Validation: ', values);
          if (!values) {
            return {};
          }
          const errors = {};
          // Validate Email Address
          const validEmailFormat = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
          const isValidEmail = value =>
            value.match(validEmailFormat) ? true : false;

          if (values.email && !isValidEmail(values.email)) {
            errors.email = i18n.t(`${packageNS}:tr000455`);
          }
        
          return errors;
        }}

        render={({ handleSubmit, form, submitting, pristine, values }) => (
          <form onSubmit={handleSubmit}>
            <Container>
              <Row>
                <Col sm="12">
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
                  <h4>{i18n.t(`${packageNS}:tr000452`)}</h4>
                  <br />
                  <FormGroup row>
                    <Label for="profilePic" sm={3}>
                      {i18n.t(`${packageNS}:tr000454`)}
                    </Label>
                    <Col sm={9}>
                      <Row style={{ marginBottom: "10px" }}>
                        <UserProfilePicFile
                          profilePicImage={
                            <img
                              src={(user && user.profilePic) || uploadedProfilePic || defaultProfilePic}
                              className="rounded-circle"
                              alt="Profile Picture"
                              style={{ width: "100px", height: "100px" }}
                            />
                          }
                          onChange={this.handleUploadedProfilePic}
                        />
                        {successMessageUploadingProfilePic}
                        {errorMessageUploadingProfilePic}
                      </Row>
                      <Field name="profilePic">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="profilePic"
                              name="profilePic"
                              type="hidden"
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
                    <Label for="username" sm={3}>
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
                  <FormGroup row>
                    <Label for="caCert" sm={3}>
                      {i18n.t(`${packageNS}:tr000129`)}
                    </Label>
                    <Col sm={9}>
                      <Input
                        id="note"
                        name="note"
                        multiline="true"
                        onChange={this.onChange}
                        rows="4"
                        type="textarea"
                      />
                      <FormText color="muted">
                        {i18n.t(`${packageNS}:tr000130`)}
                      </FormText>
                    </Col>
                  </FormGroup>
                  {
                    user.id === undefined &&
                      <>
                        <br />
                        <FormGroup row>
                          <Label for="password" sm={3}>
                            {i18n.t(`${packageNS}:tr000004`)}
                          </Label>
                          <Col sm={9}>
                            <Field name="password">
                              {({ input, meta }) => (
                                <div>
                                  <Input
                                    {...input}
                                    id="password"
                                    name="password"
                                    type="password"
                                    invalid={meta.error && meta.touched}
                                  />
                                  {meta.error && meta.touched &&
                                    <FormFeedback>{meta.error}</FormFeedback>
                                  }
                                </div>
                              )}
                            </Field>
                            <FormText color="muted">
                              {i18n.t(`${packageNS}:tr000130`)}
                            </FormText>
                          </Col>
                        </FormGroup>
                      </>
                  }
                  <FormGroup check>
                    <Field name="isAdmin" type="checkbox">
                      {({ input, meta }) => (
                        <Label check for="isAdmin">
                          <Input
                            {...input}
                            id="isAdmin"
                            name="isAdmin"
                            type="checkbox"
                            invalid={meta.error && meta.touched}
                            onClick={this.onChange}
                          />
                          {' '}
                          {i18n.t(`${packageNS}:tr000133`)}
                          {meta.error && meta.touched &&
                            <FormFeedback>{meta.error}</FormFeedback>
                          }
                        </Label>
                      )}
                    </Field>
                  </FormGroup>
                  {
                    submitLabel ? (
                      <FormGroup row>
                        <Col sm="12">
                          <br />
                          {submitButton(submitting, submitLabel)}
                        </Col>
                      </FormGroup>
                    ) : null
                  }
                </Col>
              </Row>
            </Container>
          </form>
        )}
      />
    );
  }
}

export default UserForm;
