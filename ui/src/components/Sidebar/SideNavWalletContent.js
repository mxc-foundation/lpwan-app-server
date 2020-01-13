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
                    <Link to={`/stake/${props.orgId}`} className="waves-effect side-nav-link-ref" onClick={() => props.switchSidebar(DEFAULT)}>
                        <span className="mdi mdi-arrow-left-bold"></span>
                        <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>
                        <span>{i18n.t(`${packageNS}:tr000463`)}</span>
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

export default SideNavWalletContent;
