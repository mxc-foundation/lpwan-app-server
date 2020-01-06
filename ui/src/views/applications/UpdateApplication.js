import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import ApplicationStore from "../../stores/ApplicationStore";
import ApplicationForm from "./ApplicationForm";


const styles = {
  card: {
    overflow: "visible",
  },
};


class UpdateApplication extends Component {
  constructor() {
    super();
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(application) {
    ApplicationStore.update(application, resp => {
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${application.id}`);
    });
  }

  render() {
    return(
      <React.Fragment>
        <TitleBar>
          <TitleBarTitle title="Update Application" />
        </TitleBar>
        <ApplicationForm
          submitLabel={i18n.t(`${packageNS}:tr000066`)}
          object={this.props.application}
          onSubmit={this.onSubmit}
          update={true}
        />
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(UpdateApplication));
