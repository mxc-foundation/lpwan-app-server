import React, { Component } from "react";

import Grid from '@material-ui/core/Grid';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';

import Check from "mdi-material-ui/Check";
import Close from "mdi-material-ui/Close";
import Plus from "mdi-material-ui/Plus";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellLink from "../../components/TableCellLink";
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";

import OrganizationStore from "../../stores/OrganizationStore";


class ListOrganizations extends Component {
  getPage(limit, offset, callbackFunc) {
    OrganizationStore.list("", limit, offset, callbackFunc);
  }

  getRow(obj) {
    let icon = null;

    if (obj.canHaveGateways) {
      icon = <Check />;
    } else {
      icon = <Close />;
    }

    return(
      <TableRow key={obj.id}>
        <TableCellLink to={`/organizations/${obj.id}`}>{obj.name}</TableCellLink>
        <TableCell>{obj.displayName}</TableCell>
        <TableCell>{icon}</TableCell>
        <TableCellLink to={`/organizations/${obj.id}/service-profiles/create`}>Add Service Profile</TableCellLink>
      </TableRow>
    );
  }

  render() {
    return(
      <Grid container spacing={4}>
        <TitleBar
          buttons={[
            <TitleBarButton
              key={1}
              label={i18n.t(`${packageNS}:tr000277`)}
              icon={<Plus />}
              to={`/organizations/create`}
            />,
          ]}
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000049`)} />
        </TitleBar>
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000126`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000380`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000078`)}</TableCell>
              </TableRow>
            }
            getPage={this.getPage}
            getRow={this.getRow}
          />
        </Grid>
      </Grid>
    );
  }
}

export default ListOrganizations;
