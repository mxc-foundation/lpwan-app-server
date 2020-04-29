import classNames from "classnames";
import React, { Component } from "react";
import { Link } from "react-router-dom";


class TitleBarTitle extends Component {
  render() {
    return (
      <h4 className={classNames("mt-0", "mb-1", this.props.classes)}>
        {this.props.to && <Link to={this.props.to}>{this.props.title}</Link>}
        {!this.props.to && <React.Fragment>{this.props.title}</React.Fragment>}
      </h4>
    );
  }
}

export default TitleBarTitle;
