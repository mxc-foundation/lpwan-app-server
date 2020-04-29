import React, { Component } from "react";
import { connect, Field } from 'formik';
import MaskedInput from "react-text-mask";
import { Button as RButton, UncontrolledTooltip } from 'reactstrap';
import i18n, { packageNS } from '../i18n';
import { ReactstrapInputGroup } from './FormInputs';




class AESKeyHEXMask extends Component {
  render() {
    const { inputRef, inputComponent, helpText, ...other } = this.props;

    return(
      <MaskedInput
        className="form-control"
        {...other}
        ref={inputRef}
        mask={[
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
          ' ',
          /[A-Fa-f0-9]/,
          /[A-Fa-f0-9]/,
        ]}
      />
    );
  }
}


class AESKeyField extends Component {
  constructor() {
    super();

    this.state = {
      showKey: false,
      msb: true,
      value: "",
    };
  }

  toggleShowPassword = () => {
    this.setState({
      showKey: !this.state.showKey,
    });
  }

  toggleByteOrder = () => {
    this.setState({
      msb: !this.state.msb,
    });

    const bytes = this.state.value.match(/[A-Fa-f0-9]{2}/g);
    if (bytes !== null) {
      this.setState({
        value: bytes.reverse().join(" ").replace(/\s/g, ""),
      });
    }
  }

  randomKey = () => {
    let key = "";
    const possible = 'abcdef0123456789';

    for(let i = 0; i < 32; i++){
      key += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    this.setState({
      value: key,
    });
  }

  onChange = (e) => {
    this.setState({
      value: e.target.value.replace(/\s/g, ""),
    });
  }

  componentDidMount() {
    this.setState({
      value: this.props.value.replace(/\s/g, "") || "",
      showKey: this.props.value === "" ? true : false,
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.state !== prevState) {
      const { name, formik: { setFieldValue } } = this.props;
      setFieldValue(name, this.state.value);
    }
  }

  render() {

    const controls = <React.Fragment>
      {this.state.showKey && <React.Fragment><RButton color="primary" type="button"
        onClick={this.toggleByteOrder} id={`${this.props.id}-AEStoggleBtn`}>{this.state.msb ? i18n.t(`${packageNS}:tr000220`) : i18n.t(`${packageNS}:tr000221`)}</RButton>
        <UncontrolledTooltip placement="bottom" target={`${this.props.id}-AEStoggleBtn`}>{i18n.t(`${packageNS}:tr000373`)}</UncontrolledTooltip></React.Fragment>}

      {this.props.random && this.state.showKey && !this.props.disabled && <React.Fragment><RButton color="secondary" type="button" onClick={this.randomKey} 
        id={`${this.props.id}-AESRefreshBtn`}>
        <i className="mdi mdi-refresh"></i></RButton>
        <UncontrolledTooltip placement="bottom" target={`${this.props.id}-AESRefreshBtn`}>{i18n.t(`${packageNS}:tr000376`)}</UncontrolledTooltip></React.Fragment>}

      <RButton color="secondary" type="button" onClick={this.toggleShowPassword} id={`${this.props.id}-showHideBtn`}>
        {this.state.showKey ? <i className="mdi mdi-eye-off"></i> : <i className="mdi mdi-eye"></i>}</RButton>
    </React.Fragment>;
    
    return <React.Fragment>
      <Field
        type={this.state.showKey ? "text" : "password"}
        label={this.props.label}
        name={this.props.name}
        id={this.props.id}
        helpText={this.props.helpText || null}
        append={controls}
        inputComponent={AESKeyHEXMask}
        onChange={this.onChange}
        defaultValue={this.state.value}
        disabled={this.props.disabled || !this.state.showKey}
        required={this.props.required}
        component={ReactstrapInputGroup}
      />
    </React.Fragment>;
  }
}

export default connect(AESKeyField);
