import React, { Component } from "react";
import { withRouter, Link } from "react-router-dom";

import i18n, { packageNS } from '../../i18n';
import { MAX_DATA_LIMIT } from '../../util/pagination';
import TopupStore from "../../stores/TopupStore";

import ExtLink from '../../components/ExtLink';
import AdvancedTable from "../../components/AdvancedTable";
import Loader from "../../components/Loader";

import LinkVariant from "mdi-material-ui/LinkVariant";

const TXHashColumn = (cell, row, index, extraData) => {
  const url = process.env.REACT_APP_ETHERSCAN_HOST + `/tx/${row.txHash}`;
  return <ExtLink to={url} />;
}

const getColumns = () => (
  [{
    dataField: 'from',
    text: i18n.t(`${packageNS}:menu.history.from`),
    sort: false,
  }, {
    dataField: 'to',
    text: i18n.t(`${packageNS}:menu.history.to`),
    sort: false,
  }, {
    dataField: 'txHash',
    text: i18n.t(`${packageNS}:menu.history.tx_hash`),
    formatter: TXHashColumn,
    sort: false,
  }, {
    dataField: 'moneyAbbr',
    text: i18n.t(`${packageNS}:menu.history.type`),
    sort: false,
  }, {
    dataField: 'amount',
    text: i18n.t(`${packageNS}:menu.history.mxc_amount`),
    sort: false,
  }, {
    dataField: 'lastUpdateTime',
    text: i18n.t(`${packageNS}:menu.history.update_date`),
    sort: false,
  },{
    dataField: 'status',
    text: i18n.t(`${packageNS}:menu.history.status`),
    sort: false,
  }]
);

class Transactions extends Component {
  constructor(props) {
    super(props);

    this.handleTableChange = this.handleTableChange.bind(this);
    this.getPage = this.getPage.bind(this);
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
    const offset = (page - 1) * sizePerPage + 1;

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
    this.setState({ loading: true });
    TopupStore.getTransactionsHistory(this.props.organizationID, offset, limit, res => {
      const object = this.state;
      object.totalSize = res.count;
      object.data = res.transactionHistory;
      object.loading = false;
      this.setState({object});
    }); 
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  render() {
    return(
      <div className="position-relative">
        {this.state.loading && <Loader />}
        <AdvancedTable data={this.state.data} columns={getColumns()}
          keyField="id" onTableChange={this.handleTableChange} searchEnabled={false} totalSize={this.state.totalSize} rowsPerPage={10}></AdvancedTable>
      </div>
    );
  }
}

export default withRouter(Transactions);