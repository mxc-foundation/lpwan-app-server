import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Toastr from 'toastr2';

import NotificationStore from "../stores/NotificationStore";
import dispatcher from "../dispatcher";


class Item extends Component {
  state = {};
  componentDidMount() {
    const {id, notification, hideNotification} = this.props;
    const toastr = new Toastr({
      positionClass: "toast-bottom-left",
      closeButton: true,
      preventDuplicates: true,
      
      onShown: () => {
        setTimeout(() => {
          hideNotification(id);
        });
      }
    });

    this.setState({toastr});

    if (notification) {
      switch(notification.type) {
        case 'error':
          toastr.error(notification.message);
          break;
        case 'success':
          toastr.success(notification.message);
          break;
        case 'warning':
          toastr.warning(notification.message);
          break;
        case 'info':
          toastr.info(notification.message);
          break;
        default:
          break;
      }
    }
  }

  render() {
    const {id, notification, hideNotification} = this.props;
    const {toastr} = this.state;

    return <React.Fragment></React.Fragment>;
  }
}

class Notifications extends Component {
  constructor() {
    super();

    this.state = {
      notifications: NotificationStore.getAll(),
    };
  }

  /**
   * Clears the notification
   * @param {*} notificationId 
   */
  onClose(notificationId) {
    dispatcher.dispatch({
      type: "DELETE_NOTIFICATION",
      id: notificationId,
    });
  }

  componentDidMount() {
    NotificationStore.on("change", () => {
      this.setState({
        notifications: NotificationStore.getAll(),
      });
    });
  }

  render() {
    const items = this.state.notifications.map((n, i) => <Item key={n.id} id={n.id} notification={n} hideNotification={this.onClose} />);

    return (items);
  }
}

export default withRouter(Notifications);
