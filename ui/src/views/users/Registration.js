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
import { Divider } from "@material-ui/core";


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
    console.log("Captcha value:", value);
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

    this.onSubmit = this.onSubmit.bind(this);
  }
  
  

  onSubmit(user) {
    if(isEmail(user.username)){
      SessionStore.register(user, () => {
        this.props.history.push("/");
      });
    }else{
      alert("Please, enter a valid email address to use.");
    }
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
