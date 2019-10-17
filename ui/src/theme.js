import { createMuiTheme } from "@material-ui/core/styles";

const blueHighLight = '#4D89E5';
const blueHighLight20 = '#4D89E520';
//"#206CDF"
const tealHighLight = '#00FFD9';
const tealHighLight20 = '#00FFD920';
const blueMxcBrand = '#09006E';
const blueBG = '#090046';
const overlayBG = '#0C027060';
const white = '#F9FAFC';
const linkTextColor = '#216CDF';

const theme = createMuiTheme({
    palette: {
      primary: { main: tealHighLight, secondary: tealHighLight20 }, 
      secondary: { main: blueMxcBrand, secondary: overlayBG }, 
      darkBG: { main: blueBG }, 
      textPrimary: { main: white }, 
      textSecondary: { main: linkTextColor } 
    },
    MuiListItemIcon: {
      root: {
        color: white
      }
    },
    //tab 
    MuiTypography: {
      root: {
        color: white
      },
    },
    typography: {
      subheading: {
        color: white
      },
      title: {
        color: white
      },
      fontFamily: [
        'Montserrat',
      ].join(','),
    },
    overrides: {
      MuiTypography: {
        root: {
          color: white,
        },
        body1: {
          color: white,
        },
        body2: {
          color: white,
          fontSize: '0.7rem'
        },
        colorTextSecondary: {
          color: white,
        },
        overline: {
          color: white
        },
        subtitle1: {
          color: white
        },        
        subtitle2: {
          color: white
        },
        caption: {
          color: white
        },
      },
      MuiListItem: {
        root: {
          padding: '5px 0',
        }
      },
      MuiInput: {
        root: {
          color: white,
        },
        underline: {
          "&:before": {
            borderBottom: `1px solid #F9FAFC`
          },
          "&:hover": {
            borderBottom: `1px solid #00FFD9`
          }
        },
      },
      MuiInputLabel: {
        root: {
          color: white,
          "&$focused": {
            color: white
          }
        }
      },
      MuiSelect: {
        icon: {
          color: white,
          right: 0,
          position: 'absolute',
          pointerEvents: 'none',
        }
      },
      MuiIconButton: {
        root: {
          color: white,
        }
      },
      MuiInputBase: {
        input: {
          color: white,
        }
      },
      MuiDivider: {
        root: {
          backgroundColor: '#00000040',
          margin: '5px 10px',
        },
        light: {
          backgroundColor: '#FFFFFF40',
        }
      },
      MuiTableCell: {
        head: {
          color: white,
          fontWeight: 'bold'
        },
        body: {
          color: white,
        },
      },
      MuiPaper: {
        root: {
          backgroundColor: overlayBG,
          padding: 10,
        }
      },
      MuiTablePagination: {
        root: {
          color: white,
        }
      },
      MuiButton: { 
        root: {
          background: tealHighLight,
          color: blueMxcBrand,
          width: 160,
          height: 50,
          fontWeight: 'bolder',
          marginRight: 5,
          boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
          "&:hover": {
            backgroundColor: "#00CCAE",
            fontSize: '1rem'
          },
        },
        outlined: {
          backgroundColor: 'transparent',
          color: tealHighLight,
          //padding: 30,
          fontWeight: 900,
          lineHeight: 1.5,
          borderWidth: 2,
          borderColor: tealHighLight,
          "&:hover": {
            backgroundColor: tealHighLight20,
            borderColor: "#00CCAE",
            color: "#00CCAE",
          },
        },
        /*         link: {
          color: tealHighLight,
          //padding: 30,
          fontWeight: 900,
          lineHeight: 1.5,
          "&:hover": {
            color: "#206CDF",
          },
        }, */
        text: { 
          color: blueMxcBrand, 
          marginBottom: 24,
        },
        textPrimary: {
          color: blueMxcBrand,
          "&:hover": {
            backgroundColor: "#206CDF",
          },
        },
      },
      MuiFormControlLabel: {
        root: { 
          color: white, 
        },
      },
      MuiFormLabel: { 
        root: { 
          color: white, 
        },
      },
      MuiFormHelperText: { 
        root: { 
          color: white, 
        },
      },
      MuiPrivateTabScrollButton:{
        root: {
          width: 0
        }
      },
      MuiTab: {
        root: {
          textColor: white,
        },
        textColorPrimary: {
          color: white
        },
        label: {
          color: white,
        },
        textColorInherit: {
          color: white,
        }
      },
      MuiSvgIcon: {
        root: {
          fill: '#F9FAFC80',
        },
      },
      MuiDialog: {
        color: white,
        root: {
          color: white,
          boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
        },
        paper:{
          backgroundColor: '#191980'
        }
      },
      MuiMenu: {
        paper: {
          color: white,
          backgroundColor: blueBG,
          marginTop: '50px'
        }
      },
      MuiMenuItem: {
        root: {
          "&:hover": {
            backgroundColor: "#1a2d6e",
          },
        }
      }
    },
});
  
export default theme;
