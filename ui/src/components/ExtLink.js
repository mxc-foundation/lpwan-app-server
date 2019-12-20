import React, { Component } from "react";
import { Link } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import theme from "../theme";

const styles = {
  link: {
    textDecoration: "none",
    color: theme.palette.primary.main,
    cursor: "pointer",
    padding: 0,
    fontWeight: "bold",
    fontSize: 14,
    opacity: 0.7,
      "&:hover": {
        opacity: 1,
      }
  },
};

class ExtLink extends Component {
  onClick = () => {
    const url = this.props.to;
    if(this.props.for === 'lora'){
      window.location.replace(url)
    }else if(this.props.for === 'local'){
      this.props.dismissOn();
    }else{
      window.open(url, '_blank');
    } 
  }
  
  render() {
    return(
      <Typography className={this.props.classes.link} onClick={this.onClick} gutterBottom>
        {this.props.context}
      </Typography>
    );
  }
}

export default withStyles(styles)(ExtLink);
