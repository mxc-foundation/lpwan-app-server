import { withStyles } from "@material-ui/core/styles";
import React, { Component } from "react";
import { Link, withRouter } from "react-router-dom";
import { Breadcrumb, BreadcrumbItem, Button, Card, Col, Row } from 'reactstrap';
import AdvancedTable from "../../../components/AdvancedTable";
import Loader from "../../../components/Loader";
import TitleBar from "../../../components/TitleBar";
import i18n, { packageNS } from "../../../i18n";
import WithdrawStore from "../../../stores/WithdrawStore";
import { MAX_DATA_LIMIT } from '../../../util/pagination';
import breadcrumbStyles from "../../common/BreadcrumbStyles";
import localStyles from "../../withdraw/WithdrawStyle";
import Modal from './Modal';



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

        /* let searchQuery = null;
        if (type === 'search' && searchText && searchText.length) {
            searchQuery = searchText;
        } */
        // TODO - how can I pass search query to server?
        this.getPage(sizePerPage, offset);
    }

    /**
     * Fetches data from server
     */
    getPage = (limit, offset) => {
        limit = MAX_DATA_LIMIT;
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
        req.denyComment = (this.state.value === undefined) ? '' : this.state.value;
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
        this.setState({ value: event.target.value });
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
                    {/* <Row style={{width: '100%', height: '700px'}}></Row> */}
                    {this.state.nsDialog && <Modal
                        title={i18n.t(`${packageNS}:menu.withdraw.confirm_modal_title`)}
                        context={(this.state.status) ? i18n.t(`${packageNS}:menu.withdraw.confirm_text`) : i18n.t(`${packageNS}:menu.withdraw.deny_text`)}
                        status={this.state.status}
                        row={this.state.row}
                        handleChange={this.handleChange}
                        closeModal={() => this.setState({ nsDialog: false })}
                        callback={() => { this.confirm(this.state.row, this.state.status) }} />}
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
                            </Card>
                        </Col>
                    </Row>
                </div>{/* will be taken off  */}
            </React.Fragment>
        );
    }
}

export default withStyles(styles)(withRouter(SuperAdminWithdraw));