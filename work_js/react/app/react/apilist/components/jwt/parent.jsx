import React     from 'react'

import JwtInput  from './input.jsx'
import JwtCode   from './code.jsx'


export default class JwtParent extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
        inputEM: '',
        inputPW: '',
        error: {
          email: '',
          pass: ''
        },
        code: ''
    }

    this.getBtnEvt = this.getBtnEvt.bind(this)
    this.emailChange = this.emailChange.bind(this)
    this.pwChange = this.pwChange.bind(this)
  }

  emailChange(e) {
    console.log("jwt:email change / ", e.target.value)
    this.setState({
      inputEM: e.target.value
    })
  }

  pwChange(e) {
    console.log("jwt:pw change / ", e.target.value)
    this.setState({
      inputPW: e.target.value
    })
  }

  // click btn
  getBtnEvt(e) {
    console.log("jwt:get")
    //check input value
    if (this.state.inputEM == "" || this.state.inputPW == ""){
        let em, pw
        //TODO:comment out now
        //swal("warning!", "blank filed is not allowed.", "warning")
        if (this.state.inputEM == "") em = 'email is invalid'
        if (this.state.inputPW == "") pw = 'password is invalid'

        this.setState({
          error: {email: em, pass: pw}
        })
    }else{
        //TODO:send data by Ajax
        let sendData = new Object()
        sendData.inputEmail = this.state.inputEM
        sendData.inputPassword = this.state.inputPW

        let url = '/json/userIDs.json'

        //Only this API can Access without jwt
        //hy.sendAjax(url, "post", "form", sendData) //This can't use
        $.ajax({
          url: url,
          dataType: 'json',
          cache: false,
          success: (data) => {
            this.setState({
              code: '12345'
            })
            this.setState({ids: data.ids})
          },
          error: (xhr, status, err) => {
            console.error(url, status, err.toString())
          }
        })
    }

    console.log("email:", this.state.inputEM)
    console.log("pw:", this.state.inputPW)

  }

  render() {
    return (
      <div className='row'>
        <JwtInput btn={this.getBtnEvt} em={this.emailChange} pw={this.pwChange} error={this.state.error} />
        <JwtCode code={this.state.code} />
      </div>
    )
  }
}


/*
(function (){
  var getBtn = document.getElementById("jwtBtn");
  jwtBtn.addEventListener("click", getJWT, false);

  //
  function getJWT(evt){
    //
    var url = "/api/jwt";
    var sendData = new Object();

    //create data
    var errFlg = 0;
    //create data
    if (document.getElementById("jwtEM").value != ""){
      sendData.inputEmail = document.getElementById("jwtEM").value;
    }else{
        errFlg=1;
    }
    if (document.getElementById("jwtPW").value != ""){
      sendData.inputPassword = document.getElementById("jwtPW").value;
    }else{
        errFlg=1;
    }

    if (errFlg==1){
        swal("warning!", "blank filed is not allowed.", "warning");
        return;
    }

    //send
    hy.sendAjax(url, "post", "form", sendData)
  }

})();
*/
