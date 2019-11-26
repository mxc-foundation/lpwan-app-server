import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { isEmail } from 'validator';

import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";
import ALiYunCaptcha from 'react-aliyun-captcha';

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

  onCallback = (value) => {
    const req = {
      token: value.token,
      sessionId : value.csessionid,
      sig: value.sig,
      remoteIp: window.location.origin
    }

    if(value.value === 'pass'){
      SessionStore.getVerifyingRecaptcha(req, resp => {
        console.log('ali resp: ', resp);
        this.state.object.isVerified = resp.success;
      });
    }else{
      console.log("can't pass verify process.");
    }
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
          type="email"
          value={this.state.object.username || ""}
          onChange={this.onChange}
          fullWidth
          required
        />
        <ALiYunCaptcha
              appKey="FFFF0N000000000087AA"
              scene="nc_login"
              onCallback={this.onCallback}
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
      alert(i18n.t(`${packageNS}:tr000021`));
      return false;
    }

    if(SessionStore.getLanguage() && SessionStore.getLanguage().id){
      user.language = SessionStore.getLanguage().id.toLowerCase();
    }else {
      user.language = 'en';
    }

    if(isEmail(user.username)){
      SessionStore.register(user, () => {
        this.props.history.push("/");
      });
    }else{
      alert(i18n.t(`${packageNS}:tr000024`));
    }
  }

  render() {
    return(
      <Grid container justify="center">
        <Grid item xs={6} lg={4}>
          <Card>
            <CardHeader
              title={i18n.t(`${packageNS}:tr000019`)}
            />
            <CardContent>
              <RegistrationForm
                submitLabel={i18n.t(`${packageNS}:tr000020`)}
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
