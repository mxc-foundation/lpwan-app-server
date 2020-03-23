import React from "react";

import { Field } from 'formik';

import Grid from "@material-ui/core/Grid";
import TextField from '@material-ui/core/TextField';
import IconButton from '@material-ui/core/IconButton';
import Delete from "mdi-material-ui/Delete";

import { ReactstrapInput } from '../../../../components/FormInputs';
import FormComponent from "../../../../classes/FormComponent";

class HTTPIntegrationHeaderForm extends FormComponent {
  constructor() {
    super();

    this.onDelete = this.onDelete.bind(this);
  }

  onChange(e) {
    super.onChange(e);
    this.props.onChange(this.props.index, this.state.object);
  }

  onDelete(e) {
    e.preventDefault();
    this.props.onDelete(this.props.index);
  }

  render() {
    if (this.state.object === undefined) {
      return(<div></div>);
    }

    return(
      <Grid container spacing={4}>
        <Grid item xs={5}>
          <Field
            autoFocus
            component={ReactstrapInput}
            type="text"
            label="Header name"
            name="key"
            id="key"
            helpText=""
            placeholder=""
            value={this.state.object.key || ""}
            onChange={this.onChange}
          />
        </Grid>
        <Grid item xs={5}>
          <Field
            component={ReactstrapInput}
            type="text"
            label="Header value"
            name="value"
            id="value"
            helpText=""
            placeholder=""
            value={this.state.object.value || ""}
            onChange={this.onChange}
          />
        </Grid>
        <Grid item xs={2} className={this.props.classes.delete}>
          <IconButton aria-label="delete" onClick={this.onDelete}>
            <Delete />
          </IconButton>
        </Grid>
      </Grid>
    );    
  }
}


export default HTTPIntegrationHeaderForm;
