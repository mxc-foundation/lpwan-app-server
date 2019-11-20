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

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import MulticastGroupForm from "./MulticastGroupForm";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import MulticastGroupStore from "../../stores/MulticastGroupStore";


const styles = {
  card: {
    overflow: "visible",
  },
};


class CreateMulticastGroup extends Component {
  constructor() {
    super();
    this.state = {
      spDialog: false,
    };
    this.onSubmit = this.onSubmit.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
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

  closeDialog() {
    this.setState({
      spDialog: false,
    });
  }

  onSubmit(multicastGroup) {
    MulticastGroupStore.create(multicastGroup, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/multicast-groups`);
    });
  }

  render() {
    return(
      <Grid container spacing={4}>
        <Dialog
          open={this.state.spDialog}
          onClose={this.closeDialog}
        >
          <DialogTitle>{i18n.t(`${packageNS}:tr000164`)}</DialogTitle>
          <DialogContent>
            <DialogContentText paragraph>
              {i18n.t(`${packageNS}:tr000165`)}
              {i18n.t(`${packageNS}:tr000326`)}
            </DialogContentText>
            <DialogContentText>
              {i18n.t(`${packageNS}:tr000327`)}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button color="primary.main" component={Link} to={`/organizations/${this.props.match.params.organizationID}/service-profiles/create`} onClick={this.closeDialog}>{i18n.t(`${packageNS}:tr000277`)}</Button>
            <Button color="primary.main" onClick={this.closeDialog}>{i18n.t(`${packageNS}:tr000166`)}</Button>
          </DialogActions>
        </Dialog>

        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000083`)} to={`/organizations/${this.props.match.params.organizationID}/multicast-groups`} />
          <TitleBarTitle title="/" />
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000277`)} />
        </TitleBar>

        <Grid item xs={12}>
          <Card className={this.props.classes.card}>
            <CardContent>
              <MulticastGroupForm
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

export default withStyles(styles)(withRouter(CreateMulticastGroup));
