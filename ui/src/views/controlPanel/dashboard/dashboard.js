import React, { Component } from "react";

import {Grid,Card,Table,TableBody,TextField} from "@material-ui/core";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import { withRouter, Link } from 'react-router-dom';
import { withStyles } from "@material-ui/core/styles";
import i18n, { packageNS } from '../../../i18n';
import HistoryStore from "../../../stores/HistoryStore";
import TitleBar from "../../../components/TitleBar";
import TitleBarTitle from "../../../components/TitleBarTitle";
import TitleBarButton from "../../../components/TitleBarButton";
import DataTable from "../../../components/DataTable";
import i18n, { packageNS } from '../i18n';
import styles from "./dashboardStyle"

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {

   

  }

  getRow(obj, index) {
    return(
      <TableRow key={index}>
        <TableCell>{obj.org}</TableCell>
        <TableCell>{obj.timestamp}</TableCell>
        <TableCell>{obj.value}</TableCell>
        <TableCell>{obj.type}</TableCell>
        <TableCell>{obj.income}</TableCell>
      </TableRow>
    );
  }

  render() {
    return(
      <Grid container spacing={3} className={this.props.classes.root}>
      <Grid item xs={12}>
        <TitleBar>
         <TitleBarTitle title={`${i18n.t(`${packageNS}:menu.settings.welcome`)} SuperAdmin`} />
        </TitleBar>
        </Grid>
   
        <Grid item xs={8}>
        
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:menu.settings.organization`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.settings.timestamp`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.settings.value`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.settings.type`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:menu.settings.income`)}</TableCell>
              </TableRow>
            }
            getPage={this.getPage}
            getRow={this.getRow}
          />
        </Grid>

        <Grid item xs={4}>
        <Grid container direction="column" spacing={10}>
        <Grid item xs={12}>
        <Card  className={this.props.classes.card}>
        <Table className={this.props.classes.cardTable}>
          <TableBody>
            <TableRow >
              <TableCell>{i18n.t(`${packageNS}:menu.settings.today_income`)}</TableCell>
              <TableCell align="right">1.244 MXC</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>{i18n.t(`${packageNS}:menu.settings.monthly_balance`)}</TableCell>
              <TableCell align="right"><span>1.244 MXC</span></TableCell>
            </TableRow>
            <TableRow>
              <TableCell></TableCell>
              <TableCell align="right"><b>{i18n.t(`${packageNS}:menu.settings.set_alert`)}</b></TableCell>
            </TableRow>
          </TableBody>
        </Table>

        </Card>
        </Grid>
        <Grid item container direction="column" xs={12}>

        <h4>{i18n.t(`${packageNS}:menu.settings.general_settings`)}</h4>
          <TextField
            id="standard-number"
            label={i18n.t(`${packageNS}:menu.settings.withdraw_fee`)}
            className={this.props.classes.TextField}
            variant="filled"
            type="number"
    
            InputLabelProps={{
              shrink: true,
            }}
            margin="normal"
          />


      <TextField
        id="standard-number"
        label={i18n.t(`${packageNS}:menu.settings.downlink_price`)}
        className={this.props.classes.TextField}
        variant="filled"
        type="number"
 
        InputLabelProps={{
          shrink: true,
        }}
        margin="normal"
      />
      
      <h4>System</h4>
        <Table className={this.props.classes.cardTable}>
          <TableBody>
            <TableRow>
              <TableCell>{i18n.t(`${packageNS}:menu.settings.monthly_downtime`)}</TableCell>
              <TableCell align="right">3 {i18n.t(`${packageNS}:menu.settings.hours`)}</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>{i18n.t(`${packageNS}:menu.settings.tickets_opened`)}</TableCell>
              <TableCell align="right"><b>1</b></TableCell>
            </TableRow>
            <TableRow>
              <TableCell>{i18n.t(`${packageNS}:menu.settings.tickets_closed`)}</TableCell>
              <TableCell align="right"><b>5</b></TableCell>
            </TableRow>
          </TableBody>
        </Table>

        </Grid>
        </Grid>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles) (withRouter(Dashboard));
