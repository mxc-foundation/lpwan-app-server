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
import i18n, { packageNS } from '../../i18n';


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
          label={i18n.t(`${packageNS}:tr000003`)}
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
          label={i18n.t(`${packageNS}:tr000004`)}
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
          label={i18n.t(`${packageNS}:tr000023`)}
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
          label={i18n.t(`${packageNS}:tr000030`)}
          pattern="[\w-]+"
          margin="normal"
          value={this.state.object.organizationName || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <TextField
          id="organizationDisplayName"
          label={i18n.t(`${packageNS}:tr000031`)}
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
                title={i18n.t(`${packageNS}:tr000019_confirmation`)}
              />
              <CardContent>
                {this.state.isPwdMatch !== null && this.state.isPwdMatch === false &&
                  <p style={{color: 'Red', textAlign: 'center'}}>{i18n.t(`${packageNS}:tr000025`)}</p>
                }
                <RegistrationConfirmForm
                  submitLabel={i18n.t(`${packageNS}:tr000022`)}
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
