import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import {
    Breadcrumb, BreadcrumbItem, Row, Col, UncontrolledButtonDropdown, DropdownMenu, DropdownItem, DropdownToggle,
    Media, Badge
} from 'reactstrap';

import i18n, { packageNS } from '../../i18n';
import SessionStore from "../../stores/SessionStore";
import TitleBar from "../../components/TitleBar";
import Loader from "../../components/Loader";
import defaultProfilePic from "../../assets/images/users/profile-icon.png";

import StatWidget from "./StatWidget";
import MXCAmountChart from "./MXCAmountChart";
import StakingAmountChart from "./StakingAmountChart";
import EarnedAmountChart from "./EarnedAmountChart";
import DataPacketChart from "./DataPacketChart";
import DataMap from "./DataMap";


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
 * Topup
 * @param {*} props 
 */
const Topup = (props) => {
    const data = props.data || {};
    const { user: { id, profilePic, username } } = props;
    const formattedVal = (data.amount || 0).toLocaleString(navigator.language, { minimumFractionDigits: 4 });

    const userRole = SessionStore.isOrganizationAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleOrgAdmin`) :
        SessionStore.isOrganizationDeviceAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleDeviceAdmin`) :
            SessionStore.isOrganizationGatewayAdmin() ? i18n.t(`${packageNS}:menu.dashboard.roleGatewayAdmin`) : "";

    const orgId = SessionStore.getOrganizationID();

    return <div className="card-box">
        <div className="float-right">
            <ChartActions />
        </div>

        <Media className="align-items-center">
            <Media left className="avatar-xl">
                <img src={profilePic || defaultProfilePic} className="img-fluid rounded-circle" alt="user" />
            </Media>
            <Media body>
                <div className="ml-2">
                    <h3 className="font-weight-normal mt-0">{username}</h3>
                    <h5 className="text-primary mb-0">{userRole}</h5>
                </div>
            </Media>
        </Media>

        <Row className="mt-3">
            <Col className="text-right mb-0">
                {data.growth ? <h5><Badge className="px-1">
                    {data.growth} <i className="mdi mdi-arrow-up"></i></Badge></h5> : null}
                <h2 className="my-2 font-2rem">{formattedVal} MXC</h2>
                <Link className="btn btn-primary" to={`/topup/${orgId}`}>{i18n.t(`${packageNS}:menu.dashboard.topupButton`)}</Link>
            </Col>
        </Row>
    </div>;
}


class UserDashboard extends Component {
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
                "topup": {
                    "amount": 1235.09,
                    "growth": "15%"
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
                            <Col lg={6}>
                                <Topup data={this.state.data.topup} user={this.props.user} />
                            </Col>
                            <Col lg={6}>
                                <Row>
                                    <Col className="mb-0" lg={6}>
                                        <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalUsers`)}
                                            value={this.state.data.totalUsers} formatNum={true}></StatWidget>
                                    </Col>
                                    <Col className="mb-0" lg={6}>
                                        <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalGateways`)}
                                            value={this.state.data.totalGateways} formatNum={true}></StatWidget>
                                    </Col>
                                    <Col className="mb-0" lg={6}>
                                        <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalDevices`)}
                                            value={this.state.data.totalDevices} formatNum={true}></StatWidget>
                                    </Col>
                                    <Col className="mb-0" lg={6}>
                                        <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalApps`)}
                                            value={this.state.data.totalApplications} formatNum={true}></StatWidget>
                                    </Col>
                                </Row>
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

export default withRouter(UserDashboard);
