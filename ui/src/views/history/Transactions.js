import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import i18n, { packageNS } from '../../i18n';
import TopupStore from "../../stores/TopupStore";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellExtLink from '../../components/TableCellExtLink';
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";
import LinkVariant from "mdi-material-ui/LinkVariant";
import Admin from "../../components/Admin";
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

const styles = {
  maxW140: {
    maxWidth: 140,
    //backgroundColor: "#0C0270",
    whiteSpace: 'nowrap', 
    overflow: 'hidden',
    textOverflow: 'ellipsis'
  },
  flex:{
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center'

  }
};

class Transactions extends Component {
  constructor(props) {
    super(props);
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    TopupStore.getTransactionsHistory(this.props.organizationID, offset, limit, data => {
        callbackFunc({
            totalCount: parseInt(data.count),
            result: data.transactionHistory
          });
      }); 
  }
  
  getRow(obj, index) {
    const url = process.env.REACT_APP_ETHERSCAN_HOST + `/tx/${obj.txHash}`;
    return(
      <TableRow key={index}>
        <TableCell align={'center'} className={this.props.classes.maxW140} >{obj.from}</TableCell>
        <TableCell align={'center'} className={this.props.classes.maxW140}>{obj.to}</TableCell>
        <TableCellExtLink to={url} ><div className={this.props.classes.flex}><LinkVariant /></div></TableCellExtLink>
        <TableCell align={'center'}>{obj.transactionType}</TableCell>
        <TableCell align={'right'}>{obj.amount}</TableCell>
        <TableCell align={'center'}>{obj.lastUpdateTime.substring(0, 19)}</TableCell>
        <TableCell align={'center'}>{obj.status}</TableCell>
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
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.from`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.to`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.tx_hash`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.type`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.mxc_amount`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.update_date`)}</TableCell>
                <TableCell align={'center'}>{i18n.t(`${packageNS}:menu.history.status`)}</TableCell>
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

export default withStyles(styles)(withRouter(Transactions));