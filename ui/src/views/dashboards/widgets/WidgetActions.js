import React from "react";
import { Link } from "react-router-dom";

import i18n, { packageNS } from '../../../i18n';

/**
 * Widget Actions
 */
const WidgetActions = ({ widget, onDelete }) => {
    return <UncontrolledButtonDropdown>
        <DropdownToggle className="arrow-none card-drop p-0" color="link"><i className="mdi mdi-dots-vertical"></i> </DropdownToggle>
        <DropdownMenu right>
            <DropdownItem>Week</DropdownItem>
            <DropdownItem>Month</DropdownItem>
            <DropdownItem className="">
                <Link to='#' className="text-warning" onClick={(e) => { onDelete ? onDelete(widget) : '' }}>{i18n.t(`${packageNS}:menu.dashboard.remove`)}</Link>
            </DropdownItem>
        </DropdownMenu>
    </UncontrolledButtonDropdown>
}

export default WidgetActions;