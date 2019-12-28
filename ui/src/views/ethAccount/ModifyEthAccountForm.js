import React, { Component, useEffect } from "react";

import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { ReactstrapInput } from '../../components/FormInputs';
import * as Yup from 'yup';

import i18n, { packageNS } from '../../i18n';



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


  // onChange = (event) => {
  //   const { id, value } = event.target;

  //   this.setState({
  //     object: { [id]: value }
  //   });
  // }
  reset = () => {
    this.setState({
      object: {
        username: '',
        password: '',
        newaccount: ''
      }
    })
  }

  submit = () => {
    this.props.onSubmit({
      action: 'modifyAccount',
      currentAccount: this.state.object.newaccount,
      createAccount: this.state.object.newaccount,
      username: this.state.object.username,
      password: this.state.object.password
    })

    this.reset();
  }

  render() {
    if (this.props.activeAccount == '0') {
      return 'loading...';
    }
    let fieldsSchema = {
      activeAccount: Yup.string(),
      newaccount: Yup.string(),
      username: Yup.string(),
      password: Yup.string(),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    console.log('ModfyEthAccountForm.render', this.state.object);
      
    return (
      <React.Fragment>
        <Formik
          enableReinitialize
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={this.props.onSubmit}>
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
                  label={i18n.t(`${packageNS}:menu.withdraw.username`)}
                  name="newaccount"
                  id="newaccount"
                  value={this.state.object.newaccount || ""}
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

export default ModifyEthAccountForm;
