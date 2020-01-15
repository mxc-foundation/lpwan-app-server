import React, { Component } from 'react';
import AsyncSelect from 'react-select/async';
//import ProfileStore from '../stores/ProfileStore'
import SessionStore from "../stores/SessionStore";
import { SUPER_ADMIN } from "../util/M2mUtil";

const customStyles = {
    control: (base, state) => ({
        ...base,
        margin: 5,
        // match with the menu
        borderRadius: state.isFocused ? "3px 3px 0 0" : 3,
        // Overwrittes the different states of border
        borderColor: state.isFocused ? "#00FFD9" : "white",
        // Removes weird border around container
        boxShadow: state.isFocused ? null : null,
        "&:hover": {
            // Overwrittes the different states of border
            borderColor: state.isFocused ? "#00FFD9" : "white"
        }
    }),
    menu: base => ({
        ...base,
        //background:'#101c4a', [edit] 191126
        background: 'white',
        // override border radius to match the box
        borderRadius: 0,
        // kill the gap
        marginTop: 0,
        paddingLeft: 20,
        paddingRight: 20,
    }),
    menuList: base => ({
        ...base,
        //background: '#1a2d6e', [edit] 191126
        background: 'white',
        // kill the white space on first and last option
        paddingTop: 0,

    }),
    option: base => ({
        ...base,
        // kill the white space on first and last option
        padding: '10px',
        maxWidth: 229,
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis'
    }),
};

const getOrgList = (organizations) => {
    let organizationList = null;
    if (organizations) {
        organizationList = organizations.filter(function (org) {
            if (org.organizationID === SUPER_ADMIN) {
                return false
            }
            return true
        }).map((o, i) => {
            return { label: o.organizationName, value: o.organizationID };
        });
    }

    return organizationList;
};

const promiseOptions = () =>
    new Promise((resolve, reject) => {
        SessionStore.fetchProfile(
            resp => {
                resolve(getOrgList(resp.body.organizations));
            })
    });

export default class WithPromises extends Component {
    constructor() {
        super();
        this.state = {
            selectedValue: null,
            options: [],
            dOptions: {}
        };
    }

    componentDidMount() {
        promiseOptions().then(options => {
            this.setState({
                options,
                dOptions: { label: options[0].label, value: options[0].value }
            })
        })
    }

    onChange = (v) => {
        let value = null;
        if (v !== null) {
            value = v.value;
        }

        this.props.onChange({
            target: {
                id: this.props.id,
                value
            },
        });
    }
    onClick = (v) => {
    }
    render() {
        const dValue = { label: SessionStore.getOrganizations()[0].organizationName, value: SessionStore.getOrganizations()[0].organizationID }

        return (
            <AsyncSelect
                cacheOptions
                defaultOptions
                styles={customStyles}
                theme={(theme) => ({
                    ...theme,
                    borderRadius: 4,
                    colors: {
                        primary25: '#00FFD950',
                        primary: '#00FFD950',
                    },
                })}
                defaultValue={dValue}
                onClick={this.onClick}
                //defaultValue={this.state.value}
                //inputValue={this.state.value}
                onChange={this.onChange}
                loadOptions={promiseOptions}
            //options={this.state.options}
            />
        );
    }
}