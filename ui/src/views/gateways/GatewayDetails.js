import React, { Component } from "react";

import { Row, Col, Button, CustomInput, Modal, ModalBody, ModalFooter, Table, Card, CardTitle, CardBody } from 'reactstrap';
import moment from "moment";
import { Map, Marker } from 'react-leaflet';
import { Line } from "react-chartjs-2";

import i18n, { packageNS } from '../../i18n';
import MapTileLayer from "../../components/MapTileLayer";
import GatewayStore from "../../stores/GatewayStore";


class GatewayDetails extends Component {
  constructor() {
    super();
    this.state = {
      config: {},
      modal: false
    };
    this.loadStats = this.loadStats.bind(this);
  }

  componentDidMount() {
    this.loadStats();
    this.loadConfig();
  }

  toggle = () => {
    const { modal } = this.state;
    this.setState({modal: !this.state.modal});
  }

  resetGateway = () => {
    this.setState({modal: !this.state.modal});
    // submit all data for reset
  }

  loadConfig() {
    // GatewayStore.getConfig(this.props.match.params.gatewayID, resp => {
    //   this.setState({
    //     config: resp
    //   });
    // });
    this.setState({
      config: {
        manufacturer: "manufacturer",
        id: "507f298d3b1a444b",
        discoveryEnabled: false,
        name: "Test",
        altitude: "Altitude",
        description: "People can change, this is like a note",
        GPSCoordinates: "coordinates",
        lastSeen: "10 seconds ago",
        WiFiSSID: "123",
        WiFiSSIDStatus: true,
        WiFiModel: "Model",
        LANIPAddress: "255.255.255.0",
        OSVersion: "MXOS 1.1",
        firmwareVersion: "2.0",
        processes: [
          {process: "monit", status: "OK", uptime: "1m", cpuTotal: "0.6%", memoryTotal: "3.0% [7.2 MB]", read: "0.0/s", write: "0.6/s"},
          {process: "sshd", status: "OK", uptime: "224d 2h 42m", cpuTotal: "0.0%", memoryTotal: "2.8% [6.8 MB]", read: "-", write: "-"},
          {process: "postfix", status: "OK", uptime: "224d 2h 42m", cpuTotal: "0.0%", memoryTotal: "1.8% [4.3 MB]", read: "0.0/s", write: "-"},
          {process: "cron", status: "OK", uptime: "224d 2h 42m", cpuTotal: "0.0%", memoryTotal: "0.3% [740 KB]", read: "0.0/s", write: "-"},
          {process: "devd", status: "OK", uptime: "224d 2h 42m", cpuTotal: "0.0%", memoryTotal: "0.1% [268 KB]", read: "0.0/s", write: "-"},
          {process: "ntpd", status: "OK", uptime: "224d 2h 42m", cpuTotal: "0.0%", memoryTotal: "0.4% [1.0 MB]", read: "0.0/s", write: "-"},
        ],
        networks: [
          {net: "rm0", status: "OK", upload: "1.5 kB/s", download: "1.0kB/s"},
          {net: "lo0", status: "OK", upload: "651 B/s", download: "651 B/s"},
        ],
        hosts: [
          {host: "tildeslash2", status: "OK", upload: "", protocol: "[ping] [HTTP] at port 80"}
        ],
      }
    });
  }

  loadStats() {
    const end = moment().toISOString()
    const start = moment().subtract(30, "days").toISOString()

    GatewayStore.getStats(this.props.match.params.gatewayID, start, end, resp => {
      let statsDown = {
        labels: [],
        datasets: [
          {
            label: "rx received",
            borderColor: "rgba(33, 150, 243, 1)",
            backgroundColor: "rgba(0, 0, 0, 0)",
            lineTension: 0,
            pointBackgroundColor: "rgba(33, 150, 243, 1)",
            data: [],
          },
        ],
      }

      let statsUp = {
        labels: [],
        datasets: [
          {
            label: "tx emitted",
            borderColor: "rgba(33, 150, 243, 1)",
            backgroundColor: "rgba(0, 0, 0, 0)",
            lineTension: 0,
            pointBackgroundColor: "rgba(33, 150, 243, 1)",
            data: [],
          },
        ],
      }

      for (const row of resp.result) {
        statsUp.labels.push(moment(row.timestamp).format("Do"));
        statsDown.labels.push(moment(row.timestamp).format("Do"));
        statsUp.datasets[0].data.push(row.txPacketsEmitted);
        statsDown.datasets[0].data.push(row.rxPacketsReceivedOK);
      }

      this.setState({
        statsUp: statsUp,
        statsDown: statsDown,
      });
    });
  }

