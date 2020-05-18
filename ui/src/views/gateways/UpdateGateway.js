import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import Loader from "../../components/Loader";
import i18n, { packageNS } from "../../i18n";
import GatewayStore from "../../stores/GatewayStore";
import GatewayForm from "./GatewayForm";


class UpdateGateway extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit = this.onSubmit.bind(this);
    this.state = {
      loading: false
    };
  }

  onSubmit = async (gateway, config, classBConfig) => {
    this.setState({ loading: true });
    config.gateway_conf.server_address = gateway.server_address;
    config.gateway_conf.keepalive_interval = gateway.keepalive_interval;
    config.gateway_conf.stat_interval = gateway.stat_interval;
    config.gateway_conf.push_timeout_ms = gateway.push_timeout_ms;
    config.gateway_conf.serv_port_up = gateway.serv_port_up;
    config.gateway_conf.gps_tty_path = gateway.gps_tty_path;
    config.gateway_conf.serv_port_down = gateway.serv_port_down;
    config.gateway_conf.forward_crc_disabled = gateway.forward_crc_disabled;
    config.gateway_conf.forward_crc_error = gateway.forward_crc_error;
    config.gateway_conf.forward_crc_valid = gateway.forward_crc_valid;

    gateway.beacon_period = classBConfig.beacon_period;
    gateway.beacon_freq_hz = classBConfig.beacon_freq_hz;
    gateway.beacon_datarate = classBConfig.beacon_datarate;
    gateway.beacon_bw_hz = classBConfig.beacon_bw_hz ;
    gateway.beacon_power = classBConfig.beacon_power;
    gateway.beacon_infodesc = classBConfig.beacon_infodesc;

    console.log('gateway.id', gateway.id);
    console.log('gateway.autoUpdate', gateway.autoUpdate);
    
    const setAutoUpdateFirmwareRes = await GatewayStore.setAutoUpdateFirmware( gateway.id, gateway.autoUpdate);
    const updateRes = await GatewayStore.update(gateway);
    const updateConfigRes = await GatewayStore.updateConfig(gateway, config);
    this.setState({ loading: false });
    this.props.history.push(
      `/organizations/${this.props.match.params.organizationID}/gateways`
    ); 
  }

  render() {
    return (
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <GatewayForm
          submitLabel={i18n.t(`${packageNS}:tr000614`)}
          object={this.props.gateway}
          onSubmit={this.onSubmit}
          update={true}
          match={this.props.match}
        ></GatewayForm>
      </div>
    );
  }
}

export default withRouter(UpdateGateway);
