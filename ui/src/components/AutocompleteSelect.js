import React, { Component } from "react";

import Input from '@material-ui/core/Input';
import MenuItem from '@material-ui/core/MenuItem';
import Chip from '@material-ui/core/Chip';
import FormControl from "@material-ui/core/FormControl";

import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

import MenuDown from "mdi-material-ui/MenuDown";
import Cancel from "mdi-material-ui/Cancel";
import MenuUp from "mdi-material-ui/MenuUp";
import Close from "mdi-material-ui/Close";
import AsyncSelect from 'react-select/async';
const inputStyles = {
  marginB: {
    marginBottom: 24,
  },
} 
const customStyles = {
  control: (base, state) => ({
    ...base,
    // match with the menu
    borderRadius: state.isFocused ? "3px 3px 0 0" : 3,
    // Overwrittes the different states of border
    borderColor: state.isFocused ? "#00FFD9" : "white",
    // Removes weird border around container
    boxShadow: state.isFocused ? null : null,
    "&:hover": {
      // Overwrittes the different states of border
      borderColor: state.isFocused ? "#00FFD9" : "white"
    }
  }),
  menu: base => ({
    ...base,
    background:'#1a2d6e',
    // override border radius to match the box
    borderRadius: 0,
    // kill the gap
    marginTop: 28,
  }),
  menuList: base => ({
    ...base,
    background:'#1a2d6e',
    // kill the white space on first and last option
    padding: 0,
  }),
  option: base => ({
    ...base,
    // kill the white space on first and last option
    padding: '10px',
    maxWidth: 221,
    whiteSpace: 'nowrap', 
    overflow: 'hidden',
    textOverflow: 'ellipsis'
  }),
};
// taken from react-select example
// https://material-ui.com/demos/autocomplete/

class Option extends Component {
  handleClick = event => {
    this.props.onSelect(this.props.option, event);
  };

  render() {
    const { children, isFocused, isSelected, onFocus } = this.props;

    return (
      <MenuItem
        onFocus={onFocus}
        selected={isFocused}
        onClick={this.handleClick}
        component="div"
        style={{
          fontWeight: isSelected ? 500 : 400,
        }}
      >
        {children}
      </MenuItem>
    );
  }
}

function SelectWrapped(props) {
  const { classes, inputRef, ...other } = props;

  React.useImperativeHandle(inputRef, () => ({
    focus: () => {
    },
  }));
  
  const components = {
    option: Option,
    value: (valueProps) => {
      const { value, children, onRemove } = valueProps;
      const onDelete = event => {
        event.preventDefault();
        event.stopPropagation();
        onRemove(value);
      };

      if (onRemove) {
        return (
          <Chip
            tabIndex={-1}
            label={children}
            className={classes.chip}
            deleteIcon={<Cancel onTouchEnd={onDelete} />}
            onDelete={onDelete}
          />
        );
      }

      return <div className="Select-value">{children}</div>;
    }
  };

  return (
     <AsyncSelect
      components={components}
      styles={customStyles}
      theme={(theme) => ({
        ...theme,
        borderRadius: 4,
        colors: {
          primary25: '#00FFD950',
          primary: '#00FFD950',
        },
      })}
      //noOptionsMessage={<Typography>{'No results found'}</Typography>}
      arrowRenderer={arrowProps => {
        return arrowProps.isOpen ? <MenuUp /> : <MenuDown />;
      }}
      clearRenderer={() => <Close />}
      {...other}
    />
  );
}


class AutocompleteSelect extends Component {
  constructor() {
    super();

    this.state = {
      options: [],
    };

    this.setInitialOptions = this.setInitialOptions.bind(this);
    this.setSelectedOption = this.setSelectedOption.bind(this);
    this.onAutocomplete = this.onAutocomplete.bind(this);
    this.onChange = this.onChange.bind(this);
  }

  componentDidMount() {
    this.setInitialOptions(this.setSelectedOption);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.value === this.props.value && prevProps.triggerReload === this.props.triggerReload) {
      return;
    }

    this.setInitialOptions(this.setSelectedOption);
  }

  setInitialOptions(callbackFunc) {
    this.props.getOptions("", options => {
      
      this.setState({
        options: options,
      }, callbackFunc);
    });
  }

  setSelectedOption() {
    if (this.props.getOption !== undefined) {
      if (this.props.value !== undefined && this.props.value !== "" && this.props.value !== null) {
        this.props.getOption(this.props.value, resp => {
          this.setState({
            selectedOption: resp,
          });
        });
      } else {
        this.setState({
          selectedOption: "",
        });
      }
    } else {
      if (this.props.value !== undefined && this.props.value !== "" && this.props.value !== null) {
        for (const opt of this.state.options) {
          if (this.props.value === opt.value) {
            this.setState({
              selectedOption: opt,
            });
          }
        }
      } else {
        this.setState({
          selectedOption: "",
        });
      }
    }
  }

  onAutocomplete(input) {
    return new Promise((resolve, reject) => {
      this.props.getOptions(input, options => {
        
        this.setState({
          options: options,
        });

        resolve(options);
      });
    });
  }

  onChange(v) {
    let value = null;
    let label = null;
    if (v !== null) {
      value = v.value;
      label = v.label;
    }

    this.setState({
      selectedOption: v,
    });

    this.props.onChange({
      target: {
        id: this.props.id,
        value,
        label
      },
    });
  }

  render() {
    const inputProps = this.props.inputProps || {};
    return(
      <FormControl margin={this.props.margin || ""}  fullWidth={true} 
        className={this.props.className}>
        <Input
          fullWidth
          className={this.props.classes.marginB}
          inputComponent={SelectWrapped}
          placeholder={this.props.label}
          id={this.props.id}
          onChange={this.onChange}
          disableUnderline
          inputProps={{...{
            instanceId: this.props.id,
            clearable: false,
            defaultOptions: this.state.options,
            loadOptions: this.onAutocomplete,
            value: this.state.selectedOption || "",
          }, ...inputProps}}
        />
      </FormControl>
    );
  }
}

export default withStyles(inputStyles)(withRouter(AutocompleteSelect));