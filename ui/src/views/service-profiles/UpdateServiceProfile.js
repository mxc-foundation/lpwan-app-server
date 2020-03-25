import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import ServiceProfileStore from "../../stores/ServiceProfileStore";
import ServiceProfileForm from "./ServiceProfileForm";



class UpdateServiceProfile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false
    };

    this.onSubmit = this.onSubmit.bind(this);

  }

  onSubmit(serviceProfile) {
    this.setState({ loading: true });
    ServiceProfileStore.update(serviceProfile, resp => {
      this.setState({ loading: false });
      this.props.history.push(`/organizations/${this.props.match.params.organizationID}/service-profiles`);
    }, error => { this.setState({ loading: false }) });
  }

  render() {
    return (<div className="position-relative">
      
      {this.state.loading && <Loader />}
      <ServiceProfileForm
        submitLabel={i18n.t(`${packageNS}:tr000066`)}
        object={this.props.serviceProfile}
        onSubmit={this.onSubmit}
        update={true}
        match={this.props.match}
        disabled={!this.props.admin}
      ></ServiceProfileForm>
    </div>
    );
  }
}

export default withRouter(UpdateServiceProfile);
