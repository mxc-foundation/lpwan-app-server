import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";

import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";


const styles = {
  textField: {
    width: "100%",
  },
  link: {
    "& a": {
      color: theme.palette.primary.main,
      textDecoration: "none",
    },
  },
};


class RegistrationForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return null;
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        onSubmit={this.onSubmit}
      >
        <TextField
          id="username"
          label="Email"
          margin="normal"
          value={this.state.object.username || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
      </Form>
    );
  }
}


class Registration extends Component {
  constructor() {
    super();

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(user) {
    SessionStore.register(user, () => {
      this.props.history.push("/");
    });
  }

  render() {
    return(
      <Grid container justify="center">
        <Grid item xs={6} lg={4}>
          <Card>
            <CardHeader
              title="Registration"
            />
            <CardContent>
              <RegistrationForm
                submitLabel="Register"
                onSubmit={this.onSubmit}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Registration));
