import React, { Component } from "react";

import Grid from "@material-ui/core/Grid";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import Card from '@material-ui/core/Card';
import CardContent from "@material-ui/core/CardContent";

import OrganizationStore from "../../stores/OrganizationStore";
import WithdrawStore from "../../stores/WithdrawStore";
import WithdrawForm from "./WithdrawForm";
import { withRouter } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";
import Divider from '@material-ui/core/Divider';
import Typography from "@material-ui/core/Typography";
import theme from "../../theme";
//import { FormHelperText } from "@material-ui/core";
//import { endianness } from "os";

const styles = {
  card: {
    minWidth: 180,
    width: 220,
    backgroundColor: "#0C0270",
  },
  flex: {
    display: 'flex',
    flexDirection: 'column'
  },
  title: {
    color: '#FFFFFF',
    fontSize: 14,
    padding: 6,
  },
  balance: {
    fontSize: 24,
    color: '#FFFFFF',
    textAlign: 'center',
  },
  newBalance: {
    fontSize: 24,
    textAlign: 'center',
    color: theme.palette.primary.main,
  },
  navText: {
    fontSize: 14,
  },
  pos: {
    marginBottom: 12,
    color: '#FFFFFF',
    textAlign: 'right',
  },
  TitleBar: {
    height: 115,
    width: '50%',
    light: true,
    display: 'flex',
    flexDirection: 'column'
  },
  divider: {
    padding: 0,
    color: '#FFFFFF',
    width: '100%',
  },
  padding: {
    padding: 0,
  },
  between: {
    display: 'flex',
    justifyContent:'spaceBetween'
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
    OrganizationStore.update(organization, resp => {
    this.props.history.push(`/withdraw/${this.props.match.params.organizationID}`);
    }); 
  }

  render() {
    
    return(
      <Grid container spacing={24} className={this.props.classes.backgroundColor}>
        <Grid item xs={12} className={this.props.classes.divider}>
          <div className={this.props.classes.TitleBar}>
              <TitleBar className={this.props.classes.padding}>
                <TitleBarTitle title="Withdraw" />
              </TitleBar>
              <Divider light={true}/>
              <TitleBar>
                <TitleBarTitle title="M2M Wallet" className={this.props.classes.navText}/>
                <TitleBarTitle title="/" className={this.props.classes.navText}/>
                <TitleBarTitle title="Withdraw" className={this.props.classes.navText}/>
              </TitleBar>
          </div>
        </Grid>
        <Grid item xs={6} className={this.props.classes.divider}></Grid>
        <Grid item xs={12} className={this.props.classes.divider}>
          
        </Grid>
        <Grid item xs={6}>
          <WithdrawForm
            submitLabel="Withdraw"
            organization={this.state.organization} {...this.props}
            onSubmit={this.onSubmit}
          />
        </Grid>
        <Grid item xs={2}>
          
        </Grid>
        <Grid item xs={3}>
          <Card className={this.props.classes.card}>
            <CardContent className="space-between" >
              <Typography  className={this.props.classes.title} gutterBottom>
                Balance
              </Typography>
              <Typography className={this.props.classes.title} gutterBottom>
                Tokens
              </Typography>
            </CardContent>
            <CardContent>    
              <Typography className={this.props.classes.balance} variant="h5" component="h2">
                1234.00098
              </Typography>
              <Typography className={this.props.classes.pos} color="textSecondary">
                MXC
              </Typography>
              <Typography className={this.props.classes.title} color="textSecondary" gutterBottom>
                New Balance
              </Typography>
              <Typography className={this.props.classes.newBalance} variant="h5" component="h2">
                1234.00098
              </Typography>
              <Typography className={this.props.classes.pos} color="textSecondary">
                MXC
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));