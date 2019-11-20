import React, { Component } from "react";
import { Link } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import DataTable from "../../components/DataTable";
import SessionStore from "../../stores/SessionStore";
import theme from "../../theme";


const styles = {
  link: {
    color: theme.palette.primary.main,
    textDecoration: "none",
  },

  type: {
    fontWeight: "bold",
  },
};


class ApplicationResult extends Component {
  render() {
    return(
      <TableRow>
        <TableCell className={this.props.classes.type}>{i18n.t(`${packageNS}:tr000407`)}</TableCell>
        <TableCell><Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}/applications/${this.props.result.applicationID}`}>{this.props.result.applicationName}</Link></TableCell>
        <TableCell>{i18n.t(`${packageNS}:tr000409`)}: <Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}`}>{this.props.result.organizationName}</Link></TableCell>
        <TableCell>{this.props.result.applicationID}</TableCell>
      </TableRow>
    );
  }
}

ApplicationResult = withStyles(styles)(ApplicationResult);


class OrganizationResult extends Component {
  render() {
    return(
      <TableRow>
        <TableCell className={this.props.classes.type}>{i18n.t(`${packageNS}:tr000409`)}</TableCell>
        <TableCell><Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}`}>{this.props.result.organizationName}</Link></TableCell>
        <TableCell></TableCell>
        <TableCell>{this.props.result.organizationID}</TableCell>
      </TableRow>
    );
  }
}

OrganizationResult = withStyles(styles)(OrganizationResult);

class DeviceResult extends Component {
  render() {
    return(
      <TableRow>
        <TableCell className={this.props.classes.type}>{i18n.t(`${packageNS}:tr000408`)}</TableCell>
        <TableCell><Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}/applications/${this.props.result.applicationID}/devices/${this.props.result.deviceDevEUI}`}>{this.props.result.deviceName}</Link></TableCell>
        <TableCell>{i18n.t(`${packageNS}:tr000409`)}: <Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}`}>{this.props.result.organizationName}</Link>, {i18n.t(`${packageNS}:tr000407`)}: <Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}/applications/${this.props.result.applicationID}`}>{this.props.result.applicationName}</Link></TableCell>
        <TableCell>{this.props.result.deviceDevEUI}</TableCell>
      </TableRow>
    );
  }
}

DeviceResult = withStyles(styles)(DeviceResult);

class GatewayResult extends Component {
  render() {
    return(
      <TableRow>
        <TableCell className={this.props.classes.type}>{i18n.t(`${packageNS}:tr000410`)}</TableCell>
        <TableCell><Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}/gateways/${this.props.result.gatewayMAC}`}>{this.props.result.gatewayName}</Link></TableCell>
        <TableCell>{i18n.t(`${packageNS}:tr000409`)}: <Link className={this.props.classes.link} to={`/organizations/${this.props.result.organizationID}`}>{this.props.result.organizationName}</Link></TableCell>
        <TableCell>{this.props.result.gatewayMAC}</TableCell>
      </TableRow>
    );
  }
}

GatewayResult = withStyles(styles)(GatewayResult);


class Search extends Component {
  constructor() {
    super();
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    const query = new URLSearchParams(this.props.location.search);
    const search = (query.get("search") === null) ? "" : query.get("search");

    if (search === "") {
      callbackFunc({result: [], totalCount: 0});
      return;
    }

    SessionStore.globalSearch(search, limit, offset, resp => {
      let r = resp;
      r.totalCount = r.result.length;
      callbackFunc(r);
    });
  }

  getRow(obj) {
    switch (obj.kind) {
      case "application":
        return <ApplicationResult result={obj} />
      case "organization":
        return <OrganizationResult result={obj} />
      case "device":
        return <DeviceResult result={obj} />
      case "gateway":
        return <GatewayResult result={obj} />
      default:
        break;
    }
  }

  render() {
    return(
      <Grid container spacing={4}>
        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000411`)} />
        </TitleBar>
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell>{i18n.t(`${packageNS}:tr000412`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
                <TableCell></TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000077`)}</TableCell>
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

export default Search;
