import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Doughnut, Bar, defaults as ChartJsDefaults } from "react-chartjs-2";
import { Breadcrumb, BreadcrumbItem, Row, Col, UncontrolledButtonDropdown, DropdownMenu, DropdownItem, DropdownToggle, Progress } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';

import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import StatWidget from "./StatWidget";
import MXCAmountChart from "./MXCAmountChart";
import StakingAmountChart from "./StakingAmountChart";
import EarnedAmountChart from "./EarnedAmountChart";
import DataPacketChart from "./DataPacketChart";
import DataMap from "./DataMap";


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
            <Link className="text-muted" to='#'>{i18n.t(`${packageNS}:menu.dashboard.tickets.view_history`)}</Link>
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.tickets.title`)}</h4>

        <div className="widget-chart mt-3">
            <Row>
                <Col lg={6}>
                    <Doughnut data={chartData} options={donutOpts} height={160} />
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

/**
 * Chart Actions
 */
const ChartActions = () => {
    return <UncontrolledButtonDropdown>
        <DropdownToggle className="arrow-none card-drop p-0" color="link"><i className="mdi mdi-dots-vertical"></i> </DropdownToggle>
        <DropdownMenu right>
            <DropdownItem>Week</DropdownItem>
            <DropdownItem>Month</DropdownItem>
        </DropdownMenu>
    </UncontrolledButtonDropdown>
}


/**
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
            <ChartActions />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.withdrawal.title`)}</h4>

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Bar data={chartData} options={barOpts} height={160} />
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


/**
 * Topup
 * @param {*} props 
 */
const Topup = (props) => {
    const topup = props.data || {};
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
    for (const v of (topup.data || [])) {
        labels.push(v.month);
        series.push(v.amount);
        hoverColors.push('#10c469');
        colors.push('rgba(16,196,105,0.5)');
    }

    const chartData = {
        labels: labels,
        datasets: [{
            label: i18n.t(`${packageNS}:menu.dashboard.topup.title`),
            data: series,
            backgroundColor: colors,
            hoverBackgroundColor: hoverColors,
            barPercentage: 0.65,
            categoryPercentage: 0.5,
        }]
    };


    return <div className="card-box">
        <div className="float-right">
            <ChartActions />
        </div>

        <h4 className="header-title mt-0">{i18n.t(`${packageNS}:menu.dashboard.topup.title`)}</h4>

        <div className="widget-chart mt-3">
            <Row>
                <Col className="mb-0">
                    <Bar data={chartData} options={barOpts} height={160} />
                </Col>
            </Row>
            <Row>
                <Col className="text-right mb-0">
                    <h2 className="mb-1">{topup.total ? topup.total / 1000 : 0}k MXC</h2>
                    <p className="mb-0">{i18n.t(`${packageNS}:menu.dashboard.topup.subtext`)}</p>
                </Col>
            </Row>
        </div>
    </div>;
}


class AdminDashboard extends Component {
    constructor() {
        super();


        this.state = {
            data: {},
            loading: false
        };
    }

    componentDidMount() {
        // TODO - call api to get the data
        this.setState({ loading: true });
        // mimiking the loading - should reverted later when we integrate api
        setTimeout(() => {
            this.setState({ loading: false });
        }, 1000);

        // Dummy data generator - remove this when you remove below sample data
        let packetsData = [];
        for (let idx = 30; idx > 0; idx--) {
            let day = new Date();
            day.setDate(day.getDate() - idx);
            packetsData.push({ "day": day.getDate(), "packets": Math.floor(Math.random() * 120) + 10 })
        }

        this.setState({
            data: {
                "tickets": {
                    "total": 80,
                    "approved": 34,
                    "pending": 46,
                },
                "withdrawal": {
                    "total": 50000,
                    "data": [
                        { "day": "M", "amount": 92000 },
                        { "day": "T", "amount": 220000 },
                        { "day": "W", "amount": 242000 },
                        { "day": "T", "amount": 34000 },
                        { "day": "F", "amount": 155000 },
                        { "day": "S", "amount": 172050 },
                        { "day": "S", "amount": 47500 },
                    ]
                },
                "topup": {
                    "total": 200000,
                    "data": [
                        { "month": "Jun", "amount": 92000 },
                        { "month": "Jul", "amount": 220000 },
                        { "month": "Aug", "amount": 242000 },
                        { "month": "Sep", "amount": 34000 },
                        { "month": "Oct", "amount": 155000 },
                        { "month": "Nov", "amount": 172050 },
                        { "month": "Dec", "amount": 48500 },
                    ]
                },
                "supernodeAmount": {
                    "total": 545000,
                    "data": [
                        { "day": "M", "amount": 205000 },
                        { "day": "T", "amount": 185000 },
                        { "day": "W", "amount": 220500 },
                        { "day": "T", "amount": 162050 },
                        { "day": "F", "amount": 187500 },
                        { "day": "S", "amount": 215000 },
                        { "day": "S", "amount": 145000 },
                    ]
                },
                "stakingAmount": {
                    "total": 845000,
                    "data": [
                        { "day": "M", "amount": 205000 },
                        { "day": "T", "amount": 185000 },
                        { "day": "W", "amount": 220500 },
                        { "day": "T", "amount": 162050 },
                        { "day": "F", "amount": 187500 },
                        { "day": "S", "amount": 215000 },
                        { "day": "S", "amount": 145000 },
                    ]
                },
                "earnedAmount": {
                    "total": 125,
                    "data": [
                        { "day": "M", "amount": 205, "amount2": 105 },
                        { "day": "T", "amount": 185, "amount2": 145 },
                        { "day": "W", "amount": 220, "amount2": 125 },
                        { "day": "T", "amount": 162, "amount2": 205 },
                        { "day": "F", "amount": 187, "amount2": 250 },
                        { "day": "S", "amount": 215, "amount2": 115 },
                        { "day": "S", "amount": 145, "amount2": 65 },
                    ]
                },
                "packetsSent": [...packetsData],
                "packetsReceived": [...packetsData],
                "packetsByChannel": [
                    { "channel": "864.7MHZ", "packets": 92000 },
                    { "channel": "864.9MHZ", "packets": 220000 },
                    { "channel": "866.4MHZ", "packets": 242000 },
                    { "channel": "867.2MHZ", "packets": 34000 },
                    { "channel": "869.8MHZ", "packets": 155000 },
                    { "channel": "870.1MHZ", "packets": 172050 },
                    { "channel": "872.2MHZ", "packets": 47500 },
                ],
                "packetsBySpreadFactor": [
                    { "spreadFactor": "7", "packets": 92000 },
                    { "spreadFactor": "8", "packets": 220000 },
                    { "spreadFactor": "9", "packets": 242000 },
                    { "spreadFactor": "10", "packets": 34000 },
                    { "spreadFactor": "11", "packets": 155000 },
                    { "spreadFactor": "12", "packets": 172050 }
                ],
                "totalUsers": 1230,
                "totalOrganizations": 45,
                "totalGateways": 90,
                "totalDevices": 260,
                "totalApplications": 260,
            }
        })
    }


