import React, { Component } from "react";

import { Field, connect } from 'formik';
import { Button, UncontrolledTooltip } from 'reactstrap';
import { ReactstrapInputGroup } from './FormInputs';
import MaskedInput from "react-text-mask";

import i18n, { packageNS } from '../i18n';


class EUI64HEXMask extends Component {
  render() {
    const { inputRef, inputComponent, helpText, ...other } = this.props;

    return (
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
        ]}
      />
    );
  }
}


class EUI64Field extends Component {
  constructor() {
    super();

    this.state = {
      msb: true,
      value: "",
    };
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

    for (let i = 0; i < 16; i++) {
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
      <Button color="primary" type="button"
        onClick={this.toggleByteOrder} id="toggleBtn">{this.state.msb ? i18n.t(`${packageNS}:tr000220`) : i18n.t(`${packageNS}:tr000221`)}</Button>
      {this.props.random && !this.props.disabled && <Button color="secondary" type="button" onClick={this.randomKey} id="generateBtn">
        <i className="mdi mdi-refresh"></i>
      </Button>}
      <UncontrolledTooltip placement="bottom" target="toggleBtn">{i18n.t(`${packageNS}:tr000373`)}</UncontrolledTooltip>
      {this.props.random && !this.props.disabled && <UncontrolledTooltip placement="bottom" target="generateBtn">{i18n.t(`${packageNS}:tr000391`)}</UncontrolledTooltip>}
    </React.Fragment>;

    return (<React.Fragment>

      <Field
        type="text"
        label={this.props.label}
        name={this.props.name}
        id={this.props.id}
        helpText={this.props.helpText || null}
        append={controls}
        inputComponent={EUI64HEXMask}
        onChange={this.onChange}
        defaultValue={this.state.value}
        disabled={this.props.disabled}
        required={this.props.required}
        component={ReactstrapInputGroup}
      />
    </React.Fragment>
    );
  }
}

export default connect(EUI64Field);
