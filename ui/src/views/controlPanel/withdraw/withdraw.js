import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Card, Button, Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "../../withdraw/WithdrawStyle"
import i18n, { packageNS } from "../../../i18n";

import breadcrumbStyles from "../../common/BreadcrumbStyles";
import Modal from './Modal';
import { MAX_DATA_LIMIT } from '../../../util/pagination';
import TitleBar from "../../../components/TitleBar";
import AdvancedTable from "../../../components/AdvancedTable";
import Loader from "../../../components/Loader";
import WithdrawStore from "../../../stores/WithdrawStore";

const styles = {
    ...breadcrumbStyles,
    ...localStyles
};

class SuperAdminWithdraw extends Component {
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
        //this.setState({ loading: true });//commented out by Namgyeong 

        WithdrawStore.getWithdrawRequestList(limit, offset, (res) => {
            console.log('res', res);
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

    confirm = (row, confirmStatus) => {
        if (!row.hasOwnProperty('withdrawId')) {
            return;
        }
        let req = {};
        req.orgId = 1;
        req.confirmStatus = confirmStatus;
        req.denyComment = (this.state.value === undefined)?'':this.state.value;
        req.withdrawId = row.withdrawId;

        WithdrawStore.confirmWithdraw(req, (res) => {
            const object = this.state;
            object.loading = false;
            this.props.history.push(`/control-panel/withdraw`);
        });
    }

    ConfirmationColumn = (cell, row, index, extraData) => {
        return <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button style={{ width: 120, marginRight: 10 }} color="primary" onClick={() => { this.openModal(row, true) }}>
                {i18n.t(`${packageNS}:menu.withdraw.confirm`)}
            </Button>
            <Button outline style={{ width: 120 }} color="primary" onClick={() => { this.openModal(row, false) }}>
                {i18n.t(`${packageNS}:menu.withdraw.deny`)}
            </Button>
        </div>;
    }

    openModal = (row, status) => {
        this.setState({
            nsDialog: true,
            row,
            status
        });
    };

    AvailableTokenColumn = (cell, row, index, extraData) => {
        return <div>{row.availableToken} MXC</div>;
    }

    AmountColumn = (cell, row, index, extraData) => {
        return <div>{row.amount} MXC</div>;
    }

    getColumns = () => (
        [{
            dataField: 'userName',
            text: i18n.t(`${packageNS}:menu.withdraw.username`),
            sort: false
        }, {
            dataField: 'availableToken',
            text: i18n.t(`${packageNS}:menu.withdraw.total_token_available`),
            sort: false,
            formatter: this.AvailableTokenColumn
        }, {
            dataField: 'amount',
            text: i18n.t(`${packageNS}:menu.withdraw.amount`),
            sort: false,
            formatter: this.AmountColumn
        }, {
            dataField: '',
            text: '',
            sort: false,
            formatter: this.ConfirmationColumn
        }]
    );

    handleChange = (event) => {
        this.setState({value: event.target.value});
    }

    render() {
        const { classes } = this.props;
        const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

        return (

            <React.Fragment>
                <div className="position-relative">{/* will be taken off  */}
                <div className="card-coming-soon-2">{/* will be taken off  */}
                    <h1 className="title">{i18n.t(`${packageNS}:menu.dashboard.coming_soon`)}</h1>{/* will be taken off  */}
                </div>{/* will be taken off  */}
                <Row style={{width: '100%', height: '700px'}}>

                </Row>
                </div>{/* will be taken off  */}
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(SuperAdminWithdraw));