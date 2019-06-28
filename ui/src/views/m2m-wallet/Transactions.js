import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import ApplicationStore from "../../stores/ApplicationStore";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellLink from "../../components/TableCellLink";
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";
import Admin from "../../components/Admin";

class Transactions extends Component {
  constructor() {
    super();
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    ApplicationStore.list("", this.props.match.params.organizationID, limit, offset, callbackFunc);
  }

  getRow(obj) {
    return(
      <TableRow key={obj.id}>
        <TableCell>{obj.id}</TableCell>
        <TableCellLink to={`/organizations/${this.props.match.params.organizationID}/applications/${obj.id}`}>{obj.name}</TableCellLink>
        <TableCellLink to={`/organizations/${this.props.match.params.organizationID}/service-profiles/${obj.serviceProfileID}`}>{obj.serviceProfileName}</TableCellLink>
        <TableCell>{obj.description}</TableCell>
      </TableRow>
    );
  }

  render() {
    return(
      <Grid container spacing={24}>
        <TitleBar
          buttons={
            <Admin organizationID={this.props.match.params.organizationID}>
              <TitleBarButton
                label="Filter"
                //icon={<Plus />}
              />
            </Admin>
          }
        >
          <TitleBarTitle title="Transactions" />
        </TitleBar>
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell>From</TableCell>
                <TableCell>To</TableCell>
                <TableCell>Value</TableCell>
                <TableCell>Date</TableCell>
                <TableCell>Balance</TableCell>
                <TableCell>Tx hash</TableCell>
                <TableCell>Status</TableCell>
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

export default Transactions;
