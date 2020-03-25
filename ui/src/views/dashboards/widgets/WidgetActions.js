import React from "react";
import { Link } from "react-router-dom";
import { DropdownItem, DropdownMenu, DropdownToggle, UncontrolledButtonDropdown } from "reactstrap";
import i18n, { packageNS } from '../../../i18n';


/**
 * Widget Actions
 */

const ActionTarget = ({ item }) => {
    if (item.to)
        return <Link to={item.to} className="dropdown-item">{item.label}</Link>
    else if (item.onClick)
        return <Link to='#'  className="dropdown-item" onClick={item.onClick}>{item.label}</Link>
    else
        return <div className="dropdown-item">{item.label}</div>
}


const WidgetActions = ({ widget, actionItems, onDelete }) => {
    const actions = actionItems || [];

    return <UncontrolledButtonDropdown>
        <DropdownToggle className="arrow-none card-drop p-0" color="link"><i className="mdi mdi-dots-vertical"></i> </DropdownToggle>
        <DropdownMenu right>
            {actions.map((item, i) => {
                return <ActionTarget item={item} key={i} />
            })}
            <DropdownItem divider />
                <Link to='#' className="dropdown-item text-danger" onClick={(e) => { if (onDelete) onDelete(widget); }}>{i18n.t(`${packageNS}:menu.dashboard.remove`)}</Link>
        </DropdownMenu>
    </UncontrolledButtonDropdown>
}

export default WidgetActions;