import classNames from "classnames";
import React, { Component } from "react";
import { Button, CustomInput, Modal, ModalBody, ModalHeader } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import MneMonicPhrase from '../users/MneMonicPhrase';
import MneMonicPhraseConfirm from '../users/MneMonicPhraseConfirm';



/**
 * First intro step for new feature
 */
class Intro extends Component {
  constructor(props) {
    super(props);

    this.state = {
      agree: false
    }
  }

  render() {
    return <React.Fragment>
      <div className="text-center">
        <i className="mdi mdi-information text-primary display-3"></i>

        <h3>{i18n.t(`${packageNS}:menu.dashboard_2fa.title`)}</h3>
        <p>{i18n.t(`${packageNS}:menu.dashboard_2fa.description`)}</p>

        <div className="mt-4 text-left">
          <CustomInput type="checkbox" id="agree-check" checked={this.state.agree} label={i18n.t(`${packageNS}:menu.dashboard_2fa.agree_instruction`)}
            onChange={(e) => this.setState({ agree: e.target.checked })} />
        </div>
        <div className="mt-3">
          <Button color="primary" disabled={!this.state.agree} onClick={this.props.next}>{i18n.t(`${packageNS}:menu.dashboard_2fa.next_button`)}</Button>
          <Button color="link" className="btn-block pt-0" onClick={this.props.skip}>{i18n.t(`${packageNS}:menu.dashboard_2fa.skip_button`)}</Button>
        </div>
      </div>
    </React.Fragment>
  }
}

/**
 * Confirmed step
 */
const Confirmed = ({ close }) => {
  return <React.Fragment>
    <div className="text-center">
      <i className="mdi mdi-check-circle-outline text-success display-3"></i>

      <h3>{i18n.t(`${packageNS}:menu.dashboard_2fa.menmonic_step_3_title`)}</h3>

      <div className="mt-3">
        <Button color="primary" onClick={close}>{i18n.t(`${packageNS}:menu.dashboard_2fa.done_button`)}</Button>
      </div>
    </div>
  </React.Fragment>
}

/**
 * Simple steps
 * @param {*} param0 
 */
const Steps = ({ completedSteps = [] }) => {
  return <div className="mb-3 px-5">
    <ul className="stepper">
      <li className={classNames({ "active": completedSteps.indexOf(1) !== -1 })}></li>
      <li className={classNames({ "active": completedSteps.indexOf(2) !== -1 })}></li>
      <li className={classNames({ "active": completedSteps.indexOf(3) !== -1 })}></li>
    </ul>
  </div>
}


class Feature extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showModal: true,
      showIntro: true,
      showMneMonicList: false,
      showMneMonicListConfirm: false,
      showFinal: false,
      completedSteps: []
    };

    this.showMnemonicList = this.showMnemonicList.bind(this);
    this.showMnemonicListConfirm = this.showMnemonicListConfirm.bind(this);
    this.showConfirmed = this.showConfirmed.bind(this);
    this.closeModal = this.closeModal.bind(this);
  }

  showMnemonicList() {
    // TODO - Fetch phrase - for now setting up dummy
    const phrases = ["Simba", "Sweetie", "Ziggy", "Midnight", "Kiki", "Peanut", "Midday", "Buddy", "Bently", "Gray", "Rocky", "Madison", "Bella", "Baxter"];
    this.setState({ showIntro: false, showMneMonicList: true, phrases: phrases, showMneMonicListConfirm: false, completedSteps: [1] });
  }

  showMnemonicListConfirm() {
    this.setState({ showIntro: false, showMneMonicList: false, showMneMonicListConfirm: true, completedSteps: [1, 2] });
  }

  showConfirmed() {
    // TODO - Make an api call to confirm order
    this.setState({ showIntro: false, showMneMonicList: false, showMneMonicListConfirm: false, showFinal: true, completedSteps: [1, 2, 3] });
  }

  closeModal() {
    this.setState({ showModal: false });
  }

  render() {

    return (
      <React.Fragment>

        <Modal isOpen={this.state.showModal} toggle={this.closeModal} centered={true}>
          <ModalHeader toggle={this.closeModal} className="border-0"></ModalHeader>
          <ModalBody className="pb-4">
            {!this.state.showIntro ? <Steps completedSteps={this.state.completedSteps} /> : null}

            {this.state.showIntro ? <Intro next={this.showMnemonicList} skip={this.closeModal} /> : null}

            {this.state.showMneMonicList ? <MneMonicPhrase title={i18n.t(`${packageNS}:menu.dashboard_2fa.menmonic_step_1_title`)}
              phrase={this.state.phrases} showSkip={true} titleClass="h3"
              next={this.showMnemonicListConfirm} close={this.closeModal} /> : null}

            {this.state.showMneMonicListConfirm ? <MneMonicPhraseConfirm title={i18n.t(`${packageNS}:menu.dashboard_2fa.menmonic_step_2_title`)}
              phrase={this.state.phrases} showSkipButton={true} showBackButton={true} titleClass="h3"
              next={this.showConfirmed} back={this.showMnemonicList} skip={this.closeModal} /> : null}

            {this.state.showFinal ? <Confirmed close={this.closeModal} /> : null}

          </ModalBody>
        </Modal>
      </React.Fragment>
    );
  }
}

export default Feature;
