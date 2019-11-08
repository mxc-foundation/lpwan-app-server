import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";
import { isEmail } from 'validator';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";
import ReCAPTCHA from "react-google-recaptcha";
import { PASSWORD_RECOVERY_DESCRIPTION_001 } from "../../util/Messages";
import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";

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


class PasswordRecoveryForm extends FormComponent {

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

    const extraButtons = [
        <Button 
          variant="outlined"
          color="inherit"
          component={Link} 
          to={`/login`} 
          type="button" 
          disabled={false}>Canceled</Button>
      ]

    return(
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
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


class PasswordRecovery extends Component {
  constructor() {
    super();
    this.state = {
      isVerified: false
    };

    this.onSubmit = this.onSubmit.bind(this);
  }
  
  

  onSubmit(user) {
    if(!user.isVerified){
      alert("Are you a human, please verify yourself.");
      return false;
    }

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
              title="Password Recovery"
            />
            <CardContent>
                <Typography variant="body1" className={this.props.classes.title}>
                    {PASSWORD_RECOVERY_DESCRIPTION_001}
                </Typography>
                <PasswordRecoveryForm
                    submitLabel="Reset Password"
                    onSubmit={this.onSubmit}
                    style={this.props.classes}
                    className={this.props.classes.formWidth}
                />
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6} lg={4}></Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(PasswordRecovery));
