import React, { Component } from "react";

import TextField from "@material-ui/core/TextField";
import InputAdornment from '@material-ui/core/InputAdornment';
import IconButton from '@material-ui/core/IconButton';
import Button from "@material-ui/core/Button";
import Tooltip from '@material-ui/core/Tooltip';

import Refresh from "mdi-material-ui/Refresh";

import MaskedInput from "react-text-mask";

import i18n, { packageNS } from '../i18n';


class DevAddrMask extends Component {
  render() {
    const { inputRef, ...other } = this.props;

    return(
      <MaskedInput
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
        value: bytes.reverse().join(" "),
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
        value: str,
        type: "text",
        id: this.props.id,
      },
    });
  }

  componentDidMount() {
    this.setState({
      value: this.props.value || "",
    });
  }

  render() {
    return(
      <TextField
        type="text"
        InputProps={{
          inputComponent: DevAddrMask,
          endAdornment: <InputAdornment position="end">
            <Tooltip title={i18n.t(`${packageNS}:tr000373`)}>
              <Button
                aria-label={i18n.t(`${packageNS}:tr000374`)}
                onClick={this.toggleByteOrder}
              >
                {this.state.msb ? i18n.t(`${packageNS}:tr000220`): i18n.t(`${packageNS}:tr000221`)}
              </Button>
            </Tooltip>
            {this.props.random && !this.props.disabled && <Tooltip title={i18n.t(`${packageNS}:tr000375`)}>
              <IconButton
                aria-label={i18n.t(`${packageNS}:tr000376`)}
                onClick={this.randomKey}
              >
                <Refresh />
              </IconButton>
            </Tooltip>}
          </InputAdornment>
        }}
        {...this.props}
        onChange={this.onChange}
        value={this.state.value}
        disabled={this.props.disabled}
      />
    );
  }
}

export default DevAddrField;
