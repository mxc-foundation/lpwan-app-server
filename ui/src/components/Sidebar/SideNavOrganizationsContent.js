import React from 'react';
import { Link } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import { DEFAULT, WALLET, SETTING, BACK_TO_CONTROL } from '../../util/Data';
import DropdownMenu2 from '../DropdownMenu';

const SideNavOrganizationsContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <li>
                    <Link to={BACK_TO_CONTROL} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                        <span className="mdi mdi-arrow-left-bold"></span>
                        <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
                        <span>{i18n.t(`${packageNS}:tr000463`)}</span>
                    </Link>
                </li>

                <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>
                <li>
                    {/* <DropdownMenu default={ this.state.default } onChange={this.onChange} /> [edit] */}
                    <DropdownMenu2 default={props.default} onChange={props.onChange} />
                </li>

                {/* <li>
                    <Link to={`/organizations`} className="waves-effect side-nav-link-ref" >
                        <i className="mdi mdi-domain"></i>
                        <span> {i18n.t(`${packageNS}:tr000049`)} </span>
                    </Link>
                </li> */}

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

                <li>
                    <Link to={`/organizations/${props.orgId}/service-profiles`}>
                        <i className="mdi mdi-folder-account"></i>
                        <span> {i18n.t(`${packageNS}:tr000078`)} </span>
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
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

export default SideNavOrganizationsContent;
