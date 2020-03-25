import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import Grid from "@material-ui/core/Grid";
import { withStyles } from "@material-ui/core/styles";
import Typography from '@material-ui/core/Typography';
import fileDownload from "js-file-download";
import ChevronDown from "mdi-material-ui/ChevronDown";
import moment from "moment";
import React, { Component } from "react";
import { Alert, Button, Col, Modal, ModalBody, ModalFooter, ModalHeader, Row } from 'reactstrap';
import mockDeviceData from '../../api/data/mockDeviceData';
import JSONTree from "../../components/JSONTree";
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import DeviceStore from "../../stores/DeviceStore";
import theme from "../../theme";
import isDev from '../../util/isDev';






const styles = {
  center: {
    textAlign: "center",
  },
  progress: {
    marginTop: 4 * theme.spacing(1),
  },
  headerColumn: {
    marginRight: 6 * theme.spacing(1),
  },
  headerColumnFixedSmall: {
    width: 145,
  },
  headerColumnFixedWide: {
    width: 175,
  },
  treeStyle: {
    paddingTop: 0,
    paddingBottom: 0,
    fontSize: 12,
    lineHeight: 1.5,
  },
};


class DeviceDataItem extends Component {
  render() {
    const receivedAt = moment(this.props.data.receivedAt).format("LTS");
    
    return(
      <ExpansionPanel>
        <ExpansionPanelSummary expandIcon={<ChevronDown />}>
          <div className={this.props.classes.headerColumnFixedSmall}><Typography variant="body2">{receivedAt}</Typography></div>
          <div className={this.props.classes.headerColumnFixedSmall}><Typography variant="body2">{this.props.data.type}</Typography></div>
        </ExpansionPanelSummary>
        <ExpansionPanelDetails>
          <Grid container spacing={4}>
            <Grid item xs className={this.props.classes.treeStyle}>
              <JSONTree data={this.props.data.payload} />
            </Grid>
          </Grid>
        </ExpansionPanelDetails>
      </ExpansionPanel>
    );
  }
}

DeviceDataItem = withStyles(styles)(DeviceDataItem);

class DeviceData extends Component {
  constructor() {
    super();

    this.state = {
      paused: false,
      connected: false,
      data: isDev ? mockDeviceData : [],
      dialogOpen: false,
    };
  }

  componentDidMount() {
    const conn = DeviceStore.getDataLogsConnection(this.props.match.params.devEUI, this.onData);
    this.setState({
      wsConn: conn,
    });

    DeviceStore.on("ws.status.change", this.setConnected);
    this.setConnected();
  }

  componentWillUnmount() {
    this.state.wsConn.close();
    DeviceStore.removeListener("ws.status.change", this.setConnected);
  }

  onDownload = () => {
    const dl = this.state.data.map((data, i) => {
      return {
        type: data.type,
        payload: data.payload,
      };
    });

    fileDownload(JSON.stringify(dl, null, 4), `device-${this.props.match.params.devEUI}.json`);
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
      data: [],
    });
  }

  setConnected = () => {
    this.setState({
      connected: DeviceStore.getWSDataStatus(),
    });
  }

  onData = (d) => {
    if (this.state.paused) {
      return;
    }

    let data = this.state.data;
    const now = new Date();

    data.unshift({
      id: now.getTime(),
      receivedAt: now,
      type: d.type,
      payload: JSON.parse(d.payloadJSON),
    });

    this.setState({
      data: data,
    });
  }

  render() {
    const { dialogOpen } = this.state;
    const data = this.state.data.map((d, i) => <DeviceDataItem key={d.id} data={d} />);
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
            {i18n.t(`${packageNS}:tr000354`)}
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
            {(this.state.connected && data.length === 0 && !this.state.paused) &&
              <div className={this.props.classes.center}>
                <Loader light />
              </div>
            }
            {data.length > 0 && data}
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(DeviceData);
