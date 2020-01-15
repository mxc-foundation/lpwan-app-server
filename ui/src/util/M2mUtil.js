import React, { Component } from 'react';
import NumberFormat from 'react-number-format';

export const SUPER_ADMIN = '0';

export function BackToLora() {
    window.location.replace(getLoraHost());
}

export function getLoraHost(){
    let host = process.env.REACT_APP_LORA_APP_SERVER;
    const origin = window.location.origin;
    
    if(origin.includes(process.env.REACT_APP_SUBDOM_M2M)){
        host = origin.replace(process.env.REACT_APP_SUBDOM_M2M, process.env.REACT_APP_SUBDOM_LORA);
    }

    return host;
}

export const NumberFormatMXC=(props)=> {
	const { inputRef, onChange, ...other } = props;

	return (
		<NumberFormat
			{...other}
			getInputRef={inputRef}
			onValueChange={(values) => {
				onChange({
					target: {
						value: values.value
					}
				});
			}}
			suffix=" MXC"
		/>
	);
}

export const NumberFormatPerc =(props) =>{
	const { inputRef, onChange, id, ...other } = props;
    let suffix = " %";
    if(id === 'revRate'){
        suffix += " Monthly";
    }
	return (
		<NumberFormat
			{...other}
			getInputRef={inputRef}
			onValueChange={(values) => {
				onChange({
					target: {
						value: values.value
					}
				});
			}}
			suffix={suffix}
		/>
	);
}