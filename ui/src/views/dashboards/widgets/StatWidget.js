import React from "react";

import WidgetActions from './WidgetActions';


const StatWidget = (props) => {
    const { label, data, formatNum } = props;
    const formattedVal = formatNum ? (data || 0).toLocaleString(navigator.language, { minimumFractionDigits: 0 }) : data;

    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0 mb-3">{label}</h4>

        <div className="text-right">
            <h2 className="text-primary pt-2 mb-0">{formattedVal}</h2>
        </div>
    </div>
}

export default StatWidget;