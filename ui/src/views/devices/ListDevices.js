import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, Route, Switch } from "react-router-dom";
import DeviceAdmin from "../../components/DeviceAdmin";
import TitleBar from "../../components/TitleBar";
import TitleBarButton from "../../components/TitleBarButton";
import TitleBarTitle from "../../components/TitleBarTitle";
import i18n, { packageNS } from "../../i18n";
import theme from "../../theme";
import { default as ListDevicesMap, default as ListDevicesTable } from "./ListDevicesTable";




const styles = {
  buttons: {
    textAlign: "right",
  },
  button: {
    marginLeft: 2 * theme.spacing(1),
  },
  icon: {
    marginRight: theme.spacing(1),
  },
};

class ListDevices extends Component {
  constructor() {
    super();

    this.switchToList = this.switchToList.bind(this);
    this.locationToTab = this.locationToTab.bind(this);
    this.state = {
      viewMode: 'list'
    };
  }

  componentDidMount() {
    this.locationToTab();
  }

  locationToTab = () => {
    if (window.location.href.endsWith("/map")) {
      this.setState({ viewMode: 'map' });
    }
  }

  /**
   * Switch to list
   */
  switchToList() {
    this.setState({ viewMode: 'list' });
  }

  render() {
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;
    const currentApplicationID = this.props.applicationID || this.props.match.params.applicationID;
    let currentServiceProfileID = "";
    if(this.props.application){
      currentServiceProfileID = this.props.application.serviceProfileID;
    }

    return(
      <React.Fragment>
        <TitleBar
          buttons={<DeviceAdmin organizationID={currentOrgID}>
            <TitleBarButton
              key={1}
              label={i18n.t(`${packageNS}:tr000277`)}
              icon={<i className="mdi mdi-plus mr-1 align-middle"></i>}
              to={
                currentApplicationID
                ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices/create`
                : `/organizations/${currentOrgID}/devices/create`
              }
            />
          </DeviceAdmin>}
        >
          <TitleBarTitle title={i18n.t(`${packageNS}:tr000278`)} />
        </TitleBar>

        {
          this.state.viewMode === 'map' && (
            <Link
              to={
                currentApplicationID
                ? `/organizations/${currentOrgID}/applications/${currentApplicationID}/devices`
                : `/organizations/${currentOrgID}/devices`
              }
              className="btn btn-primary mb-3"
              onClick={this.switchToList}>
                Show List
            </Link>
          )
        }
        <Switch>
          <Route exact path={this.props.match.path} render={props =>
            <ListDevicesTable {...props} organizationID={currentOrgID} applicationID={currentApplicationID} serviceProfileID={currentServiceProfileID} />} />
          <Route exact path={`${this.props.match.path}/map`} render={props =>
            <ListDevicesMap {...props} organizationID={currentOrgID} applicationID={currentApplicationID} />} />
        </Switch>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(ListDevices);
