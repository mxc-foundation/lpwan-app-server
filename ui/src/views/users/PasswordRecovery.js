import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";
import { isEmail } from 'validator';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Divider from '@material-ui/core/Divider';
import TitleBarTitle from "../../components/TitleBarTitle";
import { withStyles } from "@material-ui/core/styles";
import ReCAPTCHA from "react-google-recaptcha";
import { PASSWORD_RECOVERY_DESCRIPTION_001 } from "../../util/Messages";
import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import theme from "../../theme";

const styles = {
  textField: {
    width: "100%",
    display: 'flex',
    justifyContent: 'center'

  },
  TitleBar:{
    padding: 0
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
  flexCol: {
    display: 'flex',
    flexDirection: 'column', 
  },
  logoSection: {
    display: 'flex'
  },
  root: {
    flexGrow: 1,
    position: 'absolute',
    top: 0,
    height: 84, 
    left: 0,
    right: 0,
    zIndex: 1
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
    height: 200,
    paddingTop: 25 
  },
  appBar: {
    backgroundColor: theme.palette.secondary.main,
  },
  logo:{
    height: 50 
  },
  divider: {
    padding: 0,
    color: '#FFFFFF',
    width: '100%',
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
        {/* <ReCAPTCHA
                sitekey={process.env.REACT_APP_PUBLIC_KEY}
                onChange={this.onReCapChange}
                className={this.props.style.textField}
              /> */}
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

  onSubmit(email) {
    console.log('pass word recovery: ', email);
  }

  render() {
    return(
      <>
        <div className={this.props.classes.root}>
          <AppBar position="static" className={this.props.classes.appBar}>
            <Toolbar>
              <div className={this.props.logoSection}>
                <img src="/logo/logo_mx.png" className={this.props.classes.logo} alt="LPWAN Server" />
              </div>
              <IconButton edge="start" className={this.props.classes.menuButton} color="inherit" aria-label="menu">
                {/* <MenuIcon /> */}
              </IconButton>
              <Typography variant="h6"></Typography>
            </Toolbar>
          </AppBar>
        </div>
        <Grid container justify="center">
          <Grid item xs={12}>
          </Grid>
          <Grid item xs={12} lg={3} className={this.props.classes.flexCol}>
            <div className={this.props.classes.TitleBar}>
              <TitleBarTitle title="Password Recovery" />
            </div>
            <Divider light={true}/>
            <Typography variant="body1" className={this.props.classes.title}>
                {PASSWORD_RECOVERY_DESCRIPTION_001}
            </Typography>
            <PasswordRecoveryForm
                submitLabel="Reset Password"
                onSubmit={this.onSubmit}
                style={this.props.classes}
                className={this.props.classes.formWidth}
            />
          </Grid>
          <Grid item xs={2}>
          </Grid>
        </Grid>
      </>
    );
  }
}

export default withStyles(styles)(withRouter(PasswordRecovery));
