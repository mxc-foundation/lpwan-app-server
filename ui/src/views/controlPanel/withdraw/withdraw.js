import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Card, Button, Breadcrumb, BreadcrumbItem } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "../../withdraw/WithdrawStyle"
import i18n, { packageNS } from "../../../i18n";

import breadcrumbStyles from "../../common/BreadcrumbStyles";

import moment from "moment";

import { MAX_DATA_LIMIT } from '../../../util/pagination';
import TitleBar from "../../../components/TitleBar";
import OrgBreadCumb from '../../../components/OrgBreadcrumb';
import AdvancedTable from "../../../components/AdvancedTable";
import Loader from "../../../components/Loader";
import WithdrawStore from "../../../stores/WithdrawStore";
import NetworkServerStore from "../../../stores/NetworkServerStore";

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
            totalSize: 0
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
        if(!row.hasOwnProperty('withdrawId')){
            return;
        }
        let req = {};
        req.userId = row.userId;  
        req.confirmStatus = confirmStatus;  
        req.moneyAbbr = row.moneyAbbr;  
        req.amount = row.amount;  
        req.denyComment = "";  
        req.withdrawId = row.withdrawId;  
        req.orgId = this.props.match.params.organizationID;  

        WithdrawStore.confirmWithdraw(req, (res) => {
            const object = this.state;
            object.loading = false;
            this.props.history.push(`/control-panel/withdraw`);
        }); 
    }

    ConfirmationColumn = (cell, row, index, extraData) => {
        return <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button style={{ width: 120, marginRight: 10 }} color="primary" onClick={() => { this.confirm(row, true) }}>
                {i18n.t(`${packageNS}:menu.withdraw.confirm`)}
            </Button>
            <Button outline style={{ width: 120 }} color="primary" onClick={() => { this.confirm(row, false) }}>
                {i18n.t(`${packageNS}:menu.withdraw.deny`)}
            </Button>
        </div>;
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
            sort: false
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

    render() {
        const { classes } = this.props;
        const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

        return (

            <React.Fragment>
                <TitleBar>
                <Breadcrumb className={classes.breadcrumb}>
                    <BreadcrumbItem>
                    <Link
                        className={classes.breadcrumbItemLink}
                        to={`/organizations`}
                        onClick={() => {
                        // Change the sidebar content
                        this.props.switchToSidebarId('DEFAULT');
                        }}
                    >
                        {i18n.t(`${packageNS}:menu.control_panel`)}
                    </Link>
                    </BreadcrumbItem>
                    <BreadcrumbItem className={classes.breadcrumbItem}>{i18n.t(`${packageNS}:tr000084`)}</BreadcrumbItem>
                    <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.withdraw.withdraw`)}</BreadcrumbItem>
                </Breadcrumb>
                </TitleBar>

                <Row>
                    <Col>
                        <Card className="card-box shadow-sm">
                            {/* <CardBody className="position-relative"> */}
                            {this.state.loading && <Loader />}
                            <AdvancedTable
                                data={this.state.data}
                                columns={this.getColumns()}
                                keyField="id"
                                onTableChange={this.handleTableChange}
                                rowsPerPage={10}
                                totalSize={this.state.totalSize}
                                searchEnabled={false}
                            />
                            {/* </CardBody> */}
                        </Card>
                    </Col>
                </Row>
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(SuperAdminWithdraw));