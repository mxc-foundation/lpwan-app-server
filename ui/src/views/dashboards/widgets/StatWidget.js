import React from "react";
import { DropdownItem, DropdownMenu, DropdownToggle, UncontrolledButtonDropdown } from 'reactstrap';

const StatWidget = ({ label, value, formatNum }) => {

    const formattedVal = formatNum ? (value || 0).toLocaleString(navigator.language, { minimumFractionDigits: 0 }) : value;

    return <div className="card-box">
        <div className="float-right">
            <UncontrolledButtonDropdown>
                <DropdownToggle className="arrow-none card-drop p-0" color="link"><i className="mdi mdi-dots-vertical"></i> </DropdownToggle>
                <DropdownMenu right>
                    <DropdownItem>Week</DropdownItem>
                    <DropdownItem>Month</DropdownItem>
                </DropdownMenu>
            </UncontrolledButtonDropdown>
        </div>

        <h4 className="header-title mt-0 mb-3">{label}</h4>

        <div className="text-right">
            <h2 className="text-primary pt-2 mb-0">{formattedVal}</h2>
        </div>
    </div>
}

export default StatWidget;