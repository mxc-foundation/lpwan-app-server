import React, { Component } from "react";
import { Route, Redirect, Switch, withRouter, Link } from "react-router-dom";

import { Button } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import TitleBarTitle from "../../components/TitleBarTitle";
import OrganizationStore from "../../stores/OrganizationStore";
import UpdateOrganization from "./UpdateOrganization";


class OrganizationLayout extends Component {
  constructor() {
    super();
    this.state = {};
    this.loadData = this.loadData.bind(this);
    this.deleteOrganization = this.deleteOrganization.bind(this);
  }

  componentDidMount() {
    this.loadData();
  }

  componentDidUpdate(prevProps) {
    if (prevProps === this.props) {
      return;
    }

    this.loadData();
  }

  loadData() {
    OrganizationStore.get(this.props.match.params.organizationID, resp => {
      this.setState({
        organization: resp,
      });
    });
  }

  deleteOrganization() {
    if (window.confirm("Are you sure you want to delete this organization?")) {
      OrganizationStore.delete(this.props.match.params.organizationID, () => {
        this.props.history.push("/organizations");
      });
    }
  }

  render() {
    if (this.state.organization === undefined) {
      return (<div></div>);
    }


    return (
      <React.Fragment>
        <TitleBar
          buttons={[
            <Button color="danger"
              onClick={this.deleteOrganization}
              className="">{i18n.t(`${packageNS}:tr000061`)}
            </Button>,
          ]}
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000049`)} />
        </TitleBar>

        <Switch>
          <Route exact path={this.props.match.path} render={() => <Redirect to={`${this.props.match.url}/edit`} />} />
          <Route exact path={`${this.props.match.path}/edit`} render={props => <UpdateOrganization organization={this.state.organization} {...props} />} />
        </Switch>
      </React.Fragment>
    );
  }
}


export default withRouter(OrganizationLayout);
