import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Card, CardBody, Button, FormGroup } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

import ReCAPTCHA from "react-google-recaptcha";
import { Map } from 'react-leaflet';

import FoundLocationMap from "../../components/FoundLocationMap"
import { ReactstrapInput, ReactstrapPasswordInput } from '../../components/FormInputs';
import DropdownMenuLanguage from "../../components/DropdownMenuLanguage";
import Loader from "../../components/Loader";
import SessionStore from "../../stores/SessionStore";
import i18n, { packageNS } from '../../i18n';

const VERIFY_ERROR_MESSAGE = i18n.t(`${packageNS}:tr000021`);


class LoginForm extends Component {
  constructor(props) {
    super(props);
    this.onChangeLanguage = this.onChangeLanguage.bind(this);

    let object = this.props.object || { username: "", password: "" };

    if (window.location.origin.includes(process.env.REACT_APP_DEMO_HOST_SERVER)) {
      object['username'] = process.env.REACT_APP_DEMO_USER;
      object['password'] = process.env.REACT_APP_DEMO_USER_PASSWORD;
      object['helpText'] = i18n.t(`${packageNS}:tr000010`);
    }

    this.state = {
      object: object,
      isVerified: false
    }
  }

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
      secret: process.env.REACT_APP_PUBLIC_KEY,
      response: value,
      remoteip: window.location.origin
    }

    SessionStore.getVerifyingGoogleRecaptcha(req, resp => {
      this.setState({ isVerified: resp.success });
    });
  }

  render() {

    return (<React.Fragment>
      <Formik
        initialValues={this.state.object}
        validationSchema={
          Yup.object().shape({
            username: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
            password: Yup.string().required(i18n.t(`${packageNS}:tr000431`)),
          })
        }
        onSubmit={(values) => {
          this.props.onSubmit({ isVerified: this.state.isVerified, ...values })
        }}>
        {({
          handleSubmit,
          handleBlur
        }) => (
            <Form onSubmit={handleSubmit} noValidate>
              <Field
                type="text"
                label={i18n.t(`${packageNS}:tr000003`)}
                name="username"
                id="username"
                component={ReactstrapInput}
                onBlur={handleBlur}
              />

              <Field
                helpText={this.state.object.helpText}
                label={i18n.t(`${packageNS}:tr000004`)}
                name="password"
                id="password"
                component={ReactstrapPasswordInput}
                onBlur={handleBlur}
              />

              <FormGroup className="mt-2">
                <ReCAPTCHA
                  sitekey={process.env.REACT_APP_PUBLIC_KEY}
                  onChange={this.onReCapChange}
                />
              </FormGroup>

              <div className="mt-1">
                <Button type="submit" color="primary" className="btn-block" disabled={!this.state.isVerified}>{i18n.t(`${packageNS}:tr000011`)}</Button>
                <Link to={`/registration`} className="btn btn-outline-primary btn-block mt-2">{i18n.t(`${packageNS}:tr000020`)}</Link>
                {/* <Link to={`/password-recovery`} className="btn btn-link btn-block text-muted mt-0">{i18n.t(`${packageNS}:tr000009`)}</Link> */}
              </div>
            </Form>
          )}
      </Formik>
    </React.Fragment>
    );
  }
}

function GetBranding() {
  return new Promise((resolve, reject) => {
    SessionStore.getBranding(resp => {
      return resolve(resp);
    });
  });
}

class Login extends Component {
  constructor() {
    super();

    this.state = {
      registration: null,
      open: true,
      accessOn: false,
      isVerified: false,
      logoPath: "/logo/MATCHX-SUPERNODE2.png",
      loading: false,
      showLoginContainer: true
    };

    this.onSubmit = this.onSubmit.bind(this);
    this.showLoginContainer = this.showLoginContainer.bind(this);
    this.hideLoginContainer = this.hideLoginContainer.bind(this);
  }


  componentDidMount() {
    this.loadData();
  }
  
  loadData = async () => {
    try {
      let result = await GetBranding();

      this.setState({
        registration: result.registration,
        logoPath: result.logoPath
      });
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }


  componentDidUpdate(oldProps) {
    if (this.props.logoPath === oldProps.logoPath) {
      return;
    }
    
    this.loadData();
  }

  onChangeLanguage = (newLanguageState) => {
    this.props.onChangeLanguage(newLanguageState);
  }

  hideLoginContainer = () => {
    this.setState({ showLoginContainer: false })
  }

  showLoginContainer = () => {
    this.setState({ showLoginContainer: true })
  }

  onSubmit(login) {
    if (login.hasOwnProperty('isVerified')) {
      if (!login.isVerified) {
        alert(VERIFY_ERROR_MESSAGE);
        return false;
      }

      SessionStore.login(login, () => {
        this.setState({ loading: false });

        const orgs = SessionStore.getOrganizations();

        if (SessionStore.getToken() && orgs.length > 0) {
          this.props.history.push(`/`);
        } else {
          console.log('User has no organisations. Redirecting to login');
          this.props.history.push("/");
        }
      });
    } else {
      alert(VERIFY_ERROR_MESSAGE);
      return false;
    }
  }

  onClick = () => {
    this.setState(function (prevState) {
      return { accessOn: !prevState.accessOn };
    });
  }

  render() {

    let position = [];

    position = [51, 13];

    return (<React.Fragment>
      <div>
        <Map center={position} zoom={6} className="map-container" animate={true} scrollWheelZoom={false}>
          <FoundLocationMap />

          {!this.state.showLoginContainer && <Button type="button" color="primary" className="back-to-login-btn" onClick={this.showLoginContainer}>
            <i className="mdi mdi-arrow-left mr-1"></i>{i18n.t(`${packageNS}:tr000462`)}</Button>}
        </Map>

        {this.state.showLoginContainer && <div className="login-form-container">
          <div className="d-flex align-items-center w-100 h-100 p-2 p-sm-3 mx-auto">
            <div className="w-100">
              <img src={this.state.logoPath} className="mx-auto d-block img-fluid logo" alt={i18n.t(`${packageNS}:tr000051`)} height="54" />

              <div className="mt-2">
                <Card className="shadow-sm">
                  <CardBody>
                    <div className="position-relative">
                      {this.state.loading && <Loader />}
                      <LoginForm
                        onSubmit={this.onSubmit}
                      />
                    </div>

                    <Row>
                      <Col>
                        <Button type="button" color="link" className="p-0 btn-block text-muted align-middle" onClick={this.hideLoginContainer}>
                          <i className="mdi mdi-arrow-left mr-1"></i>{i18n.t(`${packageNS}:tr000461`)}</Button>
                      </Col>
                      <Col className="text-right">
                        <DropdownMenuLanguage onChangeLanguage={this.onChangeLanguage} />
                      </Col>
                    </Row>

                    {this.state.registration &&
                      <Row className="mt-2">
                        <Col>
                          <h6 dangerouslySetInnerHTML={{ __html: this.state.registration }}></h6>
                        </Col>
                      </Row>}
                  </CardBody>
                </Card>
              </div>
            </div>
          </div>
        </div>}
      </div>
    </React.Fragment>
    );
  }
}

export default withRouter(Login);
