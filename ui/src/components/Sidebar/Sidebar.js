import React, { Component } from 'react';
import MetisMenu from 'metismenujs/dist/metismenujs';
import PerfectScrollbar from 'react-perfect-scrollbar';
import 'react-perfect-scrollbar/dist/css/styles.css';
import { withRouter } from 'react-router-dom';
import OrganizationStore from '../../stores/OrganizationStore';
import ServerInfoStore from '../../stores/ServerInfoStore';
import SessionStore from '../../stores/SessionStore';
import { DEFAULT, ORGANIZATIONS, SETTING, SUPERNODE_SETTING, SUPERNODE_WALLET, WALLET } from '../../util/Data';
import SideNavContent from './SideNavContent';
import SideNavOrganizationsContent from './SideNavOrganizationsContent';
import SideNavSettingContent from './SideNavSettingContent';
import SideNavSupernodeSettingContent from './SideNavSupernodeSettingContent';
import SideNavSupernodeWalletContent from './SideNavSupernodeWalletContent';
import SideNavWalletContent from './SideNavWalletContent';




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
}] */

class Sidebar extends Component {
    constructor(props) {
        super(props);

        this.state = {
            currentSidebarId: props.currentSidebarId || DEFAULT,
            open: true,
            //organization: {},
            organization: null,
            organizationID: '',
            cacheCounter: 0,
            version: ''
        };

        this.handleOtherClick = this.handleOtherClick.bind(this);
        this.initMenu = this.initMenu.bind(this);
    }
    loadData = async () => {
        try {
          const organizationIDs = SessionStore.getOrganizations();
          var data = await ServerInfoStore.getAppserverVersion();
          const serverInfo = JSON.parse(data);
          
          this.setState({
            organizationID: organizationIDs[0].organizationID,
            version: serverInfo.version
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
        SessionStore.on("organization.change", async () => {
            let organization = await OrganizationStore.get(SessionStore.getOrganizationID());
            this.setState({ organization: organization.organization });
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
    
      getOrganizationOption = async (id, callbackFunc) =>  {
        let resp = await OrganizationStore.get(id);
        callbackFunc({label: resp.organization.name, value: resp.organization.id, color:"black"});
      }
    
      getOrganizationOptions = async (search, callbackFunc) => {
        const resp = await OrganizationStore.list(search, 10, 0);
        const options = resp.result.map((o, i) => {return {label: o.name, value: o.id, color:'black'}});
        callbackFunc(options);
      }
    
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

        // on route change - switch to default
        const innerExcludePaths = ['/history/', '/topup/', '/withdraw/', '/control-panel/history', '/control-panel/withdraw', '/control-panel/system-settings'];
        
        if (this.props.location !== prevProps.location && ([SUPERNODE_WALLET, WALLET, SUPERNODE_SETTING, SETTING].indexOf(prevProps.currentSidebarId) !== -1)) {
          if (this.props.location) {
            const match = innerExcludePaths.findIndex(p => this.props.location.pathname.startsWith(p));
            if (match === -1)
              this.switchSidebar(DEFAULT);  
          }
          
        }


        if (this.props === prevProps) {
            return;
        }

        if (this.props.currentSidebarId !== this.state.currentSidebarId) {
          this.setState({
            currentSidebarId: this.props.currentSidebarId
          })
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

    switchSidebar = (newSidebarId) => {
        const { switchToSidebarId } = this.props;
        switchToSidebarId(newSidebarId);
    }

    render() {
        const { version } = this.state;
        const { currentSidebarId } = this.props;

        const isCondensed = this.props.isCondensed || false;
        const user = SessionStore.getUser();

        let sidebarComponent;
        switch (currentSidebarId) {
            case SUPERNODE_WALLET:
                sidebarComponent = <SideNavSupernodeWalletContent orgId={this.state.organizationID} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case SUPERNODE_SETTING:
                sidebarComponent = <SideNavSupernodeSettingContent orgId={this.state.organizationID} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case ORGANIZATIONS:
                sidebarComponent = <SideNavOrganizationsContent orgId={this.state.organizationID} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case WALLET:
                sidebarComponent = <SideNavWalletContent orgId={this.state.organizationID} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            case SETTING:
                sidebarComponent = <SideNavSettingContent orgId={this.state.organizationID} user={user} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
            default:
                sidebarComponent = <SideNavContent orgId={this.state.organizationID} version={version} onChange={this.onChange} switchSidebar={this.switchSidebar} />;
                break;
        }

        return (
            <React.Fragment>
                <div className='left-side-menu' ref={node => this.menuNodeRef = node}>
                    {!isCondensed && <PerfectScrollbar>{sidebarComponent}</PerfectScrollbar>}
                    {isCondensed && <PerfectScrollbar>{sidebarComponent}</PerfectScrollbar>}
                </div>
            </React.Fragment>
        );
    }
}

export default withRouter(Sidebar);
