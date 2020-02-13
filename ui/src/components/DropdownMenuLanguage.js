import React, { Component } from "react";
import Select, {components} from "react-select";
import { withStyles } from "@material-ui/core/styles";
import classNames from "classnames";
import { DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from "../i18n";
import SessionStore from "../stores/SessionStore";
import FlagIcon from "./FlagIcon";
import DropdownMenuLanguageStyle from "./DropdownMenuLanguageStyle";
import DropdownMenuLanguageMobileStyle from "./DropdownMenuLanguageMobileStyle";

const styles = {
  // languageWrapper: {
  //   display: "inline-flex"
  // },
  languageIcon: {
    display: "inline-block"
  },
  languageSelection: {
    display: "inline-block"
  }
};


const customSelectComponents = {
  SingleValue: ({ children, ...props }) => {
    
    const {code} = props.data || {};
    return (<components.SingleValue {...props}>
      {<FlagIcon
              code={code}
              // size='1x'
            />}
    </components.SingleValue>);
  }
};

class WithPromises extends Component {
  constructor() {
    super();

    this.state = {
      selectedOption: null,
      options: []
    };
  } 

  componentDidMount() {
    let selectedOption = null;

    const language = SessionStore.getLanguage();

    if (!language || !language.id) {
      selectedOption = DEFAULT_LANGUAGE;
    } else if (language.label && language.label && language.value && language.code) {
      selectedOption = {
        id: language.id,
        label: language.label,
        value: language.value,
        code: language.code
      };
    }

    this.setState({
      selectedOption,
      options: SUPPORTED_LANGUAGES
    });
  }

  onChangeLanguage = selectedOption => {
    if (
      selectedOption !== null && selectedOption.id !== null && selectedOption.label !== null &&
      selectedOption.value !== null && selectedOption.code !== null
    ) {
      this.setState({
        selectedOption
      });
  
      this.props.onChangeLanguage({
        id: selectedOption.id,
        label: selectedOption.label,
        value: selectedOption.value,
        code: selectedOption.code
      });
    }
  }

  render() {
    const { selectedOption } = this.state;
    let isMobile = this.props.isMobile;
    let customStyle = DropdownMenuLanguageStyle;
    if(isMobile){
      customStyle = DropdownMenuLanguageMobileStyle;
    }
    
    return (
      <div className={classNames(this.props.classes.languageWrapper)}>
        {/* {
          selectedOption && selectedOption.code
          ? (
            <FlagIcon
              className={classNames(this.props.classes.languageIcon)}
              code={selectedOption.code}
              size='2x'
            />
          ) : null
        } */}
        <Select
          className={classNames('react-select', this.props.classes.languageSelection)}
          menuPlacement="auto"
          classNamePrefix="react-select"
          styles={customStyle}
          theme={(theme) => ({
            ...theme,
            borderRadius: 4,
            colors: {
              primary25: "#00FFD950",
              primary: "#00FFD950",
            },
          })}
          isSearchable={false}
          placeholder="Select Language"
          onChange={this.onChangeLanguage}
          options={SUPPORTED_LANGUAGES}
          value={selectedOption}
          components={customSelectComponents}
          {...(this.props.extraSelectOpts || {})}
        />
      </div>
    );
  }
}

export default withStyles(styles)(WithPromises);
