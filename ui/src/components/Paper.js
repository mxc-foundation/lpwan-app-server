import Paper from '@material-ui/core/Paper';
import React, { Component } from "react";


class PaperComponent extends Component {
  render() {
    return(
      <Paper>
        {this.props.children}
      </Paper>
    );
  }
}

export default PaperComponent;