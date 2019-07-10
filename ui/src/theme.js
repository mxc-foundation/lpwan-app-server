import { createMuiTheme } from "@material-ui/core/styles";

const blueHighLight = '#4D89E5';
const blueMxcBrand = '#09006E';
const white = '#F9FAFC';
const linkTextColor = '#CAFCF5';

const theme = createMuiTheme({
    palette: {
      primary: { main: blueHighLight }, 
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
    MuiTypography: {
      root: {
        color: white
      },
    },
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
        }
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
          width: 135,
          height: 50,
          fontWeight: 'bolder',
          marginRight: 5,
          boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
        },
        text: { 
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
          color: white,
        },
        textColorPrimary: {
          color: white
        }
      },
      MuiSvgIcon: {
        root: {
          fill: white,
        },
      },
    },
});
  
export default theme;
