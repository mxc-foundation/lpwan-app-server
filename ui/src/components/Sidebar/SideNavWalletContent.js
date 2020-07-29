import React from 'react';
import { Link } from 'react-router-dom';
import i18n, { packageNS } from '../../i18n';
import { DEFAULT } from '../../util/Data';
import DropdownMenu2 from '../DropdownMenu';


const SideNavWalletContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">
                <li>
                    <Link to={`/stake/${props.orgId}/set-stake`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                        <span className="mdi mdi-arrow-left-bold"></span>
                        <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
                        <span>{i18n.t(`${packageNS}:tr000463`)}</span>
                    </Link>
                </li>

                <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>

                <li>
                    <DropdownMenu2 default={props.default} onChange={props.onChange} />
                </li>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

export default SideNavWalletContent;
