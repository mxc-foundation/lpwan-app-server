import ReCAPTCHA from "react-google-recaptcha";
 
function onChange(value) {
  console.log("Captcha value:", value);
}
 
ReactDOM.render(
  <ReCAPTCHA
    sitekey="6LfUDLgUAAAAAAiBf12rhXVWmdW4gZojyCHNK7Oa"
    onChange={onChange}
  />,
  document.body
);