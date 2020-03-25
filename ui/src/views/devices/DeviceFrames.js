import { withStyles } from "@material-ui/core/styles";
import fileDownload from "js-file-download";
import React, { Component } from "react";
import { Alert, Button, Col, Modal, ModalBody, ModalFooter, ModalHeader, Row } from 'reactstrap';
import mockDeviceFrame from '../../api/data/mockDeviceFrame';
import Loader from "../../components/Loader";
import LoRaWANFrameLog from "../../components/LoRaWANFrameLog";
import i18n, { packageNS } from '../../i18n';
import DeviceStore from "../../stores/DeviceStore";
import theme from "../../theme";
import isDev from '../../util/isDev';






const styles = {
  buttons: {
    textAlign: "right",
  },
  button: {
    marginLeft: 2 * theme.spacing(1),
  },
  icon: {
    marginRight: theme.spacing(1),
  },
  center: {
    textAlign: "center",
  },
  progress: {
    marginTop: 4 * theme.spacing(1),
  },
};

class DeviceFrames extends Component {
  constructor() {
    super();

    this.state = {
      paused: false,
      connected: false,
      frames: [],
      dialogOpen: false,
    };
  }

  componentDidMount() {
    const conn = DeviceStore.getFrameLogsConnection(this.props.match.params.devEUI, this.onFrame);
    /* if (isDev) {
      this.onFrame(mockDeviceFrame);
    } */
    this.setState({
      wsConn: conn,
    });

    DeviceStore.on("ws.status.change", this.setConnected);
    this.setConnected();
  }

  componentWillUnmount() {
    /* if (isDev) {
      return;
    } */
    this.state.wsConn.close();

    DeviceStore.removeListener("ws.status.change", this.setConnected);
  }

  onDownload = () => {
    const dl = this.state.frames.map((frame, i) => {
      return {
        uplinkMetaData: frame.uplinkMetaData,
        downlinkMetaData: frame.downlinkMetaData,
        phyPayload: frame.phyPayload,
      };
    });

    fileDownload(JSON.stringify(dl, null, 4), `gateway-${this.props.match.params.gatewayID}.json`);
  }

  togglePause = () => {
    this.setState({
      paused: !this.state.paused,
    });
  }

  toggleHelpDialog = () => {
    this.setState({
      dialogOpen: !this.state.dialogOpen,
    });
  }

  onClear = () => {
    this.setState({
      frames: [],
    });
  }

  setConnected = () => {
    this.setState({
      connected: DeviceStore.getWSFramesStatus(),
    });
  }

  onFrame = (frame) => {
    let _frame = isDev ? mockDeviceFrame : frame;

    if (this.state.paused) {
      return;
    }

    let frames = this.state.frames;
    const now = new Date();

    if (_frame.uplinkFrame !== undefined) {
      frames.unshift({
        id: now.getTime(),
        receivedAt: now,
        uplinkMetaData: {
          rxInfo: _frame.uplinkFrame.rxInfo,
          txInfo: _frame.uplinkFrame.txInfo,
        },
        phyPayload: JSON.parse(_frame.uplinkFrame.phyPayloadJSON),
      });
    }

    if (_frame.downlinkFrame !== undefined) {
      frames.unshift({
        id: now.getTime(),
        receivedAt: now,
        downlinkMetaData: {
          txInfo: _frame.downlinkFrame.txInfo,
        },
        phyPayload: JSON.parse(_frame.downlinkFrame.phyPayloadJSON),
      });
    }

    this.setState({
      frames: frames,
    });
  }

  render() {
    const { dialogOpen } = this.state;
    const frames = this.state.frames.map((frame, i) => <LoRaWANFrameLog key={frame.id} frame={frame} />);
    const closeBtn = <button className="close" onClick={this.toggleHelpDialog}>&times;</button>;

    return(
      <React.Fragment>
        <Modal
          isOpen={dialogOpen}
          toggle={this.toggleHelpDialog}
          aria-labelledby="help-dialog-title"
          aria-describedby="help-dialog-description"
        >
          <ModalHeader
            toggle={this.toggleHelpDialog}
            close={closeBtn}
            id="help-dialog-title"
          >
            {i18n.t(`${packageNS}:tr000248`)}
          </ModalHeader>
          <ModalBody id="help-dialog-description">
            {i18n.t(`${packageNS}:tr000249`)}
          </ModalBody>
          <ModalFooter>
            <Button color="primary" onClick={this.toggleHelpDialog}>Close</Button>{' '}
          </ModalFooter>
        </Modal>

        <Row xs={1}>
          <Col xs={6} sm={{ size: 3, offset: 3 }}>
            <Button
              variant="outlined"
              className={this.props.classes.button}
              onClick={this.toggleHelpDialog}
            >
              <span style={{ display: "flex" }}>
                <i className="mdi mdi-help"></i>&nbsp;{i18n.t(`${packageNS}:tr000248`)}
              </span>
            </Button>
          </Col>
          <Col xs={6} sm={{ size: 3, offset: 0 }}>
            {!this.state.paused &&
              <Button variant="outlined" className={this.props.classes.button} onClick={this.togglePause}>
                <span style={{ display: "flex" }}>
                  <i className="mdi mdi-pause"></i>&nbsp;{i18n.t(`${packageNS}:tr000250`)}
                </span>
              </Button>
            }
            {this.state.paused &&
              <Button variant="outlined" className={this.props.classes.button} onClick={this.togglePause}>
                <span style={{ display: "flex" }}>
                  <i className="mdi mdi-play"></i>&nbsp;{i18n.t(`${packageNS}:tr000355`)}
                </span>
              </Button>
            }
          </Col>
          <Col xs={12} sm={0}><br /></Col>
          <Col xs={6} sm={{ size: 3, offset: 3 }}>
            <Button variant="outlined" className={this.props.classes.button} onClick={this.onDownload}>
              <span style={{ display: "flex" }}>
                <i className="mdi mdi-download"></i>&nbsp;{i18n.t(`${packageNS}:tr000251`)}
              </span>
            </Button>
          </Col>
          <Col xs={6} sm={{ size: 3, offset: 0 }}>
            <Button variant="outlined" className={this.props.classes.button} color="secondary" onClick={this.onClear}>
              <span style={{ display: "flex" }}>
                <i className="mdi mdi-delete"></i>&nbsp;{i18n.t(`${packageNS}:tr000252`)}
              </span>
            </Button>
          </Col>
        </Row>

        <Row xs={1}>
          <Col xs={12} sm={0}><br /></Col>
          <Col xs={12}>
            {!this.state.connected &&
              <div className={this.props.classes.center}>
                <Alert color="info" style={{ fontSize: "1.25em" }}>
                <i className="mdi mdi-information-outline mr-1"></i>&nbsp;{i18n.t(`${packageNS}:tr000392`)}
                </Alert>
              </div>
            }
            <br />
            {(this.state.connected && frames.length === 0 && !this.state.paused) &&
              <div className={this.props.classes.center}>
                <Loader light />
              </div>
            }
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

export default withStyles(styles)(DeviceFrames);
