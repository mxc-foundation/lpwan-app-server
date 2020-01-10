import React, { Component } from "react";
import { Link, withRouter } from 'react-router-dom';

import ProfileDropdown from './ProfileDropdown';
import DropdownMenuLanguage from "./DropdownMenuLanguage";

import i18n, { packageNS } from '../i18n';
import SessionStore from "../stores/SessionStore";
import WithdrawStore from "../stores/WithdrawStore";
import WalletStore from "../stores/WalletStore";
import logoSm from '../assets/images/logo-sm.png';
import logo from '../assets/images/MATCHX-SUPERNODE2.png';

function getWalletBalance() {
  var organizationId = SessionStore.getOrganizationID();
  if (organizationId) {
    return 0;
  }

  /* if (SessionStore.isAdmin()) {
    organizationId = '0';
  } */

  return new Promise((resolve, reject) => {
    WalletStore.getWalletBalance(organizationId, resp => {
      return resolve(resp);
    });
  });
} 

class Topbar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      balance: 0,
      ProfileMenus: [{
        label: 'Logout',
        icon: 'fe-log-out',
        redirectTo: '/logout',
        hasDivider: true
      }]
    };
  }

  componentDidMount() {
    this.loadData();

    SessionStore.on("organization.change", () => {
      this.loadData();
    });
    WithdrawStore.on("withdraw", () => {
      this.loadData();
    });
  }

  loadData = async () => {
    try {
      var result = await getWalletBalance();
      this.setState({ balance: result.balance });

    } catch (error) {
      console.error(error);
      this.setState({ error });
    }
  }

  onChangeLanguage = (newLanguageState) => {
    this.props.onChangeLanguage(newLanguageState);

    const obj = this.state.ProfileMenus.filter((obj, i, b) => {
      if (obj.hasOwnProperty('redirectTo')) return obj;
    });

    if (obj) {
      this.state.ProfileMenus[0].label = i18n.t(`${packageNS}:menu.settings.logout`);
    }
  }

  render() {
    const { balance } = this.state;

    // let searchbar;
    // searchbar = <Input
    //               placeholder={i18n.t(`${packageNS}:tr000033`)}
    //               className={this.props.classes.search}
    //               disableUnderline={true}
    //               value={this.state.search || ""}
    //               onChange={this.onSearchChange}
    //               startAdornment={
    //                 <InputAdornment position="start">
    //                 <Magnify />
    //                 </InputAdornment>
    //               }
    //             />

    const balanceEl = balance === null ?
      <span className="color-gray">(no org selected)</span> :
      balance + " MXC";
    let user = null;
    if (SessionStore.getUser()) {
      user = SessionStore.getUser();
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
              <ProfileDropdown menuItems={this.state.ProfileMenus} user={user} />
            </li>
          </ul>

          <div className="logo-box">
            <div to="/" className="logo text-center">
              <span className="logo-lg">
                <img src={logo} alt="" height="53" />
              </span>
              <span className="logo-sm">
                <img src={logoSm} alt="" height="16" />
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
