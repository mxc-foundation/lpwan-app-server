import React, { Component } from "react";
import { connect, Field } from 'formik';
import MaskedInput from "react-text-mask";
import { Button as RButton, UncontrolledTooltip } from 'reactstrap';
import i18n, { packageNS } from '../i18n';
import { ReactstrapInputGroup } from './FormInputs';




class DevAddrMask extends Component {
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
        ]}
      />
    );
  }
}


class DevAddrField extends Component {
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
    this.props.randomFunc((k) => {
      let key = k;
      const bytes = key.match(/[\w]{2}/g);
      if(!this.state.msb && bytes !== null) {
        key = bytes.reverse().join("");
      }

      this.setState({
        value: key,
      });

      this.props.onChange({
        target: {
          value: k,
          type: "text",
          id: this.props.id,
        },
      });
    });
  }

  onChange = (e) => {
    this.setState({
      value: e.target.value,
    });

    let str = "";

    const bytes = e.target.value.match(/[\w]{2}/g);
    if (!this.state.msb && bytes !== null) {
      str = bytes.reverse().join("");
    } else if (bytes !== null) {
      str = bytes.join("");
    } else {
      str = "";
    }

    this.props.onChange({
      target: {
        value: str.replace(/\s/g, ""),
        type: "text",
        id: this.props.id,
      },
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

    const controls = (
      <React.Fragment>
        <React.Fragment>
          <RButton color="primary" type="button" onClick={this.toggleByteOrder} id={`${this.props.id}-DevAddrToggleBtn`}>
            {this.state.msb ? i18n.t(`${packageNS}:tr000220`) : i18n.t(`${packageNS}:tr000221`)}
          </RButton>
          <UncontrolledTooltip placement="bottom" target={`${this.props.id}-DevAddrToggleBtn`}>
            {i18n.t(`${packageNS}:tr000373`)}
          </UncontrolledTooltip>
        </React.Fragment>

        {this.props.random && !this.props.disabled &&
          <React.Fragment>
            <RButton color="secondary" type="button" onClick={this.randomKey} id={`${this.props.id}-DevAddrRefreshBtn`}>
              <i className="mdi mdi-refresh"></i>
            </RButton>
            <UncontrolledTooltip placement="bottom" target={`${this.props.id}-DevAddrRefreshBtn`}>
              {i18n.t(`${packageNS}:tr000375`)}
            </UncontrolledTooltip>
          </React.Fragment>
        }
      </React.Fragment>
    );
    
    return (
      <React.Fragment>
        <Field
          type="text"
          label={this.props.label}
          name={this.props.name}
          id={this.props.id}
          helpText={this.props.helpText || null}
          append={controls}
          inputComponent={DevAddrMask}
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

export default connect(DevAddrField);
