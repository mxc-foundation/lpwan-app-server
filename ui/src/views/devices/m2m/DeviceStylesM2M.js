import theme from "../../../theme";

const DeviceStylesM2M = {
  title: {
    color: '#FFFFFF',
    fontSize: 14,
    padding: 6,
  },
  pos: {
    marginBottom: 12,
    color: '#FFFFFF',
    textAlign: 'right',
  },
  between: {
    display: 'flex',
    justifyContent:'spaceBetween'
  },
  flex: {
    display: 'flex',
    flexDirection: 'column'
  },
  navText: {
    fontSize: '0.85rem !important',
  },
  TitleBar: {
    width: '100%',
    light: true,
    display: 'flex',
    justifyContent:'spaceBetween',
    padding: '0px 0px 50px 0px'
  },
  divider: {
    padding: 0,
    color: '#FFFFFF',
    width: '100%',
  },
  padding: {
    padding: 0,
  },
  link: {
    textDecoration: "none",
    fontWeight: "bold",
    fontSize: '1rem',
    color: theme.palette.textSecondary.main,
    opacity: 0.7,
      "&:hover": {
        opacity: 1,
      }
  },
  subTitle:{
    textDecoration: "none",
    paddingTop: 20,
    fontWeight: "bold",
    fontSize: 12,
    color: theme.palette.textPrimary.main,
    cursor: "pointer",
    opacity: 0.7,
      "&:hover": {
        opacity: 1,
      } 
    }       
};
  
export default DeviceStylesM2M;
