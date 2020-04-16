import theme from "../../theme";

const Styles = {
  num: {
    width: '7vw',
    height: '7vw',
    margin: '2%',
    borderRadius: 5,
    fontSize: '7vw',
    textAlign: 'center'
  },
  numLayout: {
    display: 'flex',
    justifyContent: 'between'
  },
  numLayout: {
    display: 'flex',
    justifyContent: 'center'
  },
  card: {
    minWidth: 180,
    width: 220,
    backgroundColor: "#0C0270",
  },
  title: {
    color: '#FFFFFF',
    fontSize: 14,
    padding: 6,
  },
  balance: {
    fontSize: 24,
    color: '#FFFFFF',
    textAlign: 'center',
  },
  newBalance: {
    fontSize: 24,
    textAlign: 'center',
    color: theme.palette.primary.main,
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
    fontSize: 14,
  },
  TitleBar: {
    height: 115,
    width: '50%',
    light: true,
    display: 'flex',
    flexDirection: 'column'
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
    fontSize: 12,
    color: theme.palette.textSecondary.main,
    opacity: 0.7,
      "&:hover": {
        opacity: 1,
      }
  },
  s_input: {
    width: "100%",
    height: "calc(1.5em + 0.9rem + 2px)",
    padding: "0.45rem 0.9rem",
    fontSize: "0.9rem",
    fontWeight: "400",
    lineHeight: "1.5",
    color: "#6c757d",
    backgroundColor: "#fff",
    backgroundClip: "padding-box",
    border: "1px solid #ced4da",
    borderRadius: "0.2rem"
  },
  t_input: {
    width: "100%",
    height: "calc(1.5em + 0.9rem + 2px)",
    padding: "0.45rem 0.9rem",
    fontSize: "1.5rem",
    fontWeight: "400",
    lineHeight: "1.5",
    //color: "#6c757d",
    backgroundColor: "#ebeff2",
    border: "0px",
  },
  arrowLeft: {
    //position: 'absolute',
    top: '50%',
    width: '3vmin',
    height: '3vmin',
    margin: 10,
    background: 'transparent',
    borderTop: '1vmin solid white',
    borderRight: '1vmin solid white',
    boxShadow: '0 0 0 lightgray',
    transition: 'all 200ms ease',
    left: 0,
    transform: 'translate3d(0,-50%,0) rotate(-135deg)',
  
    "&.right": {
      right: 0,
      transform: 'translate3d(0,-50%,0) rotate(45deg)',
    },
    
    "&:hover": {
      borderColor: 'orange',
      boxShadow: '0.5vmin -0.5vmin 0 white',
    },
    
    "&:before": { 
      content: '',
      //position: 'absolute',
      top: '50%',
      left: '50%',
      transform: 'translate(-40%,-60%) rotate(45deg)',
      width: '200%',
      height: '200%',
    }
  },
  arrowRight: {
    //position: 'absolute',
    top: '50%',
    width: '3vmin',
    height: '3vmin',
    margin: 10,
    background: 'transparent',
    borderTop: '1vmin solid white',
    borderRight: '1vmin solid white',
    boxShadow: '0 0 0 lightgray',
    transition: 'all 200ms ease',
    left: 0,
    transform: 'translate3d(0,-50%,0) rotate(-135deg)',
    right: 0,
    transform: 'translate3d(0,-50%,0) rotate(45deg)',
    
    "&:hover": {
      borderColor: 'orange',
      boxShadow: '0.5vmin -0.5vmin 0 white',
    },
    
    "&:before": { 
      content: '',
      //position: 'absolute',
      top: '50%',
      left: '50%',
      transform: 'translate(-40%,-60%) rotate(45deg)',
      width: '200%',
      height: '200%',
    }
  },
};
  
export default Styles;
