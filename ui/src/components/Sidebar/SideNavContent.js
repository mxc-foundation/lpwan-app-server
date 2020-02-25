import React from 'react';
import { Link } from 'react-router-dom';
import Divider from '@material-ui/core/Divider';

import i18n, { packageNS } from '../../i18n';
import mxcLogo from '../../assets/images/mxc_logo-social_2.png';
import { SUPERNODE_WALLET, SUPERNODE_SETTING, WALLET, SETTING, ORGANIZATIONS } from '../../util/Data';
import Admin from '../Admin';
import NonAdmin from '../NonAdmin';
import DropdownMenu2 from '../DropdownMenu';

const SideNavContent = (props) => {
    return <React.Fragment>

        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                {/* <li>
                    <Link to="/dashboard" className="waves-effect side-nav-link-ref">
                        <i className="mdi mdi-view-dashboard"></i>
                        <span> {i18n.t(`${packageNS}:menu.dashboard.title`)} </span>
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
                        <Link to={`/organizations`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(ORGANIZATIONS)}>
                            <i className="mdi mdi-domain"></i>
                            <span> {i18n.t(`${packageNS}:tr000049`)} </span>
                            <span className="menu-arrow"></span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/users`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-account-multiple"></i>
                            <span> {i18n.t(`${packageNS}:tr000055`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/service-profiles`}>
                            <i className="mdi mdi-folder-account"></i>
                            <span> {i18n.t(`${packageNS}:tr000078`)} </span>
                        </Link>
                    </li>

                    {/*<li>
                        <Link to={`/organizations/${props.orgId}/edit`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-domain"></i>
                            <span> {i18n.t(`${packageNS}:tr000418`)} </span>
                        </Link>
                    </li>*/}

                    <li>
                        <Link to="/control-panel/modify-account/" className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(SUPERNODE_SETTING)}>
                            <i className="mdi mdi-settings"></i>
                            <span> {i18n.t(`${packageNS}:tr000451`)} </span>
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

                    <li className="menu-title">{i18n.t(`${packageNS}:menu.lpwan_management`)}</li>

                    <li>
                        <Link to={`/topup/${props.orgId}`} className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(WALLET)}>
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
                        <Link to={`/stake/${props.orgId}/set-stake`} className="waves-effect side-nav-link-ref">
                            {/* <i className="mdi mdi-cloud-print-outline"></i> */}
                            <i className="mdi mdi-vote"></i>
                            <span> {i18n.t(`${packageNS}:menu.staking.staking`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/gateway-profiles`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-remote"></i>
                            <span> {i18n.t(`${packageNS}:tr000046`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/gateways`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-remote"></i>
                            <span> {i18n.t(`${packageNS}:menu.gateways.gateways`)} </span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/organizations/${props.orgId}/device-profiles`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-memory"></i>
                            <span> {i18n.t(`${packageNS}:tr000278`)} </span>
                        </Link>
                    </li>

                    {/* <li>
                        <Link to={`/organizations/${props.orgId}/multicast-groups`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-podcast"></i>
                            <span> {i18n.t(`${packageNS}:tr000083`)} </span>
                        </Link>
                    </li> */}

                    <li>
                        <Link to={`/organizations/${props.orgId}`} className="waves-effect" aria-expanded="false" onClick={() => props.switchSidebar(SETTING)}>
                            <i className="mdi mdi-settings"></i>
                            <span> {i18n.t(`${packageNS}:tr000451`)} </span>
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


export default SideNavContent;
