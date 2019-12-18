import React from "react";

import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
//import FormGroup from "@material-ui/core/FormGroup";
import Checkbox from '@material-ui/core/Checkbox';

import { Button, Form, FormGroup, Label, Input, FormText, Row, Col } from 'reactstrap';
import FormComponent from "../../classes/FormComponent";
import FormControl from "../../components/FormControl";
import FormSubmit from "../../components/Form";
import i18n, { packageNS } from '../../i18n';


class UserForm extends FormComponent {

  render() {
    if (this.state.object === undefined) {
      return (<div></div>);
    }

    return (
      <Form>
        <FormGroup row>
          <Label for="username" sm={2}>{i18n.t(`${packageNS}:tr000056`)}</Label>
          <Col sm={10}>
            <Input type="text" name="username" id="username" value={this.state.object.username || ""} onChange={this.onChange} />
          </Col>
        </FormGroup>
        <FormGroup row>
          <Label for="email" sm={2}>{i18n.t(`${packageNS}:tr000147`)}</Label>
          <Col sm={10}>
            <Input type="email" name="email" id="email" value={this.state.object.email || ""} onChange={this.onChange} />
          </Col>
        </FormGroup>
        <FormGroup row>
          <Label for="note" sm={2}>{i18n.t(`${packageNS}:tr000129`)}</Label>
          <Col sm={10}>
            <Input type="textarea" name="note" id="note" value={this.state.object.note || ""} onChange={this.onChange} />
            <FormText color="muted">{i18n.t(`${packageNS}:tr000130`)}</FormText>
          </Col>
        </FormGroup>
        {this.state.object.id === undefined && <FormGroup row>
          <Label for="password" sm={2}>{i18n.t(`${packageNS}:tr000004`)}</Label>
          <Col sm={10}>
            <Input type="password" name="password" id="password" value={this.state.object.password || ""} onChange={this.onChange} />
            <FormText color="muted">{i18n.t(`${packageNS}:tr000130`)}</FormText>
          </Col>
        </FormGroup>}
        <FormGroup check>
          <Label check>
            <Input type="checkbox"
              id="isAdmin"
              checked={!!this.state.object.isAdmin}
              onChange={this.onChange}
              color="primary"
            />{' '}
            {i18n.t(`${packageNS}:tr000133`)}
          </Label>
        </FormGroup>

        {this.props.submitLabel && <Button color="primary"
          onClick={this.onSubmit}
          disabled={this.props.disabled}
          className="btn-block">{this.props.submitLabel}
        </Button>}

      </Form>
    );
  }
}

export default UserForm;
