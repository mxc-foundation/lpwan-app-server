import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import { Card } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import NetworkServerStore from "../../stores/NetworkServerStore";
import NetworkServerForm from "./NetworkServerForm";


class UpdateNetworkServer extends Component {
  constructor() {
    super();

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(networkServer) {
    NetworkServerStore.update(networkServer, resp => {
      this.props.history.push("/network-servers");
    });
  }

  render() {
    return(
      <React.Fragment>
        <Card className="card-box shadow-sm" style={{ minWidth: "25rem" }}>
          <NetworkServerForm
            object={this.props.networkServer}
            onSubmit={this.onSubmit}
            submitLabel={i18n.t(`${packageNS}:tr000066`)}
          />
        </Card>
      </React.Fragment>
    );
  }
}

export default withRouter(UpdateNetworkServer);
