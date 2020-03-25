import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem, Card } from 'reactstrap';
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";
import NetworkServerForm from "./NetworkServerForm";



class CreateNetworkServer extends Component {
  constructor() {
    super();
    this.state = {
      loading: false,
    }
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(networkServer) {
    this.setState({ loading: true });
    NetworkServerStore.create(networkServer, resp => {
      this.setState({ loading: false });
      this.props.history.push("/network-servers");
    }, error => { this.setState({ loading: false }) });
  }

  render() {
    return (
      <React.Fragment>
        <TitleBar>
          <Breadcrumb>
            <BreadcrumbItem>{i18n.t(`${packageNS}:menu.control_panel`)}</BreadcrumbItem>
            <BreadcrumbItem><Link to={`/network-servers`}>{i18n.t(`${packageNS}:tr000040`)}</Link></BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>

        <Card className="card-box shadow-sm">
          {this.state.loading && <Loader />}
          <NetworkServerForm
            onSubmit={this.onSubmit}
            submitLabel={i18n.t(`${packageNS}:tr000041`)}
          />
        </Card>
      </React.Fragment>
    );
  }
}

export default withRouter(CreateNetworkServer);
