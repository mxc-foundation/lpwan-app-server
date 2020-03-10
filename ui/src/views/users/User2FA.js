import React, { Component } from "react";
import { Button, Row, Col, Card, CardBody, Modal, ModalHeader, ModalBody, CustomInput } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import Google2FA from './Google2FA';
import MneMonicPhraseConfirm from './MneMonicPhraseConfirm';
import speakeasy from 'speakeasy';


class User2FA extends Component {
  constructor(props) {
    super(props);
    this.state = {
      object: props.object || {},
      showSetup2FA: false,
      twofa_enabled: false,
      showReset2FA: false,
      showMnemonicPhraseConfirm: false,
    };

    this.showSetup2FA = this.showSetup2FA.bind(this);
    this.confirm2fa = this.confirm2fa.bind(this);
    this.skip2fa = this.skip2fa.bind(this);
    this.changeSetting = this.changeSetting.bind(this);
    this.confirmReset2fa = this.confirmReset2fa.bind(this);
    this.skipReset2fa = this.skipReset2fa.bind(this);
    this.confirmMnemonicPhraseList = this.confirmMnemonicPhraseList.bind(this);
  }

  showSetup2FA() {
    // TODO - API Call to fetch the initial code
    var secret = speakeasy.generateSecret({ length: 20 });
    console.log(secret.base32); //save to user db
    var qr = secret.google_auth_qr;
    this.setState({ showSetup2FA: true, auth_2fa_code: secret, qr: qr });
  }

  confirm2fa(confirmCode) {
    // TODO  - API call to confirm
    this.setState({ showSetup2FA: false, twofa_enabled: true });
  }

  skip2fa() {
    this.setState({ showSetup2FA: false });
  }

  confirmReset2fa(confirmCode) {
    // TODO  - API call to confirm
    // TODO - Fetch phrase - for now setting up dummy
    const phrases = ["Simba", "Sweetie", "Ziggy", "Midnight", "Kiki", "Peanut", "Midday", "Buddy", "Bently", "Gray", "Rocky", "Madison", "Bella", "Baxter"];
    this.setState({ showMnemonicPhraseConfirm: true, phrases: phrases });
  }

  skipReset2fa() {
    this.setState({ showReset2FA: false, showMnemonicPhraseConfirm: false });
  }

  changeSetting(e) {
    if (!e.target.checked) {
      // TODO - API Call to fetch the initial code
      this.setState({ showSetup2FA: false, showReset2FA: true, auth_2fa_code: '12345678' });
    }
  }

  confirmMnemonicPhraseList(phrases) {
    // TODO - API call to confirm order of phrase
    this.setState({ showSetup2FA: false, showReset2FA: false, twofa_enabled: false });
  }

  render() {
    const { object } = this.state;

    if (object === undefined) {
      return null;
    }

    return (
      <React.Fragment>
        <Card className="h-auto">
          <CardBody className="pb-0">
            <h5>{i18n.t(`${packageNS}:menu.profile_2fa.title`)}</h5>
            <Row>
              <Col>
                <p className="mt-2">{i18n.t(`${packageNS}:menu.profile_2fa.google.title`)}</p>
              </Col>
              <Col className="text-right">
                {!this.state.twofa_enabled ?
                  <Button color="primary" outline onClick={this.showSetup2FA}>{i18n.t(`${packageNS}:menu.profile_2fa.google.setup_button`)}</Button> :
                  <CustomInput type="switch" id="disable-2fa" name="2fa-setting" label=""
                    defaultChecked={this.state.twofa_enabled} checked={!this.state.showReset2FA}
                    onChange={this.changeSetting} />
                }
              </Col>
            </Row>

            {this.state.showSetup2FA ?
              <Modal isOpen={this.state.showSetup2FA} toggle={this.skip2fa} centered={true}>
                <ModalHeader toggle={this.skip2fa}>{i18n.t(`${packageNS}:menu.profile_2fa.google.2fa_title`)}</ModalHeader>
                <ModalBody>
                  <img src={this.state.qr} />
                  <Google2FA
                    title={i18n.t(`${packageNS}:menu.profile_2fa.google.2fa_instruction`)}
                    titleClass="font-weight-normal"
                    code={this.state.auth_2fa_code}
                    confirm={this.confirm2fa} skip={this.skip2fa} />
                </ModalBody>
              </Modal> : null}

            {this.state.showReset2FA ?
              <Modal isOpen={this.state.showReset2FA} toggle={this.skipReset2fa} centered={true}>
                <ModalHeader toggle={this.skipReset2fa}>{i18n.t(`${packageNS}:menu.profile_2fa.google.reset_2fa_title`)}</ModalHeader>
                <ModalBody>

                  {this.state.showMnemonicPhraseConfirm ?
                    <MneMonicPhraseConfirm
                      title={i18n.t(`${packageNS}:menu.profile_2fa.google.mnemonic_phrase_confirm_title`)}
                      titleClass="font-weight-normal" showSkipButton={true} showBackButton={false}
                      phrase={this.state.phrases} next={this.confirmMnemonicPhraseList} skip={this.skipReset2fa} /> :

                    <Google2FA
                      title={i18n.t(`${packageNS}:menu.profile_2fa.google.reset_2fa_instruction`)}
                      titleClass="font-weight-normal"
                      code={this.state.auth_2fa_code}
                      confirm={this.confirmReset2fa} skip={this.skipReset2fa} />}
                </ModalBody>
              </Modal> : null}

          </CardBody>
        </Card>
      </React.Fragment>
    );
  }
}

export default User2FA;
