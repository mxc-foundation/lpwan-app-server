import React from "react";
import { defaults as ChartJsDefaults, Doughnut } from "react-chartjs-2";
import { Col, Progress, Row } from "reactstrap";
import i18n, { packageNS } from '../../../i18n';
import WidgetActions from './WidgetActions';


// default
ChartJsDefaults.global.defaultFontColor = 'rgba(0, 0, 0, 0.65)';
ChartJsDefaults.global.defaultFontSize = 12;
ChartJsDefaults.global.defaultFontFamily = 'Karla, Microsoft YaHei';


/**
 * Tickets
 * @param {*} param0 
 */
const Tickets = (props) => {
    const data = props.data || {};
    const donutOpts = {
        maintainAspectRatio: false,
        cutoutPercentage: 80,
        legend: {
            display: false
        }
    };

    const chartData = {
        labels: [i18n.t(`${packageNS}:menu.dashboard.tickets.approved`), i18n.t(`${packageNS}:menu.dashboard.tickets.pending`)],
        datasets: [{
            data: [data.approved, data.pending],
            backgroundColor: ['#10c469', '#ff5b5b'],
            hoverBackgroundColor: ['#10c469', '#ff5b5b']
        }]
    };

    const approvedPer = (data && data.approved ? (data.approved / data.total) * 100 : 0).toFixed(2);
    const pendingPer = (data && data.pending ? (data.pending / data.total) * 100 : 0).toFixed(2);


    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: i18n.t(`${packageNS}:menu.dashboard.tickets.view_history`) }]}
                onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.tickets.title`)}</h4>
        <p>&nbsp;</p>

        <div className="widget-chart mt-3">
            <Row>
                <Col lg={6}>
                    <Doughnut data={chartData} options={donutOpts} height={200} />
                </Col>
                <Col lg={6} className="">
                    <div className="pl-2">
                        <label className="mb-1">{i18n.t(`${packageNS}:menu.dashboard.tickets.approved`)}</label>
                        <Row className="align-items-center no-gutters">
                            <Col lg={7}>
                                <Progress value={approvedPer} color="success" className="mt-0" />
                            </Col>
                            <Col lg={2}><span className="pl-2">{approvedPer}%</span></Col>
                        </Row>
                        <hr />
                        <label className="mb-1">{i18n.t(`${packageNS}:menu.dashboard.tickets.pending`)}</label>
                        <Row className="align-items-center no-gutters">
                            <Col lg={7}>
                                <Progress value={pendingPer} color="danger" className="mt-0" />
                            </Col>
                            <Col lg={2}><span className="pl-2">{pendingPer}%</span></Col>
                        </Row>
                    </div>
                </Col>
            </Row>
            <Row>
                <Col className="text-right mb-0">
                    <h2 className="mb-1">{data.total}</h2>
                    <p className="mb-0">{i18n.t(`${packageNS}:menu.dashboard.tickets.subtext`)}</p>
                </Col>
            </Row>
        </div>
    </div>;
}

export default Tickets;