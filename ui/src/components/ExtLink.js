import React, { Component } from "react";
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import LinkVariant from "mdi-material-ui/LinkVariant";
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
    let context = this.props.context;
    if(!context){
      context = <LinkVariant />;
    } 
    return(
      <Typography className={this.props.classes.link} onClick={this.onClick} gutterBottom>
        {context}
      </Typography>
    );
  }
}

export default withStyles(styles)(ExtLink);
