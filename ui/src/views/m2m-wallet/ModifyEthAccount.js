import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import Button from "@material-ui/core/Button";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellLink from "../../components/TableCellLink";
import Divider from '@material-ui/core/Divider';

import ApplicationStore from "../../stores/ApplicationStore";
import ModifyEthAccountForm from "./ModifyEthAccountForm";
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

const styles = {
  navText: {
    fontSize: 14,
  },
  TitleBar: {
    height: 115,
    width: '50%',
    light: true,
    display: 'flex',
    flexDirection: 'column'
  },
  card: {
    backgroundColor: "#090046",
    display: 'flex',
    justifyContent: 'flex-end',
    boxShadow: 0
  },
  divider: {
    padding: 0,
    color: '#FFFFFF',
    width: '100%',
  },
  padding: {
    paddingTop: 13,
  },
};

class ModifyEthAccount extends Component {
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
        <Grid item xs={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
              <TitleBar className={this.props.classes.padding}>
                <TitleBarTitle title="ETH Account" />
              </TitleBar>
              <Divider light={true}/>
              <TitleBar>
                <TitleBarTitle title="M2M Wallet" className={this.props.classes.navText}/>
                <TitleBarTitle title="/" className={this.props.classes.navText}/>
                <TitleBarTitle title="ETH Account" className={this.props.classes.navText}/>
              </TitleBar>
          </div>
        </Grid>
        <Grid item xs={6}>
          <ModifyEthAccountForm
            submitLabel="Confirm"
            //object={this.state.organization} {...props}
            //onSubmit={this.onSubmit}
          />
        </Grid>
        <Grid item xs={6} className={this.props.classes.card}>
          <Button color="primary" type="submit" disabled={this.props.disabled}>ADD TOKEN</Button>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(ModifyEthAccount));
