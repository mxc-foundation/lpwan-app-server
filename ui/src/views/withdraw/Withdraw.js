import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import { Row, Col, Card, Button, Breadcrumb, BreadcrumbItem, FormGroup, Label, Input } from 'reactstrap';
import { withStyles } from "@material-ui/core/styles";
import localStyles from "./WithdrawStyle"
import i18n, { packageNS } from "../../i18n";
import NumberFormat from 'react-number-format';


import breadcrumbStyles from "../common/BreadcrumbStyles";
import Modal from './Modal';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TitleBar from "../../components/TitleBar";
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";
import WithdrawStore from "../../stores/WithdrawStore";
import MoneyStore from "../../stores/MoneyStore";

const styles = {
  ...breadcrumbStyles,
  ...localStyles
};

const NumberFormatMXC = (props) => {
  const { inputRef, onChange, ...other } = props;

  return (
    <NumberFormat
      {...other}
      getInputRef={inputRef}
      onValueChange={(values) => {
        onChange({
          target: {
            value: values.value
          }
        });
      }}
      suffix=" MXC"
    />
  );
}

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
      console.log('res', res);
      const object = this.state;
      object.totalSize = Number(res.count);
      object.data = res.withdrawRequest;
      object.loading = false;
      this.setState({ object });
    });
  }


  loadData = () => {

    const orgId = this.props.match.params.organizationID;

    MoneyStore.getActiveMoneyAccount(0, orgId, resp => {
      console.log('resp', resp);
      this.setState({
        activeAccount: resp.activeAccount,
      });
    });
  }

  componentDidMount() {
    this.loadData();
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

  handleChange = (event) => {
    this.setState({ value: event.target.value });
  }

  submitWithdrawReq = () => {
    let req = {};
    
    WithdrawStore.WithdrawReq(req, (res) => {
      const object = this.state;
      object.loading = false;
      this.props.history.push(`/withdraw/${this.props.match.params.organizationID}`);
    });
  }

  render() {
    const { classes } = this.props;
    const currentOrgID = this.props.organizationID || this.props.match.params.organizationID;

    return (

      <React.Fragment>
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
                {i18n.t(`${packageNS}:tr000049`)}
              </Link>
            </BreadcrumbItem>
            <BreadcrumbItem className={classes.breadcrumbItem}>{i18n.t(`${packageNS}:tr000084`)}</BreadcrumbItem>
            <BreadcrumbItem active>{i18n.t(`${packageNS}:menu.withdraw.withdraw`)}</BreadcrumbItem>
          </Breadcrumb>
        </TitleBar>
        <Row>
          <Col xs="4">
            <FormGroup>
              <Label for="exampleEmail">Destination</Label>
              <NumberFormat className={classes.s_input} value={"asdasdasd"} />
            </FormGroup>
            <FormGroup>
              <Label for="amount">Amount</Label>
              <NumberFormat className={classes.s_input} value={246} thousandSeparator={true} suffix={' MXC'} />
            </FormGroup>
            <FormGroup>
              <Label for="exampleEmail">Available MXC</Label>
              <NumberFormat className={classes.s_input} value={6981} thousandSeparator={true} suffix={' MXC'} />
            </FormGroup>
          </Col>
          <Col xs="4"></Col>
          <Col xs="4"></Col>
        </Row>
        <Row>
          <Col>
            <FormGroup>
              <Button onClick={this.submitWithdrawReq} color="primary">Submit Withdrawal Request</Button>
            </FormGroup>
          </Col>
        </Row>
        <Row>
          <Col>
            <Card className="card-box shadow-sm">
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
            </Card>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default withStyles(styles)(withRouter(Withdraw));
