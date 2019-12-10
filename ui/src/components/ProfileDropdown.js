import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Dropdown, DropdownMenu, DropdownToggle, DropdownItem } from 'reactstrap';
import i18n, { packageNS } from '../i18n';

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
        const profilePic = this.props.profilePic || null;

        return (
            <Dropdown isOpen={this.state.dropdownOpen} toggle={this.toggleDropdown} className="notification-list">
                <DropdownToggle
                    data-toggle="dropdown"
                    tag="button"
                    className="btn btn-link nav-link dropdown-toggle nav-user mr-0 waves-effect waves-light"
                    onClick={this.toggleDropdown} aria-expanded={this.state.dropdownOpen}>
                    <img src={profilePic} className="rounded-circle" alt="user" />
                    <span className="pro-user-name ml-1">{this.props.username}  <i className="mdi mdi-chevron-down"></i> </span>
                </DropdownToggle>
                <DropdownMenu right className="profile-dropdown">
                    <div onClick={this.toggleDropdown}>
                        <div className="dropdown-header noti-title">
                            <h6 className="text-overflow m-0">{i18n.t(`${packageNS}:menu.settings.welcome`)}</h6>
                        </div>
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