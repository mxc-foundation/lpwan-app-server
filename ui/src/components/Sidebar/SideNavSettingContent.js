import React from 'react';
import { Link } from 'react-router-dom';

import i18n, { packageNS } from '../../i18n';
import { DEFAULT } from '../../util/Data';
import DropdownMenu2 from '../DropdownMenu';

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

export default SideNavSettingContent;
