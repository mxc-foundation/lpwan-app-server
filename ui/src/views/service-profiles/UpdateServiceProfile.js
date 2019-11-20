import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";

import i18n, { packageNS } from '../../i18n';
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import ServiceProfileForm from "./ServiceProfileForm";


class UpdateServiceProfile extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(serviceProfile) {
    ServiceProfileStore.update(serviceProfile, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <ServiceProfileForm
                submitLabel={i18n.t(`${packageNS}:tr000066`)}
                object={this.props.serviceProfile}
                onSubmit={this.onSubmit}
                match={this.props.match}
                disabled={!this.props.admin}
                update={true}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(UpdateServiceProfile);
