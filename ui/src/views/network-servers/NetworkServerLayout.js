import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import Modal from "../../components/Modal";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";
import UpdateNetworkServer from "./UpdateNetworkServer";




class NetworkServerLayout extends Component {
  constructor() {
    super();

    this.state = {
      nsDialog: false,
    };

    this.deleteNetworkServer = this.deleteNetworkServer.bind(this);
    this.openConfirmModal = this.openConfirmModal.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  loadData = async () => {
    let networkServer = await NetworkServerStore.get(this.props.match.params.networkServerID);
    
    this.setState({
      networkServer
    });
  }

  deleteNetworkServer() {
    NetworkServerStore.delete(this.props.match.params.networkServerID, () => {
      this.props.history.push("/network-servers");
    });
    this.setState({ nsDialog: false });
  }

  openConfirmModal = () => {
    this.setState({
      nsDialog: true,
    });
  };

  render() {

    if (this.state.networkServer === undefined) {
      return (<div></div>);
    }

    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <TitleBarButton
              color="danger"
              key={1}
              icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
              label={i18n.t(`${packageNS}:tr000061`)}
              onClick={this.openConfirmModal}
            />,
          ]}
        >
          <Breadcrumb>
            <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
            <BreadcrumbItem><Link to={
              `/network-servers`}>{i18n.t(`${packageNS}:tr000040`)
              }</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000066`)}</BreadcrumbItem>
            <BreadcrumbItem active>{this.state.networkServer.networkServer.id}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        {this.state.nsDialog && <Modal
          title={""}
          closeModal={() => this.setState({ nsDialog: false })}
          context={i18n.t(`${packageNS}:lpwan.network_servers.delete_server`)}
          callback={this.deleteNetworkServer} />}

        <UpdateNetworkServer networkServer={this.state.networkServer.networkServer} version={this.state.networkServer.version}
          region={this.state.networkServer.region} />
      </React.Fragment>
    );
  }
}

export default withRouter(NetworkServerLayout);
