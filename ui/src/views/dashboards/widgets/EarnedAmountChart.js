import React from "react";
import { defaults as ChartJsDefaults, Line } from "react-chartjs-2";
import { Col, Row } from "reactstrap";
import i18n, { packageNS } from '../../../i18n';
import WidgetActions from './WidgetActions';



// default
ChartJsDefaults.global.defaultFontColor = 'rgba(0, 0, 0, 0.65)';
ChartJsDefaults.global.defaultFontSize = 12;
ChartJsDefaults.global.defaultFontFamily = 'Karla, Microsoft YaHei';


const EarnedAmountChart = (props) => {
    const data = props.data || {};

    const lineOpts = {
        maintainAspectRatio: false,
        legend: {
            display: false
        },
        scales: {
            yAxes: [{
                gridLines: {
                    color: "#ebeff2"
                },
                stacked: false,
                ticks: {
                    min: 0
                },
            }],
            xAxes: [{
                stacked: false,
                gridLines: {
                    display: false,
                    zeroLineColor: '#ebeff2'
                },
                zeroLineColor: '#ebeff2'
            }]
        }
    };

    let labels = [];
    let series = [];
    let series2 = [];
    for (const v of (data.data || [])) {
        labels.push(v.day);
        series.push(v.amount);
        series2.push(v.amount2)
    }

    const chartData = {
        labels: labels,
        datasets: [{
            label: 'Amount',
            data: series,
            backgroundColor: "transparent",
            borderColor: "#f9c851",
            borderWidth: 3,
            pointBorderWidth: 2,
            pointBackgroundColor: "#ffffff",
            pointHoverBackgroundColor: "#ffffff",
            pointHoverBorderColor: "#f9c851",
        },
        {
            label: "Amount 2",
            data: series2,
            backgroundColor: "transparent",
            borderColor: "#5b69bc",
            borderWidth: 3,
            pointBorderWidth: 2,
            pointBackgroundColor: "#ffffff",
            pointHoverBackgroundColor: "#ffffff",
            pointHoverBorderColor: "#5b69bc",
        }]
    };


    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.title`)}</h4>
        <p>&nbsp;</p>
        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Line data={chartData} options={lineOpts} height={200} />
                </Col>
            </Row>
            <Row>
                <Col className="text-right mb-0">
                    <h2 className="mb-1">{data.total ? data.total / 1000 : 0}k MXC</h2>
                    <p className="mb-0">{i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.statSubText`)}</p>
                </Col>
            </Row>
        </div>
    </div>;
}

export default EarnedAmountChart;
