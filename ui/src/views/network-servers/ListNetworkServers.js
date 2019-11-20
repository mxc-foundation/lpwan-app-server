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

import NetworkServerStore from "../../stores/NetworkServerStore";


class ListNetworkServers extends Component {
  getPage(limit, offset, callbackFunc) {
    NetworkServerStore.list(0, limit, offset, callbackFunc);
  }

  getRow(obj) {
    return(
      <TableRow key={obj.id}>
        <TableCellLink to={`/network-servers/${obj.id}`}>{obj.name}</TableCellLink>
        <TableCell>{obj.server}</TableCell>
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
              icon={<Plus />}
              label={i18n.t(`${packageNS}:tr000041`)}
              to={`/network-servers/create`}
            />,
          ]}
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000040`)} />
        </TitleBar>
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000043`)}</TableCell>
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

export default ListNetworkServers;
