import React, { Component } from "react";
import { withRouter } from "react-router-dom";
import { Card, CardBody, Col, Row } from 'reactstrap';
import AdvancedTable from "../../components/AdvancedTable";
import ExtLink from "../../components/ExtLink";
import Loader from "../../components/Loader";
import i18n, { packageNS } from '../../i18n';
import TopupStore from "../../stores/TopupStore";
import { MAX_DATA_LIMIT } from '../../util/pagination';


const tableCols = [
  {
    dataField: 'amount',
    text: i18n.t(`${packageNS}:menu.topup.history.amount`),
    sort: false,
    formatter: (cell, row, rowIndex, formatExtraData) => {
      return <React.Fragment>{cell} MXC</React.Fragment>
    },
  },
  {
    dataField: 'createdAt',
    text: i18n.t(`${packageNS}:menu.topup.history.date`),
    sort: false,
    formatter: (cell, row, rowIndex, formatExtraData) => {
      return row.createdAt.substring(0, 10);
    },
  },
  {
    dataField: 'txHash',
    text: i18n.t(`${packageNS}:menu.topup.history.tx_hash`),
    sort: false,
    formatter: (cell, row, rowIndex, formatExtraData) => {
      const url = process.env.REACT_APP_ETHERSCAN_HOST + `/tx/${row.txHash}`;
      return <ExtLink to={url} />;
    },
  }
]

class TopupHistory extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: false,
      data: [],
      totalSize: 0
    };
  }

  componentDidMount() {
    this.getPage(MAX_DATA_LIMIT);
  }

  /**
   * Handles table changes including pagination, sorting, etc
   */
  handleTableChange = (type, { page, sizePerPage, filters, sortField, sortOrder }) => {
    const offset = (page - 1) * sizePerPage;
    
    /* let searchQuery = null;
    if (type === 'search' && searchText && searchText.length) {
      searchQuery = searchText;
    } */

    // TODO - how can I pass search query to server?
    this.getPage(sizePerPage, offset);
  };


  getPage = async (limit, offset) => {
    this.setState({ loading: true });
    const res = await TopupStore.getTopUpHistory(this.props.organizationID, offset, limit);
    const object = this.state;

    object.totalSize = Number(res.count);
    object.data = res.topupHistory;
    object.loading = false;
    this.setState({ object });
  }

  render() {


    return (<React.Fragment>
      <Row>
        <Col>
          <Card>
            <CardBody className="pb-0">
              <div className="position-relative">
                {this.state.loading ? <Loader /> : null}
                <AdvancedTable data={this.state.data} columns={tableCols} keyField="id" totalSize={this.state.totalSize}
                  onTableChange={this.handleTableChange} searchEnabled={true}></AdvancedTable>
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>
    </React.Fragment>
    );
  }
}

export default withRouter(TopupHistory);
