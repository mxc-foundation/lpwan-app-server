import React, { Component } from "react";

import Paper from '@material-ui/core/Paper';

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