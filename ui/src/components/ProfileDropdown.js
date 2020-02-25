import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Dropdown, DropdownMenu, DropdownToggle, DropdownItem } from 'reactstrap';
import i18n, { packageNS } from '../i18n';
import defaultProfilePic from '../assets/images/users/profile-icon.png';

class ProfileDropdown extends Component {
    constructor(props) {
        super(props);

        this.toggleDropdown = this.toggleDropdown.bind(this);
        this.state = {
            dropdownOpen: false
        };
    }

    /*:: toggleDropdown: () => void */
    toggleDropdown() {
        this.setState({
            dropdownOpen: !this.state.dropdownOpen
        });
    }

    render() {
        if (!this.props.user) {
            return null;
        }
        const { user: { id, profilePic, username } } = this.props;

        return (
            <Dropdown isOpen={this.state.dropdownOpen} toggle={this.toggleDropdown} className="notification-list">
                <DropdownToggle
                    data-toggle="dropdown"
                    tag="button"
                    className="btn btn-link nav-link dropdown-toggle nav-user mr-0 waves-effect waves-light"
                    onClick={this.toggleDropdown} aria-expanded={this.state.dropdownOpen}>
                    <img src={profilePic || defaultProfilePic} className="rounded-circle" alt="user" />
                    <span className="pro-user-name ml-1">{username}  <i className="mdi mdi-chevron-down"></i> </span>
                </DropdownToggle>
                <DropdownMenu right className="profile-dropdown">
                    <div onClick={this.toggleDropdown}>
                        {
                            this.props.user ? (
                                <Link to={`/users/${id}`} className="dropdown-item notify-item side-nav-link-ref">
                                    <i className="mdi mdi-account-circle"></i>
                                    <span> {i18n.t(`${packageNS}:tr000452`)} </span>
                                </Link>
                            ) : null
                        }
                        {this.props.menuItems.map((item, i) => {
                            return <React.Fragment key={i + "-profile-menu"}>
                                {item.hasDivider ? <DropdownItem divider /> : null}
                            <Link to={item.redirectTo} className="dropdown-item notify-item">
                                    <i className={`${item.icon} mr-1`}></i>
                                    <span>{item.label}</span>
                                </Link>
                            </React.Fragment>
                        })}
                    </div>
                </DropdownMenu>
            </Dropdown>
        );
    }
}

export default ProfileDropdown;