import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import { Breadcrumb, BreadcrumbItem, Row, Col, Button } from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";

import AddWidget from './AddWidget';
import { adminWidgetCatalog } from './widgets/';

class AdminDashboard extends Component {
    constructor() {
        super();

        this.state = {
            data: {},
            loading: false,
            openAddWidget: false,
            widgets: []
        }

        this.openAddWidget = this.openAddWidget.bind(this);
        this.closeAddWidget = this.closeAddWidget.bind(this);
        this.onAddWidget = this.onAddWidget.bind(this);
        this.getData = this.getData.bind(this);
    }

    openAddWidget() {
        this.setState({ openAddWidget: true });
    }

    closeAddWidget() {
        this.setState({ openAddWidget: false });
    }

    onAddWidget(widget) {
        let widgets = [...this.state.widgets];
        widget.push(widget);
        this.setState({widgets: widgets, openAddWidget: false});
        this.getData();
    }

    componentDidMount() {
        this.getData();
    }

    getData() {
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
        });
    }


    render() {

        return (<React.Fragment>

            <TitleBar buttons={[
                <Button color="primary" onClick={this.openAddWidget}><i className="mdi mdi-plus"></i></Button>
            ]}>
                <Breadcrumb>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.dashboard.title`)}</BreadcrumbItem>
                </Breadcrumb>
            </TitleBar>

            {this.state.openAddWidget ? <AddWidget availableWidgets={adminWidgetCatalog} closeModal={this.closeAddWidget} /> : null}

            <Row>
                <Col>
                    <div className="position-relative">
                        {this.state.loading ? <Loader /> : null}

                        {/* <Row>
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
                        </Row> */}
                    </div>
                </Col>
            </Row>
        </React.Fragment>
        );
    }
}

export default withRouter(AdminDashboard);
