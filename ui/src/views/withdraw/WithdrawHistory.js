import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import { withStyles } from "@material-ui/core/styles";
import localStyles from "./WithdrawStyle"


import breadcrumbStyles from "../common/BreadcrumbStyles";
import { MAX_DATA_LIMIT } from '../../util/pagination';
import AdvancedTable from "../../components/AdvancedTable";
import WithdrawStore from "../../stores/WithdrawStore";

const styles = {
    ...breadcrumbStyles,
    ...localStyles
};

class Withdraw extends Component {
    constructor(props) {
        super(props);
        this.state = {
            data: [],
            stats: {},
            totalSize: 0,
            nsDialog: false
        }
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

        WithdrawStore.getWithdrawHistory(moneyAbbr, orgId, limit, offset, (res) => {
            const object = this.state;
            object.totalSize = Number(res.count);
            object.data = res.withdrawRequest;
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

    DateRequestedColumn = (cell, row, index, extraData) => {
        return <div>{row.txSentTime.substring(0, 10)}</div>;
    }

    AmountColumn = (cell, row, index, extraData) => {
        return <div>{row.amount} MXC</div>;
    }

    getColumns = () => (
        [{
            dataField: 'txSentTime',
            //text: i18n.t(`${packageNS}:menu.withdraw.username`),
            text: 'Date Requested',
            sort: false,
            formatter: this.DateRequestedColumn
        }, {
            dataField: 'txStatus',
            //text: i18n.t(`${packageNS}:menu.withdraw.total_token_available`),
            text: 'Status',
            sort: false
        }, {
            dataField: 'amount',
            //text: i18n.t(`${packageNS}:menu.withdraw.amount`),
            text: 'Amount',
            sort: false,
            formatter: this.AmountColumn
        }, {
            dataField: 'denyComment',
            //text: i18n.t(`${packageNS}:menu.withdraw.amount`),
            text: 'Comment',
            sort: false,
            formatter: this.AmountColumn
        }, {
            dataField: 'txHash',
            //text: i18n.t(`${packageNS}:menu.withdraw.amount`),
            text: 'Transaction Hash',
            sort: false,
            formatter: this.AmountColumn
        }]
    );

    render() {
        const { classes } = this.props;
        const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

        return (

            <React.Fragment>
                {/* {this.state.loading && <Loader />} */}
                <AdvancedTable
                    data={this.state.data}
                    columns={this.getColumns()}
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

export default withStyles(styles)(withRouter(Withdraw));