  render() {
    if (this.props.gateway === undefined || this.state.statsDown === undefined || this.state.statsUp === undefined) {
      return (<div></div>);
    }
    const { config, modal } = this.state;
    const style = {
      height: 322,
      zIndex: 1
    };

    const statsOptions = {
      legend: {
        display: false,
      },
      maintainAspectRatio: false,
      scales: {
        yAxes: [{
          ticks: {
            beginAtZero: true,
          },
        }],
      },
    }

    let position = [];
    if (typeof (this.props.gateway.location.latitude) !== "undefined" && typeof (this.props.gateway.location.longitude !== "undefined")) {
      position = [this.props.gateway.location.latitude, this.props.gateway.location.longitude];
    } else {
      position = [0, 0];
    }

    let lastseen = "";
    if (this.props.lastSeenAt !== undefined) {
      lastseen = moment(this.props.lastSeenAt).fromNow();
    }

    return (<React.Fragment>
      <Row>
        <Col lg={12}>
          <Card className="border shadow-none">
            <CardBody>
              <CardTitle tag="h4">{i18n.t(`${packageNS}:tr000423`)}</CardTitle>
              {config && 
               <>
              <Row>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000571`)}
                  </h6>
                  <p>
                    {config.manufacturer}
                  </p>
                </Col>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                  {i18n.t(`${packageNS}:tr000572`)} / {i18n.t(`${packageNS}:tr000074`)}
                  </h6>
                  <Row>
                    <Col lg={6}>
                      <p>
                      {config.id}
                      </p>
                    </Col>
                    <Col lg={6}>
                      <p>
                      <CustomInput type="switch" id="discoveryEnabled" label={i18n.t(`${packageNS}:tr000228`)} disabled checked={config.discoveryEnabled}/>
                      </p>
                    </Col>
                  </Row>
                </Col>
              </Row>

              <Row>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000218`)}
                  </h6>
                  <p>
                    {config.name}
                  </p>
                </Col>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                  {i18n.t(`${packageNS}:tr000573`)}
                  </h6>
                  <p>
                    {config.altitude}
                  </p>
                </Col>
              </Row>

              <Row>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000219`)}
                  </h6>
                  <p>
                    {config.description}
                  </p>
                </Col>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                  {i18n.t(`${packageNS}:tr000241`)}
                  </h6>
                  <p>
                    {config.GPSCoordinates}
                  </p>
                </Col>
              </Row>

              <Row>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000242`)}
                  </h6>
                  <p>
                    {config.lastSeen}
                  </p>
                </Col>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                  {i18n.t(`${packageNS}:tr000574`)}
                  </h6>
                  <p>
                    {config.WiFiSSID}
                    <CustomInput type="switch" id="WiFiSSIDStatus" label={""} disabled defaultChecked={config.WiFiSSIDStatus} className="ml-1 d-inline" />
                  </p>
                </Col>
              </Row>

              <Row>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000575`)}
                  </h6>
                  <p>
                    {config.WiFiModel}
                  </p>
                </Col>
                <Col lg={6}>
                  <h6 className="text-primary font-16">
                  {i18n.t(`${packageNS}:tr000576`)}
                  </h6>
                  <p>
                    {config.LANIPAddress}
                  </p>
                </Col>
              </Row>

              <Row>
                <Col lg={4}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000577`)}
                  </h6>
                  <p>
                    {config.OSVersion}
                    <Button
                      type="button"
                      color="primary"
                      className="ml-2 d-inline" 
                    >
                      {i18n.t(`${packageNS}:tr000579`)}
                    </Button>
                  </p>
                </Col>
                <Col lg={4}>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000578`)}
                  </h6>
                  <p>
                    {config.firmwareVersion}
                    <Button
                      type="button"
                      color="primary"
                      className="ml-2 d-inline" 
                    >
                      {i18n.t(`${packageNS}:tr000579`)}
                    </Button>
                  </p>
                </Col>
                <Col>
                  <h6 className="text-primary font-16">
                    {i18n.t(`${packageNS}:tr000242`)} ( {config.lastSeen} )
                  </h6>
                </Col>
              </Row>
            </>
            }
            </CardBody>
          </Card>
        </Col>
      </Row>

      <Row className="mt-2">
        <Col>
          <Card className="border shadow-none">
            <CardBody className="p-1">
              {config.processes ? <Table striped size="sm">
                <thead>
                  <tr>
                    <th>{i18n.t(`${packageNS}:tr000580`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000282`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000581`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000582`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000583`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000584`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000585`)}</th>
                  </tr>
                </thead>
                <tbody>
                {config.processes.map((process, index) => {
                  return <tr key={index}>
                    <td>{process.process}</td>
                    <td>{process.status}</td>
                    <td>{process.uptime}</td>
                    <td>{process.cpuTotal}</td>
                    <td>{process.memoryTotal}</td>
                    <td>{process.read}</td>
                    <td>{process.write}</td>
                  </tr>
                })}
                </tbody>
                </Table>: null}
            </CardBody>
          </Card>
        </Col>
      </Row>

      <Row className="mt-2">
        <Col>
          <Card className="border shadow-none">
            <CardBody className="p-1">
              {config.networks ? <Table striped size="sm">
                <thead>
                  <tr>
                    <th>{i18n.t(`${packageNS}:tr000586`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000282`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000587`)}</th>
                    <th className="text-right">{i18n.t(`${packageNS}:tr000251`)}</th>
                  </tr>
                </thead>
                <tbody>
                {config.networks.map((network, index) => {
                  return <tr key={index}>
                    <td>{network.net}</td>
                    <td>{network.status}</td>
                    <td>{network.upload}</td>
                    <td className="text-right">{network.download}</td>
                  </tr>
                })}
                </tbody>
                </Table>: null}

              {config.hosts ? <Table striped size="sm">
                <thead>
                  <tr>
                    <th>{i18n.t(`${packageNS}:tr000588`)}</th>
                    <th>{i18n.t(`${packageNS}:tr000282`)}</th>
                    <th className="text-right">{i18n.t(`${packageNS}:tr000589`)}</th>
                  </tr>
                </thead>
                <tbody>
                {config.hosts.map((host, index) => {
                  return <tr key={index}>
                    <td>{host.host}</td>
                    <td>{host.status}</td>
                    <td className="text-right">{host.protocol}</td>
                  </tr>
                })}
                </tbody>
                </Table>: null}
            </CardBody>
          </Card>
        </Col>
      </Row>

      <Row className="mt-2">
        <Col>
          <Card className="border shadow-none">
            <CardBody className="p-1">
              <Map center={position} zoom={15} style={style} animate={true} scrollWheelZoom={false}>
                <MapTileLayer />
                <Marker position={position} />
              </Map>
            </CardBody>
          </Card>
        </Col>
      </Row>

      <Row className="mt-2">
        <Col lg={12}>
          <Card className="border shadow-none">
            <CardBody>
              <CardTitle tag="h4">{i18n.t(`${packageNS}:tr000243`)}</CardTitle>

              <div style={{height: '300px'}}>
                <Line height={75} options={statsOptions} data={this.state.statsDown} redraw />
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>

      <Row className="mt-2">
        <Col lg={12}>
          <Card className="border shadow-none">
            <CardBody>
              <CardTitle tag="h4">{i18n.t(`${packageNS}:tr000244`)}</CardTitle>
              <div style={{ height: '300px' }}>
                <Line height={75} options={statsOptions} data={this.state.statsUp} redraw />
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>
      <Row className="mt-2">
        <Col>
          <h6>{i18n.t(`${packageNS}:tr000590`)}</h6>

          <Row>
            <Col className="col-auto align-self-center">
              <CustomInput type="switch" id="discoveryEnabled" label={i18n.t(`${packageNS}:tr000228`)} disabled checked={config.discoveryEnabled} />
            </Col>
            <Col className="col-auto align-self-center">
              <CustomInput type="checkbox" id="agreeTerms" label={i18n.t(`${packageNS}:tr000591`)} />
            </Col>
            <Col className="col-auto">
              <Button
              type="button"
              color="primary"
              className="ml-2 d-inline" 
            >
              {i18n.t(`${packageNS}:tr000592`)}
            </Button>
            <Button
              type="button"
              color="primary"
              className="ml-2 d-inline" 
            >
              {i18n.t(`${packageNS}:tr000593`)}
            </Button>
            <Button
              type="button"
              color="primary"
              onClick={this.toggle}
              className="ml-2 d-inline mt-2 mt-sm-0" 
            >
              {i18n.t(`${packageNS}:tr000594`)}
            </Button>
            </Col>
          </Row>
        </Col>
      </Row>
      <Modal isOpen={modal} toggle={this.toggle}>
        <ModalBody className="text-center">
          {i18n.t(`${packageNS}:tr000595`)} <br />
          <small>{i18n.t(`${packageNS}:tr000596`)}</small>
        </ModalBody>
        <ModalFooter>
          <Button color="secondary" onClick={this.toggle}>Cancel</Button>
          <Button color="primary" onClick={this.resetGateway}>Confirm</Button>
        </ModalFooter>
      </Modal>
    </React.Fragment>
    );
  }
}

export default GatewayDetails;
