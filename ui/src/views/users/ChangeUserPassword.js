import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Formik, Form, Field, FieldArray } from 'formik';
import { ReactstrapInput } from '../../components/FormInputs';
import * as Yup from 'yup';
import { Breadcrumb, BreadcrumbItem, FormGroup, Label, Input, Button, Container, Row, Col, Card, CardBody } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";

import TitleBar from "../../components/TitleBar";
import UserStore from "../../stores/UserStore";
import FormComponent from "../../classes/FormComponent";
import i18n, { packageNS } from '../../i18n';

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class PasswordForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }
    const object = this.state.object;
    const id = this.props.userId;

    let fieldsSchema = {
      id: Yup.string().trim(),
      password: Yup.string().trim().matches(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[/\W/])[A-Za-z\d/\W/]{8,}$/g, i18n.t(`${packageNS}:menu.messages.format_unmatch`)).required(i18n.t(`${packageNS}:tr000431`)),
      passwordConfirmation: Yup.string().oneOf([Yup.ref('password'), null], i18n.t(`${packageNS}:menu.registration.confirm_password_match_error`))
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <Formik
        enableReinitialize
        initialValues={{
          id: id,
          password: object.password || '',
          passwordConfirmation: object.passwordConfirmation || ''
        }}
        validationSchema={formSchema}
        onSubmit={this.props.onSubmit}>
        {({
          handleSubmit,
          setFieldValue,
          handleChange,
          handleBlur,
          values
        }) => (
            <Form onSubmit={handleSubmit} noValidate>
              <Field
                type="password"
                label={i18n.t(`${packageNS}:tr000004`) + "*"}
                name="password"
                id="password"
                value={values.password}
                onChange={handleChange}
                helpText={i18n.t(`${packageNS}:menu.registration.password_hint`)}
                component={ReactstrapInput}
                onBlur={handleBlur}
                required
              />

              <Field
                type="password"
                label={i18n.t(`${packageNS}:tr000023`) + "*"}
                name="passwordConfirmation"
                id="passwordConfirmation"
                value={values.passwordConfirmation}
                onChange={handleChange}
                helpText={i18n.t(`${packageNS}:menu.registration.password_hint`)}
                component={ReactstrapInput}
                onBlur={handleBlur}
                required
              />
              <Button type="submit" color="primary">{this.props.submitLabel}</Button>
            </Form>
          )}
      </Formik>
    );
  }
}


class ChangeUserPassword extends Component {
  constructor() {
    super();
    this.state = {};

    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    UserStore.get(this.props.match.params.userID, resp => {
      this.setState({
        user: resp,
      });
    });
  }

  onSubmit(password) {
    UserStore.updatePassword(this.props.match.params.userID, password.password, resp => {
      this.props.history.push(`/users/${this.props.match.params.userID}`);
    });
  }

  render() {
    const { classes } = this.props;

    if (this.state.user === undefined) {
      return (<div></div>);
    }

    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={`/users`}>{i18n.t(`${packageNS}:tr000036`)}</Link></BreadcrumbItem>
            <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={`/users/${this.state.user.user.id}`}>{this.state.user.user.username}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000038`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Container fluid>
          <Row xs="1" lg="1">
            <Card>
              <CardBody>
                <PasswordForm
                  submitLabel={i18n.t(`${packageNS}:tr000022`)}
                  userId={this.state.user.user.id}
                  onSubmit={this.onSubmit}
                />
              </CardBody>
            </Card>
          </Row>
        </Container>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(ChangeUserPassword));
