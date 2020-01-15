import React from "react";

const NonAuthLayout = (props) => {
    const children = props.children || null;
    return (<React.Fragment>
        <div className="app">
            <div id="wrapper">
                {children}
            </div>
        </div>
    </React.Fragment>
    );
}

export default NonAuthLayout;