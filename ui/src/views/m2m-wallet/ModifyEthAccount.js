import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";

import Plus from "mdi-material-ui/Plus";

import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import TableCellLink from "../../components/TableCellLink";
import TitleBarButton from "../../components/TitleBarButton";
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";

import Admin from "../../components/Admin";
import ApplicationStore from "../../stores/ApplicationStore";
import ModifyEthAccountForm from "./ModifyEthAccountForm";

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
        <TitleBar
          buttons={
            <Admin organizationID={this.props.match.params.organizationID}>
              <TitleBarButton
                label="Create"
                icon={<Plus />}
                to={`/organizations/${this.props.match.params.organizationID}/applications/create`}
              />
            </Admin>
          }
        >
          <TitleBarTitle title="ModifyEthAccount" />
        </TitleBar>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <ModifyEthAccountForm
                submitLabel="Confirm"
                //object={this.state.organization} {...props}
                //onSubmit={this.onSubmit}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default ModifyEthAccount;
