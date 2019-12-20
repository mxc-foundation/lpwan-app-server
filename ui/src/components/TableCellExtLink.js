import React, { Component } from "react";
import { Link } from 'react-router-dom';

import TableCell from '@material-ui/core/TableCell';
import { withStyles } from '@material-ui/core/styles';

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


class TableCellExtLink extends Component {
  onClick = () => {
    const url = this.props.to;
    if(this.props.for === 'lora'){
      window.location.replace(url)
    }else{
      window.open(url, '_blank');
    } 
  }
  
  render() {
    return(
      <TableCell align={this.props.align}>
        <span className={this.props.classes.link} onClick={this.onClick}>{this.props.children}</span>
      </TableCell>
    );
  }
}

export default withStyles(styles)(TableCellExtLink);
