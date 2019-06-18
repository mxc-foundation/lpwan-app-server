import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
//import Admin from "../../components/Admin";
import OrganizationStore from "../../stores/OrganizationStore";
//import WithdrawStore from "../../stores/WithdrawStore";
import WithdrawForm from "./WithdrawForm";
import { withRouter } from "react-router-dom";

class Withdraw extends Component {
  constructor() {
    super();
    this.state = {};
    this.loadData = this.loadData.bind(this);
    //this.deleteOrganization = this.deleteOrganization.bind(this);
  }
  
  componentDidMount() {
    //console.log("componentDidMount this.props")
    //console.log(this.props)
    this.loadData();
  }

  formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }

  loadData() {
    //console.log("loadData")
    //console.log(this.props)
    OrganizationStore.get(this.props.match.params.organizationID, resp => {
      //console.log(resp)
      resp.balance = 1000000.23
      
      Object.keys(resp).forEach(attr => {
        const value = resp[attr];
    
        if (typeof value === 'number') {
          resp[attr] = this.formatNumber(value);
        }
        });
      

      this.setState({
        organization: resp,
      });
    });
  }
  
  componentDidUpdate(prevProps) {
    //console.log("prevProps")
    //console.log(prevProps)
    if (prevProps === this.props) {
      return;
    }

    this.loadData();
  }

  deleteOrganization() {
    if (window.confirm("Are you sure you want to delete this organization?")) {
      OrganizationStore.delete(this.props.match.params.organizationID, () => {
        this.props.history.push("/withdraw");
      });
    }
  }

  onSubmit(organization) {
    OrganizationStore.update(organization, resp => {
    this.props.history.push("/withdraw");
    });
  }

  render() {
    
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
                object={this.state.organization} {...this.props}
                onSubmit={this.onSubmit}
                
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withRouter(Withdraw);