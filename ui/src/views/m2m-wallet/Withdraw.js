import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import OrganizationStore from "../../stores/OrganizationStore";
import WithdrawForm from "./WithdrawForm";


class Withdraw extends Component {
    constructor() {
        super();
        this.state = {};
        this.loadData = this.loadData.bind(this);
        //this.deleteOrganization = this.deleteOrganization.bind(this);
      }
    
      /* componentDidMount() {
        console.log(this.props)
        this.loadData();
      }
    
      componentDidUpdate(prevProps) {
        if (prevProps === this.props) {
          return;
        }
    
        this.loadData();
      }*/
    
      loadData() {
        console.log("this.props")
        console.log(this.props)
        OrganizationStore.get(this.props.match.params.organizationID, resp => {
          console.log(resp)
          this.setState({
            organization: resp,
          });
        });
      }
    /*
      deleteOrganization() {
        if (window.confirm("Are you sure you want to delete this organization?")) {
          OrganizationStore.delete(this.props.match.params.organizationID, () => {
            this.props.history.push("/topup");
          });
        }
      }

    onSubmit(organization) {
        OrganzationStore.update(organization, resp => {
        this.props.history.push("/topup");
        });
    } */
    onSubmit(organization) {
      console.log("withdraw")
      console.log(organization)
  }

  render() {
      console.log(this.state)
    return(
      <Grid container spacing={24}>
        <TitleBar>
          <TitleBarTitle title="Withdraw" />
        </TitleBar>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <WithdrawForm
                submitLabel="Withdraw"
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

export default withRouter(Withdraw);
