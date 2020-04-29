import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Button, Row } from 'reactstrap';
import Modal from '../../components/Modal';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import UpdateGatewayProfile from "./UpdateGatewayProfile";




class GatewayProfileLayout extends Component {
  constructor() {
    super();

    this.state = {
      nsDialog: false
    };

    this.deleteGatewayProfile = this.deleteGatewayProfile.bind(this);
  }

  componentDidMount() {
    GatewayProfileStore.get(this.props.match.params.gatewayProfileID, resp => {
      this.setState({
        gatewayProfile: resp,
      });
    });
  }

  deleteGatewayProfile = () => {
    GatewayProfileStore.delete(this.props.match.params.gatewayProfileID, () => {
      this.props.history.push("/gateway-profiles");
    });
    this.setState({ nsDialog: false });
  }

  openModal = () => {
    this.setState({
      nsDialog: true,
    });
  }

  render() {
  
    if (this.state.gatewayProfile === undefined) {
      return (<div></div>);
    }

    return (
      <React.Fragment>
        {this.state.nsDialog && <Modal
          title={""}
          context={i18n.t(`${packageNS}:tr000426`)}
          closeModal={() => this.setState({ nsDialog: false })}
          callback={this.deleteGatewayProfile} />}
        <TitleBar
          buttons={[
            <Button color="danger"
              key={1}
              onClick={this.openModal}
              className=""><i className="mdi mdi-delete"></i>{' '}{i18n.t(`${packageNS}:tr000401`)}
            </Button>
          ]}
        >
          <Breadcrumb>
            <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
            <BreadcrumbItem><Link to={`/gateway-profiles`}>{i18n.t(`${packageNS}:tr000046`)}</Link></BreadcrumbItem>
            <BreadcrumbItem>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
            <BreadcrumbItem active>{`${this.state.gatewayProfile.gatewayProfile.name}`}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <UpdateGatewayProfile gatewayProfile={this.state.gatewayProfile.gatewayProfile} />
        </Row>
      </React.Fragment>
    );
  }
}

export default withRouter(GatewayProfileLayout);
