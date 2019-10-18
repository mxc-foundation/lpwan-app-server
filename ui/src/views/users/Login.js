import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from "@material-ui/core/Typography";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import ReCAPTCHA from "react-google-recaptcha";

import Form from "../../components/Form";
import FormComponent from "../../classes/FormComponent";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";

const VERIFY_ERROR_MESSAGE = "Are you a human, please verify yourself.";
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
    paddingTop: 230,
  },
  padd: {
    position: 'absolute',
    width: '100%',
    height: '100%',
    top: 0,
    left: 0,
    backgroundImage: 'url("/img/world-map.png")',
  },
  cbody: {
    width: 388,
  },
  bbody: {
    display: 'flex',
    justifyContent: 'center'
  }
};



class LoginForm extends FormComponent {
  
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
      isVerified: false
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
    if(login.hasOwnProperty('isVerified')){
      if(!login.isVerified){
        alert(VERIFY_ERROR_MESSAGE);
        return false;
      }
      
      SessionStore.login(login, () => {
        this.props.history.push("/");
      });
    }else{
      alert(VERIFY_ERROR_MESSAGE);
      return false;
    }
  }

  render() {
    return(
      <div className={this.props.classes.padd}>
      <Grid container justify="center" className={this.props.classes.padding}>
        <Grid item xs={6} lg={4} className={this.props.classes.bbody}>
          <Card className={this.props.classes.cbody}>
            <CardContent>
              <LoginForm
                submitLabel="Login"
                onSubmit={this.onSubmit}
                style={this.props.classes}
              />
            </CardContent>
            {this.state.registration && <CardContent>
              <Typography className={this.props.classes.link} dangerouslySetInnerHTML={{__html: this.state.registration}}></Typography>
             </CardContent>}
          </Card>
        </Grid>
      </Grid>
      </div>
    );
  }
}

export default withStyles(styles)(withRouter(Login));
