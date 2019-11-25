import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import ReCAPTCHA from "react-google-recaptcha";

import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
//import MenuIcon from 'mdi-material-ui/Server';

import DropdownMenuLanguage from "../../components/DropdownMenuLanguage";
import TitleBarTitle from "../../components/TitleBarTitle";
import Password from '../../components/TextfileForPassword'
import { 
  Map,
  Marker,
  Popup,
  LayersControl,
  LayerGroup
 } from 'react-leaflet';

import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";
import MapTileLayerCluster from "../../components/MapTileLayerCluster";
//import { isAbsolute } from "path";
//import { NONAME } from "dns";
//import { relative } from "path";
//const DURATION = 550;
//const COLOR = 'rgba(121,244,218,0.5)';
import i18n, { packageNS } from '../../i18n';

const VERIFY_ERROR_MESSAGE = i18n.t(`${packageNS}:tr000021`);
const styles = {
  textField: {
    width: "100%",
    display: 'flex',
    justifyContent: 'center'
  },
  link: {
    "& a": {
      color: theme.palette.textSecondary.main,
      textDecoration: "none",
    },
  },
  padding: {
    paddingTop: 115,
  },
  z1000: {
    zIndex: 1000
  },
  loginFormStyle: {
    backgroundColor: '#10337b50',
    padding: '24px',
    position: 'absolute',
    width: 380,
    top: '0',
    right: '0',
    background:'linear-gradient(rgba(121,244,218,0.5),transparent)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '100%',
  },
  logo: {
    height: 90,
    marginLeft: 'auto',
    opacity: '0.7',
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
  },
  appBar: {
    backgroundColor: theme.palette.secondary.main,
  },
};

class LoginForm extends FormComponent {
  constructor(props) {
    super(props);

    this.handleChange = this.handleChange.bind(this);
  }

  handleChange = (event) => {
    this.state.object.password = event;
  };

  onChangeLanguage = e => {
    const newLanguage = {
      id: e.id,
      label: e.label,
      value: e.value,
      code: e.code
    }

    this.props.onChangeLanguage(newLanguage);
  }

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
        to={`/registration`} 
        type="button" 
        disabled={false}>{i18n.t(`${packageNS}:tr000020`)}</Button>
    ]
    let demoUsername = "";
    let demoPassword = "";
    let helpText = "";
    if(window.location.origin.includes(process.env.REACT_APP_DEMO_HOST_SERVER)){
      demoUsername = process.env.REACT_APP_DEMO_USER;
      demoPassword = process.env.REACT_APP_DEMO_USER_PASSWORD;
      helpText = i18n.t(`${packageNS}:tr000010`);
    }

    return(
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
        onSubmit={this.onSubmit}
      >
        <div className={this.props.style.logoSection}>
          <img src="/logo/mxc_logo-social.png" className={this.props.style.logo} alt={i18n.t(`${packageNS}:tr000051`)} />
        </div>

        <DropdownMenuLanguage onChangeLanguage={this.onChangeLanguage} />

        <TextField
          id="username"
          label={i18n.t(`${packageNS}:tr000003`)}
          margin="normal"
          value={this.state.object.username === undefined 
            ? this.state.object.username = demoUsername 
            : this.state.object.username }
          autoComplete='off'
          onChange={this.onChange}
          fullWidth
        />
        <Password handleChange={this.handleChange} demoPassword={demoPassword} helpText={helpText} label={i18n.t(`${packageNS}:tr000004`)} />
        {/* <TitleBarTitle component={Link} to={`/password-recovery`} title={i18n.t(`${packageNS}:tr000009`)} /> */}
        <ReCAPTCHA
          sitekey={process.env.REACT_APP_PUBLIC_KEY}
          onChange={this.onReCapChange}
          className={this.props.style.textField}
        />
      </Form>
    );
  }
}


class Login extends Component {
  constructor() {
    super();

    this.state = {
      registration: null,
      open: true,
      accessOn: false,
      isVerified: false
    };
    
    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    SessionStore.getBranding(resp => {
      if (resp.registration !== "") {
        this.setState({
          registration: resp.registration,
        });
      }
    });
  }

  onChangeLanguage = (newLanguageState) => {
    this.props.onChangeLanguage(newLanguageState);
  }

  onSubmit(login) {
    if(login.hasOwnProperty('isVerified')){
      if(!login.isVerified){
        alert(VERIFY_ERROR_MESSAGE);
        return false;
      }

      SessionStore.login(login, () => {
        const orgs = SessionStore.getOrganizations();
        console.log('Organizations: ', orgs);

        if (SessionStore.getToken() && orgs.length > 0) {
          this.props.history.push(`/organizations/${orgs[0]}`);
        } else {
          console.log('User has no organisations. Redirecting to login');
          this.props.history.push("/");
        }
      });
    }else{
      alert(VERIFY_ERROR_MESSAGE);
      return false;
    }
  }

  onClick = () => {
    this.setState(function(prevState) {
			return {accessOn: !prevState.accessOn};
		});
  }

  render() {
    const style = {
      position: 'absolute',
      top: 0,
      bottom: 0,
      left: 0,
      right: 0,
      zIndex: 1
    };

    let position = [];
    
    position = [51,13];
    
    return(
      <>
        <Map center={position} zoom={6} style={style} animate={true} scrollWheelZoom={false}>
          <MapTileLayerCluster />
        </Map>
        <div className={this.props.classes.padding + ' ' + this.props.classes.z1000}>
          <div className={this.props.classes.loginFormStyle}>
            <LoginForm
              onChangeLanguage={this.onChangeLanguage}
              onSubmit={this.onSubmit}
              submitLabel={i18n.t(`${packageNS}:tr000011`)}
              style={this.props.classes}
            />
          </div>
          {this.state.registration && <div>
            <Typography className={this.props.classes.link} dangerouslySetInnerHTML={{__html: this.state.registration}}></Typography>
          </div>}
        </div>
      </>
    );
  }
}

export default withStyles(styles)(withRouter(Login));
