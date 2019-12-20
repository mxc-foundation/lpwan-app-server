import React from 'react';
import { css } from '@emotion/core';
// First way to import
import { ScaleLoader } from 'react-spinners';

const override = css`
    display: block;
    border-color: red;
    position: fixed;
    top: 50%;
    left: 50%;
    margin-top: -50px;
    margin-left: -20px;
`;

const Spinner = (props) => (
    <div className='sweet-loading'>
        <ScaleLoader
        css={override}
        sizeUnit={"px"}
        height={80}
        width={8}
        redius={20}
        margin={'4px'}
        color={'#09006E'}  
        loading={props.on}
        />
    </div>
)

export default Spinner;