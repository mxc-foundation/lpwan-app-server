import Button from "@material-ui/core/Button";
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';
import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import MulticastGroupStore from "../../stores/MulticastGroupStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import breadcrumbStyles from "../common/BreadcrumbStyles";
import MulticastGroupForm from "./MulticastGroupForm";





const localStyles = {
  card: {
    overflow: "visible",
  },
};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
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
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

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
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
              >
                  Organizations
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}`}
              >
                {currentOrgID}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}/multicast-groups`}
              >
                {i18n.t(`${packageNS}:tr000083`)}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
          </Breadcrumb>
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
