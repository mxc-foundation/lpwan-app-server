import React, { Component } from 'react';
import Select from 'react-select';
import { DV_FREE_GATEWAYS_LIMITED } from "../../util/Data";

const MXC_TURQUOISE_RGB_COLOR = "#00FFD950";

const customStyles = {
    control: (base, state) => ({
      ...base,
      marginTop: 0,
      color: state.isDisabled ? 'grey' : 'white',
      marginBottom: 0,
      // match with the menu
      borderRadius: state.isFocused ? "3px 3px 0 0" : 3,
      // Overwrites the different states of border
      border: "1px solid #DDDDDD",
      // Removes weird border around container
      boxShadow: state.isFocused ? null : null,
      "&:hover": {
        borderColor: MXC_TURQUOISE_RGB_COLOR
      },
      "&:focused": {
        borderColor: MXC_TURQUOISE_RGB_COLOR
      }, 
      //left: 'calc(100%/3.3)',
      //width: '100%',
      maxWidth: 120
    }),
    menu: base => ({
      ...base,
      // background:'#101c4a',
      // override border radius to match the box
      borderRadius: "0 0 3px 3px",
      // kill the gap
      //left: 'calc(100%/3.3)',
      marginTop: 0,
      //width: '100%',
      maxWidth: 120
    }),
    menuList: base => ({
      ...base,
      // background: '#1a2d6e',
      // kill the white space on first and last option
      paddingTop: 0,
      zIndex: 999,
      //width: '100%',
      //left: 'calc(100%/3.3)',
      maxWidth: 120
    }),
    // Customize the selection box value that is shown
    singleValue: (base, state) => ({
      color: '#666666'
    }),
    // Customize the selection box down arrow
    indicatorsContainer: (base, state) => ({
      color: '#CCCCCC'
    }),
  };
 
export default class SelectPlain extends Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedValue: null,
            options:[],
            //isDisabled: (this.props.gwId)?true:false,
            //haveGateway:false
        };
    } 

    componentDidMount() {
      
    }
      
    onChange = (v) => {
        let value = null;
        if (v !== null) {
            value = v.value;
        }
        
        this.props.onSelectChange({
            target: {
                id: this.props.id,
                value: value,
            },...this.props
        });
    }
    
    
    render() {
        let dValue = this.props.defaultValue;
        let options = this.props.options;

        if(!this.props.haveGateway){
          options = options.filter(function(value, index, arr){
            return value.value !== DV_FREE_GATEWAYS_LIMITED;//private
          });
        }

        return (
            <Select 
                //cacheOptions
                defaultOptions
                styles={customStyles}
                theme={(theme) => ({
                    ...theme,
                    borderRadius: 4,
                    colors: {
                        primary25: MXC_TURQUOISE_RGB_COLOR,
                        primary: MXC_TURQUOISE_RGB_COLOR,
                    },
                })}
                isDisabled={this.props.isDisabled}
                value={dValue}
                onChange={this.onChange}
                options={options}
            />
      );
    }
}
