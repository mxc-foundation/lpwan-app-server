import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Button from "@material-ui/core/Button";

import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import GatewayForm from "./GatewayForm";
import GatewayStore from "../../stores/GatewayStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";


const styles = {
  card: {
    overflow: "visible",
  },
};


class CreateGateway extends Component {
  constructor() {
    super();

    this.state = {
      spDialog: false,
    };

    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    ServiceProfileStore.list(this.props.match.params.organizationID, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          spDialog: true,
        });
      }
    });
  }

  closeDialog = () => {
    this.setState({
      spDialog: false,
    });
  }

  onSubmit(gateway) {
    let gw = gateway;
    gw.organizationID = this.props.match.params.organizationID;

    GatewayStore.create(gateway, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/gateways`);
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Dialog
          open={this.state.spDialog}
          onClose={this.closeDialog}
        >
          <DialogTitle>Add a service-profile?</DialogTitle>
          <DialogContent>
            <DialogContentText paragraph>
              The selected organization does not have a service-profile yet.
              A service-profile connects an organization to a network-server and defines the features that an organization can use on this network-server.
            </DialogContentText>
            <DialogContentText>
              Would you like to create a service-profile?
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button color="primary" component={Link} to={`/organizations/${this.props.match.params.organizationID}/service-profiles/create`} onClick={this.closeDialog}>Create service-profile</Button>
            <Button color="primary" onClick={this.closeDialog}>Dismiss</Button>
          </DialogActions>
        </Dialog>
        <TitleBar>
          <TitleBarTitle title="Gateways" to={`/organizations/${this.props.match.params.organizationID}/gateways`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title="Create" />
        </TitleBar>
        <Grid item xs={12}>
          <Card className={this.props.classes.card}>
            <CardContent>
              <GatewayForm
                match={this.props.match}
                submitLabel="Create"
                onSubmit={this.onSubmit}
                object={{location: {}}}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(withStyles(styles)(CreateGateway));
