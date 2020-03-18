import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import localStyles from "./StakeStyle"
import i18n, { packageNS } from "../../i18n";
import breadcrumbStyles from "../common/BreadcrumbStyles";
import { MAX_DATA_LIMIT } from '../../util/pagination';
import NumberFormat from 'react-number-format';
import AdvancedTable from "../../components/AdvancedTable";
import StakeStore from "../../stores/StakeStore";

const styles = {
    ...breadcrumbStyles,
    ...localStyles
};

class StakeHistory extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
            stats: {},
            totalSize: 0,
            nsDialog: false
        }
        const { classes } = this.props;
    }
    /**
       * Handles table changes including pagination, sorting, etc
       */
    handleTableChange = (type, { page, sizePerPage, searchText, sortField, sortOrder, searchField }) => {
        const offset = (page - 1) * sizePerPage;

        let searchQuery = null;
        if (type === 'search' && searchText && searchText.length) {
            searchQuery = searchText;
        }
        // TODO - how can I pass search query to server?
        this.getPage(sizePerPage, offset);
    }

    /**
     * Fetches data from server
     */
    getPage = (limit, offset) => {
        limit = MAX_DATA_LIMIT;
        const defaultOrgId = 0;
        this.setState({ loading: true });
        const moneyAbbr = 2;
        const orgId = this.props.match.params.organizationID;
        console.log('paramter:', orgId, offset, limit);
        StakeStore.getStakingHistory(orgId, offset, limit, (res) => {
            const object = this.state;
            object.totalSize = Number(res.count);
            object.data = res.stakingHist;
            object.loading = false;
            this.setState({ object });
        });
    }

    componentDidMount() {
        this.getPage(MAX_DATA_LIMIT);
    }

    componentDidUpdate(prevProps, prevState) {
        if (prevState !== this.state && prevState.data !== this.state.data) {

        }
    }

    StartColumn = (cell, row, index, extraData) => {
        return row.start.substring(0, 10);
    }

    EndColumn = (cell, row, index, extraData) => {
        return row.end.substring(0, 10);
    }

    getColumns = () => (
        [{
            dataField: 'stakeAmount',
            text: i18n.t(`${packageNS}:menu.staking.stake_amount`),
            sort: false
        }, {
            dataField: 'start',
            text: i18n.t(`${packageNS}:menu.staking.start`),
            formatter: this.StartColumn,
            sort: false,
        }, {
            dataField: 'end',
            text: i18n.t(`${packageNS}:menu.staking.end`),
            formatter: this.EndColumn,
            sort: false,
        }, {
            dataField: 'revMonth',
            text: i18n.t(`${packageNS}:menu.staking.revenue_month`),
            sort: false,
        }, {
            dataField: 'networkIncome',
            text: i18n.t(`${packageNS}:menu.staking.network_income`),
            sort: false,
        }, {
            dataField: 'monthlyRate',
            text: i18n.t(`${packageNS}:menu.staking.monthly_rate`),
            sort: false,
        }, {
            dataField: 'revenue',
            text: i18n.t(`${packageNS}:menu.staking.revenue`),
            sort: false,
        }, {
            dataField: 'balance',
            text: i18n.t(`${packageNS}:menu.staking.balance`),
            sort: false,
        }]
    );

    render() {
        const { classes } = this.props;
        return (
            <React.Fragment>
                {/* {this.state.loading && <Loader />} */}
                <AdvancedTable
                    data={this.state.data}
                    columns={this.getColumns(classes)}
                    keyField="id"
                    onTableChange={this.handleTableChange}
                    rowsPerPage={10}
                    totalSize={this.state.totalSize}
                    searchEnabled={false}
                />
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(StakeHistory));
