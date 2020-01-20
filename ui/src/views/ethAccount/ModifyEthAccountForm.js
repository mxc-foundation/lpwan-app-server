import React, { Component, useEffect } from "react";

import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { ReactstrapInput } from '../../components/FormInputs';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../i18n';


const fieldsSchema = {
  activeAccount: Yup.string().trim(),
  newAccount: Yup.string().trim(),
  username: Yup.string().trim(),
  password: Yup.string(),
}

class ModifyEthAccountForm extends Component {

  constructor(props) {
    super(props);

    this.state = {
      object: {
        activeAccount: this.props.activeAccount
      }
    };
  }

  componentDidUpdate(oldProps) {
    if (this.props.activeAccount ===  oldProps.activeAccount){
      return;
    }

    this.setState({
      object: {
        activeAccount: this.props.activeAccount
      }
    })
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

  render() {
    if (this.props.activeAccount == '0') {
      return 'loading...';
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
                  label={i18n.t(`${packageNS}:menu.eth_account.current_account`)}
                  name="activeAccount"
                  id="activeAccount"
                  value={this.state.object.activeAccount || ""}
                  placeholder="0x0000000000000000000000000000000000000000"
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  readOnly
                  inputProps={{
                    clearable: true,
                    cache: false,
                  }}
                />
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.eth_account.new_account`)}
                  name="newAccount"
                  id="newAccount"
                  value={this.state.object.newAccount || ""}
                  placeholder="0x0000000000000000000000000000000000000000"
                  component={ReactstrapInput}
                  placeholder={i18n.t(`${packageNS}:menu.eth_account.new_account`)}
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
                  placeholder={i18n.t(`${packageNS}:menu.eth_account.type_here`)}
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
                  placeholder={i18n.t(`${packageNS}:menu.eth_account.type_here`)}
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

export default ModifyEthAccountForm;
