import React, { Component } from "react";
import { Input } from "reactstrap";

export default class EditableTextInput extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: props.value,
      showEditableInput: false
    };

    this.onInputChange = this.onInputChange.bind(this);
    this.onBlur = this.onBlur.bind(this);
    this.handleKeyDown = this.handleKeyDown.bind(this);
    this.toggleEditable = this.toggleEditable.bind(this);
  }

  onInputChange = e => {
    this.setState({ value: e.target.value });
    this.props.onDataChanged(e.target);
    /* if (this.props.onDataChanged) {
      this.props.onDataChanged(this.props.id, this.props.field, e.target.value);
    } */
  };

  onBlur(e) {
    this.onInputChange(e);
    this.toggleEditable();
  }

  handleKeyDown(target) {
    if (target.keyCode == 13) {
      this.onInputChange(target);
      this.toggleEditable();
    } else if (target.keyCode == 27) {
      this.toggleEditable();
    }
  }

  toggleEditable() {
    this.setState({ showEditableInput: !this.state.showEditableInput });
  }

  render() {
    const inputType = this.props.inputType || "text";

    return (
      <React.Fragment>
        {this.state.showEditableInput ? (
          <Input
            type={inputType}
            id={this.props.id}
            placeholder={this.props.placeholder}
            value={this.state.value || ""}
            onChange={this.onInputChange}
            onBlur={this.onBlur}
            onKeyDown={this.handleKeyDown}
          />
        ) : (
          <span className="d-block w-auto h-100" onClick={this.toggleEditable}>
            {this.props.valueFormatter
              ? this.props.valueFormatter(this.state.value)
              : this.state.value}
          </span>
        )}
      </React.Fragment>
    );
  }
}
