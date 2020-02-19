import React, { Component } from "react";
import { Row, Col } from 'reactstrap';
import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider, { Search, CSVExport } from 'react-bootstrap-table2-toolkit';
import paginationFactory from 'react-bootstrap-table2-paginator';

class AdvancedTable extends Component {
  render() {
    const { SearchBar } = Search;
    const { ExportCSVButton } = CSVExport;
    const rowsPerPage = this.props.rowsPerPage || 10;
    const totalSize = this.props.totalSize;

    return (
      <React.Fragment>
        <ToolkitProvider
          remote
          bootstrap4
          keyField={this.props.keyField}
          data={this.props.data}
          columns={this.props.columns}
          search>
          {props => (
            <React.Fragment>

              {(this.props.searchEnabled || this.props.exportEnabled) && <Row>
                {this.props.searchEnabled && <Col>
                  <SearchBar {...props.searchProps} />
                </Col>}
                {this.props.exportEnabled && <Col className="text-right">
                  <ExportCSVButton {...props.csvProps} className="btn btn-primary">{this.props.exportButtonLabel}</ExportCSVButton>
                </Col>}
              </Row>}
              <BootstrapTable
                {...props.baseProps}
                remote
                onTableChange={this.props.onTableChange}
                wrapperClasses="table-responsive"
                bordered={false}
                pagination={paginationFactory({ sizePerPage: rowsPerPage, hideSizePerPage: true, totalSize })}
                wrapperClasses="table-responsive"
              />

            </React.Fragment>
          )}
        </ToolkitProvider>
      </React.Fragment>
    );
  }
}

export default AdvancedTable;
