import React from 'react';
import i18n, { packageNS } from '../../../i18n';

import DataMap from './DataMap';
import DataPacketChart from './DataPacketChart';
import EarnedAmountChart from './EarnedAmountChart';
import MXCAmountChart from './MXCAmountChart';
import StakingAmountChart from './StakingAmountChart';
import StatWidget from './StatWidget';
import Tickets from './Tickets';
import Topup from './Topup';
import UserTopup from './UserTopup';
import Withdrawal from './Withdrawal';

import ticketImg from './images/tickets.png';

const WIDGET_TYPE_GRAPH = "graph";
const WIDGET_TYPE_STAT = "stat";
const WIDGET_TYPE_MAP = "map";


const TotalUsers = (props) => {
    return <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalUsers.title`)}
        {...props} formatNum={true}></StatWidget>
}

const TotalOrgs = (props) => {
    return <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalOrgs.title`)}
        {...props} formatNum={true}></StatWidget>
}

const TotalGateway = (props) => {
    return <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalGateways.title`)}
        {...props} formatNum={true}></StatWidget>
}

const TotalDevices = (props) => {
    return <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalDevices.title`)}
        {...props} formatNum={true}></StatWidget>
}

const TotalApps = (props) => {
    return <StatWidget label={i18n.t(`${packageNS}:menu.dashboard.totalApps.title`)}
        {...props} formatNum={true}></StatWidget>
}

const DataPacketsReceived = (props) => {
    return <DataPacketChart {...props} color="#10c469"
        title={i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.title`)}
        subTitle={i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.subTitle`)}
        subTitleClass="text-success" labelField="day" />
}

const DataPacketsSent = (props) => {
    return <DataPacketChart {...props} color="#71b6f9"
        title={i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.title`)}
        subTitle={i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.subTitle`)}
        subTitleClass="text-primary" labelField="day" />
}

const DataPacketsByChannel = (props) => {
    return <DataPacketChart {...props} color="#71b6f9"
        title={i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.title`)}
        labelField="channel" showYAxis={true} />
}

const DataPacketsBySpread = (props) => {
    return <DataPacketChart {...props} color="#71b6f9"
        title={i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.title`)}
        labelField="spreadFactor" showYAxis={true} />
}

const adminWidgetCatalog = [
    {
        type: WIDGET_TYPE_GRAPH, name: 'tickets', label: i18n.t(`${packageNS}:menu.dashboard.tickets.title`),
        component: Tickets, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.tickets.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'withdrawal', label: i18n.t(`${packageNS}:menu.dashboard.withdrawal.title`),
        component: Withdrawal, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.withdrawal.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'topup', label: i18n.t(`${packageNS}:menu.dashboard.topup.title`),
        component: Topup, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.topup.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalUsers', label: i18n.t(`${packageNS}:menu.dashboard.totalUsers.title`),
        component: TotalUsers, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalUsers.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalOrganizations', label: i18n.t(`${packageNS}:menu.dashboard.totalOrgs.title`),
        component: TotalOrgs, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalOrgs.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalGateways', label: i18n.t(`${packageNS}:menu.dashboard.totalGateways.title`),
        component: TotalGateway, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalOrgs.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalDevices', label: i18n.t(`${packageNS}:menu.dashboard.totalDevices.title`),
        component: TotalDevices, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalDevices.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalApplications', label: i18n.t(`${packageNS}:menu.dashboard.totalApps.title`),
        component: TotalApps, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalApps.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'supernodeAmount', label: i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.title`),
        component: MXCAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'stakingAmount', label: i18n.t(`${packageNS}:menu.dashboard.stakingAmountChart.title`),
        component: StakingAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.stakingAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'earnedAmount', label: i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.title`),
        component: EarnedAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsReceived', label: i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.title`),
        component: DataPacketsReceived, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsSent', label: i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.title`),
        component: DataPacketsSent, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.description`),
    },
    {
        type: WIDGET_TYPE_MAP, name: 'dataMap', label: i18n.t(`${packageNS}:menu.dashboard.dataMap.title`),
        component: DataMap, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.dataMap.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsByChannel', label: i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.title`),
        component: DataPacketsByChannel, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsBySpreadFactor', label: i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.title`),
        component: DataPacketsBySpread, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.description`),
    },
];

const userWidgetCatalog = [
    {
        type: WIDGET_TYPE_GRAPH, name: 'topup', label: i18n.t(`${packageNS}:menu.dashboard.topup.title`),
        component: UserTopup, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.topup.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalUsers', label: i18n.t(`${packageNS}:menu.dashboard.totalUsers.title`),
        component: TotalUsers, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalUsers.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalGateways', label: i18n.t(`${packageNS}:menu.dashboard.totalGateways.title`),
        component: TotalGateway, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalOrgs.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalDevices', label: i18n.t(`${packageNS}:menu.dashboard.totalDevices.title`),
        component: TotalDevices, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalDevices.description`),
    },
    {
        type: WIDGET_TYPE_STAT, name: 'totalApplications', label: i18n.t(`${packageNS}:menu.dashboard.totalApps.title`),
        component: TotalApps, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.totalApps.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'supernodeAmount', label: i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.title`),
        component: MXCAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.mxcAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'stakingAmount', label: i18n.t(`${packageNS}:menu.dashboard.stakingAmountChart.title`),
        component: StakingAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.stakingAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'earnedAmount', label: i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.title`),
        component: EarnedAmountChart, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.earnedAmountChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsReceived', label: i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.title`),
        component: DataPacketsReceived, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsReceivedChart.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsSent', label: i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.title`),
        component: DataPacketsSent, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsSentChart.description`),
    },
    {
        type: WIDGET_TYPE_MAP, name: 'dataMap', label: i18n.t(`${packageNS}:menu.dashboard.dataMap.title`),
        component: DataMap, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.dataMap.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsByChannel', label: i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.title`),
        component: DataPacketsByChannel, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsByChannel.description`),
    },
    {
        type: WIDGET_TYPE_GRAPH, name: 'packetsBySpreadFactor', label: i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.title`),
        component: DataPacketsBySpread, avatar: ticketImg,
        description: i18n.t(`${packageNS}:menu.dashboard.packetsBySpreadFactor.description`),
    },
];

export { WIDGET_TYPE_GRAPH, WIDGET_TYPE_MAP, WIDGET_TYPE_STAT, adminWidgetCatalog, userWidgetCatalog };

