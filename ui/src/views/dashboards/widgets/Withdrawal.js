import React from "react";
import { Bar, defaults as ChartJsDefaults } from "react-chartjs-2";
import { Row, Col } from "reactstrap";

import i18n, { packageNS } from '../../../i18n';
import WidgetActions from './WidgetActions';

// default
ChartJsDefaults.global.defaultFontColor = 'rgba(0, 0, 0, 0.65)';
ChartJsDefaults.global.defaultFontSize = 12;
ChartJsDefaults.global.defaultFontFamily = 'Karla, Microsoft YaHei';

/**
 * 
 * Withdrawal
 * @param {*} props 
 */
const Withdrawal = (props) => {
    const withdrawal = props.data || {};
    const barOpts = {
        maintainAspectRatio: false,
        legend: {
            display: false
        },
        tooltips: {
            callbacks: {
                label: function (tooltipItems, data) {
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
                    }
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
    let colors = [];
    let hoverColors = [];
    for (const v of (withdrawal.data || [])) {
        labels.push(v.day);
        series.push(v.amount);
        hoverColors.push('#ff5b5b');
        colors.push('rgba(255,91,91,0.65)');
    }

    const chartData = {
        labels: labels,
        datasets: [{
            label: i18n.t(`${packageNS}:menu.dashboard.withdrawal.title`),
            data: series,
            backgroundColor: colors,
            hoverBackgroundColor: hoverColors,
            barPercentage: 0.65,
            categoryPercentage: 0.5,
        }]
    };


    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.withdrawal.title`)}</h4>
        <p>&nbsp;</p>

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Bar data={chartData} options={barOpts} height={200} />
                </Col>
            </Row>
            <Row>
                <Col className="text-right mb-0">
                    <h2 className="mb-1">{withdrawal.total ? withdrawal.total / 1000 : 0}k MXC</h2>
                    <p className="mb-0">{i18n.t(`${packageNS}:menu.dashboard.withdrawal.subtext`)}</p>
                </Col>
            </Row>
        </div>
    </div>;
}

export default Withdrawal;