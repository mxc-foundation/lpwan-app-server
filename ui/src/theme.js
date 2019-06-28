import { createMuiTheme } from "@material-ui/core/styles";
import teal from "@material-ui/core/colors/teal";


const theme = createMuiTheme({
    palette: {
      primary: { main: teal['A200'] },
      secondary: { main: '#11cb5f' }, 
      textPrimary: {main: '#FFFFFF'},
      textSecondary: {main: '#FFFFFF'}
    },
    MuiListItemIcon: {
      root: {
        color: '#FFFFFF'
      }
    },
    //tab 
    MuiTypography: {
      root: {
        color: '#FFFFFF'
      },
    },
    typography: {
      subheading: {
        color: '#FFFFFF'
      },
      title: {
        color: '#FFFFFF'
      },
      fontFamily: [
        'Montserrat',
      ].join(','),
    },
    overrides: {
      MuiInput: {
        root: {
          color: '#FFFFFF'
        },
        underline: {
          "&:before": {
            borderBottom: `1px solid white`
          }
        }
      },
      MuiSelect: {
        icon: {
          color: '#FFFFFF',
          right: 0,
          position: 'absolute',
          pointerEvents: 'none',
        }
      },
      MuiIconButton: {
        root: {
          color: '#FFFFFF',
        }
      },
      MuiInputBase: {
        input: {
          color: '#FFFFFF',
        }
      },
      MuiTable: {
        root: {
          background: '#0C0270',
        }
      },
      MuiDivider: {
        root: {
          backgroundColor: '#FFFFFF50',
        },
        light: {
          backgroundColor: '#FFFFFF50',
        }
      },
      MuiTableCell: {
        head: {
          background: '#0C0270',
          color: 'white',
          fontWeight: 'bold'
        }
      },
      MuiPaper: {
        root: {
          backgroundColor: '#0C0270',
          padding: 10,
        }
      },
      MuiTablePagination: {
        root: {
          color: 'white',
          background: '#0C0270',
        }
      },
      MuiButton: { 
        root: {
          //background: teal['A200'],
          background: '#311b92',
          width: 135,
          height: 50,
          fontWeight: 'bold',
          marginRight: 5,
        },
        text: { 
          color: 'white', 
        },
      },
      MuiFormLabel: { 
        root: { 
          color: 'white', 
        },
      },
      MuiFormHelperText: { 
        root: { 
          color: 'white', 
        },
      },
      MuiPrivateTabScrollButton:{
        root: {
          width: 0
        }
      },
      MuiTab: {
        root: {
          color: 'white',
        },
        textColorPrimary: {
          color: 'white'
        }
      },
    },
});
  
export default theme;
