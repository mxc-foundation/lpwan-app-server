import React, { Component } from 'react';

/**
 * Renders the preloader
 */
class Loader extends Component {
    render() {
        const { light } = this.props;
        return (
            <div className={light ? `preloader-light` : `preloader`}>
                <div className="status">
                    <div className="spinner-border avatar-sm text-primary m-2" role="status"></div>
                </div>
            </div>
        )
    }
}

export default Loader;