    render() {

        return (<React.Fragment>

            <TitleBar buttons={[]}>
                <Breadcrumb>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.dashboard.title`)}</BreadcrumbItem>
                </Breadcrumb>
            </TitleBar>

            <Row>
                <Col>
                    <div className="position-relative">
                        {this.state.loading ? <Loader /> : null}

                        <Row>
                            <Col lg={4}>
                                <Tickets data={this.state.data.tickets} />
                            </Col>
                            <Col lg={4}>
                                <Withdrawal data={this.state.data.withdrawal} />
                            </Col>
                            <Col lg={4}>
                                <Topup data={this.state.data.topup} />
                            </Col>
                        </Row>

                        <Row>
                            <Col className="mb-0">
                                <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalUsers`)}
                                    value={this.state.data.totalUsers} formatNum={true}></StatWidget>
                            </Col>
                            <Col className="mb-0">
                                <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalOrgs`)}
                                    value={this.state.data.totalOrganizations} formatNum={true}></StatWidget>
                            </Col>
                            <Col className="mb-0">
                                <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalGateways`)}
                                    value={this.state.data.totalGateways} formatNum={true}></StatWidget>
                            </Col>
                            <Col className="mb-0">
                                <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalDevices`)}
                                    value={this.state.data.totalDevices} formatNum={true}></StatWidget>
                            </Col>
                            <Col className="mb-0">
                                <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalApps`)}
                                    value={this.state.data.totalApplications} formatNum={true}></StatWidget>
                            </Col>
                        </Row>

                        <Row>
                            <Col lg={4} className="mb-0">
                                <MXCAmountChart data={this.state.data.supernodeAmount} />
                            </Col>
                            <Col lg={4} className="mb-0">
                                <StakingAmountChart data={this.state.data.stakingAmount} />
                            </Col>
                            <Col lg={4} className="mb-0">
                                <EarnedAmountChart data={this.state.data.earnedAmount} />
                            </Col>
                        </Row>

                        <Row>
                            <Col className="mb-0">
                                <DataPacketChart data={this.state.data.packetsReceived} color="#10c469"
                                    title={i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.title`)}
                                    subTitle={i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.subTitle`)}
                                    subTitleClass="text-success" labelField="day" />
                            </Col>
                        </Row>

                        <Row>
                            <Col className="mb-0">
                                <DataPacketChart data={this.state.data.packetsSent} color="#71b6f9"
                                    title={i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.title`)}
                                    subTitle={i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.subTitle`)}
                                    subTitleClass="text-primary" labelField="day" />
                            </Col>
                        </Row>

                        <Row>
                            <Col className="mb-0">
                                <DataMap position={[51, 13]} />
                            </Col>
                        </Row>

                        <Row>
                            <Col lg={6} className="mb-0">
                                <DataPacketChart data={this.state.data.packetsByChannel} color="#71b6f9"
                                    title={i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.title`)}
                                    labelField="channel" showYAxis={true} />
                            </Col>
                            <Col lg={6} className="mb-0">
                                <DataPacketChart data={this.state.data.packetsBySpreadFactor} color="#71b6f9"
                                    title={i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.title`)}
                                    labelField="spreadFactor" showYAxis={true} />
                            </Col>
                        </Row>
                    </div>
                </Col>
            </Row>
        </React.Fragment>
        );
    }
}

export default withRouter(AdminDashboard);
