import React, { Component } from 'react';
import { withRouter } from "react-router-dom";
import i18n, { packageNS } from '../i18n';
import { Button, Modal, ModalHeader, ModalBody, ModalFooter, Progress } from 'reactstrap';

let timer = null;

class ModalWithProgress extends Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: true,
            completed: 0
        };
    }

    progress = () => {
        if (this.state.completed === 100) {
            clearInterval(timer);
            this.props.handleProgress(this.state.completed);
            return;
        }
        this.setState({ completed: this.state.completed + 10 })

    }

    componentDidMount() {
        timer = setInterval(this.progress, 100);
    }

    toggle = () => {
        this.setState({ modal: !this.state.modal });
    };

    proc = () => {
        clearInterval(timer);
        this.props.handleProgress();
    }

    render() {
        return (
            <React.Fragment>
                <Modal isOpen={this.state.modal} toggle={this.toggle} centered={true}>
                    <ModalHeader toggle={this.toggle}>{i18n.t(`${packageNS}:menu.messages.stake_proc_tit`)}</ModalHeader>
                    <ModalBody>
                        <div>
                            <Progress value={this.state.completed} >{i18n.t(`${packageNS}:menu.staking.staking_is_being_processed`)}</Progress>
                        </div>
                        <div style={{paddingTop: 16}}>
                            {i18n.t(`${packageNS}:menu.messages.stake_proc_desc`)}
                        </div>
                    </ModalBody>
                    <ModalFooter>
                        <Button color="primary" onClick={this.proc}>{i18n.t(`${packageNS}:menu.withdraw.cancel`)}</Button>
                    </ModalFooter>
                </Modal>
            </React.Fragment>
        );
    }
}

export default withRouter(ModalWithProgress);
/* const ModalWithProgress = (props) => {
    const {
        className,
        showConfirmButton = true,
        show = true,
    } = props;

    let [modal, setModal, completed=0] = useState(show);

    const toggle = () => setModal(!modal);
    const proc = () => {
        setModal(!modal);
        props.callback();
    }

    React.useEffect(() => {

        function progress() {
            if (completed === 100) {
                clearInterval(timer);
                props.onProgress(completed);
            }

            //const diff = Math.random() * 10;

            return completed += 10;
        }

        const timer = setInterval(progress, 800);
        return () => {
            clearInterval(timer);
        };
    }, []);


    return (
        <div>
            <Modal isOpen={modal} toggle={toggle} className={className} centered={true}>
                <ModalHeader toggle={toggle}>{props.title}</ModalHeader>
                <ModalBody>
                    <div>
                        <Progress value={completed} >{i18n.t(`${packageNS}:staking.staking_is_being_processed`)}</Progress>
                    </div>
                </ModalBody>
                <ModalFooter>
                    {showConfirmButton && <Button color="primary" onClick={proc}>{i18n.t(`${packageNS}:menu.withdraw.cancel`)}</Button>}
                </ModalFooter>
            </Modal>
        </div>
    );
}

export default ModalWithProgress; */