import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Button, Col, Row } from 'reactstrap';
import Loader from "../../components/Loader";
import TitleBar from "../../components/TitleBar";
import i18n, { packageNS } from '../../i18n';
import AddWidget from './AddWidget';
import { userWidgetCatalog, WIDGET_TYPE_GRAPH, WIDGET_TYPE_MAP, WIDGET_TYPE_STAT } from './widgets/';
import WalletStore from "../../stores/WalletStore";
import SessionStore from "../../stores/SessionStore";



class UserDashboard extends Component {
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
        this.onDeletewidget = this.onDeletewidget.bind(this);
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
        widgets.push(widget);
        this.setState({ widgets: widgets, openAddWidget: false });
        this.getData();
    }

    onDeletewidget(widget) {
        let widgets = [...this.state.widgets];
        widgets = widgets.filter(w => w.name !== widget.name);
        this.setState({ widgets: widgets });
    }

    componentDidMount() {
        this.getData();

        // showing dummy widgets on load - remove this when API is available
        let widgets = [...userWidgetCatalog];
        this.setState({ widgets });
    }

    getData = async () => {
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

        const user = await SessionStore.getUser();
        const orgId = await SessionStore.getOrganizationID();
        const topup = await  WalletStore.getMiningInfo(user.id, orgId);

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
        });
    }

    /**
     * Gets the widgets by type
     * @param {*} type 
     * @param {*} startIdx 
     * @param {*} size 
     */
    getWidgets(type, startIdx, size) {
        let typeWiseWidgets = [];
        for (const widget of this.state.widgets) {
            if (widget['type'] === type)
                typeWiseWidgets.push({ meta: widget, component: widget.component, data: this.state.data[widget.name] });
        }
        return typeWiseWidgets.slice(startIdx, startIdx + size) || [];
    }

    render() {
        const statWidgets = this.getWidgets(WIDGET_TYPE_STAT, 0, 4) || [];

        return (<React.Fragment>

            <TitleBar buttons={[
                <Button color="primary" onClick={this.openAddWidget}><i className="mdi mdi-plus"></i></Button>
            ]}>
                <Breadcrumb>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.dashboard.title`)}</BreadcrumbItem>
                </Breadcrumb>
            </TitleBar>

            {this.state.openAddWidget ? <AddWidget availableWidgets={userWidgetCatalog} closeModal={this.closeAddWidget}
                addWidget={this.onAddWidget} addedWidgets={this.state.widgets} /> : null}

            <Row>
                <Col>
                    <div className="position-relative">
                        {this.state.loading ? <Loader /> : null}

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_GRAPH, 0, 1).map((widget, idx) => {
                                if (idx < 2) {
                                    return <Col key={idx} className="mb-0">
                                        <div className="position-relative">
                                            <div className="card-coming-soon-2"></div>
                                            <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                        </div>
                                    </Col>
                                } else {
                                    return <Col key={idx} className="mb-0">
                                        <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                    </Col>

                                }
                            })}
                            <Col>
                                <Row>
                                    {statWidgets.map((widget, idx) => {
                                        return <Col lg={statWidgets.length > 2 ? 6 : 12} key={idx} className="mb-0">
                                            <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                        </Col>
                                    })}
                                </Row>
                            </Col>
                        </Row>

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_GRAPH, 1, 3).map((widget, idx) => {
                                return <Col key={idx} className="mb-0">
                                    <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                </Col>
                            })}
                        </Row>

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_GRAPH, 4, 1).map((widget, idx) => {
                                return <Col key={idx} className="mb-0">
                                    <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                </Col>
                            })}
                        </Row>

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_GRAPH, 5, 1).map((widget, idx) => {
                                return <Col key={idx} className="mb-0">
                                    <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                </Col>
                            })}
                        </Row>

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_MAP, 0, 1).map((widget, idx) => {
                                return <Col key={idx} className="mb-0">
                                    <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                </Col>
                            })}
                        </Row>

                        <Row>
                            {this.getWidgets(WIDGET_TYPE_GRAPH, 6, 2).map((widget, idx) => {
                                return <Col key={idx} className="mb-0">
                                    <widget.component data={widget.data} widget={widget.meta} onDelete={this.onDeletewidget} />
                                </Col>
                            })}
                        </Row>
                    </div>
                </Col>
            </Row>
        </React.Fragment>
        );
    }
}

export default withRouter(UserDashboard);