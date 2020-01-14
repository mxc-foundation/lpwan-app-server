import React from 'react';
import { Link } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import { DEFAULT } from '../../util/Data';
import Admin from '../Admin';

const SideNavSupernodeSettingContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <Admin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.control_panel`)}</li>

                    <li>
                        <Link to={`/organizations`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
                            <span>{i18n.t(`${packageNS}:tr000450`)}</span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/modify-account`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-ethereum"></i>
                            <span> {i18n.t(`${packageNS}:menu.eth_account.eth_account`)} </span>
                        </Link>
                    </li>

{/*                    <li>
                        <Link to={`/control-panel/system-settings`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-settings"></i>
                            <span> {i18n.t(`${packageNS}:tr000417`)} </span>
                        </Link>
                    </li>*/}
                </Admin>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

export default SideNavSupernodeSettingContent;
