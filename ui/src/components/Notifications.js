import React, { Component } from "react";
import { withRouter } from "react-router-dom";

import Toastr from 'toastr2';

import NotificationStore from "../stores/NotificationStore";
import dispatcher from "../dispatcher";


const Item = ({id, notification, hideNotification}) => {
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

  return <React.Fragment></React.Fragment>;
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
