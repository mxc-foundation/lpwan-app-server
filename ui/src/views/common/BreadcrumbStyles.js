import theme from "../../theme";

// Reference: https://material-ui.com/customization/breakpoints/#breakpoints
const breadcrumbStyles = {
  [theme.breakpoints.down('xs')]: {
    breadcrumb: {
      fontSize: "0.8rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.between('sm', 'md')]: {
    breadcrumb: {
      fontSize: "1.1rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  [theme.breakpoints.up('md')]: {
    breadcrumb: {
      fontSize: "1.25rem",
      margin: "0rem",
      padding: "0rem"
    },
  },
  breadcrumbItem: {
    color: "#666 !important"
  },
  breadcrumbItemLink: {
    color: "#71b6f9 !important"
  },
};

export default breadcrumbStyles;
