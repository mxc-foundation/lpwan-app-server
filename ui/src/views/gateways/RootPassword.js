import React, { Component } from "react";

import { Button, CustomInput, Input, InputGroup, InputGroupAddon, Modal, ModalBody } from 'reactstrap';
import i18n, { packageNS } from '../../i18n';


class RootPassword extends Component {
  constructor() {
    super();
    this.state = {
      showPassword: false,
      password: '',
      modal: true,
      agreedToTerms: false
    };
    this.getPassword = this.getPassword.bind(this);
    this.close = this.close.bind(this);
    this.onAgree = this.onAgree.bind(this);
    this.togglePassword = this.togglePassword.bind(this);
  }

  close = () => {
    this.setState({modal: !this.state.modal});
    if (this.props.onClose)
      this.props.onClose();
  }

  /**
   * On checked
   */
  onAgree = (e) => {
    this.setState({agreedToTerms: e.target.checked});
  }

  /**
   * toggle visibility of password
   */
  togglePassword = () => {
    this.setState({showPassword: !this.state.showPassword });
  }

  /**
   * Fetches password 
   */
  getPassword() {
    // TODO - Map to real api when it's available
    this.setState({
      password: '12345678'
    });
  }

  render() {
    
    return (<React.Fragment>
      <Modal isOpen={this.state.modal} toggle={this.close} size="lg">
        <ModalBody className="text-center">
          {!this.state.password ? <div className="px-3">
            <h1 className="text-danger display-4"><i className="mdi mdi-alert-circle"></i></h1>
            <h4>{i18n.t(`${packageNS}:tr000620`)}</h4>
            <p className='text-left p-3'>{i18n.t(`${packageNS}:tr000621`)}</p>

            <div className="px-3 text-left">
              <CustomInput type="checkbox" id="agree" label={i18n.t(`${packageNS}:tr000622`)} checked={this.state.agreedToTerms} onChange={this.onAgree} />
            </div>
            <div className='my-2'>
              <Button color="danger" className="mr-3 d-inline" disabled={!this.state.agreedToTerms} 
                onClick={this.getPassword}>{i18n.t(`${packageNS}:tr000623`)}</Button>
              <Button color="success" className="d-inline" onClick={this.close}>{i18n.t(`${packageNS}:tr000624`)}</Button>
            </div>
          </div>: <div className="px-3">
              <h1 className="text-danger display-4"><i className="mdi mdi-alert-circle"></i></h1>
              <h4 className="mx-3">{i18n.t(`${packageNS}:tr000625`)}</h4>

              <div className="py-2 px-3">
                <InputGroup>
                  <Input className=""
                    type={this.state.showPassword ? "text" : "password"}
                    defaultValue={this.state.password} />
                    <InputGroupAddon addonType="append">
                    <Button onClick={this.togglePassword}>{!this.state.showPassword ? <i className="mdi mdi-eye"></i>: <i className="mdi mdi-eye-off"></i>}</Button>
                    </InputGroupAddon>
                </InputGroup>
                <Button color="success" className="mt-2" onClick={this.close}>{i18n.t(`${packageNS}:tr000430`)}</Button>
              </div>
              
            </div>}

          
        </ModalBody>
      </Modal>
    </React.Fragment>
    );
  }
}

export default RootPassword;
