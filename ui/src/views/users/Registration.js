import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";
import { isEmail } from 'validator';

import { Row, Col, Container, Card, CardBody, Button, FormGroup } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import ReCAPTCHA from "react-google-recaptcha";

import { ReactstrapInput, ReactstrapPasswordInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import SessionStore from "../../stores/SessionStore";
import i18n, { packageNS } from '../../i18n';

import MneMonicPhrase from './MneMonicPhrase';
import MneMonicPhraseConfirm from './MneMonicPhraseConfirm';
import Google2FA from './Google2FA';


const regSchema = Yup.object().shape({
  username: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
  password: Yup.string().trim().min(8).required(i18n.t(`${packageNS}:tr000431`)).matches(
    /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/,
    i18n.t(`${packageNS}:menu.registration.password_match_error`)
  ),
  confirm_password: Yup.string().trim().oneOf([Yup.ref('password')], i18n.t(`${packageNS}:menu.registration.confirm_password_match_error`)).required(i18n.t(`${packageNS}:tr000431`)),
  org_name: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
  org_display_name: Yup.string().trim()
})

class RegistrationForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || { username: "", password: "", confirm_password: "", org_name: "", org_display_name: "" },
      isVerified: false,
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

  render() {

    return (
      <React.Fragment>
        <Formik
          initialValues={this.state.object}
          validationSchema={regSchema}
          onSubmit={(values) => {
            const castValues = regSchema.cast(values);
            this.props.onSubmit({ isVerified: this.state.isVerified, ...castValues })
          }}>
          {({
            handleSubmit,
            handleBlur,
          }) => (
              <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.registration.email`) + '*'}
                  name="username"
                  id="username"
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                />

                <Field
                  helpText={this.state.object.helpText}
                  label={i18n.t(`${packageNS}:menu.registration.password`) + '* ' + i18n.t(`${packageNS}:menu.registration.password_hint`)}
                  name="password"
                  id="password"
                  component={ReactstrapPasswordInput}
                  onBlur={handleBlur}
                />

                <Field
                  helpText={this.state.object.helpText}
                  label={i18n.t(`${packageNS}:menu.registration.confirm_password`) + '*'}
                  name="confirm_password"
                  id="confirm_password"
                  component={ReactstrapPasswordInput}
                  onBlur={handleBlur}
                />

                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.registration.org_name`) + '*'}
                  name="org_name"
                  id="org_name"
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                />

                <Field
                  type="text"
                  label={i18n.t(`${packageNS}:menu.registration.org_display_name`)}
                  name="org_display_name"
                  id="org_display_name"
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
                  <Button type="submit" color="primary" className="btn-block" disabled={(!this.state.bypassCaptcha) && (!this.state.isVerified)}>{i18n.t(`${packageNS}:menu.registration.next`)}</Button>
                  <Link to={`/login`} className="btn btn-link btn-block text-muted mt-0">{i18n.t(`${packageNS}:tr000462`)}</Link>
                </div>
              </Form>
            )}
        </Formik>
      </React.Fragment>
    );
  }
}

const MneMonicPhraseStep1 = ({ next }) => {
  return <React.Fragment>
    <Row>
      <Col>
        <p className="pb-4 pt-3">{i18n.t(`${packageNS}:menu.registration.mnemonic_phrase_instruction`)}</p>

        <Button type="submit" color="primary" className="btn-block mt-5" onClick={next}>{i18n.t(`${packageNS}:menu.registration.next`)}</Button>
      </Col>
    </Row>
  </React.Fragment>
}

function GetBranding() {
  return new Promise((resolve, reject) => {
    SessionStore.getBranding(resp => {
      return resolve(resp);
    });
  });
}

