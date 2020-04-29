import React, { Component } from "react";
//import WithdrawBalanceInfo from "./WithdrawBalanceInfo";
import { withRouter } from "react-router-dom";
import { Card, CardBody, Col, Row } from 'reactstrap';
import TitleBar from "../../../components/TitleBar";
import TitleBarTitle from "../../../components/TitleBarTitle";
import i18n, { packageNS } from '../../../i18n';
import DeviceStore from "../../../stores/DeviceStore.js";
import GatewayStore from "../../../stores/GatewayStore.js";
import WalletStore from "../../../stores/WalletStore.js";
import { DV_FREE_GATEWAYS_LIMITED, DV_INACTIVE, DV_WHOLE_NETWORK } from "../../../util/Data";
import OrganizationDevices from "../OrganizationDevices";
import DeviceForm from "./DeviceFormM2M";


function doIHaveGateway(orgId) {
  return new Promise((resolve, reject) => {
    GatewayStore.getGatewayList(orgId, 0, 1, data => {
      resolve(parseInt(data.count));
    });
  });
}

function getDlPrice(orgId) {
  return new Promise((resolve, reject) => {
    WalletStore.getDlPrice(orgId, resp => {
      resolve(resp.downLinkPrice);
    });
  });
}

class DeviceLayoutM2M extends Component {
  constructor(props) {
    super(props);

    this._isMounted = false;

    this.state = {
      loading: true,
      mod: null,
      haveGateway: false,
      downlinkFee: 0
    };
  }

  loadData = async () => {
    try {
      const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
      if (this._isMounted) {
        this.setState({ loading: true });
      }

      let res = await doIHaveGateway(currentOrgID);
      let downlinkFee = await getDlPrice(currentOrgID);

      let haveGateway = (res > 0) ? true : false;

      if (this._isMounted) {
        this.setState({
          downlinkFee,
          haveGateway,
          loading: false
        });
      }
    } catch (error) {
      if (this._isMounted) {
        console.error(error);
        this.setState({
          error,
          loading: false
        });
      }
    }
  }

  componentDidMount() {
    /*window.analytics.page();*/
    this._isMounted = true;
    this.loadData();
  }

  componentDidUpdate(oldProps) {
    if (this.props === oldProps) {
      return;
    }
    if (this._isMounted) {
      this.loadData();
    }
  }

  componentWillUnmount() {
    this._isMounted = false;
  }

  onSubmit = (e, apiWithdrawReqRequest) => {
    e.preventDefault();
  }

  handleCloseModal = () => {
    if (this._isMounted) {
      this.setState({
        modal: null
      })
    }
  }

  onSelectChange = (device) => {
    const { dvId, dvMode } = device;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    //console.log('device', device);
    DeviceStore.setDeviceMode(currentOrgID, dvId, dvMode, data => {
      this.props.history.push(`/device/${currentOrgID}`);
    });
  }

  onSwitchChange = (device, e) => {
    e.preventDefault();
    const { dvId, available } = device;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    //console.log('onSwitchChange', device);
    let mod = DV_FREE_GATEWAYS_LIMITED;
    if (!this.state.haveGateway) {
      mod = DV_WHOLE_NETWORK;
    }
    if (!available) {
      mod = DV_INACTIVE;
    }
    //console.log('onSwitchChange', mod);
    DeviceStore.setDeviceMode(currentOrgID, dvId, mod, data => {
    });
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const { loading } = this.state;

    return (
      <OrganizationDevices
        mainTabIndex={0}
        organizationID={currentOrgID}
        loading={this.state.loading}
      >
        <TitleBar
          buttons={
            <div>{i18n.t(`${packageNS}:menu.devices.downlink_fee_mxc`)} {this.state.downlinkFee} MXC</div>
          }
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:menu.devices.devices`)} />
        </TitleBar>
        <Row>
          <Col>
            <Card className="shadow-sm">
              <CardBody className="position-relative">
                <DeviceForm
                  submitLabel={i18n.t(`${packageNS}:menu.devices.devices`)}
                  onSubmit={this.onSubmit}
                  downlinkFee={this.state.downlinkFee}
                  haveGateway={this.state.haveGateway}
                  loading={loading}
                  onSelectChange={this.onSelectChange}
                  onSwitchChange={this.onSwitchChange}
                />
              </CardBody>
            </Card>
          </Col>
        </Row>
      </OrganizationDevices>
    );
  }
}

export default withRouter(DeviceLayoutM2M);
