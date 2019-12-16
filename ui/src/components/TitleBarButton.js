import React, { Component } from "react";
import classNames from "classnames";
import { Link } from "react-router-dom";
import { Button } from 'reactstrap';


class TitleBarButton extends Component {
  render() {
    const color = this.props.color || "primary";
    const icon = this.props.icon || null;

    return (<React.Fragment>
      {this.props.to && <Link to={this.props.to} className={classNames("btn", `btn-${color}`)}>
        {icon}
        {this.props.label}
      </Link>}
      {!this.props.to && <Button to={this.props.to} color={color} onClick={this.props.onClick}>{icon}{this.props.label}</Button>}
    </React.Fragment>
    );
  }
}

export default TitleBarButton;
