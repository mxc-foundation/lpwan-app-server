import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import TextField from '@material-ui/core/TextField';
import Typography from "@material-ui/core/Typography";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import Drawer from '@material-ui/core/Drawer';

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
import { isAbsolute } from "path";
import { NONAME } from "dns";
//import { relative } from "path";
const DURATION = 550;
const COLOR = 'rgba(121,244,218,0.5)';

const styles = {
  textField: {
    width: "100%",
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
    width: '20%',
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
  }
};

class LoginForm extends FormComponent {
  render() {
    if (this.state.object === undefined) {
      return null;
    }
    
    const extraButtons = [
      <Button color="primary.main" component={Link} to={`/registration`} type="button" disabled={false}>Register</Button>
    ]
    let demoUsername = "";
    let demoPassword = "";
    let helpText = "";
    if(window.location.origin.includes(process.env.REACT_APP_DEMO_HOST_SERVER)){
      demoUsername = process.env.REACT_APP_DEMO_USER;
      demoPassword = process.env.REACT_APP_DEMO_USER_PASSWORD;
      helpText = "You can access with this account right now as a demo user.";
    }
    return(
      <Form
        submitLabel={this.props.submitLabel}
        extraButtons={extraButtons}
        onSubmit={this.onSubmit}
      >
        <div className={this.props.logoSection}>
          <img src="/logo/mxc_logo-social.png" className={this.props.logo} alt="LoRa Server" />
        </div>
        <TextField
          id="username"
          label="Username"
          margin="normal"
          value={this.state.object.username === undefined 
            ? this.state.object.username = demoUsername 
            : this.state.object.username }
          autoComplete='off'
          onChange={this.onChange}
          fullWidth
          required
        />
        <TextField
          id="password"
          label="Password"
          type="password"
          margin="normal"
          value={this.state.object.password === undefined 
            ? this.state.object.password = demoPassword 
            : this.state.object.password }
          helperText={helpText}  
          onChange={this.onChange}
          fullWidth
          required
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
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    SessionStore.logout(() => {});

    SessionStore.getBranding(resp => {
      if (resp.registration !== "") {
        this.setState({
          registration: resp.registration,
        });
      }
    });
  }

  onSubmit(login) {
    SessionStore.login(login, () => {
      this.props.history.push("/");
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
    const markers = [];
    

    return(
    <>
      <Map center={position} zoom={6} style={style} animate={true} scrollWheelZoom={false}>
        <MapTileLayerCluster />
        <LayersControl position="bottomleft">
          <LayersControl.Overlay name="Markers" checked>
            <LayerGroup>
              {markers.map((position, index) => {
                return (
                  <Marker key={index} position={position} radius={10}>
                    <Popup>{position.text}</Popup>
                  </Marker>
                );
              })}
            </LayerGroup>
          </LayersControl.Overlay>
        </LayersControl>
      </Map>
          <div className={this.props.classes.padding + ' ' + this.props.classes.z1000}>
            <div className={this.props.classes.loginFormStyle}>
              <LoginForm
                submitLabel="Login"
                onSubmit={this.onSubmit}
                logo={this.props.classes.logo}
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
