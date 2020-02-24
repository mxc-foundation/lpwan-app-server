import React from "react";
import { Line, defaults as ChartJsDefaults } from "react-chartjs-2";
import { Row, Col } from "reactstrap";

import i18n, { packageNS } from '../../i18n';


// default
ChartJsDefaults.global.defaultFontColor = 'rgba(0, 0, 0, 0.65)';
ChartJsDefaults.global.defaultFontSize = 12;
ChartJsDefaults.global.defaultFontFamily = 'Karla, Microsoft YaHei';


const MXCAmountChart = (props) => {
    const data = props.data || {};

    const lineOpts = {
        maintainAspectRatio: false,
        legend: {
            display: false
        },
        tooltips: {
            callbacks: {
                label: function (tooltipItems) {
                    return tooltipItems.yLabel / 1000 + 'k';
                }
            }
        },
        scales: {
            yAxes: [{
                gridLines: {
                    color: "#ebeff2"
                },
                stacked: false,
                ticks: {
                    callback: function (label, index, labels) {
                        return label / 1000 + 'k';
                    },
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
    for (const v of (data.data || [])) {
        labels.push(v.day);
        series.push(v.amount);
    }

    const chartData = {
        labels: labels,
        datasets: [{
            label: i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.title`),
            data: series,
            backgroundColor: "rgba(249,200,81,0.1)",
            borderColor: "#f9c851",
            borderWidth: 3,
            pointBorderWidth: 2,
            pointBackgroundColor: "#ffffff",
            pointHoverBackgroundColor: "#ffffff",
            pointHoverBorderColor: "#f9c851",
        }]
    };


    return <div className="card-box">
        <div className="float-right"></div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.title`)}</h4>
        <p className="mt-0 text-warning">{i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.subtext`)}</p>

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Line data={chartData} options={lineOpts} height={200} />
                </Col>
            </Row>
            <Row>
                <Col className="text-right mb-0">
                    <h2 className="mb-1">{data.total ? data.total / 1000 : 0}k MXC</h2>
                    <p className="mb-0">{i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.statSubText`)}</p>
                </Col>
            </Row>
        </div>
    </div>;
}

export default MXCAmountChart;
