import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { isEmail } from 'validator';

import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";
import ReCAPTCHA from "react-google-recaptcha";


import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import i18n, { packageNS } from '../../i18n';

const styles = {
  textField: {
    width: "100%",
    display: 'flex',
    justifyContent: 'center'

  },
  formWidth: {
    width: 352,
  },
  link: {
    "& a": {
      color: theme.palette.primary.main,
      textDecoration: "none",
    },
  },
};


class RegistrationForm extends FormComponent {

  onReCapChange = (value) => {
    const req = {
      secret : process.env.REACT_APP_PUBLIC_KEY,
      response: value,
      remoteip: window.location.origin
    }

    SessionStore.getVerifyingGoogleRecaptcha(req, resp => {
      this.state.object.isVerified = resp.success;
    }); 
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
          label={i18n.t(`${packageNS}:common.email`)}
          margin="normal"
          type="email"
          value={this.state.object.username || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <ReCAPTCHA
                sitekey={process.env.REACT_APP_PUBLIC_KEY}
                onChange={this.onReCapChange}
                className={this.props.style.textField}
              />
      </Form>
    );
  }
}


class Registration extends Component {
  constructor() {
    super();
    this.state = {
      isVerified: false
    };

    this.onSubmit = this.onSubmit.bind(this);
  }
  
  

  onSubmit(user) {
    if(!user.isVerified){
      alert(i18n.t(`${packageNS}:common.human`));
      return false;
    }

    if(SessionStore.getLanguage().label){
      user.language = 'en';
    }else {
      user.language = SessionStore.getLanguage().label.toLowerCase();
    }
    
    if(isEmail(user.username)){
      SessionStore.register(user, () => {
        this.props.history.push("/");
      });
    }else{
      alert(i18n.t(`${packageNS}:registration.valid_email`));
    }
  }

  render() {
    return(
      <Grid container justify="center">
        <Grid item xs={6} lg={4}>
          <Card>
            <CardHeader
              title={i18n.t(`${packageNS}:registration.registration`)}
            />
            <CardContent>
              <RegistrationForm
                submitLabel={i18n.t(`${packageNS}:registration.register`)}
                onSubmit={this.onSubmit}
                style={this.props.classes}
                className={this.props.classes.formWidth}
              />
              
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Registration));
