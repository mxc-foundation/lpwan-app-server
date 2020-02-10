import React, { Component } from "react";

import { withStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import moment from "moment";
import SessionStore from "../stores/SessionStore";
import theme from "../theme";

const styles = {
  footer: {
    paddingBottom: theme.spacing(1),
    "& a": {
      color: theme.palette.primary.main,
      textDecoration: "none",
    },
  },
};

const footerHtml = `<!-- Footer -->
<footer class="page-footer font-small blue pt-4">
  <!-- Copyright -->
  <div class="footer-copyright text-center py-3">© ${moment().format('YYYY')} Powered by 
    <a href="https://www.mxc.org/" target="_blank"> MXC</a>
  </div>
  <!-- Copyright -->
</footer>
</footer>`

class Footer extends Component {
  constructor() {
    super();
    this.state = {
      footer: footerHtml,
    };
  }

  componentDidMount() {
    SessionStore.getBranding(resp => {
      if (resp.footer !== "") {
        this.setState({
          footer: resp.footer,
        });
      }
    });
  }

  render() {
    /* if (this.state.footer === null) {
      return(null);
    } */

    return(
      <footer className={this.props.classes.footer}>
        <span dangerouslySetInnerHTML={{__html: footerHtml}}></span>
      </footer>
    );
  }
}

export default withStyles(styles)(Footer);
