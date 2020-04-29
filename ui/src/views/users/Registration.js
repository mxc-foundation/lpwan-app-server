import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import ReCAPTCHA from "react-google-recaptcha";
import { Link, withRouter } from "react-router-dom";
import { Button, Card, CardBody, Col, Container, CustomInput, FormGroup, Row, Modal, ModalBody, ModalHeader } from 'reactstrap';
import { isEmail } from 'validator';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import ServerInfoStore from "../../stores/ServerInfoStore";
import UserStore from "../../stores/UserStore";
import Google2FA from './Google2FA';

const regSchema = Yup.object().shape({
  username: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
})

class RegistrationFormAverage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      object: this.props.object || { username: "" },
      isVerified: false,
      showSetup2FA: false,
      twofa_enabled: false,
      isVerifiedDuplicate: null,
      bypassCaptcha: this.props.bypassCaptcha
    }
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

  verificationUser = async (userName) => {
    
    if (userName === '') {
      this.setState({isVerifiedDuplicate:false});
      return false;
    }

    if (!isEmail(userName)) {
      this.setState({isVerifiedDuplicate:false});
      return false;
    }

    const resp = await UserStore.getUserEmail(userName);
    if(resp !== undefined){
      if(resp.status){
        this.setState({isVerifiedDuplicate:true});
        return true;
      }else{
        this.setState({isVerifiedDuplicate:false});
        return false;
      }
    }else{
      this.setState({isVerifiedDuplicate:false});
      return false;
    }
  }

  showSetup2FA = () => {
    // TODO - API Call to fetch the initial code
    this.setState({ showSetup2FA: true, auth_2fa_code: '12345678' });
  }

  confirm2fa = (confirmCode) => {
    // TODO  - API call to confirm
    this.setState({ showSetup2FA: false, twofa_enabled: true });
  }

  skip2fa = () => {
    this.setState({ showSetup2FA: false });
  }

  confirmReset2fa = (confirmCode) => {
    // TODO  - API call to confirm
    // TODO - Fetch phrase - for now setting up dummy
    const phrases = ["Simba", "Sweetie", "Ziggy", "Midnight", "Kiki", "Peanut", "Midday", "Buddy", "Bently", "Gray", "Rocky", "Madison", "Bella", "Baxter"];
    this.setState({ showMnemonicPhraseConfirm: true, phrases: phrases });
  }

  skipReset2fa = () => {
    this.setState({ showReset2FA: false, showMnemonicPhraseConfirm: false });
  }

  render() {
    let fieldsSchema = {
      username: Yup.string().trim().test('testname', i18n.t(`${packageNS}:menu.registration.user_validation`), this.verificationUser
      ).required(i18n.t(`${packageNS}:tr000431`)),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <React.Fragment>
        <Formik
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={(values) => {
            const castValues = formSchema.cast(values);
            const isVerified = this.state.isVerified;
            const isVerifiedDuplicate = this.state.isVerifiedDuplicate;
            this.props.onSubmit({ isVerified, isVerifiedDuplicate, ...castValues })
          }}>
          {({
            handleSubmit,
            handleBlur,
            values,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:tr000003`)}
                  name="username"
                  id="username"
                  value={values.username}
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                />

                <FormGroup className="mt-2">
                  {!this.state.bypassCaptcha && <ReCAPTCHA
                    sitekey={process.env.REACT_APP_PUBLIC_KEY}
                    onChange={this.onReCapChange}
                  />}
                </FormGroup>
                <div className="mt-1">
                  <Button type="submit" color="primary" className="btn-block" disabled={(!this.state.bypassCaptcha) && (!this.state.isVerified)}>{i18n.t(`${packageNS}:tr000020`)}</Button>
                  <Link to={`/login`} className="btn btn-link btn-block text-muted mt-0">{i18n.t(`${packageNS}:tr000462`)}</Link>
                </div>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

class RegistrationFormRestricted extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || { username: "" },
      isVerified: false,
      isVerifiedDuplicate: null,
      bypassCaptcha: this.props.bypassCaptcha
    }
  }

  verificationUser = async (userName) => {
    
    if (userName === '') {
      this.setState({isVerifiedDuplicate:false});
      return false;
    }

    if (!isEmail(userName)) {
      this.setState({isVerifiedDuplicate:false});
      return false;
    }

    const resp = await UserStore.getUserEmail(userName);
    if(resp !== undefined){
      if(resp.status){
        this.setState({isVerifiedDuplicate:true});
        return true;
      }else{
        this.setState({isVerifiedDuplicate:false});
        return false;
      }
    }else{
      this.setState({isVerifiedDuplicate:false});
      return false;
    }
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
    let fieldsSchema = {
      username: Yup.string().trim().test('testname', i18n.t(`${packageNS}:menu.registration.user_validation`), this.verificationUser
      ).required(i18n.t(`${packageNS}:tr000431`)),
    }

    const formSchema = Yup.object().shape(fieldsSchema);

    return (
      <React.Fragment>
        <Formik
          initialValues={this.state.object}
          validationSchema={formSchema}
          onSubmit={(values) => {
            const castValues = regSchema.cast(values);
            const isVerified = this.state.isVerified;
            const isVerifiedDuplicate = this.state.isVerifiedDuplicate;
            this.props.onSubmit({ isVerified, isVerifiedDuplicate, ...castValues })
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

                <div className="mt-1">
                  <Button type="submit" color="primary" className="btn-block" >{i18n.t(`${packageNS}:tr000020`)}</Button>
                  <Link to={`/login`} className="btn btn-link btn-block text-muted mt-0">{i18n.t(`${packageNS}:tr000462`)}</Link>
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

function LoadServerRegion() {
  return new Promise((resolve, reject) => {
    ServerInfoStore.getServerRegion(resp => {
      return resolve(resp);
    });
  });
}

class Registration extends Component {
  constructor() {
    super();

    let bypassCaptcha = false;
    if (window.location.origin.includes("https://lora.demo") || window.location.origin.includes("http://localhost")) {
      bypassCaptcha = true;
    }

    this.state = {
      isVerified: false,
      bypassCaptcha: bypassCaptcha,
      serverRegion : ""
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  loadData = async () => {
    try {
      let result = await GetBranding();
      let serverRegion = await ServerInfoStore.getServerRegion();

      this.setState({
        logoPath: result.logoPath,
        serverRegion: serverRegion.serverRegion.trim()
      });
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  onSubmit = (user) => {
    if (this.state.serverRegion === "NOT_DEFINED" || this.state.serverRegion === "AVERAGE") {
      if (this.state.bypassCaptcha) {
        user.isVerified = true;
      }

      if (!user.isVerified) {
        alert(i18n.t(`${packageNS}:tr000021`));
        return false;
      }

    }else{
      if (this.state.bypassCaptcha) {
        user.isVerified = true;
      }
    }

    
    
    if (!user.isVerifiedDuplicate) {
      return false;
    }

    if (SessionStore.getLanguage() && SessionStore.getLanguage().id) {
      user.language = SessionStore.getLanguage().id;
    } else {
      user.language = 'en';
    }

    if (isEmail(user.username)) {
      this.setState({ loading: true });
      SessionStore.register(user, () => {
        this.setState({ loading: false });
        SessionStore.setUsernameTemp(user.username);
        this.props.history.push("/registration-confirm/");
      });
    } else {
      alert(i18n.t(`${packageNS}:tr000024`));
    }
  }

  render() {
    if(!this.state.serverRegion){
      return <div></div>;
    }
    return (<React.Fragment>
      <div className="account-pages mt-5 mb-5">
        <Container>
          <Row className="justify-content-center">
            <Col md={8} lg={6} xl={5}>
              <div className="text-center mb-3">
                <Link to="/">
                  <span><img src={this.state.logoPath} alt="" height="54" /></span>
                </Link>
              </div>

              <Card>
                <CardBody className="p-4">
                  <div className="text-center mb-4">
                    <h4 className="text-uppercase mt-0">{i18n.t(`${packageNS}:tr000019`)}</h4>
                  </div>

                  <div className="position-relative">
                    {this.state.loading && <Loader />}
                    {
                      this.state.serverRegion === "RESTRICTED"
                        ? <RegistrationFormRestricted
                          onSubmit={this.onSubmit}
                          bypassCaptcha={this.state.bypassCaptcha}
                        />
                        : <RegistrationFormAverage
                          onSubmit={this.onSubmit}
                          bypassCaptcha={this.state.bypassCaptcha}
                        />
                    }

                  </div>
                </CardBody>
              </Card>
            </Col>
          </Row>
        </Container>
      </div>
    </React.Fragment>
    );
  }
}

export default withRouter(Registration);
