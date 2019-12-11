import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';

//import WithdrawStore from "../stores/WithdrawStore";
//import { SUPER_ADMIN } from "../util/M2mUtil";
//import WalletStore from "../stores/WalletStore";
//import NotificationDropdown from './NotificationDropdown';
import ProfileDropdown from './ProfileDropdown';
import DropdownMenuLanguage from "./DropdownMenuLanguage";

import i18n, { packageNS } from '../i18n';
import SessionStore from "../stores/SessionStore";

import logoSm from '../assets/images/logo-sm.png';
import logo from '../assets/images/logos_wallet_light.png';
import profilePic from '../assets/images/users/profile-icon.png'; 

const Notifications = [{
  id: 1,
  text: 'Caleb Flakelar commented on Admin',
  subText: '1 min ago',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'primary'
},
{
  id: 2,
  text: 'New user registered.',
  subText: '5 min ago',
  icon: 'mdi mdi-account-plus',
  bgColor: 'info'
},
{
  id: 3,
  text: 'Cristina Pride',
  subText: 'Hi, How are you? What about our next meeting',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'success'
},
{
  id: 4,
  text: 'Caleb Flakelar commented on Admin',
  subText: '2 days ago',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'danger'
},
{
  id: 5,
  text: 'Caleb Flakelar commented on Admin',
  subText: '1 min ago',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'primary'
},
{
  id: 6,
  text: 'New user registered.',
  subText: '5 min ago',
  icon: 'mdi mdi-account-plus',
  bgColor: 'info'
},
{
  id: 7,
  text: 'Cristina Pride',
  subText: 'Hi, How are you? What about our next meeting',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'success'
},
{
  id: 8,
  text: 'Caleb Flakelar commented on Admin',
  subText: '2 days ago',
  icon: 'mdi mdi-comment-account-outline',
  bgColor: 'danger'
}];

/* const ProfileMenus = [{
  label: 'My Account',
  icon: 'fe-user',
  redirectTo: "/",
},
{
  label: 'Settings',
  icon: 'fe-settings',
  redirectTo: "/"
},
{
  label: 'Lock Screen',
  icon: 'fe-lock',
  redirectTo: "/"
},
{
  label: 'Logout',
  icon: 'fe-log-out',
  redirectTo: "/logout",
  hasDivider: true
}] [edit] 191126 */

const ProfileMenus = [{
  label: 'Change Password',
  icon: 'mdi mdi-key-change',
  redirectTo: "/",
},
{
  label: 'Log out',
  icon: 'fe-log-out',
  redirectTo: '/logout',
  hasDivider: true
}]

/* function getWalletBalance() {
  var organizationId = SessionStore.getOrganizationID();
  if (organizationId === undefined) {
    return null;
  }

  if (SessionStore.isAdmin()) {
    organizationId = SUPER_ADMIN
  }

  return new Promise((resolve, reject) => {
    WalletStore.getWalletBalance(organizationId, resp => {
      return resolve(resp);
    });
  });
} */

class Topbar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      balance: 0,
      ProfileMenus : [{
        label: 'Logout',
        icon: 'fe-log-out',
        redirectTo: '/logout',
        hasDivider: true
      }]
    };
  }

  componentDidMount() {
    /* this.loadData();

    SessionStore.on("organization.change", () => {
      this.loadData();
    });
    WithdrawStore.on("withdraw", () => {
      this.loadData();
    }); */
  }

  /* loadData = async () => {
    try {
      var result = await getWalletBalance();
      this.setState({ balance: result.balance });

    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  } */

  onChangeLanguage = (newLanguageState) => {
    this.props.onChangeLanguage(newLanguageState);

    const obj = this.state.ProfileMenus.filter((obj, i, b)=>{
      if(obj.hasOwnProperty('redirectTo')) return obj;
    });
    
    if(obj){
      this.state.ProfileMenus[0].label = i18n.t(`${packageNS}:menu.settings.logout`);
    }
  }

  render() {
    const { balance } = this.state;

    const balanceEl = balance === null ? 
      <span className="color-gray">(no org selected)</span> : 
      balance + " MXC";
      let username = null;
      if(SessionStore.getUser()){
        username = SessionStore.getUser().username;
      }
     
    return (
      <React.Fragment>
        <div className="navbar-custom">
          <ul className="list-unstyled topnav-menu float-right mb-0">

            <li className="d-none d-sm-block">
              {/* <form className="app-search">
                <div className="app-search-box">
                  <div className="input-group">
                    <input type="text" className="form-control" placeholder="Search..." />
                    <div className="input-group-append">
                      <button className="btn" type="submit">
                        <i className="fe-search"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </form> */}
            </li>

            <li>
              {/* <NotificationDropdown notifications={Notifications} /> */}
            </li>

            <li className="dropdown notification-list">
              <button className="btn btn-link nav-link right-bar-toggle waves-effect waves-light" onClick={this.props.rightSidebarToggle}>
                <i className="mdi mdi-wallet-outline"></i>
                <span> {balanceEl}</span>
              </button>
            </li>

            <li>
              <DropdownMenuLanguage onChangeLanguage={this.onChangeLanguage} />
            </li>

            <li>
              <ProfileDropdown profilePic={profilePic} menuItems={this.state.ProfileMenus} username={username} />
            </li>

            {/* <li className="dropdown notification-list">
              <button className="btn btn-link nav-link right-bar-toggle waves-effect waves-light" onClick={this.props.rightSidebarToggle}>
                <i className="mdi mdi-help-circle-outline"></i>
              </button>
            </li> */}
          </ul>

          <div className="logo-box">
            <div to="/" className="logo text-center">
              <span className="logo-lg">
                <img src={logo} alt="" height="16" />
              </span>
              <span className="logo-sm">
                <img src={logoSm} alt="" height="24" />
              </span>
            </div>
          </div>

          <ul className="list-unstyled topnav-menu topnav-menu-left m-0">
            <li>
              <button className="button-menu-mobile disable-btn waves-effect" onClick={this.props.menuToggle}>
                <i className="fe-menu"></i>
              </button>
            </li>

            <li>
              <h4 className="page-title-main">{this.props.title}</h4>
            </li>
          </ul>
            
        </div>
      </React.Fragment >
    );
  }
}

export default withRouter(Topbar);
//export default connect()(Topbar); [edit]

