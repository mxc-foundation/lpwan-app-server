import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from "@material-ui/core/styles";
import Typography from '@material-ui/core/Typography';
import React from "react";
import { withRouter } from "react-router-dom";
import FormComponent from "../../classes/FormComponent";
import styles from "./WithdrawStyle";


class WithdrawBalanceInfo extends FormComponent {
    
  render() {
    if (this.props.txinfo === undefined) {
      return(<div>loading...</div>);
    }

    return(
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
            {this.props.txinfo.balance || ""}
          </Typography>
          <Typography className={this.props.classes.pos} color="textSecondary">
            MXC
          </Typography>
          <Typography className={this.props.classes.title} color="textSecondary" gutterBottom>
            New Balance
          </Typography>
          <Typography className={this.props.classes.newBalance} variant="h5" component="h2">
            {this.props.txinfo.newBalance || ""}
          </Typography>
          <Typography className={this.props.classes.pos} color="textSecondary">
            MXC
          </Typography>
        </CardContent>
      </Card>
    );
  }
}

export default withStyles(styles)(withRouter(WithdrawBalanceInfo));;
