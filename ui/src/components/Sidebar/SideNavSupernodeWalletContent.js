import React from 'react';
import { Link } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import { DEFAULT, BACK_TO_CONTROL } from '../../util/Data';
import Admin from '../Admin';

const SideNavSupernodeWalletContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <Admin>
                    <li className="menu-title">{i18n.t(`${packageNS}:menu.control_panel`)}</li>

                    <li>
                        <Link to={BACK_TO_CONTROL} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                            <span className="mdi mdi-arrow-left-bold"></span>
                            <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
                            <span>{i18n.t(`${packageNS}:tr000450`)}</span>
                        </Link>
                    </li>

                    <li>
                        <Link to={`/control-panel/withdraw/`} className="waves-effect side-nav-link-ref">
                            <i className="mdi mdi-hammer"></i>
                            <span> {i18n.t(`${packageNS}:menu.dashboard.coming_soon`)} </span>
                        </Link>
                    </li>
                    
                </Admin>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

export default SideNavSupernodeWalletContent;
