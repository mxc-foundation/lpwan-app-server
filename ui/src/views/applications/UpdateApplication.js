import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';
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
    this.state = {
      loading: false
    }
    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(application) {
    this.setState({ loading: true });
    ApplicationStore.update(application, resp => {
      this.setState({ loading: false });
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${application.id}`);
    }, error => { this.setState({ loading: false }) });
  }

  render() {
    return(
      <React.Fragment>
        <TitleBar>
          <TitleBarTitle title="Update Application" />
        </TitleBar>

        <div className="position-relative">
          {this.state.loading && <Loader />}

          <ApplicationForm
            submitLabel={i18n.t(`${packageNS}:tr000066`)}
            object={this.props.application}
            onSubmit={this.onSubmit}
            update={true}
          />
        </div>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(UpdateApplication));
