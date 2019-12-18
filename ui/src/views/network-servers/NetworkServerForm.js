import React from "react";
import {
  Button, Col, FormFeedback, FormGroup, FormText, Input, Label, 
  TabContent, TabPane, Nav, NavItem, NavLink, Row,
} from 'reactstrap';
// Example: https://final-form.org/docs/react-final-form/examples/record-level-validation
import { Form, Field } from "react-final-form";
import classnames from 'classnames';

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";

const submitButton = (submitting, submitLabel) => {
  return (
    <Button
      aria-label={submitLabel}
      block
      color="primary"
      disabled={submitting}
      size="md"
      type="submit"
    >
      {submitLabel}
    </Button>
  );
};

class NetworkServerForm extends FormComponent {
  constructor() {
    super();

    this.state = {
      activeTab: '1'
    };
  }

  toggle = tab => {
    const { activeTab } = this.state;
    if (activeTab !== tab) {
      this.setState({
        activeTab: tab
      });
    }
  }

  render() {
    const { activeTab } = this.state;
    const { onSubmit, submitLabel } = this.props;

    if (this.state.object === undefined) {
      return(null);
    }

    return(
      <Form
        onSubmit={onSubmit}
        initialValues={{
          id: this.state.object.id,
          name: this.state.object.name,
          server: this.state.object.server,
          gatewayDiscoveryEnabled: !!this.state.object.gatewayDiscoveryEnabled,
          // Fallback to undefined, otherwise it defaults to a value of 0 even if the user hasn't entered anything
          gatewayDiscoveryInterval: this.state.object.gatewayDiscoveryInterval || undefined,
          gatewayDiscoveryTXFrequency: this.state.object.gatewayDiscoveryTXFrequency || undefined,
          gatewayDiscoveryDR: this.state.object.gatewayDiscoveryDR || undefined,
          caCert: this.state.object.caCert,
          tlsCert: this.state.object.tlsCert,
          tlsKey: this.state.object.tlsKey,
          routingProfileCACert: this.state.object.routingProfileCACert,
          routingProfileTLSCert: this.state.object.routingProfileTLSCert,
          routingProfileTLSKey: this.state.object.routingProfileTLSKey
        }}
        validate={values => {
          console.log('validateForm values/activeTab: ', values);
          if (!values) {
            return {};
          }
          const errors = {};
        
          if (activeTab == '1') {
            if (!values.name) {
              errors.name = "Required";
            }
            if (!values.server) {
              errors.server = "Required";
            }
          }
          if (activeTab == '2') {
            if (values.gatewayDiscoveryEnabled && !values.gatewayDiscoveryInterval) {
              errors.gatewayDiscoveryInterval = "Required";
            }
            if (values.gatewayDiscoveryEnabled && !values.gatewayDiscoveryTXFrequency) {
              errors.gatewayDiscoveryTXFrequency = "Required";
            }
            if (values.gatewayDiscoveryEnabled && !values.gatewayDiscoveryDR) {
              errors.gatewayDiscoveryDR = "Required";
            }
          }
        
          return errors;
        }}
        render={({ handleSubmit, form, submitting, pristine, values }) => (
          <form onSubmit={handleSubmit}>
            <Nav tabs>
              <NavItem>
                <NavLink
                  className={classnames({ active: activeTab === '1' })}
                  onClick={() => { this.toggle('1'); }}
                >
                  {i18n.t(`${packageNS}:tr000167`)}
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  className={classnames({ active: activeTab === '2' })}
                  onClick={() => { this.toggle('2'); }}
                >
                  {i18n.t(`${packageNS}:tr000095`)}
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  className={classnames({ active: activeTab === '3' })}
                  onClick={() => { this.toggle('3'); }}
                >
                  {i18n.t(`${packageNS}:tr000104`)}
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink
                  className={classnames({ active: activeTab === '4' })}
                  onClick={() => { this.toggle('4'); }}
                >
                  {i18n.t(`${packageNS}:tr000428`)}
                </NavLink>
              </NavItem>
            </Nav>
            <TabContent activeTab={activeTab}>
              <TabPane tabId="1">
                <Row>
                  <Col sm="12">
                    <FormGroup row>
                      <Field name="id">
                        {({ input, meta }) => (
                          <div>
                            <Input
                              {...input}
                              id="id"
                              name="id"
                              type="hidden"
                            />
                          </div>
                        )}
                      </Field>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="name" sm={3}>
                        {i18n.t(`${packageNS}:tr000090`)}
                      </Label>
                      <Col sm={9}>
                        <Field name="name">
                          {({ input, meta }) => (
                            <div>
                              <Input
                                {...input}
                                autoFocus
                                id="name"
                                name="name"
                                type="text"
                                invalid={meta.error && meta.touched}
                              />
                              {meta.error && meta.touched &&
                                <FormFeedback>{meta.error}</FormFeedback>
                              }
                            </div>
                          )}
                        </Field>
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000091`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="server" sm={3}>
                        {i18n.t(`${packageNS}:tr000092`)}
                      </Label>
                      <Col sm={9}>
                        <Field name="server">
                          {({ input, meta }) => (
                            <div>
                              <Input
                                {...input}
                                id="server"
                                name="server"
                                type="text"
                                invalid={meta.error && meta.touched}
                              />
                              {meta.error && meta.touched &&
                                <FormFeedback>{meta.error}</FormFeedback>
                              }
                            </div>
                          )}
                        </Field>  
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000093`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Col sm="12">
                        <br />
                        {submitButton(submitting, submitLabel)}
                      </Col>
                    </FormGroup>
                  </Col>
                </Row>
              </TabPane>
              <TabPane tabId="2">
                <Row>
                  <Col sm="12">
                    <h5>{i18n.t(`${packageNS}:tr000095`)}</h5>
                    <br />
                    <FormGroup check>      
                      <Field name="gatewayDiscoveryEnabled" type="checkbox">
                        {({ input, meta }) => (
                          <Label check for="gatewayDiscoveryEnabled">
                            <Input
                              {...input}
                              id="gatewayDiscoveryEnabled"
                              name="gatewayDiscoveryEnabled"
                              type="checkbox"
                              invalid={meta.error && meta.touched}
                              onClick={this.onChange}
                            />
                            {' '}
                            {i18n.t(`${packageNS}:tr000096`)}
                            {meta.error && meta.touched &&
                              <FormFeedback>{meta.error}</FormFeedback>
                            }
                          </Label>
                        )}
                      </Field>
                      <FormText color="muted">
                        {i18n.t(`${packageNS}:tr000097`)}
                      </FormText>
                    </FormGroup>
                    {this.state.object.gatewayDiscoveryEnabled &&
                      <>
                        <br />
                        <FormGroup row>
                          <Label for="gatewayDiscoveryInterval" sm={3}>
                            {i18n.t(`${packageNS}:tr000098`)}
                          </Label>
                          <Col sm={9}>
                            <Field name="gatewayDiscoveryInterval">
                              {({ input, meta }) => (
                                <div>
                                  <Input
                                    {...input}
                                    autoFocus
                                    id="gatewayDiscoveryInterval"
                                    invalid={meta.error && meta.touched}
                                    min="0" 
                                    name="gatewayDiscoveryInterval"
                                    type="number"
                                  />
                                  {meta.error && meta.touched &&
                                    <FormFeedback>{meta.error}</FormFeedback>
                                  }
                                </div>
                              )}
                            </Field>  
                            <FormText color="muted">
                              {i18n.t(`${packageNS}:tr000099`)}
                            </FormText>
                          </Col>
                        </FormGroup>
                        <FormGroup row>
                          <Label for="gatewayDiscoveryTXFrequency" sm={3}>
                            {i18n.t(`${packageNS}:tr000100`)}
                          </Label>
                          <Col sm={9}>
                            <Field name="gatewayDiscoveryTXFrequency">
                              {({ input, meta }) => (
                                <div>
                                  <Input
                                    {...input}
                                    id="gatewayDiscoveryTXFrequency"
                                    invalid={meta.error && meta.touched}
                                    min="0"
                                    name="gatewayDiscoveryTXFrequency"
                                    type="number"
                                  />
                                  {meta.error && meta.touched &&
                                    <FormFeedback>{meta.error}</FormFeedback>
                                  }
                                </div>
                              )}
                            </Field>
                            <FormText color="muted">
                              {i18n.t(`${packageNS}:tr000101`)}
                            </FormText>
                          </Col>
                        </FormGroup>
                        <FormGroup row>
                          <Label for="gatewayDiscoveryDR" sm={3}>
                            {i18n.t(`${packageNS}:tr000102`)}
                          </Label>
                          <Col sm={9}>
                            <Field name="gatewayDiscoveryDR">
                              {({ input, meta }) => (
                                <div>
                                  <Input
                                    {...input}
                                    id="gatewayDiscoveryDR"
                                    invalid={meta.error && meta.touched}
                                    min="0"
                                    name="gatewayDiscoveryDR"
                                    type="number"
                                  />
                                  {meta.error && meta.touched &&
                                    <FormFeedback>{meta.error}</FormFeedback>
                                  }
                                </div>
                              )}
                            </Field>
                            <FormText color="muted">
                              {i18n.t(`${packageNS}:tr000103`)}
                            </FormText>
                          </Col>
                        </FormGroup>
                      </>
                    }
                    <FormGroup row>
                      <Col sm="12">
                        <br />
                        {submitButton(submitting, submitLabel)}
                      </Col>
                    </FormGroup>
                  </Col>
                </Row>
              </TabPane>
              <TabPane tabId="3">
                <Row>
                  <Col sm="12">
                    <h5>{i18n.t(`${packageNS}:tr000105`)}</h5>
                    <br />
                    <FormGroup row>
                      <Label for="caCert" sm={3}>
                        {i18n.t(`${packageNS}:tr000106`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="caCert"
                          name="caCert"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000107`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="tlsCert" sm={3}>
                        {i18n.t(`${packageNS}:tr000110`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="tlsCert"
                          name="tlsCert"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000109`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="tlsKey" sm={3}>
                        {i18n.t(`${packageNS}:tr000108`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="tlsKey"
                          name="tlsKey"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000109`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Col sm="12">
                        <br />
                        {submitButton(submitting, submitLabel)}
                      </Col>
                    </FormGroup>
                  </Col>
                </Row>
              </TabPane>
              <TabPane tabId="4">
                <Row>
                  <Col sm="12">
                    <h5>{i18n.t(`${packageNS}:tr000427`)}</h5>
                    <br />
                    <FormGroup row>
                      <Label for="routingProfileCACert" sm={3}>
                        {i18n.t(`${packageNS}:tr000106`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="routingProfileCACert"
                          name="routingProfileCACert"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000107`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="routingProfileTLSCert" sm={3}>
                        {i18n.t(`${packageNS}:tr000110`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="routingProfileTLSCert"
                          name="routingProfileTLSCert"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000107`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Label for="routingProfileTLSKey" sm={3}>
                        {i18n.t(`${packageNS}:tr000108`)}
                      </Label>
                      <Col sm={9}>
                        <Input
                          id="routingProfileTLSKey"
                          name="routingProfileTLSKey"
                          multiline="true"
                          onChange={this.onChange}
                          rows="4"
                          type="textarea"
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000109`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Col sm="12">
                        <br />
                        {submitButton(submitting, submitLabel)}
                      </Col>
                    </FormGroup>
                  </Col>
                </Row>
              </TabPane>
            </TabContent>
          </form>
        )}
      />
    );
  }
}

export default NetworkServerForm;
