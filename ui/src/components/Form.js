import React, { Component } from "react";

import Grid from '@material-ui/core/Grid';
//import Button from "@material-ui/core/Button";
import { Button } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";


const styles = {
  formControl: {
    paddingTop: 24,
  },
}


class Form extends Component {
  render() {
    return(
      <form onSubmit={this.props.onSubmit}>
        {this.props.children}

        <Grid container justify="flex-end" className={this.props.classes.formControl}>
          {this.props.extraButtons}
          {this.props.submitLabel && <Button color="primary" 
              onClick={this.deleteOrganization}
              disabled={this.props.disabled}
              className="btn-block">{this.props.submitLabel}
          </Button>} 
        </Grid>
      </form>
    );
  }
}

export default withStyles(styles)(Form);
