import React from 'react'

export default class Post extends React.Component {
  //like getInitialState()
  constructor(props) {
    super(props)
    this.state = {
        error: {
          all: '',
          firstName: '',
          lastName: '',
          email: '',
          pass: ''
        }
    }

    this.clickBtnEvt = this.clickBtnEvt.bind(this)
  }

  //Click post button
  //For validator and call passed method
  clickBtnEvt(e) {
    console.log("[Post]:clickBtnEvt(e)")
    //check values of form elements
    if (this.refs.postFN.value != "" && this.refs.postLN.value != "" &&
      this.refs.postEM.value != "" && this.refs.postPW.value != ""){

      //reset error
      this.setState({
        error: {all: "", firstName: "", lastName: "", email: "", pass: ""}
      })

      //send data
      //let sendData = [this.refs.postFN.value, this.refs.postLN.value, this.refs.postEM.value, this.refs.postPW.value]
      let sendData = {
        firstName: this.refs.postFN.value,
        lastName:  this.refs.postLN.value,
        email:     this.refs.postEM.value,
        password:  this.refs.postPW.value
      }

      //call event for post btn click
      this.props.btn.call(this, sendData)
    }else{
      //error
      this.setState({
        error: {all: "it requires all input", firstName: "", lastName: "", email: "", pass: ""}
      })

    }
  }

  render() {
    return (
        <div className="col-md-6 col-sm-6 col-xs-12">
          <div className="panel panel-default">
            <div className="panel-heading">POST API</div>
            <div className="panel-body">
              <form role="form">
                <span className="form-err">{this.state.error.all}</span>
                <div className="form-group">
                  <label>First Name</label>
                  <input id="postFN" className="form-control" type="text" ref="postFN" />
                  <span className="form-err">{this.state.error.firstName}</span>
                </div>
                <div className="form-group">
                  <label>Last Name</label>
                  <input id="postLN" className="form-control" type="text" ref="postLN" />
                  <span className="form-err">{this.state.error.lastName}</span>
                </div>
                <div className="form-group">
                  <label>E-mail</label>
                  <input id="postEM" className="form-control" type="text" autoComplete="off" ref="postEM" />
                  <span className="form-err">{this.state.error.email}</span>
                </div>
                <div className="form-group">
                  <label>Password</label>
                  <input id="postPW" className="form-control" type="password" autoComplete="off" ref="postPW" />
                  <span className="form-err">{this.state.error.pass}</span>
                </div>
                <button id="postBtn" type="button" onClick={this.clickBtnEvt} className="btn btn-info">Register</button>
              </form>
            </div>
          </div>
        </div>
    )
  }
}





