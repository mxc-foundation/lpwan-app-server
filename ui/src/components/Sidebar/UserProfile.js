import React from 'react';
import { Link } from 'react-router-dom';
import { UncontrolledDropdown, DropdownMenu, DropdownItem, DropdownToggle } from 'reactstrap';

import profilePic from '../../assets/images/users/profile-icon.png';

const UserProfile = () => {
  return <React.Fragment>
      <div className="user-box text-center">
          <img src={profilePic} alt="user-img" title="Nik Patel" className="rounded-circle img-thumbnail avatar-lg" />
          <UncontrolledDropdown>
              <DropdownToggle caret tag="a" className="text-dark dropdown-toggle h5 mt-2 mb-1 d-block">
                  Nik Patel
              </DropdownToggle>
              <DropdownMenu className="user-pro-dropdown">
                  <DropdownItem>
                      <i className="fe-user mr-1"></i>
                      <span>My Account</span>
                  </DropdownItem>
                  <DropdownItem>
                      <i className="fe-settings mr-1"></i>
                      <span>Settings</span>
                  </DropdownItem>
                  <DropdownItem>
                      <i className="fe-lock mr-1"></i>
                      <span>Lock Screen</span>
                  </DropdownItem>
                  <DropdownItem>
                      <i className="fe-log-out mr-1"></i>
                      <span>Logout</span>
                  </DropdownItem>
              </DropdownMenu>
          </UncontrolledDropdown>

          <p className="text-muted">Admin Head</p>
          <ul className="list-inline">
              {/* <li className="list-inline-item">
                  <Link to="/" className="text-muted">
                      <i className="mdi mdi-settings"></i>
                  </Link>
              </li> [edit]*/}

              <li className="list-inline-item">
                  <Link to="/logout" className="text-custom">
                      <i className="mdi mdi-power"></i>
                  </Link>
              </li>
          </ul>
      </div>
  </React.Fragment>
}

export default UserProfile;