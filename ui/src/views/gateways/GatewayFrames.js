import React, { Component } from "react";

import { Row, Col, Button as RButton, UncontrolledAlert } from 'reactstrap';

import fileDownload from "js-file-download";

import i18n, { packageNS } from '../../i18n';
import LoRaWANFrameLog from "../../components/LoRaWANFrameLog";
import CommonModal from '../../components/Modal';
import GatewayStore from "../../stores/GatewayStore";


class GatewayFrames extends Component {
  constructor() {
    super();

    this.state = {
      connected: false,
      paused: false,
      frames: [],
      dialogOpen: false,
    };

    this.onFrame = this.onFrame.bind(this);
    this.onDownload = this.onDownload.bind(this);
    this.togglePause = this.togglePause.bind(this);
    this.onClear = this.onClear.bind(this);
    this.toggleHelpDialog = this.toggleHelpDialog.bind(this);
    this.setConnected = this.setConnected.bind(this);
  }

  componentDidMount() {
    const conn = GatewayStore.getFrameLogsConnection(this.props.gateway.id, () => { }, () => { }, this.onFrame);
    this.setState({
      wsConn: conn,
    });

    GatewayStore.on("ws.status.change", this.setConnected);
    this.setConnected();
  }

  componentWillUnmount() {
    this.state.wsConn.close();
    GatewayStore.removeListener("ws.status.change", this.setConnected);
  }

  onDownload() {
    const dl = this.state.frames.map((frame, i) => {
      return {
        uplinkMetaData: frame.uplinkMetaData,
        downlinkMetaData: frame.downlinkMetaData,
        phyPayload: frame.phyPayload,
      };
    });

    fileDownload(JSON.stringify(dl, null, 4), `gateway-${this.props.match.params.gatewayID}.json`);
  }

  togglePause() {
    this.setState({
      paused: !this.state.paused,
    });
  }

  toggleHelpDialog() {
    this.setState({
      dialogOpen: !this.state.dialogOpen,
    });
  }

  onClear() {
    this.setState({
      frames: [],
    });
  }

  setConnected() {
    this.setState({
      connected: GatewayStore.getWSStatus(),
    });
  }

  onFrame(frame) {
    if (this.state.paused) {
      return;
    }

    let frames = this.state.frames;
    const now = new Date();

    if (frame.uplinkFrame !== undefined) {
      frames.unshift({
        id: now.getTime(),
        receivedAt: now,
        uplinkMetaData: {
          rxInfo: frame.uplinkFrame.rxInfo,
          txInfo: frame.uplinkFrame.txInfo,
        },
        phyPayload: JSON.parse(frame.uplinkFrame.phyPayloadJSON),
      });
    }

    if (frame.downlinkFrame !== undefined) {
      frames.unshift({
        id: now.getTime(),
        receivedAt: now,
        downlinkMetaData: {
          txInfo: frame.downlinkFrame.txInfo,
        },
        phyPayload: JSON.parse(frame.downlinkFrame.phyPayloadJSON),
      });
    }

    this.setState({
      frames: frames,
    });
  }

  render() {
    const frames = this.state.frames.map((frame, i) => <LoRaWANFrameLog key={frame.id} frame={frame} />);
    
    return (<React.Fragment>
      <Row>
        <Col className="text-right">
          <div className="button-list">
            <CommonModal buttonLabel={<React.Fragment><i className="mdi mdi-help-circle mr-2"></i>{i18n.t(`${packageNS}:tr000248`)}</React.Fragment>}
              outline={true}
              buttonColor={"info"} callback={() => { }}
              context={i18n.t(`${packageNS}:tr000249`)} title={i18n.t(`${packageNS}:tr000248`)}
              showConfirmButton={false} left={i18n.t(`${packageNS}:tr000430`)}></CommonModal>

            {!this.state.paused && <RButton outline color="primary" onClick={this.togglePause}>
              <i className="mdi mdi-pause mr-2"></i>
              {i18n.t(`${packageNS}:tr000250`)}
            </RButton>}

            {this.state.paused && <RButton outline color="primary" onClick={this.togglePause}>
              <i className="mdi mdi-play mr-2"></i>
              {i18n.t(`${packageNS}:tr000355`)}
            </RButton>}

            <RButton outline color="secondary" onClick={this.onDownload}>
              <i className="mdi mdi-download mr-2"></i>
              {i18n.t(`${packageNS}:tr000251`)}
            </RButton>

            <RButton outline color="danger" onClick={this.onClear}>
              <i className="mdi mdi-delete mr-2"></i>
              {i18n.t(`${packageNS}:tr000252`)}
            </RButton>
          </div>
        </Col>
      </Row>

      <Row>
        <Col>
          {!this.state.connected && <div className="mt-3 mb-2">
            <UncontrolledAlert color="info">
              <p className="font-15 mb-0"><i className="mdi mdi-alert-circle mr-1"></i>
                {i18n.t(`${packageNS}:tr000392`)}</p>
            </UncontrolledAlert>
          </div>}

          {(this.state.connected && frames.length === 0 && !this.state.paused) &&
            <div className="text-center mt-2 mb-2">
              <div className="spinner-border text-primary" role="status">
                <span className="sr-only">...</span>
              </div>
            </div>}
        </Col>
      </Row>

      <Row>
        <Col className="mb-0">
          {frames.length > 0 && frames}
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default GatewayFrames;
