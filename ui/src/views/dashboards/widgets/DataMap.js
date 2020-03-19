import React from "react";
import { Row, Col } from "reactstrap";
import { Map, Marker } from 'react-leaflet';
import FoundLocationMap from "../../../components/FoundLocationMap"

import i18n, { packageNS } from '../../../i18n';
import WidgetActions from "./WidgetActions";


const DataMap = (props) => {
    const position = props.data;
    const style = {
        height: 360,
        zIndex: 1
    };

    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.dataMap.title`)}</h4>

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Map center={position} zoom={15} style={style} animate={true} scrollWheelZoom={false}>
                        <FoundLocationMap />
                        <Marker position={position} />
                    </Map>
                </Col>
            </Row>
        </div>
    </div>;
}

export default DataMap;
