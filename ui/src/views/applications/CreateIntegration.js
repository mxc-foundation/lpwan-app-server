import React, { Component } from "react";
import { withRouter } from 'react-router-dom';

import { withStyles } from "@material-ui/core/styles";

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";

import ApplicationStore from "../../stores/ApplicationStore";
import Loader from "../../components/Loader";
import IntegrationForm from "./IntegrationForm";


const styles = {
  card: {
    overflow: "visible",
  },
};


class CreateIntegration extends Component {
  constructor() {
    super();
    this.state = {
      loading: false
    };
    this.onSubmit = this.onSubmit.bind(this);
  }

  componentDidMount() {
    ApplicationStore.get(this.props.match.params.applicationID, resp => {
      this.setState({
        application: resp,
      });
    });
  }

  onSubmit(integration) {
    let integr = integration;
    integr.applicationID = this.props.match.params.applicationID;

    this.setState({ loading: true });

    switch (integr.kind) {
      case "http":
        ApplicationStore.createHTTPIntegration(integr, resp => {
          this.setState({ loading: false });
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        }, error => { this.setState({ loading: false }) });
        break;
      case "influxdb":
        ApplicationStore.createInfluxDBIntegration(integr, resp => {
          this.setState({ loading: false });
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        }, error => { this.setState({ loading: false }) });
        break;
      case "thingsboard":
        ApplicationStore.createThingsBoardIntegration(integr, resp => {
          this.setState({ loading: false });
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        }, error => { this.setState({ loading: false }) });
        break;
      default:
        break;
    }
    this.setState({ loading: false });
  }

  render() {
    const { application } = this.state;
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;

    if (application === undefined) {
      return(<div></div>);
    }

    return(
      <React.Fragment>
        <TitleBar>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000076`)} to={`/organizations/${currentOrgID}/applications`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={application.application.name} to={`/organizations/${currentOrgID}/applications/${currentApplicationID}`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000384`)} to={`/organizations/${currentOrgID}/applications/${currentApplicationID}/integrations`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000277`)} />
        </TitleBar>
        
        {this.state.loading && <Loader />}
        <IntegrationForm
          match={this.props.match}
          onSubmit={this.onSubmit}
          submitLabel={i18n.t(`${packageNS}:tr000277`)}
        />
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(CreateIntegration));
