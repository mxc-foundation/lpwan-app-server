import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";
//import Admin from "../../components/Admin";
//import OrganizationStore from "../../stores/OrganizationStore";
import WithdrawStore from "../../stores/WithdrawStore";
import WithdrawForm from "./WithdrawForm";
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

const styles = {
  backgroundColor: {
    backgroundColor: "#090046",
  },
  font: {
    color: '#FFFFFF', 
    fontFamily: 'Montserrat',
  }
};

class Withdraw extends Component {
  constructor() {
    super();
    this.state = {};
    this.loadData = this.loadData.bind(this);
  }
  
  componentDidMount() {
    this.loadData();
  }

  formatNumber(number) {
    //let balance = number.toString().replace(".", ",");
    //balance = number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }

  loadData() {
    WithdrawStore.getWithdrawFee("Ether", 
      resp => {
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
    if (prevProps === this.props) {
      return;
    }

    this.loadData();
  }

  deleteOrganization() {
    
  }

  onSubmit(organization) {
    /* OrganizationStore.update(organization, resp => {
    this.props.history.push(`/withdraw/${this.props.match.params.organizationID}`);
    }); */
  }

  render() {
    
    return(
      <Grid container spacing={24} className={this.props.classes.backgroundColor}>
        <TitleBar>
          <TitleBarTitle title="Withdraw" className={this.props.classes.font}/>
        </TitleBar>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card>
            <CardContent>
              <WithdrawForm
                submitLabel="Withdraw"
                organization={this.state.organization} {...this.props}
                onSubmit={this.onSubmit}
              />
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card>
            <CardContent>
              <WithdrawForm
                submitLabel="Withdraw"
                organization={this.state.organization} {...this.props}
                onSubmit={this.onSubmit}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));