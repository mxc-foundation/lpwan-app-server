import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Button } from 'reactstrap';
import * as Yup from 'yup';
import { ReactstrapCheckbox, ReactstrapInput } from '../../components/FormInputs';
import i18n, { packageNS } from '../../i18n';



class NonAdminChoices extends Component {
  render() {
    const { handleChange, handleBlur } = this.props;

    return (<><Field
      type="checkbox"
      label={i18n.t(`${packageNS}:tr000141`)}
      name="isDeviceAdmin"
      id="isDeviceAdmin"
      component={ReactstrapCheckbox}
      onChange={handleChange}
      onBlur={handleBlur}
      helpText={i18n.t(`${packageNS}:tr000142`)}
    />
      <Field
        type="checkbox"
        label={i18n.t(`${packageNS}:tr000143`)}
        name="isGatewayAdmin"
        id="isGatewayAdmin"
        component={ReactstrapCheckbox}
        onChange={handleChange}
        onBlur={handleBlur}
        helpText={i18n.t(`${packageNS}:tr000144`)}
      /></>);
  }
}

class OrganizationUserForm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      object: this.props.object || {},
    };

  }

  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }

    let fieldsSchema = {
      username: Yup.string().trim().required("Required"),
      isAdmin: Yup.bool(),
      isDeviceAdmin: Yup.bool(),
      isGatewayAdmin: Yup.bool(),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={(values) => {
            const castValues = formSchema.cast(values);
            this.props.onSubmit({ ...castValues })
          }}>
          {({
            handleSubmit,
            handleChange,
            setFieldValue,
            values,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:tr000056`)}
                  name="username"
                  id="username"
                  value={this.state.object.username || ""}
                  helpText={i18n.t(`${packageNS}:tr000138`)}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly={this.props.update}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />

                <Field
                  type="checkbox"
                  label={i18n.t(`${packageNS}:tr000139`)}
                  name="isAdmin"
                  id="isAdmin"
                  value={values.isAdmin}

                  component={ReactstrapCheckbox}
                  onChange={handleChange}

                  onBlur={handleBlur}
                  helpText={i18n.t(`${packageNS}:tr000140`)}
                />
                {!values.isAdmin && <NonAdminChoices handleChange={handleChange} handleBlur={handleBlur} />}
                <Button type="submit" color="primary">{this.props.submitLabel || i18n.t(`${packageNS}:tr000066`)}</Button>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default OrganizationUserForm;
