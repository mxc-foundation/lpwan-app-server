import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import QRCode from "qrcode.react";
import { Row, Col, Button, UncontrolledTooltip, Modal as RModal, ModalBody, ModalHeader } from 'reactstrap';

import Modal from "../../components/Modal";
import Loader from "../../components/Loader";
import TopupStore from "../../stores/TopupStore";
import MoneyStore from "../../stores/MoneyStore";
import SessionStorage from "../../stores/SessionStore";
import { ETHER } from "../../util/CoinType"

import i18n, { packageNS } from '../../i18n';

function loadSuperNodeActiveMoneyAccount(organizationID) {
  return new Promise((resolve, reject) => {
    TopupStore.getTopUpDestination(organizationID, resp => {
      resolve(resp.activeAccount);
    }, reject);
  });
}

function loadActiveMoneyAccount(organizationID) {
  return new Promise((resolve, reject) => {
    TopupStore.getTopUpDestination(organizationID, resp => {
      resolve(resp.activeAccount);
    }, reject);
  });
}

/*function loadActiveMoneyAccount(organizationID) {
  return new Promise((resolve, reject) => {
    MoneyStore.getActiveMoneyAccount(ETHER, organizationID, resp => {
      resolve(resp.activeAccount);
    }, reject);
  });
}*/


class TopupCrypto extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
      showCopied: false,
      showQRCode: false,
      nsDialog: false,
      description: null,
      object: this.props.object || {},
    };

    this.copyAddr = this.copyAddr.bind(this);
    this.toggleQRCodeModal = this.toggleQRCodeModal.bind(this);
    this.downloadQRCode = this.downloadQRCode.bind(this);

    this.downloadQRRef = React.createRef();
  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
    else {
      this.loadData();
    }
  }

  loadData = async () => {
    try {
      const organizationID = this.props.match.params.organizationID;

      this.setState({ loading: true });
      var superNodeAccount = await loadSuperNodeActiveMoneyAccount(organizationID);
      var account = await loadActiveMoneyAccount(organizationID);

      const accounts = {};
      accounts.superNodeAccount = superNodeAccount;
      accounts.account = account;


      let object = this.state.object;
      object.accounts = {
        superNodeAccount: superNodeAccount,
        account: account,
      }

      if (SessionStorage.getUser().isAdmin && !superNodeAccount) {
        this.showModal(true);
      }

      if (!accounts.account && !SessionStorage.getUser().isAdmin) {
        this.showModal(true);
      }

      let description = '';
      if (SessionStorage.getUser().isAdmin) {
        description = i18n.t(`${packageNS}:menu.topup.notice001`) + " " + i18n.t(`${packageNS}:menu.topup.notice003`);
      } else {
        description = i18n.t(`${packageNS}:menu.topup.notice002`) + " " + i18n.t(`${packageNS}:menu.topup.notice003`);
      }
      
      this.setState({
        object: object,
        description
      });

      this.setState({ loading: false });
    } catch (error) {
      this.setState({ loading: false, error });
    }
  }

  showModal = (nsDialog) => {
    this.setState({ nsDialog });
  }

  handleLink = () => {
    //window.location.replace(`http://wallet.mxc.org/`);
    this.props.history.push(this.props.path);
  }

  /**
   * Copy address into clipboard
   */
  copyAddr (address) {
    navigator.clipboard.writeText(address);
    this.setState({showCopied: true});
  }

  /**
   * Toggels the qr code modal
   */
  toggleQRCodeModal() {
    this.setState({ showQRCode: !this.state.showQRCode });
  }

  /**
   * Download QR code
   */
  downloadQRCode() {
    const canvas = document.querySelector('.qr-code-container > canvas');
    if (canvas) {
      var link = document.createElement('a');
      link.download = 'mxc-qr.png';
      link.href = canvas.toDataURL();
      link.click();
    }
  }

  render() {
    
    let accounts = {};
    if (this.state.object.accounts !== undefined) {
      if (this.state.object.accounts.superNodeAccount !== undefined) {
        accounts.superNodeAccount = this.state.object.accounts.superNodeAccount;
      } else {
        accounts.superNodeAccount = '';
      }
    }

    return (<React.Fragment>

      {this.state.nsDialog && <Modal
        title={i18n.t(`${packageNS}:menu.topup.notice`)}
        left={"DISMISS"}
        right={"ADD ETH ACCOUNT"}
        context={this.state.description}
        callback={this.handleLink} />}

      <div className="position-relative">
        {this.state.loading ? <Loader /> : null}

        <Row>
          <Col className="mb-0">
            <h6 className="font-weight-normal font-14 mb-3">{i18n.t(`${packageNS}:menu.topup.eth_address`)}:</h6>
            <h4>{accounts.superNodeAccount}</h4>

            <div className="mt-3">
              <Button color="primary" className='mr-2' onClick={() => this.copyAddr(accounts.superNodeAccount || "")}
                id="copy-btn">
                {i18n.t(`${packageNS}:menu.topup.copy_addr`)}
              </Button>
              
              <UncontrolledTooltip placement="bottom" target="copy-btn">
                {this.state.showCopied ? i18n.t(`${packageNS}:menu.topup.copied_notice`) : i18n.t(`${packageNS}:menu.topup.click_to_copy`)}
              </UncontrolledTooltip>

              <Button color="primary" onClick={this.toggleQRCodeModal}>{i18n.t(`${packageNS}:menu.topup.show_qr_code`)}</Button>

              <h5 className="mt-4">{i18n.t(`${packageNS}:menu.topup.instruction001`)}</h5>
              <p className="text-muted">{i18n.t(`${packageNS}:menu.topup.instruction002`)}</p>
              <h5 className="mb-0">{i18n.t(`${packageNS}:menu.topup.instruction003`)}</h5>
            </div>
          </Col>
        </Row>
      </div>

      <RModal isOpen={this.state.showQRCode} toggle={this.toggleQRCodeModal} centered={true}>
        <ModalBody className="text-center">
          <ModalHeader toggle={this.toggleQRCodeModal} className="border-0"></ModalHeader>

          <Row className="mt-2">
            <Col className="qr-code-container">
              {accounts.superNodeAccount ? <QRCode value={accounts.superNodeAccount} size={280} level={'H'} />: null}
            </Col>
          </Row>

          <Row className="mt-2">
            <Col>
              <h5>{accounts.superNodeAccount}</h5>
              <Button className="mt-2" color="primary" onClick={this.downloadQRCode}>
                {i18n.t(`${packageNS}:menu.topup.download_qr_code`)}</Button>
            </Col>
          </Row>
        </ModalBody>
      </RModal>
    </React.Fragment>
    );
  }
}

export default withRouter(TopupCrypto);
