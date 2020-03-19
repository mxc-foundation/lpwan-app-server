import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import { CardContent } from "@material-ui/core";

import i18n, { packageNS } from '../../i18n';
import GatewayProfileStore from "../../stores/GatewayProfileStore";
import Loader from "../../components/Loader";
import GatewayProfileForm from "./GatewayProfileForm";


class UpdateGatewayProfile extends Component {
  constructor() {
    super();
    this.state = {
      loading: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(gatewayProfile) {
    this.setState({ loading: true });
    GatewayProfileStore.update(gatewayProfile, resp => {
      this.setState({ loading: false });
      this.props.history.push("/gateway-profiles");
    }, error => { this.setState({ loading: false }) });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
            {this.state.loading && <Loader />}
              <GatewayProfileForm
                submitLabel={i18n.t(`${packageNS}:tr000066`)}
                object={this.props.gatewayProfile}
                onSubmit={this.onSubmit}
                update={true}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(UpdateGatewayProfile);
