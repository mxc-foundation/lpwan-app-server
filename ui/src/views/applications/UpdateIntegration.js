import React, { Component } from 'react';
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from '../../i18n';
import ApplicationStore from "../../stores/ApplicationStore";
import IntegrationForm from "./IntegrationForm";




class UpdateIntegration extends Component {
  constructor() {
    super();
    this.state = {};
    this.onSubmit = this.onSubmit.bind(this);
    this.deleteIntegration = this.deleteIntegration.bind(this);
  }

  componentDidMount() {
    ApplicationStore.get(this.props.match.params.applicationID, resp => {
      this.setState({
        application: resp,
      });
    });

    switch (this.props.match.params.kind) {
      case "http":
        ApplicationStore.getHTTPIntegration(this.props.match.params.applicationID, resp => {
          let integration = resp.integration;
          integration.kind = "http";

          this.setState({
            integration: integration,
          });
        });
        break;
      case "influxdb":
        ApplicationStore.getInfluxDBIntegration(this.props.match.params.applicationID, resp => {
          let integration = resp.integration;
          integration.kind = "influxdb";

          this.setState({
            integration: integration,
          });
        });
        break;
      case "thingsboard":
        ApplicationStore.getThingsBoardIntegration(this.props.match.params.applicationID, resp => {
          let integration = resp.integration;
          integration.kind = "thingsboard";

          this.setState({
            integration: integration,
          });
        });
        break;
      default:
        break;
    }
  }

  onSubmit(integration) {
    switch (integration.kind) {
      case "http":
        ApplicationStore.updateHTTPIntegration(integration, resp => {
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        });
        break;
      case "influxdb":
        ApplicationStore.updateInfluxDBIntegration(integration, resp => {
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        });
        break;
      case "thingsboard":
        ApplicationStore.updateThingsBoardIntegration(integration, resp => {
          this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
        });
        break;
      default:
        break;
    }
  }

  deleteIntegration() {
    if (window.confirm("Are you sure you want to delete this integration?")) {
      switch(this.props.match.params.kind) {
        case "http":
          ApplicationStore.deleteHTTPIntegration(this.props.match.params.applicationID, resp => {
            this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
          });
          break;
        case "influxdb":
          ApplicationStore.deleteInfluxDBIntegration(this.props.match.params.applicationID, resp => {
            this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
          });
          break;
        case "thingsboard":
          ApplicationStore.deleteThingsBoardIntegration(this.props.match.params.applicationID, resp => {
            this.props.history.push(`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`);
          });
          break;
        default:
          break;
      }
    }
  }

  render() {
    if (this.state.application === undefined || this.state.integration === undefined) {
      return(<div></div>);
    }

    return(
      <React.Fragment>
        <TitleBar
          buttons={[
            <TitleBarButton
              key={1}
              icon={<i className="mdi mdi-delete mr-1 align-middle"></i>}
              label={i18n.t(`${packageNS}:tr000061`)}
              onClick={this.deleteIntegration}
            />,
          ]}
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000076`)} to={`/organizations/${this.props.match.params.organizationID}/applications`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={this.state.application.application.name} to={`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000384`)} to={`/organizations/${this.props.match.params.organizationID}/applications/${this.props.match.params.applicationID}/integrations`} />
          <span>&nbsp;</span>
          <TitleBarTitle title="/" />
          <span>&nbsp;</span>
          <TitleBarTitle title={this.props.match.params.kind} />
        </TitleBar>
        <IntegrationForm
          object={this.state.integration}
          onSubmit={this.onSubmit}
          submitLabel={i18n.t(`${packageNS}:tr000066`)}
          update={true}
        />
      </React.Fragment>
    );
  }
}

export default UpdateIntegration;
