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

const submitButton = (submitting) => {
  return (
    <Button
      aria-label={i18n.t(`${packageNS}:tr000277`)}
      block
      color="primary"
      disabled={submitting}
      size="md"
      type="submit"
    >
      {i18n.t(`${packageNS}:tr000277`)}
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
    const { onSubmit } = this.props;

    if (this.state.object === undefined) {
      return(null);
    }

    return(
      <Form
        onSubmit={onSubmit}
        validate={values => {
          // console.log('validateForm values/activeTab: ', values);
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
            </Nav>
            <TabContent activeTab={activeTab}>
              <TabPane tabId="1">
                <Row>
                  <Col sm="12">
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
                                id="name"
                                name="name"
                                placeholder={i18n.t(`${packageNS}:tr000091`)}
                                type="text"
                                value={this.state.object.name}
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
                                placeholder={i18n.t(`${packageNS}:tr000093`)}
                                type="text"
                                value={this.state.object.server}
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
                        {submitButton(submitting)}
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
                              checked={!!this.state.object.gatewayDiscoveryEnabled}
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
                                    id="gatewayDiscoveryInterval"
                                    invalid={meta.error && meta.touched}
                                    min="0" 
                                    name="gatewayDiscoveryInterval"
                                    placeholder={i18n.t(`${packageNS}:tr000099`)}
                                    type="number"
                                    value={this.state.object.gatewayDiscoveryInterval}
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
                                    placeholder={i18n.t(`${packageNS}:tr000101`)}
                                    type="number"
                                    value={this.state.object.gatewayDiscoveryTXFrequency}
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
                                    placeholder={i18n.t(`${packageNS}:tr000103`)}
                                    type="number"
                                    value={this.state.object.gatewayDiscoveryDR}
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
                        {submitButton(submitting)}
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
                          placeholder={i18n.t(`${packageNS}:tr000107`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.caCert || ""}
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
                          placeholder={i18n.t(`${packageNS}:tr000109`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.tlsCert || ""}
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
                          placeholder={i18n.t(`${packageNS}:tr000109`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.tlsKey || ""}
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000109`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <br />
                    <h5>{i18n.t(`${packageNS}:tr000421`)}</h5>
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
                          placeholder={i18n.t(`${packageNS}:tr000107`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.routingProfileCACert || ""}
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
                          placeholder={i18n.t(`${packageNS}:tr000107`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.routingProfileTLSCert || ""}
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
                          placeholder={i18n.t(`${packageNS}:tr000109`)}
                          rows="4"
                          type="textarea"
                          value={this.state.object.routingProfileTLSKey || ""}
                        />
                        <FormText color="muted">
                          {i18n.t(`${packageNS}:tr000109`)}
                        </FormText>
                      </Col>
                    </FormGroup>
                    <FormGroup row>
                      <Col sm="12">
                        <br />
                        {submitButton(submitting)}
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
