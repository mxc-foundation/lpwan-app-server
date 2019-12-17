import React from "react";
import {
  Button, Col, Form, FormGroup, FormText, Input, Label, 
  TabContent, TabPane, Nav, NavItem, NavLink, Row,
} from 'reactstrap';
import classnames from 'classnames';

import i18n, { packageNS } from '../../i18n';
import FormComponent from "../../classes/FormComponent";

const submitButton = () => {
  return (
    <Button
      aria-label={i18n.t(`${packageNS}:tr000277`)}
      block
      color="primary"
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

    if (this.state.object === undefined) {
      return(null);
    }

    return(
      <Form>
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
                    <Input
                      id="name"
                      name="email"
                      onChange={this.onChange}
                      placeholder={i18n.t(`${packageNS}:tr000091`)}
                      type="text"
                      value={this.state.object.name || ""}
                    />
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
                    <Input
                      id="server"
                      name="server"
                      onChange={this.onChange}
                      placeholder={i18n.t(`${packageNS}:tr000093`)}
                      type="text"
                      value={this.state.object.server || ""}
                    />
                    <FormText color="muted">
                      {i18n.t(`${packageNS}:tr000093`)}
                    </FormText>
                  </Col>
                </FormGroup>
                <FormGroup row>
                  <Col sm="12">
                    <br />
                    {submitButton()}
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
                  <Label check for="gatewayDiscoveryEnabled">
                    <Input
                      checked={!!this.state.object.gatewayDiscoveryEnabled}
                      id="gatewayDiscoveryEnabled"
                      name="gatewayDiscoveryEnabled"
                      onChange={this.onChange}
                      placeholder={i18n.t(`${packageNS}:tr000097`)}
                      type="checkbox"
                    />
                    {i18n.t(`${packageNS}:tr000096`)}
                  </Label>
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
                        <Input
                          id="gatewayDiscoveryInterval"
                          name="gatewayDiscoveryInterval"
                          onChange={this.onChange}
                          placeholder={i18n.t(`${packageNS}:tr000099`)}
                          type="number"
                          value={this.state.object.gatewayDiscoveryInterval || 0}
                        />
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
                        <Input
                          id="gatewayDiscoveryTXFrequency"
                          name="gatewayDiscoveryTXFrequency"
                          onChange={this.onChange}
                          placeholder={i18n.t(`${packageNS}:tr000101`)}
                          type="number"
                          value={this.state.object.gatewayDiscoveryTXFrequency || 0}
                        />
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
                        <Input
                          id="gatewayDiscoveryDR"
                          name="gatewayDiscoveryDR"
                          onChange={this.onChange}
                          placeholder={i18n.t(`${packageNS}:tr000103`)}
                          type="number"
                          value={this.state.object.gatewayDiscoveryDR || 0}
                        />
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
                    {submitButton()}
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
                <h5>{i18n.t(`${packageNS}:tr000105`)}</h5>
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
                    {submitButton()}
                  </Col>
                </FormGroup>
              </Col>
            </Row>
          </TabPane>
        </TabContent>
      </Form>
    );
  }
}

export default NetworkServerForm;
