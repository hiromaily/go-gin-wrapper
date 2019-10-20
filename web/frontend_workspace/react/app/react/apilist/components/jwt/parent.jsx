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

  //change value of email form
  emailChange(e) {
    //console.log("jwt:email change / ", e.target.value)
    this.setState({
      inputEM: e.target.value
    })
  }

  //change value of password form
  pwChange(e) {
    //console.log("jwt:pw change / ", e.target.value)
    this.setState({
      inputPW: e.target.value
    })
  }

  // click btn
  getBtnEvt(e) {
    console.log("[JwtParent]:getBtnEvt()")
    //check input value
    if (this.state.inputEM == "" || this.state.inputPW == ""){
        let em, pw
        swal("warning!", "blank filed is not allowed.", "warning")
        if (this.state.inputEM == "") em = 'email is invalid'
        if (this.state.inputPW == "") pw = 'password is invalid'

        this.setState({
          error: {email: em, pass: pw}
        })
    }else{
        let sendData = new Object()
        sendData.inputEmail = this.state.inputEM
        sendData.inputPassword = this.state.inputPW

        let that = this
        let url = '/api/jwt'
        let method = 'post'
        let contentType = "application/x-www-form-urlencoded"

        //Only this API can Access without jwt
        $.ajax({
          url: encodeURI(url),
          type: method,
          beforeSend: function beforeSend(xhr) {
            xhr.setRequestHeader(hiromaily_header, hiromaily_key)
          },
          //cache    : false,
          crossDomain: false,
          contentType: contentType,
          dataType: 'json', //data type from server
          data: sendData
        }).done(function (data, textStatus, jqXHR) {
          that.setState({
            code: data.token,
            inputEM: '',
            inputPW: ''
          })
          swal("success!", "get jwt code!", "success")
          //call getUserIDs()
          that.props.funcGetUserIDs.call(that)
        }).fail(function (jqXHR, textStatus, errorThrown) {
          swal("error!", "validation error was occurred!", "error")
        })
    }
  }

  render() {
    return (
      <div className='row'>
        <JwtInput btn={this.getBtnEvt} em={this.emailChange} pw={this.pwChange} inputEM={this.state.inputEM}
          inputPW={this.state.inputPW} error={this.state.error} />
        <JwtCode code={this.state.code} />
      </div>
    )
  }
}
