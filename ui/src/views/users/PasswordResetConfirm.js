import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import { withStyles } from "@material-ui/core/styles";
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import FormComponent from "../../classes/FormComponent";
import Form from "../../components/Form";
import Password from '../../components/TextfileForPassword';
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
//import ReCAPTCHA from "react-google-recaptcha";
import { PASSWORD_RECOVERY_DESCRIPTION_002, PASSWORD_RECOVERY_ERROR_MINIMUM_LENGTH, PASSWORD_RECOVERY_ERROR_MISMATCH } from "../../util/Messages";


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


class PasswordResetConfirmForm extends FormComponent {

    handlePassword = (event) => {
        this.state.object.password = event;
    };

    handlePasswordConfirm = (event) => {
        this.state.object.passwordConfirm = event;
    };

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
            disabled={false}>{i18n.t(`${packageNS}:tr000014`)}</Button>
        ]

        return(
        <Form
            submitLabel={this.props.submitLabel}
            extraButtons={extraButtons}
            onSubmit={this.onSubmit}
        >
            <Password handleChange={this.handlePassword} label={i18n.t(`${packageNS}:tr000004`)} />
            <Password handleChange={this.handlePasswordConfirm} label={i18n.t(`${packageNS}:tr000416`)} />
        </Form>
        );
    }
}


class PasswordResetConfirm extends Component {
  constructor() {
    super();
    this.state = {
      isVerified: false
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(passwords) {
    if(passwords.password.length < 6){
        alert(PASSWORD_RECOVERY_ERROR_MINIMUM_LENGTH);
        return false;
    }

    if(passwords.password !== passwords.passwordConfirm){
        alert(PASSWORD_RECOVERY_ERROR_MISMATCH);
        return false;
    }  
  }

  render() {
    return(
      <>
        <div className={this.props.classes.root}>
          <AppBar position="static" className={this.props.classes.appBar}>
            <Toolbar>
              <div className={this.props.logoSection}>
                <img src="/logo/logo_mx.png" className={this.props.classes.logo} alt={i18n.t(`${packageNS}:tr000051`)} />
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
              <TitleBarTitle title={i18n.t(`${packageNS}:tr000012`)} />
            </div>
            <Divider light={true}/>
            <Typography variant="body1" className={this.props.classes.title}>
                {PASSWORD_RECOVERY_DESCRIPTION_002}
            </Typography>
            <PasswordResetConfirmForm
                submitLabel={i18n.t(`${packageNS}:tr000325`)}
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

export default withStyles(styles)(withRouter(PasswordResetConfirm));
