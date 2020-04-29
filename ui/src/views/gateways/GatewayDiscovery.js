import L from 'leaflet';
import React, { Component } from "react";
import ReactDOM from 'react-dom';
import { Map, MapControl, Marker, Polyline, Popup, withLeaflet } from 'react-leaflet';
import { Card, CardBody, Col, Row } from 'reactstrap';
import MapTileLayer from "../../components/MapTileLayer";
import i18n, { packageNS } from '../../i18n';
import GatewayStore from "../../stores/GatewayStore";




class GatewayDiscovery extends Component {
  constructor() {
    super();
    this.state = {};
  }

  componentDidMount() {
    GatewayStore.getLastPing(this.props.gateway.id, ping => {
      this.setState({
        ping: ping,
      });
    });
  }

  getColor(dbm) {
    if (dbm >= -100) {
      return "#FF0000";
    } else if (dbm >= -105) {
      return "#FF7F00";
    } else if (dbm >= -110) {
      return "#FFFF00";
    } else if (dbm >= -115) {
      return "#00FF00";
    } else if (dbm >= -120) {
      return "#00FFFF";
    }
    return "#0000FF";
  }

  render() {
    if (this.state.ping === undefined || this.state.ping.pingRX.length === 0) {
      return (<Row>
        <Col>
          <Card className="shadow-none border">
            <CardBody>
              <p className="font-16">{i18n.t(`${packageNS}:tr000246`)}</p>
              <ul>
                <li>{i18n.t(`${packageNS}:tr000329`)}</li>
                <li>{i18n.t(`${packageNS}:tr000330`)}</li>
                <li>{i18n.t(`${packageNS}:tr000331`)}</li>
              </ul>
            </CardBody>
          </Card>
        </Col>
      </Row>
      );
    }

    let position = [0, 0];
    if (this.props.gateway.location !== undefined && this.props.gateway.location.latitude !== undefined && this.props.gateway.location.longitude !== undefined) {
      position = [this.props.gateway.location.latitude, this.props.gateway.location.longitude];
    }

    const style = {
      height: 800,
    };

    let bounds = [];
    let markers = [];
    let lines = [];

    markers.push(
      <Marker position={position} key={`gw-${this.props.gateway.id}`}>
        <Popup>
          <span>
            {this.props.gateway.id}<br />
            Freq: {this.state.ping.frequency / 1000000} MHz<br />
            DR: {this.state.ping.dr}<br />
            Altitude: {this.props.gateway.location.altitude} meter(s)
          </span>
        </Popup>
      </Marker>
    );

    bounds.push(position);

    for (const rx of this.state.ping.pingRX) {
      const pingPos = [rx.latitude, rx.longitude];

      markers.push(
        <Marker position={pingPos} key={`gw-${rx.gatewayID}`}>
          <Popup>
            <span>
              {rx.gatewayID}<br />
              RSSI: {rx.rssi} dBm<br />
              SNR: {rx.LoRaSNR} dB<br />
              Altitude: {rx.altitude} meter(s)
            </span>
          </Popup>
        </Marker>
      );

      bounds.push(pingPos);

      lines.push(
        <Polyline
          key={`line-${rx.gatewayID}`}
          positions={[position, pingPos]}
          color={this.getColor(rx.rssi)}
          opacity={.7}
          weight={3}
        />
      );
    }

    return (
      <Row>
        <Col>
          <Card className="shadow-none border">
            <CardBody>
              <Map bounds={bounds} maxZoom={19} style={style} animate={true} scrollWheelZoom={false}>
                <MapTileLayer />
                {markers}
                {lines}
                <LegendControl className={this.props.classes.mapLegend}>
                  <ul className={this.props.classes.mapLegendList}>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-100) }}>&nbsp;</span> &gt;= -100 dBm</li>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-105) }}>&nbsp;</span> &gt;= -105 dBm</li>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-110) }}>&nbsp;</span> &gt;= -110 dBm</li>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-115) }}>&nbsp;</span> &gt;= -115 dBm</li>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-120) }}>&nbsp;</span> &gt;= -120 dBm</li>
                    <li className={this.props.classes.mapLegendListItem}><span className={this.props.classes.label} style={{ background: this.getColor(-121) }}>&nbsp;</span> &lt; -120 dBm</li>
                  </ul>
                </LegendControl>
              </Map>
            </CardBody>
          </Card>
        </Col>
      </Row>
    );
  };
}

class LegendControl extends MapControl {
  componentWillMount() {
    const legend = L.control({ position: "bottomleft" });
    const jsx = (
      <div {...this.props}>
        {this.props.children}
      </div>
    );

    legend.onAdd = function (map) {
      let div = L.DomUtil.create("div", '');
      ReactDOM.render(jsx, div);
      return div;
    };

    this.leafletElement = legend;
  }

  createLeafletElement() { }
}

LegendControl = withLeaflet(LegendControl);

export default GatewayDiscovery;

