import React, { Component } from "react";
import { withRouter, Link } from 'react-router-dom';

import { Breadcrumb, BreadcrumbItem, Button,
  Modal, ModalHeader, ModalBody, ModalFooter, NavLink,
} from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from '@material-ui/core/Grid';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";

import DeviceProfileForm from "./DeviceProfileForm";
import OrganizationDevices from "../devices/OrganizationDevices";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import ServiceProfileStore from "../../stores/ServiceProfileStore";


const styles = theme => ({
  [theme.breakpoints.down('sm')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('sm')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  },
  card: {
    overflow: "visible",
  },
});


class CreateDeviceProfile extends Component {
  constructor() {
    super();
    this.state = {
      spDialog: false,
      loading: true,
    };
  }

  componentDidMount() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    ServiceProfileStore.list(currentOrgID, 0, 0, resp => {
      if (resp.totalCount === "0") {
        this.setState({
          spDialog: true,
          loading: false
        });
      }
    });
  }

  toggleSpDialog = () => {
    this.setState({
      spDialog: !this.state.spDialog,
    });
  }

  onSubmit = (deviceProfile) => {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    let sp = deviceProfile;
    sp.organizationID = currentOrgID;

    DeviceProfileStore.create(sp, resp => {
      this.props.history.push(`/organizations/${currentOrgID}/device-profiles`);
    });
  }

  render() {
    const { loading } = this.state;
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const closeSpBtn = <button className="close" onClick={this.toggleSpDialog}>&times;</button>;

    return(
      <Grid container spacing={4}>
        <OrganizationDevices
          mainTabIndex={2}
          organizationID={currentOrgID}
        >
          <Modal
            isOpen={this.state.spDialog}
            toggle={this.toggleSpDialog}
            aria-labelledby="help-dialog-title"
            aria-describedby="help-dialog-description"
          >
            <ModalHeader
              toggle={this.toggleSpDialog}
              close={closeSpBtn}
              id="help-dialog-title"
            >
              {i18n.t(`${packageNS}:tr000164`)}
            </ModalHeader>
            <ModalBody id="help-dialog-description">
              <p>{i18n.t(`${packageNS}:tr000165`)}</p>
              <p>{i18n.t(`${packageNS}:tr000326`)}</p>
              <p>{i18n.t(`${packageNS}:tr000327`)}</p>
            </ModalBody>
            <ModalFooter>
              <Button variant="outlined">
                <NavLink
                  style={{ color: "#fff", padding: "0" }}
                  tag={Link}
                  to={`/organizations/${currentOrgID}/service-profiles/create`}
                >
                  {i18n.t(`${packageNS}:tr000277`)}
                </NavLink>
              </Button>
              <Button
                color="primary"
                onClick={this.toggleSpDialog}
              >
                {i18n.t(`${packageNS}:tr000166`)}
              </Button>{' '}
            </ModalFooter>
          </Modal>

          <TitleBar>
            <Breadcrumb className={classes.breadcrumb}>
              <BreadcrumbItem><Link className={classes.breadcrumbItemLink} to={
                `/organizations/${currentOrgID}/device-profiles`
              }>{i18n.t(`${packageNS}:tr000070`)}</Link></BreadcrumbItem>
              <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000277`)}</BreadcrumbItem>
            </Breadcrumb>
          </TitleBar>

          <Grid item xs={12}>
            <DeviceProfileForm
              submitLabel={i18n.t(`${packageNS}:tr000277`)}
              onSubmit={this.onSubmit}
              match={this.props.match}
              loading={loading}
            />
          </Grid>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(CreateDeviceProfile));
