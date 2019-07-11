import { createMuiTheme } from "@material-ui/core/styles";

const blueHighLight = '#4D89E5';
const blueHighLight40 = '#4D89E540';
const blueMxcBrand = '#09006E';
const white = '#F9FAFC';
const linkTextColor = '#216CDF';

const theme = createMuiTheme({
    palette: {
      primary: { main: blueHighLight, secondary: blueHighLight40 }, 
      secondary: { main: blueMxcBrand }, 
      textPrimary: {main: white}, 
      textSecondary: {main: linkTextColor} 
    },
    MuiListItemIcon: {
      root: {
        color: white,
      }
    },
    //tab 
    
    typography: {
      subheading: {
        color: white,
      },
      title: {
        color: white,
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
        colorTextSecondary: {
          color: white,
        },
      },
      MuiInput: {
        root: {
          color: white,
        },
        underline: {
          "&:before": {
            borderBottom: `1px solid #F9FAFC`
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
      MuiTable: {
        root: {
          background: blueMxcBrand,
        }
      },
      MuiDivider: {
        root: {
          backgroundColor: '#FFFFFF50',
          margin: 15,
        },
        light: {
          backgroundColor: '#FFFFFF50',
        }
      },
      MuiTableCell: {
        head: {
          background: blueMxcBrand,
          color: white,
          fontWeight: 'bold'
        },
        body: {
          color: white,
        },
      },
      MuiPaper: {
        root: {
          backgroundColor: blueMxcBrand,
          padding: 10,
        }
      },
      MuiTablePagination: {
        root: {
          color: white,
          background: blueMxcBrand,
        }
      },
      MuiButton: { 
        root: {
          background: blueHighLight,
          color: blueMxcBrand,
          width: 135,
          height: 50,
          fontWeight: 'bolder',
          marginRight: 5,
          boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
          "&:hover": {
            backgroundColor: "#206CDF",
          },
        outline: {
          backgroundColor: blueMxcBrand,
          color: blueMxcBrand,
        },
        },
        text: { 
          color: blueMxcBrand, 
        },
        textPrimary: {
          color: blueMxcBrand,
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
      },
      MuiSvgIcon: {
        root: {
          fill: white,
        },
      },
      MuiDialog: {
        color: white,
        root: {
          color: white,
          boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
        },
      },
    },
});
  
export default theme;
