import React, { Component } from "react";

import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { ReactstrapInput } from '../../components/FormInputs';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../i18n';
class NewEthAccountForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: {},
    };
  }


  onChange = (event) => {
    const { id, value } = event.target;

    this.setState({
      object: { [id]: value }
    });
  }

  reset = () => {
    this.setState({
      object: {
        username: '',
        password: '',
        newAccount: ''
      }
    })
  }

  onSubmit = () => {
    this.props.onSubmit({
      action: 'newAccount',
      newAccount: this.state.object.newAccount,
      currentAccount: this.state.object.newAccount,
      username: this.state.object.username,
      password: this.state.object.password
    });

    this.reset();
  }

  render() {
    let fieldsSchema = {
      newAccount: Yup.string().trim(),
      username: Yup.string().trim(),
      password: Yup.string(), 
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
                  label={i18n.t(`${packageNS}:menu.eth_account.new_account`)}
                  name="newAccount"
                  id="newAccount"
                  value={this.state.object.newAccount || ""}
                  placeholder="0x0000000000000000000000000000000000000000"
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.withdraw.username`)}
                  name="username"
                  id="username"
                  value={this.state.object.username || ""}
                  component={ReactstrapInput}
                  placeholder={i18n.t(`${packageNS}:menu.withdraw.type_here`)}
                  onBlur={handleBlur}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />
                <Field
                  type="password"
                  label={i18n.t(`${packageNS}:menu.eth_account.password`)}
                  name="password"
                  id="password"
                  value={this.state.object.password || ""}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />
                <Button className="btn-block" onClick={this.reset}>{i18n.t(`${packageNS}:common.reset`)}</Button>
                <Button type="submit" className="btn-block" color="primary">{this.props.submitLabel || i18n.t(`${packageNS}:tr000066`)}</Button>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

export default NewEthAccountForm;
