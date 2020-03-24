import { Field, Form, Formik } from 'formik';
import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Button, Card, CardBody, Col, Container, Row } from 'reactstrap';
import * as Yup from 'yup';
import { ReactstrapInput } from '../../components/FormInputs';
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import { PASSWORD_RECOVERY_DESCRIPTION_001 } from "../../util/Messages";



// validation
const emailSchema = Yup.object().shape({
  email: Yup.string().trim().required(i18n.t(`${packageNS}:tr000431`)),
})

// branding
function GetBranding() {
  return new Promise((resolve, reject) => {
    SessionStore.getBranding(resp => {
      return resolve(resp);
    });
  });
}

class PasswordRecoverForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      object: this.props.object || { email: "" },
      isVerified: false
    }
  }

  render() {

    return (
      <React.Fragment>
        <Formik
          initialValues={this.state.object}
          validationSchema={emailSchema}
          onSubmit={(values) => {
            const castValues = emailSchema.cast(values);
            this.props.onSubmit({ ...castValues });
          }}>
          {({
            handleSubmit,
            handleBlur
          }) => (
            <Form onSubmit={handleSubmit} noValidate>
                <Field
                  type="email"
                  label={i18n.t(`${packageNS}:tr000003`)}
                  name="email"
                  id="email"
                  component={ReactstrapInput}
                  onBlur={handleBlur}
                  onChange={this.onChange}
                />

                <Button type="submit" color="primary" className="btn-block">{i18n.t(`${packageNS}:tr000325`)}</Button>
                <Link to={`/login`} className="btn btn-link btn-block text-muted mt-0">{i18n.t(`${packageNS}:tr000462`)}</Link>
                
              </Form>
            )}
        </Formik>
      </React.Fragment>
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

  onSubmit(email) {
    console.log('password recovery: ', email);
  }

  render() {
    return(
      <>
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
                    <h4 className="text-uppercase mt-0">{i18n.t(`${packageNS}:tr000012`)}</h4>
                  </div>

                  <p>{PASSWORD_RECOVERY_DESCRIPTION_001}</p>

                  <div className="position-relative">
                    {this.state.loading && <Loader />}
                    <PasswordRecoverForm
                      onSubmit={this.onSubmit}
                      bypassCaptcha={this.state.bypassCaptcha}
                    />
                  </div>
                </CardBody>
              </Card>
            </Col>
          </Row>
        </Container>
      </div>
      </>
    );
  }
}

export default withRouter(PasswordRecovery);
