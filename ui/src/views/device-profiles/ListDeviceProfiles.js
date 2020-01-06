import React, { Component } from "react";

import Grid from '@material-ui/core/Grid';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import Plus from "mdi-material-ui/Plus";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellLink from "../../components/TableCellLink";
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";
import DeviceAdmin from "../../components/DeviceAdmin";
import DeviceProfileStore from "../../stores/DeviceProfileStore";
import OrganizationDevices from "../devices/OrganizationDevices";



class ListDeviceProfiles extends Component {
  constructor() {
    super();

    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    DeviceProfileStore.list(currentOrgID, 0, limit, offset, callbackFunc);
  }

  getRow(obj) {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <TableRow key={obj.id}>
        <TableCellLink to={`/organizations/${currentOrgID}/device-profiles/${obj.id}`}>{obj.name}</TableCellLink>
      </TableRow>
    );
  }

  render() {
    // TODO - refactor this into a method or store in state on page load (apply to all components where this rushed approach used)
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <Grid container spacing={4}>
        <OrganizationDevices
          mainTabIndex={2}
          organizationID={currentOrgID}
        >
          <TitleBar
            buttons={
              <DeviceAdmin organizationID={currentOrgID}>
                <TitleBarButton
                  label={i18n.t(`${packageNS}:tr000277`)}
                  icon={<Plus />}
                  to={`/organizations/${currentOrgID}/device-profiles/create`}
                />
              </DeviceAdmin>
            }
          >
            <TitleBarTitle title={i18n.t(`${packageNS}:tr000070`)} />
          </TitleBar>
          <Grid item xs={12}>
            <DataTable
              header={
                <TableRow>
                  <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                </TableRow>
              }
              getPage={this.getPage}
              getRow={this.getRow}
            />
          </Grid>
        </OrganizationDevices>
      </Grid>
    );
  }
}

export default ListDeviceProfiles;
