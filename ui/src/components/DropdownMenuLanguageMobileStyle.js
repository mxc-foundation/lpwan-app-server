const DropdownMenuLanguageMobileStyle = {
    control: (base, state) => ({
      ...base,
      //color: "#FFFFFF",
      width: "70px",
      margin: "18px 0px 18px 0px",
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
      background:'white',
      // override border radius to match the box
      borderRadius: 0,
      // kill the gap
      marginTop: 0,
      // paddingLeft: 20,
      // paddingRight: 20,
    }),
    menuList: base => ({
      ...base,
      background: 'white',
      // kill the white space on first and last option
      paddingTop: 0,
    }),
    option: base => ({
      ...base,
      // kill the white space on first and last option
      padding: "10px",
      maxWidth: 229,
      whiteSpace: "nowrap", 
      overflow: "hidden",
      textOverflow: "ellipsis"
    }),
  };
  
export default DropdownMenuLanguageMobileStyle;
