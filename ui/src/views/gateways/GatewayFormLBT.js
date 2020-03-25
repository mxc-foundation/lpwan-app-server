import React, { Component } from "react";
import { Col, CustomInput, Row, Table } from "reactstrap";
import EditableTextInput from "../../components/EditableTextInput";
import i18n, { packageNS } from "../../i18n";


class GatewayFormLBT extends Component {
  constructor(props) {
    super(props);

    this.state = {
      records: props.records,
      status: props.status
    };

    this.onDataChanged = this.onDataChanged.bind(this);
    this.onStatusChange = this.onStatusChange.bind(this);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.records !== this.props.records) {
      this.setState({ records: this.props.records });
    }
    if (prevProps.status !== this.props.status) {
      this.setState({ status: this.props.status });
    }
  }

  /**
   * On data edited
   * @param {*} object
   * @param {*} field
   * @param {*} changedValue
   */
  onDataChanged(idx, field, changedValue) {
    let records = [...this.state.records];
    records[idx][field] = changedValue;
    if (this.props.onDataChanged) {
      this.props.onDataChanged(records);
    } else {
      this.setState({ records });
    }
  }

  /**
   * on lbt status change
   */
  onStatusChange() {
    this.props.onStatusChanged(!this.state.status);
  }

  render() {
    return (
      <React.Fragment>
        <Row>
          <Col>
            <h5 className="d-inline mr-1">{i18n.t(`${packageNS}:tr000598`)}</h5>
            <CustomInput
              type="switch"
              id="lbtConfig"
              name="lbtConfig"
              className="d-inline align-middle"
              checked={this.state.status}
              onChange={this.onStatusChange}
            />
            <Table className="mt-2">
              <thead>
                <tr>
                  <th width="33%">{i18n.t(`${packageNS}:tr000610`)}</th>
                  <th width="33%">{i18n.t(`${packageNS}:tr000611`)}</th>
                  <th width="33%">{i18n.t(`${packageNS}:tr000613`)}</th>
                </tr>
              </thead>
              <tbody>
                {this.state.records.map((record, idx) => {
                  return (
                    <tr key={idx}>
                      <td>
                        <EditableTextInput
                          value={record.channel}
                          id={idx}
                          field={"channel"}
                          onDataChanged={this.onDataChanged}
                        />
                      </td>
                      <td>
                        <EditableTextInput
                          value={record.freq_hz}
                          id={idx}
                          field={"freq_hz"}
                          onDataChanged={this.onDataChanged}
                          valueFormatter={value => value / 1000000}
                          inputType="number"
                        />
                      </td>
                      <td>
                        <EditableTextInput
                          value={record.scan_time_us}
                          id={idx}
                          field={"scan_time_us"}
                          onDataChanged={this.onDataChanged}
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </Table>
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

export default GatewayFormLBT;
