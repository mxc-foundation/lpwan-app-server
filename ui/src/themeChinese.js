import { createMuiTheme } from "@material-ui/core/styles";
//import { teal } from "@material-ui/core/colors";
import MicrosoftYahei from './fonts/Microsoft Yahei.ttf';
const microsoftYahei = {
    fontFamily: 'Microsoft YaHei', 
    fontStyle: 'normal',
    fontDisplay: 'swap',
    fontWeight: 400,
    src: `
    local('Microsoft YaHei'),
    url(${MicrosoftYahei}) format('ttf')
  `,
};

const tealHighLight = '#00FFD9';
const tealHighLight20 = '#00FFD920';
const blueMxcBrand = '#09006E';
const blueMxcBrand20 = '#09006E20';
const blueBG = '#070033';
const overlayBG = '#0C027060';
const white = 'white';
const dark = '#323a46';
//const linkTextColor = '#BBE9E8';

const themeChinese = createMuiTheme({
    palette: {
        primary: { main: blueMxcBrand, secondary: blueMxcBrand20 },
        secondary: { main: blueMxcBrand, secondary: overlayBG },
        darkBG: { main: blueBG },
        textPrimary: { main: dark },
        textSecondary: { main: blueMxcBrand },
        background: {
            paper: "#fff",
            default: "#ebeff2"
        }
    },
    MuiListItemIcon: {
        root: {
            color: dark
        }
    },
    //tab 
    MuiTypography: {
        root: {
            color: dark,
        },
        body1: {
            color: dark,
        },
        colorTextSecondary: {
            color: dark,
        },
    },
    typography: {
        //useNextVariants: true,
        subheading: {
            color: dark,
            "&:hover": {
                color: 'dark',
            },
        },
        title: {
            color: dark
        },
        fontFamily: [
            'Microsoft YaHei',
        ].join(','),
    },
    overrides: {
        MuiCssBaseline: {
            '@global': {
                '@font-face': [microsoftYahei],
            },
        },
        MuiTypography: {
            root: {
                color: dark,
                fontFamily: [
                    'Microsoft YaHei',
                ].join(','),
            },
            body1: {
                color: dark,
                fontSize: '0.8rem',
                fontFamily: [
                    'Microsoft YaHei',
                ].join(','),
            },
            body2: {
                color: dark,
                fontSize: '0.7rem'
            },
            colorTextSecondary: {
                color: dark,
            },
            headline: {
                color: dark,
                fontFamily: [
                    'Microsoft YaHei',
                ].join(','),
            },
            caption: {
                color: dark,
                fontFamily: [
                    'Microsoft YaHei',
                ].join(','),
            },
            
        },
        MuiInput: {
            root: {
                color: dark
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
        MuiAppBar: {
            root: {
                //width: '1024px',
                color: dark
            },
            positionFixed: {
                left: 'inherit',
                right: 'inherit'
            }
        },
        MuiSelect: {
            icon: {
                color: dark,
                right: 0,
                position: 'absolute',
                pointerEvents: 'none',
            }
        },
        MuiIconButton: {
            root: {
                color: dark,
            }
        },
        /*       MuiInputBase: {
                input: {
                  color: '#F9FAFC',
                  fontWeight: "bolder",
                  "&:-webkit-autofill": {
                    WebkitBoxShadow: "0 0 0 1000px #F9FAFC inset"
                  }
                }
              }, */
        MuiDivider: {
            root: {
                backgroundColor: '#00000040',
                margin: '5px 0px 5px 0px',
            },
            light: {
                backgroundColor: '#FFFFFF50',
            }
        },
        MuiTable: {
            root: {
                background: 'transparent',
                //minWidth: 840,
            }
        },
        MuiTableCell: {
            head: {
                color: dark,
                fontWeight: '800',
                fontSize: '1em',
                padding: 10,
            },
            body: {
                background: 'none',
                color: dark,
                //maxWidth: 140,
                whiteSpace: 'nowrap',
                //overflow: 'hidden',
                textOverflow: 'ellipsis',
                fontWeight: '400',
            },
            root: {
                padding: '4px 5px',
                //maxWidth: 140,
                whiteSpace: 'nowrap',
                //overflow: 'hidden',
                textOverflow: 'ellipsis',
                borderBottom: 'solid 1px #070033',
                lineHeight: '40px',
                textAlign: 'left',
            }
        },
        MuiPaper: {
            root: {
                backgroundColor: white,
                padding: 10,
            }
        },
        MuiTablePagination: {
            root: {
                color: dark,
                background: 'none',
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
                },
            },
            outlined: {
                backgroundColor: 'transparent',
                color: tealHighLight,
                //padding: 30,
                fontWeight: 900,
                lineHeight: 1.5,
                borderWidth: 2,
                borderColor: white,
                "&:hover": {
                    backgroundColor: tealHighLight20,
                    borderColor: "#00CCAE",
                    color: dark,
                },
            },
            colorInherit: {
                color: dark,
                "&:hover": {
                    borderColor: white,
                    color: dark,
                },
            },
        },
        MuiFormLabel: {
            root: {
                color: dark,
            },
        },
        MuiFormHelperText: {
            root: {
                color: dark,
            },
        },
        MuiPrivateTabScrollButton: {
            root: {
                width: 0
            }
        },
        MuiTab: {
            root: {
                color: dark,
            },
            textColorPrimary: {
                color: dark
            },
        },
        MuiSvgIcon: {
            root: {
                // fill: '#F9FAFC80',
            },
        },
        MuiDialog: {
            color: dark,
            root: {
                color: dark,
                boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
            },
        },
        MuiMenu: {
            paper: {
                backgroundColor: blueBG,
                marginTop: '50px',
                color: dark
            }
        }
    },
});

export default themeChinese;