class Registration extends Component {
  constructor() {
    super();

    let bypassCaptcha = false;
    if (window.location.origin.includes("http://localhost")) {
      bypassCaptcha = true;
    }

    this.state = {
      isVerified: false,
      bypassCaptcha: bypassCaptcha,
      showRegisterForm: true,
      startMnemonicPhrase: false,
      showMnemonicPhraseList: false,
      showMnemonicPhraseConfirm: false,
      showTwoFactorAuth: false
    };

    this.onSubmit = this.onSubmit.bind(this);
    this.showMnemonicPhraseList = this.showMnemonicPhraseList.bind(this);
    this.showMnemonicPhraseListConfirm = this.showMnemonicPhraseListConfirm.bind(this);
    this.confirmMnemonicPhraseList = this.confirmMnemonicPhraseList.bind(this);
    this.confirm2fa = this.confirm2fa.bind(this);
    this.skip2fa = this.skip2fa.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  loadData = async () => {
    try {
      let result = await GetBranding();

      this.setState({
        logoPath: result.logoPath
      });
    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  onSubmit(user) {
    if (this.state.bypassCaptcha) {
      user.isVerified = true;
    }

    if (!user.isVerified) {
      alert(i18n.t(`${packageNS}:tr000021`));
      return false;
    }

    if (SessionStore.getLanguage() && SessionStore.getLanguage().id) {
      user.language = SessionStore.getLanguage().id.toLowerCase();
    } else {
      user.language = 'en';
    }

    if (isEmail(user.username)) {
      this.setState({ loading: true });
      SessionStore.register(user, () => {
        this.setState({user: user, loading: false, showRegisterForm: false, startMnemonicPhrase: true });

        // this.props.history.push("/");
      });
    } else {
      alert(i18n.t(`${packageNS}:tr000024`));
    }
  }

  showMnemonicPhraseList() {
    // TODO - Fetch phrase - for now setting up dummy
    const phrases = ["Simba", "Sweetie", "Ziggy", "Midnight", "Kiki", "Peanut", "Midday", "Buddy", "Bently", "Gray", "Rocky", "Madison", "Bella", "Baxter"];
    this.setState({ showMnemonicPhraseList: true, startMnemonicPhrase: false, phrases: phrases, showMnemonicPhraseConfirm: false });
  }

  showMnemonicPhraseListConfirm() {
    this.setState({ showMnemonicPhraseList: false, startMnemonicPhrase: false, showMnemonicPhraseConfirm: true });
  }

  confirmMnemonicPhraseList(phrases) {
    // TODO - API call to confirm order of phrase
    this.setState({ showMnemonicPhraseList: false, startMnemonicPhrase: false, showMnemonicPhraseConfirm: false, showTwoFactorAuth: true });

    // TODO - API call to get the code for 2FA
    this.setState({auth_2fa_code: '123'});
  }

  confirm2fa(confirmCode) {
    // TODO  - API call to confirm
    this.setState({ showMnemonicPhraseList: false, startMnemonicPhrase: false, showMnemonicPhraseConfirm: false, showTwoFactorAuth: false });

    this.props.history.push("/");
  }

  skip2fa() {
    // TODO - for now redirecting to login
    this.props.history.push("/");
  }

  render() {
    return (<React.Fragment>
      <div className="account-pages mt-4 mb-0">
        <Container>
          <Row className="justify-content-center">
            <Col md={8} lg={6} xl={6}>
              <div className="text-center mb-3">
                <Link to="/">
                  <span><img src={this.state.logoPath} alt="" height="54" /></span>
                </Link>
              </div>

              <Card className="h-auto">
                <CardBody className="p-4">
                  <div className="text-center mb-3">
                    <h4 className="text-uppercase mt-0">{i18n.t(`${packageNS}:tr000019`)}</h4>
                  </div>

                  <div className="position-relative">
                    {this.state.loading && <Loader />}

                    {this.state.showRegisterForm ?
                      <RegistrationForm
                        onSubmit={this.onSubmit}
                        bypassCaptcha={this.state.bypassCaptcha}
                      /> : <React.Fragment>

                        {this.state.startMnemonicPhrase? <MneMonicPhraseStep1 next={this.showMnemonicPhraseList} />: null}

                        {this.state.showMnemonicPhraseList? <MneMonicPhrase 
                          title={i18n.t(`${packageNS}:menu.registration.mnemonic_phrase_title`)} phrase={this.state.phrases} 
                            next={this.showMnemonicPhraseListConfirm} showSkip={false} />: null}

                        {this.state.showMnemonicPhraseConfirm ? <MneMonicPhraseConfirm 
                          title={i18n.t(`${packageNS}:menu.registration.mnemonic_phrase_confirm_title`)}
                          phrase={this.state.phrases} next={this.confirmMnemonicPhraseList} back={this.showMnemonicPhraseList} /> : null}

                        {this.state.showTwoFactorAuth ? <Google2FA
                          title={i18n.t(`${packageNS}:menu.registration.2fa_title`)}
                          code={this.state.auth_2fa_code}
                          confirm={this.confirm2fa} skip={this.skip2fa} /> : null}
                      </React.Fragment>}

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
