import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Button from "@material-ui/core/Button";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import ServiceProfileForm from "./ServiceProfileForm";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import NetworkServerStore from "../../stores/NetworkServerStore";


class CreateServiceProfile extends Component {
  constructor() {
    super();
    this.state = {
      nsDialog: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
  }

  componentDidMount() {
    NetworkServerStore.list(0, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          nsDialog: true,
        });
      }
    });
  }

  closeDialog() {
    this.setState({
      nsDialog: false,
    });
  }

  onSubmit(serviceProfile) {
    let sp = serviceProfile;
    sp.organizationID = this.props.match.params.organizationID;

    ServiceProfileStore.create(sp, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Dialog
          open={this.state.nsDialog}
          onClose={this.closeDialog}
        >
          <DialogTitle>Add a network-server?</DialogTitle>
          <DialogContent>
            <DialogContentText paragraph>
              LoRa App Server isn't connected to a LoRa Server network-server.
              Did you know that LoRa App Server can connect to multiple LoRa Server instances, e.g. to support multiple regions?
            </DialogContentText>
            <DialogContentText>
              Would you like to connect to a network-server now?
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button color="primary.main" component={Link} to="/network-servers/create" onClick={this.closeDialog}>{i18n.t(`${packageNS}:tr000041`)}</Button>
            <Button color="primary.main" onClick={this.closeDialog}>{i18n.t(`${packageNS}:tr000166`)}</Button>
          </DialogActions>
        </Dialog>

        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000069`)} to={`/organizations/${this.props.match.params.organizationID}/service-profiles`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000277`)} />
        </TitleBar>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <ServiceProfileForm
                submitLabel={i18n.t(`${packageNS}:tr000277`)}
                onSubmit={this.onSubmit}
                match={this.props.match}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(CreateServiceProfile);
