import React from "react";
import { Bar, defaults as ChartJsDefaults } from "react-chartjs-2";
import { Row, Col } from "reactstrap";
import classNames from "classnames";

import WidgetActions from './WidgetActions';

// default
ChartJsDefaults.global.defaultFontColor = 'rgba(0, 0, 0, 0.65)';
ChartJsDefaults.global.defaultFontSize = 12;
ChartJsDefaults.global.defaultFontFamily = 'Karla, Microsoft YaHei';


const DataPacketChart = (props) => {
    const data = props.data || [];
    const showYAxis = props.showYAxis || false;

    const barOpts = {
        maintainAspectRatio: false,
        legend: {
            display: false
        },
        scales: {
            yAxes: [{
                display: showYAxis,
                gridLines: {
                    color: "#ebeff2"
                },
                stacked: false,
            }],
            xAxes: [{
                stacked: false,
                gridLines: {
                    display: false,
                    zeroLineColor: '#ebeff2'
                },
                zeroLineColor: '#ebeff2',
            }]
        }
    };


    let labels = [];
    let series = [];
    let colors = [];
    for (const v of data) {
        labels.push(v[props.labelField]);
        series.push(v.packets);
        colors.push(props.color);
    }

    const chartData = {
        labels: labels,
        datasets: [{
            label: props.title || "Packets",
            data: series,
            backgroundColor: colors,
            hoverBackgroundColor: colors,
            barPercentage: 0.65,
            categoryPercentage: 0.5,
        }]
    };

    return <div className="card-box">
        <div className="float-right">
            <WidgetActions widget={props.widget} actionItems={[{ to: '#', label: 'Week' }]} onDelete={props.onDelete} />
        </div>

        <h4 className="header-title mt-0">{props.title}</h4>
        {props.subTitle ? <p className={classNames("mt-0", props.subTitleClass)}>{props.subTitle}</p> : null}

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Bar data={chartData} options={barOpts} height={200} />
                </Col>
            </Row>
        </div>
    </div>;
}

export default DataPacketChart;
