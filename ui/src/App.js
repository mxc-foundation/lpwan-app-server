import React, { Component } from "react";
import { Router, Route, Switch, Redirect } from 'react-router-dom';


import CssBaseline from "@material-ui/core/CssBaseline";
import { MuiThemeProvider, withStyles } from "@material-ui/core/styles";
import Modal from "./components/Modal";

import history from "./history";
import themeChinese from "./themeChinese";
import theme from "./theme";
import i18n, { packageNS, DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from "./i18n";

import './assets/scss/DefaultTheme.scss';

import Topbar from "./components/Topbar";
import Sidebar from "./components/Sidebar";
import Footer from "./components/Footer";
import AuthLayout from "./components/AuthLayout";
import NonAuthLayout from "./components/NonAuthLayout";

import Notifications from "./components/Notifications";
import SessionStore from "./stores/SessionStore";

// network-server
import ListNetworkServers from "./views/network-servers/ListNetworkServers";
import CreateNetworkServer from "./views/network-servers/CreateNetworkServer";
import NetworkServerLayout from "./views/network-servers/NetworkServerLayout";

// SMB Gateway Profile
import ListGatewayProfiles from "./views/gateway-profiles/ListGatewayProfiles";
import CreateGatewayProfile from "./views/gateway-profiles/CreateGatewayProfile";
import GatewayProfileLayout from "./views/gateway-profiles/GatewayProfileLayout";

// organization
import ListOrganizations from "./views/organizations/ListOrganizations";
import CreateOrganization from "./views/organizations/CreateOrganization";
import OrganizationLayout from "./views/organizations/OrganizationLayout";
import ListOrganizationUsers from "./views/organizations/ListOrganizationUsers";
import OrganizationUserLayout from "./views/organizations/OrganizationUserLayout";
import CreateOrganizationUser from "./views/organizations/CreateOrganizationUser";
import OrganizationRedirect from "./views/organizations/OrganizationRedirect";

// user
import Login from "./views/users/Login";
import Logout from "./views/users/Logout";
import Registration from "./views/users/Registration";
import RegistrationConfirm from "./views/users/RegistrationConfirm";

import ListUsers from "./views/users/ListUsers";
import CreateUser from "./views/users/CreateUser";
import UserLayout from "./views/users/UserLayout";
import ChangeUserPassword from "./views/users/ChangeUserPassword";
import PasswordRecovery from "./views/users/PasswordRecovery";
import PasswordResetConfirm from "./views/users/PasswordResetConfirm";

// service-profile
import ListServiceProfiles from "./views/service-profiles/ListServiceProfiles";
import CreateServiceProfile from "./views/service-profiles/CreateServiceProfile";
import ServiceProfileLayout from "./views/service-profiles/ServiceProfileLayout";

// device-profile
import ListDeviceProfiles from "./views/device-profiles/ListDeviceProfiles";
import CreateDeviceProfile from "./views/device-profiles/CreateDeviceProfile";
import DeviceProfileLayout from "./views/device-profiles/DeviceProfileLayout";

// gateways
import ListGateways from "./views/gateways/ListGateways";
import GatewayLayout from "./views/gateways/GatewayLayout";
import CreateGateway from "./views/gateways/CreateGateway";
import BrandSelect from "./views/gateways/BrandSelect";
import EnterSerialNum from "./views/gateways/EnterSerialNum";

// applications
import ListApplications from "./views/applications/ListApplications";
import CreateApplication from "./views/applications/CreateApplication";
import ApplicationLayout from "./views/applications/ApplicationLayout";

// multicast-groups
import ListMulticastGroups from "./views/multicast-groups/ListMulticastGroups";
import CreateMulticastGroup from "./views/multicast-groups/CreateMulticastGroup";
import MulticastGroupLayout from "./views/multicast-groups/MulticastGroupLayout";
import AddDeviceToMulticastGroup from "./views/multicast-groups/AddDeviceToMulticastGroup";

// device
import DeviceLayoutM2M from "./views/devices/m2m/DeviceLayoutM2M";

import OrganizationDevices from "./views/devices/OrganizationDevices";
import CreateDevice from "./views/devices/CreateDevice";
import DeviceLayout from "./views/devices/DeviceLayout";

// search
import Search from "./views/search/Search";

// FUOTA
import CreateFUOTADeploymentForDevice from "./views/fuota/CreateFUOTADeploymentForDevice";
import FUOTADeploymentLayout from "./views/fuota/FUOTADeploymentLayout";
import FUOTADeploymentDetails from "./views/fuota/FUOTADeploymentDetails";
//M2M
import ModifyEthAccount from "./views/ethAccount/ModifyEthAccount";
import Withdraw from "./views/withdraw/Withdraw";
import Topup from "./views/topup/Topup";
import HistoryLayout from "./views/history/HistoryLayout";

import StakeLayout from "./views/stake/StakeLayout";
import SetStake from "./views/stake/SetStake";
import SuperNodeEth from "./views/controlPanel/superNodeEth/superNodeEth"
import SuperAdminWithdraw from "./views/controlPanel/withdraw/withdraw"
import SupernodeHistory from "./views/controlPanel/history/History"
import SystemSettings from "./views/controlPanel/settings/Settings"

//dashboard
import Dashboard from "./views/dashboards/";
// home
import HomeComponent from './views/Home';

//Temp banner
import TopBanner from "./components/TopBanner";
import { initJWTTimer } from "./util/JWTUti";

const drawerWidth = 270;

const styles = {
  root: {
    flexGrow: 1,
    display: "flex",
    minHeight: "100vh",
    flexDirection: "column",
    backgroundColor: "#070033",
    // backgroundImage: 'url("/img/world-map.png")',
    backgroundRepeat: 'no-repeat',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    //backgroundColor: '#cccccc',
    //background: "#311b92",
    fontFamily: 'Montserrat',
  },
  input: {
    color: theme.palette.textPrimary.main.white,
  },
  paper: {
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  main: {
    width: "100%",
    padding: 2 * 24,
    paddingTop: 115,
    /* display: 'flex',
    alignItems: 'center', */
    flex: 1,
  },
  mainDrawerOpen: {
    paddingLeft: drawerWidth + (2 * 24),
  },
  footerDrawerOpen: {
    paddingLeft: drawerWidth,
  },
  color: {
    backgroundColor: theme.palette.secondary.main,
  },
};



class NotLoggedinRoute extends Component {
  route = (props) => {
    const { Comp, user, ...otherProps } = this.props;
    console.log('NotLoggedinRoute', Comp.name, user);

    return !user ?
      <Comp {...otherProps} /> :
      <Redirect to='/' />;
  }
  render() {
    return (
      <Route {...this.props}
        render={this.route}
      />
    )
  }
}
/*
class HomeRoute extends Component {
  route = () => {
    const { Comp, user, ...otherProps } = this.props;
    const orgs = SessionStore.getOrganizations();

    if (SessionStore.getToken() && orgs.length > 0) {
      return <Redirect to={`/organizations/${orgs[0].organizationID}`}></Redirect>;
    } else {
      console.log('User has no organisations. Redirecting to login');
      return <Redirect to={"/logout"}></Redirect>;
    }
  }

  render() {
    return (
      <Route {...this.props}
        render={this.route}
      />
    )
  }
}*/

class LoggedInRoutes extends Component {
  render() {
    const user = SessionStore.getUser();
    const language = SessionStore.getLanguage();

    if (!user) {
      return <Redirect to={"/login"}></Redirect>;
    }

    return (<>
      <Switch>
        <Route exact path="/" component={HomeComponent} />
        <Route exact path="/logout" component={Logout} />

        <Route exact path="/dashboard" component={Dashboard} />

        <Route exact path="/users" component={ListUsers} />
        <Route exact path="/users/create" component={CreateUser} />
        <Route exact path="/users/:userID(\d+)" component={UserLayout} />
        <Route exact path="/users/:userID(\d+)/password" component={ChangeUserPassword} />
        
        <Route exact path="/network-servers" component={ListNetworkServers} />
        <Route exact path="/network-servers/create" component={CreateNetworkServer} />
        <Route path="/network-servers/:networkServerID" component={NetworkServerLayout} />

        <Route exact path="/gateway-profiles" component={ListGatewayProfiles} />
        <Route exact path="/gateway-profiles/create" component={CreateGatewayProfile} />
        <Route path="/gateway-profiles/:gatewayProfileID([\w-]{36})" component={GatewayProfileLayout} />

        <Route exact path="/organizations/:organizationID(\d+)/service-profiles" component={ListServiceProfiles} />
        <Route exact path="/organizations/:organizationID(\d+)/service-profiles/create" component={CreateServiceProfile} />
        <Route path="/organizations/:organizationID(\d+)/service-profiles/:serviceProfileID([\w-]{36})" component={ServiceProfileLayout} />

        <Route exact path="/organizations/:organizationID(\d+)/device-profiles" component={ListDeviceProfiles} />
        <Route exact path="/organizations/:organizationID(\d+)/device-profiles/create" component={CreateDeviceProfile} />
        <Route path="/organizations/:organizationID(\d+)/device-profiles/:deviceProfileID([\w-]{36})" component={DeviceProfileLayout} />

        <Route exact path="/organizations/:organizationID(\d+)/gateways/create" component={CreateGateway} />
        <Route exact path="/organizations/:organizationID(\d+)/gateways/brand" component={BrandSelect} />
        <Route exact path="/organizations/:organizationID(\d+)/gateways/input-serial" component={EnterSerialNum} />
        <Route path="/organizations/:organizationID(\d+)/gateways/:gatewayID([\w]{16})" component={GatewayLayout} />
        <Route path="/organizations/:organizationID(\d+)/gateways" component={ListGateways} />

        <Route exact path="/organizations/:organizationID(\d+)/devices/create" component={CreateDevice} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/edit" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/keys" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/activation" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/data" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/frames" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/fuota-deployments" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})/fuota-deployments/create" component={CreateFUOTADeploymentForDevice} />
        <Route exact path="/organizations/:organizationID(\d+)/devices/:devEUI([\w]{16})" component={DeviceLayout} />
        <Route exact path="/organizations/:organizationID(\d+)/devices" component={DeviceLayoutM2M} />

        <Route exact path="/organizations/:organizationID(\d+)/applications" component={ListApplications} />
        <Route exact path="/organizations/:organizationID(\d+)/applications/create" component={CreateApplication} />
        <Route exact path="/organizations/:organizationID(\d+)/applications/:applicationID(\d+)/devices/create" component={CreateDevice} />
        <Route exact path="/organizations/:organizationID(\d+)/applications/:applicationID(\d+)/devices/:devEUI([\w]{16})/fuota-deployments/create" component={CreateFUOTADeploymentForDevice} />
        <Route path="/organizations/:organizationID(\d+)/applications/:applicationID(\d+)/fuota-deployments/:fuotaDeploymentID([\w-]{36})" component={FUOTADeploymentLayout} />
        <Route path="/organizations/:organizationID(\d+)/applications/:applicationID(\d+)/devices/:devEUI([\w]{16})" component={DeviceLayout} />
        <Route path="/organizations/:organizationID(\d+)/applications/:applicationID(\d+)" component={ApplicationLayout} />

        <Route exact path="/organizations/:organizationID(\d+)/multicast-groups" component={ListMulticastGroups} />
        <Route exact path="/organizations/:organizationID(\d+)/multicast-groups/create" component={CreateMulticastGroup} />
        <Route exact path="/organizations/:organizationID(\d+)/multicast-groups/:multicastGroupID/devices/create" component={AddDeviceToMulticastGroup} />
        <Route path="/organizations/:organizationID(\d+)/multicast-groups/:multicastGroupID([\w-]{36})" component={MulticastGroupLayout} />

        <Route exact path="/organizations" component={ListOrganizations} />
        <Route exact path="/organizations/create" component={CreateOrganization} />
        <Route exact path="/organizations/:organizationID(\d+)/users" component={ListOrganizationUsers} />
        <Route exact path="/organizations/:organizationID(\d+)/users/create" component={CreateOrganizationUser} />
        <Route exact path="/organizations/:organizationID(\d+)/users/:userID(\d+)" component={OrganizationUserLayout} />
        <Route path="/organizations/:organizationID(\d+)" component={OrganizationLayout} />

        <Route path="/modify-account/:organizationID" component={ModifyEthAccount} />
        <Route path="/withdraw/:organizationID"
          render={props =>
            <Withdraw
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route path="/topup/:organizationID"
          render={props =>
            <Topup
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route path="/history/:organizationID"
          render={props =>
            <HistoryLayout
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route exact path="/stake/:organizationID" component={StakeLayout} />
        <Route exact path="/stake/:organizationID/set-stake" component={SetStake} />
        <Route path="/control-panel/modify-account"
          render={props =>
            <SuperNodeEth
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route path="/control-panel/withdraw"
          render={props =>
            <SuperAdminWithdraw
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route path="/control-panel/history"
          render={props =>
            <SupernodeHistory
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />
        <Route path="/control-panel/system-settings"
          render={props =>
            <SystemSettings
              {...props}
              switchToSidebarId={this.props.switchToSidebarId}
            />
          }
        />

        <Route exact path="/search" component={Search} />
      </Switch>
    </>);
  }
}

class App extends Component {
  constructor() {
    super();

    this.state = {
      user: null,
      organizationId: null,
      drawerOpen: false,
      language: null,
      nsDialog: null,
      width: null,
      theme: theme
    };

    this.setDrawerOpen = this.setDrawerOpen.bind(this);
  }

  componentWillMount() {
    window.addEventListener('resize', this.handleWindowSizeChange);
  }

  // make sure to remove the listener
  // when the component is not mounted anymore
  componentWillUnmount() {
    window.removeEventListener('resize', this.handleWindowSizeChange);
  }

  handleWindowSizeChange = () => {
    this.setState({ width: window.innerWidth });
  };

  onSessionInvalid = (err) => {
    // TODO: trigger modal
    //SessionStore.logout(()=>{});
    console.log('err.message', err.message);
  }

  componentDidMount = () => {
    initJWTTimer(this.onSessionInvalid);

    SessionStore.on("change", () => {
      this.setState({
        user: SessionStore.getUser(),
        organizationId: SessionStore.getOrganizationID(),
        drawerOpen: SessionStore.getUser() != null,
        language: SessionStore.getLanguage(),
        sessionInitialized: true
      });
    });

    const storedLanguageID = SessionStore.getLanguage() && SessionStore.getLanguage().id;

    if (!storedLanguageID && !i18n.language) {
      i18n.changeLanguage(DEFAULT_LANGUAGE.id, (err, t) => {
        if (err) {
          console.error(`Error setting default language to English: `, err);
        }
      });
    }

    const i18nLanguage = SUPPORTED_LANGUAGES.find(el => el.id === i18n.language);

    // Add the saved i18n language back into Local Storage if it is lost after a page refresh on Login component
    if (!storedLanguageID && i18n.language) {
      SessionStore.setLanguage(i18nLanguage);
    }

    // Language stored in Local Storage persists and takes precedence over i18n language
    if (storedLanguageID && i18n.language !== storedLanguageID) {
      i18n.changeLanguage(storedLanguageID, (err, t) => {
        if (err) {
          console.error(`Error loading language ${storedLanguageID}: `, err);
        }
      });
    }

    this.setState({
      user: SessionStore.getUser(),
      organizationId: SessionStore.getOrganizationID(),
      drawerOpen: SessionStore.getUser() != null,
      language: storedLanguageID ? SessionStore.getLanguage() : i18nLanguage
    });
  }

  onChangeLanguage = (newLanguage) => {
    SessionStore.setLanguage(newLanguage);

    i18n.changeLanguage(newLanguage.id, (err, t) => {
      if (err) {
        console.error(`Error loading language ${newLanguage.id}: `, err);
      }
    });

    this.setState({
      language: newLanguage,
      theme: newLanguage.code === 'cn' ? themeChinese : theme
    });
  }

  setDrawerOpen(state) {
    this.setState({
      drawerOpen: state,
    });
  }

  /**
   * toggle Menu
   */
  toggleMenu = (e) => {
    e.preventDefault();
    this.setState({ isCondensed: !this.state.isCondensed });
  }

  /**
   * Toggle right side bar
   */
  toggleRightSidebar = () => {
    document.body.classList.toggle("right-bar-enabled");
  }

  switchToSidebarId = (newSidebarId) => {
    this.setState({
      currentSidebarId: newSidebarId
    })
  }

  logout = () => {
    localStorage.clear();
    /* SessionStore.logout(() => {
      this.props.history.push("/login");
    }); */
  }

  render() {


    let topNav = null;
    let sideNav = null;
    let topbanner = null;

    const { currentSidebarId, isCondensed, language, sessionInitialized, user, width } = this.state;

    let isMobile = width <= 800;

    let Layout = NonAuthLayout;

    if (user !== null) {
      /* sideNav = <SideNav open={this.state.drawerOpen} user={this.state.user} />
      topbanner = <TopBanner setDrawerOpen={this.setDrawerOpen} drawerOpen={this.state.drawerOpen} user={this.state.user} organizationId={this.state.organizationId}/>; 
      topNav = (
        <TopNav
          drawerOpen={this.state.drawerOpen}
          language={language}
          onChangeLanguage={this.onChangeLanguage}
          organizationId={this.state.organizationId}
          setDrawerOpen={this.setDrawerOpen}
          user={this.state.user}
        /> 
        ); */
      topNav = <Topbar rightSidebarToggle={this.toggleRightSidebar} isMobile={isMobile} onChangeLanguage={this.onChangeLanguage} menuToggle={this.toggleMenu} {...this.props} />;
      sideNav = <Sidebar isCondensed={isCondensed} currentSidebarId={currentSidebarId} switchToSidebarId={this.switchToSidebarId} {...this.props} />;

      // if user is logged in - set auth layout
      Layout = AuthLayout;
    }

    // if (!sessionInitialized) {
    //   return 'loading...';
    // }

    return (
      <Router history={history}>
        <React.Fragment>
          {this.state.nsDialog && <Modal
            title={i18n.t(`${packageNS}:menu.topup.notice`)}
            left={"DISMISS"}
            right={"ADD ETH ACCOUNT"}
            context={"Session was expired. Please, login again."}
            callback={this.logout} />}
          <MuiThemeProvider theme={this.state.theme}>
            <CssBaseline />
            {/* <div className={this.props.classes.root}> */}

            <Layout topBar={topNav} topBanner={topbanner} sideNav={sideNav}>

              { /* TODO: move all routing to its own file */}
              <Switch>
                <NotLoggedinRoute exact path="/login"
                  Comp={Login} user={user}
                  language={language}
                  onChangeLanguage={this.onChangeLanguage} />
                <Route exact path="/registration" component={Registration} />
                <Route exact path="/password-recovery" component={PasswordRecovery} />
                <Route exact path="/password-reset-confirm" component={PasswordResetConfirm} />
                <Route exact path="/registration-confirm/:securityToken"
                  render={props =>
                    <RegistrationConfirm {...props}
                      language={language}
                      onChangeLanguage={this.onChangeLanguage}
                    />
                  }
                />
                <LoggedInRoutes switchToSidebarId={this.switchToSidebarId}/>

              </Switch>
              <Footer />
            </Layout>

            <Notifications />
          </MuiThemeProvider>
        </React.Fragment>
      </Router>
    );
  }
}

export default withStyles(styles)(App);
