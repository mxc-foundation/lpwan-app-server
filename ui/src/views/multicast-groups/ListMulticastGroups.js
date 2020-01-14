import React, { Component } from "react";
import { Link } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import Plus from "mdi-material-ui/Plus";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TableCellLink from "../../components/TableCellLink";
import TitleBarButton from "../../components/TitleBarButton";
import DataTable from "../../components/DataTable";
import DeviceAdmin from "../../components/DeviceAdmin";
import MulticastGroupStore from "../../stores/MulticastGroupStore";
import theme from "../../theme";

import breadcrumbStyles from "../common/BreadcrumbStyles";

const localStyles = {
  idColumn: {
    width: theme.spacing(45),
    whiteSpace: "nowrap",
  },
};

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

class ListMulticastGroups extends Component {
  constructor() {
    super();
    this.getPage = this.getPage.bind(this);
    this.getRow = this.getRow.bind(this);
  }

  getPage(limit, offset, callbackFunc) {
    MulticastGroupStore.list("", this.props.match.params.organizationID, "", "", limit, offset, callbackFunc);
  }

  getRow(obj) {
    return(
      <TableRow key={obj.id}>
        <TableCell>{obj.id}</TableCell>
        <TableCellLink to={`/organizations/${this.props.match.params.organizationID}/multicast-groups/${obj.id}`}>{obj.name}</TableCellLink>
        <TableCellLink to={`/organizations/${this.props.match.params.organizationID}/service-profiles/${obj.serviceProfileID}`}>{obj.serviceProfileName}</TableCellLink>
      </TableRow>
    );
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return(
      <Grid container spacing={4}>
        <TitleBar
          buttons={
            <DeviceAdmin organizationID={this.props.match.params.organizationID}>
              <TitleBarButton
                label={i18n.t(`${packageNS}:tr000277`)}
                icon={<Plus />}
                to={`/organizations/${this.props.match.params.organizationID}/multicast-groups/create`}
              />
            </DeviceAdmin>
          }
        >
          <Breadcrumb className={classes.breadcrumb}>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations`}
              >
                  Organizations
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem>
              <Link
                className={classes.breadcrumbItemLink}
                to={`/organizations/${currentOrgID}`}
              >
                {currentOrgID}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:tr000083`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Grid item xs={12}>
          <DataTable
            header={
              <TableRow>
                <TableCell className={this.props.classes.idColumn}>{i18n.t(`${packageNS}:tr000077`)}</TableCell>
                <TableCell>{i18n.t(`${packageNS}:tr000042`)}</TableCell>
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

export default withStyles(styles)(ListMulticastGroups);
