import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import i18n, { packageNS } from '../../../i18n';
import HistoryStore from "../../../stores/HistoryStore";
import DataTable from "../../../components/DataTable";
import { ETHER } from "../../../util/CoinType"
import { SUPER_ADMIN } from "../../../util/M2mUtil";

class SuperNodeEthAccount extends Component {
  constructor(props) {
    super(props);
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    HistoryStore.getChangeMoneyAccountHistory(ETHER, SUPER_ADMIN, limit, offset, (data) => {
      callbackFunc({
        totalCount: parseInt(data.count),
        result: data.changeHistory
      });
    });
  }

  getRow(obj, index) {
    return(
      <TableRow key={index}>
        <TableCell>{obj.addr}</TableCell>
        <TableCell>{obj.status}</TableCell>
        <TableCell>{obj.createdAt.substring(0,19)}</TableCell>
      </TableRow>
    );
  }

  render() {
    return(
      <Grid container spacing={24}>
          {/*<TitleBar
          buttons={
            <Admin organizationID={this.props.match.params.organizationID}>
              <TitleBarButton
                label="Filter"
                //icon={<Plus />}
              />
            </Admin>
          }
        >
        </TitleBar>*/}
        <Grid item xs={12} md={12} lg={6} >
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:menu.staking.account`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.staking.status`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.staking.date`)}</TableCell>
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

export default SuperNodeEthAccount;
