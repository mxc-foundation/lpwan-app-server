import React from "react";

import { Button } from 'reactstrap';
import TextField from '@material-ui/core/TextField';
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";
// import Button from "@material-ui/core/Button";

import FormComponent from "../../../../classes/FormComponent";
import HTTPIntegrationHeaderForm from "./HTTPIntegrationHeaderForm";

class HTTPIntegrationForm extends FormComponent {
  constructor() {
    super();
    this.addHeader = this.addHeader.bind(this);
    this.onDeleteHeader = this.onDeleteHeader.bind(this);
    this.onChangeHeader = this.onChangeHeader.bind(this);
  }

  onChange(e) {
    super.onChange(e);
    this.props.onChange(this.state.object);
  }

  addHeader(e) {
    e.preventDefault();

    let object = this.state.object;
    if(object.headers === undefined) {
      object.headers = [{}];
    } else {
      object.headers.push({});
    }

    this.props.onChange(object);
  }

  onDeleteHeader(index) {
    let object = this.state.object;
    object.headers.splice(index, 1);
    this.props.onChange(object);
  }

  onChangeHeader(index, header) {
    let object = this.state.object;
    object.headers[index] = header;
    this.props.onChange(object);
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    let headers = [];
    if (this.state.object.headers !== undefined) {
      headers = this.state.object.headers.map((h, i) =>
        <HTTPIntegrationHeaderForm
          classes={this.props.classes}
          key={i}
          index={i}
          object={h}
          onChange={this.onChangeHeader}
          onDelete={this.onDeleteHeader}
        />);
    }

    return(
      <React.Fragment>
        <FormControl fullWidth margin="normal">
          <br />
          <h4>Headers</h4>
          {headers}
        </FormControl>
        <Button variant="outlined" onClick={this.addHeader}>Add Header</Button>
        <FormControl fullWidth margin="normal">
          <FormLabel>
            <br />
            <h4>Endpoints</h4>
          </FormLabel>
          <TextField
            id="uplinkDataURL"
            label="Uplink data URL"
            placeholder="http://example.com/uplink"
            value={this.state.object.uplinkDataURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
          <TextField
            id="joinNotificationURL"
            label="Join notification URL"
            placeholder="http://example.com/join"
            value={this.state.object.joinNotificationURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
          <TextField
            id="statusNotificationURL"
            label="Device-status notification URL"
            placeholder="http://example.com/status"
            value={this.state.object.statusNotificationURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
          <TextField
            id="locationNotificationURL"
            label="Location notification URL"
            placeholder="http://example.com/location"
            value={this.state.object.locationNotificationURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
          <TextField
            id="ackNotificationURL"
            label="ACK notification URL"
            placeholder="http://example.com/ack"
            value={this.state.object.ackNotificationURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
          <TextField
            id="errorNotificationURL"
            label="Error notification url"
            placeholder="http://example.com/error"
            value={this.state.object.errorNotificationURL || ""}
            onChange={this.onChange}
            margin="normal"
            fullWidth
          />
        </FormControl>
      </React.Fragment>
    );
  }
}

export default HTTPIntegrationForm;
