import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import i18n, { packageNS } from '../../i18n';
import HistoryStore from "../../stores/HistoryStore";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";
import Admin from "../../components/Admin";
import { ETHER } from "../../util/CoinType";

class SubScriptions extends Component {
  constructor() {
    super();
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    HistoryStore.getWithdrawHistory(ETHER, this.props.match.params.organizationID, limit, offset, (data) => {
      callbackFunc({
        totalCount: offset + 2 * limit,
        result: data.withdrawHistory
      });
    });
  }

  getRow(obj, index) {
    return(
      <TableRow key={index}>
        <TableCell>{obj.from}</TableCell>
        <TableCell>{obj.to}</TableCell>
        <TableCell>{obj.moneyType}</TableCell>
        <TableCell>{obj.amount}</TableCell>
        <TableCell>{obj.createdAt}</TableCell>
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
                label={i18n.t(`${packageNS}:menu.history.filter`)}
                //icon={<Plus />}
              />
            </Admin>
          }
        >
        </TitleBar>*/}
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:menu.history.from`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.history.to`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.history.type`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.history.amount`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.history.date`)}</TableCell>
              </TableRow>
            }
            getPage={this.getPage}
            getRow={this.getRow}
          />
        </Grid>
    );
  }
}

export default SubScriptions;
