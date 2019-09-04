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


class RegistrationConfirmForm extends FormComponent {
  componentDidMount() {
    SessionStore.confirmRegistration(this.props.securityToken, (resp) => {
      if (resp) {
        this.setState({
          object: resp,
          isTokenValid: true
        })
        SessionStore.setToken(resp.jwt)
      } else {
        this.setState({
          isTokenValid: false
        })
      }
    })
  }

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
          InputProps={{
            readOnly: true,
          }}
          fullWidth
          required
        />
        <TextField
          id="password"
          label="Password"
          type="password"
          minLength="6"
          margin="normal"
          value={this.state.object.password || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <TextField
          id="passwordConfirmation"
          label="Password repeat"
          type="password"
          minLength="6"
          margin="normal"
          value={this.state.object.passwordConfirmation || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <TextField
          id="organizationName"
          label="Organization name"
          pattern="[\w-]+"
          margin="normal"
          value={this.state.object.organizationName || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <TextField
          id="organizationDisplayName"
          label="Organization display name"
          margin="normal"
          value={this.state.object.organizationDisplayName || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
      </Form>
    );
  }
}


class RegistrationConfirm extends Component {
  constructor() {
    super();

    this.state = {
      isTokenValid: null,
      isPwdMatch: null
    }

    localStorage.setItem('jwt', '')

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(data) {
    console.log('onSubmit(', data, ')')
    if (data.password === data.passwordConfirmation) {
      this.setState({
        isPwdMatch: true
      })

      let request = {
        userId: data.id,
        password: data.password,
        organizationName: data.organizationName,
        organizationDisplayName: data.organizationDisplayName,
      }

      SessionStore.finishRegistration(request, (responseData) => {
        this.props.history.push("/");
      })
    } else {
      this.setState({
          isPwdMatch: false
      })
    }
    
  }

  render() {
    return(
        <Grid container justify="center">
          <Grid item xs={6} lg={4}>
            <Card>
              <CardHeader
                title="Registration confirmation"
              />
              <CardContent>
                {this.state.isPwdMatch !== null && this.state.isPwdMatch === false &&
                  <p style={{color: 'Red', textAlign: 'center'}}>Passwords doesn't match.</p>
                }
                <RegistrationConfirmForm
                  submitLabel="Confirm"
                  onSubmit={this.onSubmit}
                  securityToken={this.props.match.params.securityToken}
                />
              </CardContent>
            </Card>
          </Grid>
        </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(RegistrationConfirm));
