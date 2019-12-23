import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { UncontrolledDropdown, DropdownMenu, DropdownItem, DropdownToggle } from 'reactstrap';
import { getLoraHost } from "../util/M2mUtil";
import DropdownMenu2 from "./DropdownMenu";
import ProfileDropdown from './ProfileDropdown';
import PerfectScrollbar from 'react-perfect-scrollbar';
import 'react-perfect-scrollbar/dist/css/styles.css';
import MetisMenu from 'metismenujs/dist/metismenujs';
import mxcLogo from '../assets/images/mxc_logo.png';
import profilePic from '../assets/images/users/profile-icon.png';
import Divider from '@material-ui/core/Divider';
import SessionStore from '../stores/SessionStore';
import Admin from '../components/Admin';
import NonAdmin from '../components/NonAdmin';
import { SUPERNODE_WALLET, SUPERNODE_SETTING, DEFAULT, WALLET, SETTING } from '../util/Data';
import i18n, { packageNS } from '../i18n';
import OrganizationStore from "../stores/OrganizationStore";
import ServerInfoStore from '../stores/ServerInfoStore';

const ProfileMenus = [{
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
}]

const SideNavContent = (props) => {
    return <React.Fragment>

        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                {/* <li>
                    <Link to="/dashboard" className="waves-effect side-nav-link-ref">
                        <i className="mdi mdi-view-dashboard"></i>
                        <span> Dashboard </span>
                    </Link>
                </li> */}
                <Admin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.control_panel`)}</li>

                    <li>
                        <Link to="/control-panel/withdraw/" className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(SUPERNODE_WALLET)}>
                            <i className="mdi mdi-wallet"></i>
                            <span> {i18n.t(`${packageNS}:tr000084`)} </span>
                            <span className="menu-arrow"></span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/network-servers`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-server"></i>
                            <span> {i18n.t(`${packageNS}:tr000420`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/gateway-profiles`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-remote"></i>
                            <span> {i18n.t(`${packageNS}:tr000046`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-domain"></i>
                            <span> {i18n.t(`${packageNS}:tr000049`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/users`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-account-multiple"></i>
                            <span> {i18n.t(`${packageNS}:tr000055`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/edit`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-domain"></i>
                            <span> {i18n.t(`${packageNS}:tr000418`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to="/control-panel/wallet/" className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(SUPERNODE_SETTING)}>
                            <i className="mdi mdi-settings"></i>
                            <span> Setting </span>
                            <span className="menu-arrow"></span>
                        </Link>
                    </li>
                </Admin>

                <NonAdmin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>
                    <li>
                        {/* <DropdownMenu default={ this.state.default } onChange={this.onChange} /> [edit] */}
                        <DropdownMenu2 default={props.default} onChange={props.onChange} />
                    </li>

                    <li>
                        <Link to={`/withdraw/${props.orgId}`} className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(WALLET)}>
                            <i className="mdi mdi-wallet"></i>
                            <span> {i18n.t(`${packageNS}:tr000084`)} </span>
                            <span className="menu-arrow"></span>
                        </Link>
                    </li>
                    
                    <li>
                        <Link to={`/organizations/${props.orgId}/users`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-account-details"></i>
                            <span> {i18n.t(`${packageNS}:tr000067`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/withdraw/${props.orgId}`} className="waves-effect side-nav-link-ref">
                            {/* <i className="mdi mdi-cloud-print-outline"></i> */}
                            <i className="mdi mdi-vote"></i>
                            <span> {i18n.t(`${packageNS}:menu.staking.staking`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/gateways`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-remote"></i>
                            <span> {i18n.t(`${packageNS}:menu.gateways.gateway`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/device-profiles`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-memory"></i>
                            <span> {i18n.t(`${packageNS}:tr000278`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/service-profiles`}>
                            <i className="mdi mdi-folder-account"></i>
                            <span> {i18n.t(`${packageNS}:tr000078`)} </span>
                        </Link>
                    </li>
                    
                    <li>
                        <Link to={`/organizations/${props.orgId}/applications`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-apps"></i>
                            <span> {i18n.t(`${packageNS}:tr000407`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/multicast-groups`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-podcast"></i>
                            <span> {i18n.t(`${packageNS}:tr000083`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/modify-account/${props.orgId}`} className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(SETTING)}>
                            <i className="mdi mdi-settings"></i>
                            <span> Setting </span>
                            <span className="menu-arrow"></span>
                        </Link>
                    </li>

                </NonAdmin>

                <li>
                    <Divider />
                </li>

                <li>
                    <Link to={'#'} className="waves-effect side-nav-link-ref">
                        <i className="mdi mdi-view-dashboard"></i>
                        <span> {i18n.t(`${packageNS}:menu.nb_iot_server`)} </span>
                    </Link>
                </li>

                {/* <li>
                    <Link to={'/ext'} className="waves-effect side-nav-link-ref">
                        <i className="mdi mdi-view-dashboard"></i>
                        <span> {i18n.t(`${packageNS}:menu.lpwan_server`)} </span>
                    </Link>
                </li> */}

                <li>
                    <Link to={'#'} className="waves-effect side-nav-link-ref">
                        <span> {i18n.t(`${packageNS}:menu.powered_by`)} </span>&nbsp;
                        <img src={mxcLogo} className="iconStyle" alt={i18n.t(`${packageNS}:menu.lora_server`)} />
                    </Link>
                </li>

                <li>
                    <Link to={'#'} className="waves-effect side-nav-link-ref">
                        <span> {i18n.t(`${packageNS}:menu.version`)}: {props.version} </span>
                    </Link>
                </li>

            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

const SideNavSupernodeWalletContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <Admin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.control_panel`)}</li>

                    <li>
                        <Link to={`/control-panel/withdraw/`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span> {'back to control pannel'} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/withdraw/`} className="waves-effect side-nav-link-ref">
                            <i className="ti-cloud-down"></i>
                            <span> {i18n.t(`${packageNS}:menu.withdraw.withdraw`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/history/`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-history"></i>
                            <span> {i18n.t(`${packageNS}:menu.history.history`)} </span>
                        </Link>
                    </li>
                </Admin>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

const SideNavSupernodeSettingContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <Admin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.control_panel`)}</li>

                    <li>
                        <Link to={`/control-panel/withdraw/`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span> {'back to control pannel'} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/withdraw/`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-ethereum"></i>
                            <span> {i18n.t(`${packageNS}:menu.eth_account.eth_account`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/history/`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-settings"></i>
                            <span> {i18n.t(`${packageNS}:tr000417`)} </span>
                        </Link>
                    </li>
                </Admin>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

const SideNavWalletContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                    <li>
                        <Link to={`/modify-account/${props.orgId}`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span> {'back'} </span>
                        </Link>
                    </li>

                    <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>

                    <li>
                        <DropdownMenu2 default={props.default} onChange={props.onChange} />
                    </li>

                    <li>
                        <Link to={`/topup/${props.orgId}`} className="waves-effect side-nav-link-ref">
                            {/* <i className="mdi mdi-inbox-arrow-down"></i>
                        <i className="mdi mdi-basket-fill"></i> */}
                            <i className="ti-cloud-up"></i>
                            <span> {i18n.t(`${packageNS}:menu.topup.topup`)} </span>
                        </Link>
                    </li>
                    
                    <li>
                        <Link to={`/withdraw/${props.orgId}`} className="waves-effect side-nav-link-ref">
                            {/* <i className="mdi mdi-cloud-print-outline"></i> */}
                            <i className="ti-cloud-down"></i>
                            <span> {i18n.t(`${packageNS}:menu.withdraw.withdraw`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/history/${props.orgId}`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-history"></i>
                            <span> {i18n.t(`${packageNS}:menu.history.history`)} </span>
                        </Link>
                    </li>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

const SideNavSettingContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                    <li>
                        <Link to={`/modify-account/${props.orgId}`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span> {'back'} </span>
                        </Link>
                    </li>
                    
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>

                    <li>
                        <DropdownMenu2 default={props.default} onChange={props.onChange} />
                    </li>

                    <li>
                        <Link to={`/modify-account/${props.orgId}`} className="waves-effect side-nav-link-ref">
                            {/* <i className="mdi mdi-inbox-arrow-down"></i>
                        <i className="mdi mdi-basket-fill"></i> */}
                            <i className="mdi mdi-ethereum"></i>
                            <span> {i18n.t(`${packageNS}:menu.eth_account.eth_account`)} </span>
                        </Link>
                    </li>

                    
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

const UserProfile = () => {
    return <React.Fragment>
        <div className="user-box text-center">
            <img src={profilePic} alt="user-img" title="Nik Patel" className="rounded-circle img-thumbnail avatar-lg" />
            <UncontrolledDropdown>
                <DropdownToggle caret tag="a" className="text-dark dropdown-toggle h5 mt-2 mb-1 d-block">
                    Nik Patel
                </DropdownToggle>
                <DropdownMenu className="user-pro-dropdown">
                    <DropdownItem>
                        <i className="fe-user mr-1"></i>
                        <span>My Account</span>
                    </DropdownItem>
                    <DropdownItem>
                        <i className="fe-settings mr-1"></i>
                        <span>Settings</span>
                    </DropdownItem>
                    <DropdownItem>
                        <i className="fe-lock mr-1"></i>
                        <span>Lock Screen</span>
                    </DropdownItem>
                    <DropdownItem>
                        <i className="fe-log-out mr-1"></i>
                        <span>Logout</span>
                    </DropdownItem>
                </DropdownMenu>
            </UncontrolledDropdown>

            <p className="text-muted">Admin Head</p>
            <ul className="list-inline">
                {/* <li className="list-inline-item">
                    <Link to="/" className="text-muted">
                        <i className="mdi mdi-settings"></i>
                    </Link>
                </li> [edit]*/}

                <li className="list-inline-item">
                    <Link to="/logout" className="text-custom">
                        <i className="mdi mdi-power"></i>
                    </Link>
                </li>
            </ul>
        </div>
    </React.Fragment>
}

function loadServerVersion() {
    return new Promise((resolve, reject) => {
        ServerInfoStore.getVersion(data => {
            resolve(data);
        });
    });
}

class Sidebar extends Component {
    constructor(props) {
        super(props);

        this.state = {
            open: true,
            //organization: {},
            organization: null,
            organizationID: '',
            cacheCounter: 0,
            version: '1.0.0',
            sidebar: DEFAULT
        };

        this.handleOtherClick = this.handleOtherClick.bind(this);
        this.initMenu = this.initMenu.bind(this);
    }
    loadData = async () => {
        try {
          const organizationID = SessionStore.getOrganizationID();
          /* var data = await loadServerVersion();
          const serverInfo = JSON.parse(data); */
          
          this.setState({
            organizationID,
            //version: serverInfo.version
          })
    
          this.setState({loading: true})
          
        } catch (error) {
          this.setState({loading: false})
          console.error(error);
          this.setState({ error });
        }
      }

    componentDidMount = () => {
        this.initMenu();
        this.loadData();
        SessionStore.on("organization.change", () => {
            OrganizationStore.get(SessionStore.getOrganizationID(), resp => {
              this.setState({
                organization: resp.organization,
              });
            });
          });
      
          OrganizationStore.on("create", () => {
            this.setState({
              cacheCounter: this.state.cacheCounter + 1,
            });
          });
      
          OrganizationStore.on("change", (org) => {
            if (this.state.organization !== null && this.state.organization.id === org.id) {
              this.setState({
                organization: org,
              });
            }
      
            this.setState({
              cacheCounter: this.state.cacheCounter + 1,
            });
          });
      
          OrganizationStore.on("delete", id => {
            if (this.state.organization !== null && this.state.organization.id === id) {
              this.setState({
                organization: null
              });
            }
      
            this.setState({
              cacheCounter: this.state.cacheCounter + 1,
            });
          });
      
          this.getOrganizationFromLocation();
    }
    
      onChange = (e) => {
        SessionStore.setOrganizationID(e.target.value);
        this.setState({
            organizationID: e.target.value
        })
        this.props.history.push(`/organizations/${e.target.value}/applications`);
      }
    
      getOrganizationFromLocation() {
        const organizationRe = /\/organizations\/(\d+)/g;
        const match = organizationRe.exec(this.props.history.location.pathname);
    
        if (match !== null && (this.state.organization === null || this.state.organization.id !== match[1])) {
          SessionStore.setOrganizationID(match[1]);
        }
      }
    
      getOrganizationOption(id, callbackFunc) {
        OrganizationStore.get(id, resp => {
          callbackFunc({label: resp.organization.name, value: resp.organization.id, color:"black"});
        });
      }
    
      getOrganizationOptions(search, callbackFunc) {
        OrganizationStore.list(search, 10, 0, resp => {
          const options = resp.result.map((o, i) => {return {label: o.name, value: o.id, color:'black'}});
          callbackFunc(options);
        });
      }
    
      /* handlingExtLink = () => {
        const resp = SessionStore.getProfile();
        resp.then((res) => {
          let orgId = SessionStore.getOrganizationID();
          const isBelongToOrg = res.body.organizations.some(e => e.organizationID === SessionStore.getOrganizationID());
    
          OrganizationStore.get(orgId, resp => {
            openM2M(resp.organization, isBelongToOrg, '/modify-account');
          });
        })
      } */

    /**
     * Bind event
     */
    componentWillMount = () => {
        document.addEventListener('mousedown', this.handleOtherClick, false);
    }


    

    /**
     * Component did update
     */
    componentDidUpdate = (prevProps) => {
        
        if (this.props.isCondensed !== prevProps.isCondensed) {
            if (prevProps.isCondensed) {
                document.body.classList.remove("sidebar-enable");
                document.body.classList.remove("enlarged");
            } else {
                document.body.classList.add("sidebar-enable");
                const isSmallScreen = window.innerWidth < 768;
                if (!isSmallScreen) {
                    document.body.classList.add("enlarged");
                }
            }

            this.initMenu();
        }

        if (this.props === prevProps) {
            return;
        }

        this.getOrganizationFromLocation();
    }

    /**
     * Bind event
     */
    componentWillUnmount = () => {
        document.removeEventListener('mousedown', this.handleOtherClick, false);
    }

    /**
     * Handle the click anywhere in doc
     */
    handleOtherClick = (e) => {
        if (this.menuNodeRef.contains(e.target))
            return;
        // else hide the menubar
        document.body.classList.remove('sidebar-enable');
    }

    /**
     * Init the menu
     */
    initMenu = () => {
        // render menu
        new MetisMenu("#side-menu");
        var links = document.getElementsByClassName('side-nav-link-ref');
        var matchingMenuItem = null;

        for (var i = 0; i < links.length; i++) {
            if (this.props.location.pathname === links[i].pathname) {
                matchingMenuItem = links[i];
                break;
            }
        }

        if (matchingMenuItem) {
            matchingMenuItem.classList.add('active');
            var parent = matchingMenuItem.parentElement;

            /**
             * TODO: This is hard coded way of expading/activating parent menu dropdown and working till level 3. 
             * We should come up with non hard coded approach
             */
            if (parent) {
                parent.classList.add('active');
                const parent2 = parent.parentElement;
                if (parent2) {
                    parent2.classList.add('in');
                }
                const parent3 = parent2.parentElement;
                if (parent3) {
                    parent3.classList.add('active');
                    var childAnchor = parent3.querySelector('.has-dropdown');
                    if (childAnchor) childAnchor.classList.add('active');
                }

                const parent4 = parent3.parentElement;
                if (parent4)
                    parent4.classList.add('in');
                const parent5 = parent4.parentElement;
                if (parent5)
                    parent5.classList.add('active');
            }
        }
    }

    switchSidebar = (sidebarNo) => {
        this.setState({ sidebarNo })
    }

    render() {
        const isCondensed = this.props.isCondensed || false;
        const orgId = SessionStore.getOrganizationID();
        const version = this.state.version;
        let sidebar = this.state.sidebar;

        switch (this.state.sidebarNo) {
            case SUPERNODE_WALLET:
                sidebar = <SideNavSupernodeWalletContent orgId={orgId} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case SUPERNODE_SETTING:
                sidebar = <SideNavSupernodeSettingContent orgId={orgId} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case WALLET:
                sidebar = <SideNavWalletContent orgId={orgId} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case SETTING:
                sidebar = <SideNavSettingContent orgId={orgId} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            default:
                sidebar = <SideNavContent orgId={orgId} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
        }

        return (
            <React.Fragment>
                <div className='left-side-menu' ref={node => this.menuNodeRef = node}>
                    {!isCondensed && <PerfectScrollbar>{sidebar}</PerfectScrollbar>}
                    {isCondensed && <PerfectScrollbar>{sidebar}</PerfectScrollbar>}
                </div>
            </React.Fragment>
        );
    }
}

export default withRouter(Sidebar);
