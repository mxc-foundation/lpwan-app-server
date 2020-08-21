import React from 'react';
import i18n, {packageNS} from '../../i18n';
import DropdownMenu2 from '../DropdownMenu';


const SideNavWalletContent = (props) => {
    return <React.Fragment>
        <div id="sidebar-menu">
            <ul className="metismenu" id="side-menu">

                <li className="menu-title">{i18n.t(`${packageNS}:menu.organization_list`)}</li>

                <li>
                    <DropdownMenu2 default={props.default} onChange={props.onChange}/>
                </li>
            </ul>
        </div>
        <div className="clearfix"></div>
    </React.Fragment>
}

export default SideNavWalletContent;
